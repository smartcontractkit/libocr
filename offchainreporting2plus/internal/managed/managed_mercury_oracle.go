package managed

import (
	"context"
	"errors"
	"fmt"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/internal/metricshelper"
	"github.com/RoSpaceDev/libocr/internal/util"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/managed/limits"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/mercuryshim"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3/protocol"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/shim"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
	"github.com/prometheus/client_golang/prometheus"
)

// RunManagedMercuryOracle runs a "managed" version of protocol.RunOracle. It handles
// setting up telemetry, garbage collection, configuration updates, translating
// from commontypes.BinaryNetworkEndpoint to protocol.NetworkEndpoint, and
// creation/teardown of reporting plugins.
func RunManagedMercuryOracle(
	ctx context.Context,

	v2bootstrappers []commontypes.BootstrapperLocator,
	configTracker types.ContractConfigTracker,
	contractTransmitter types.ContractTransmitter,
	database ocr3types.Database,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	monitoringEndpoint commontypes.MonitoringEndpoint,
	netEndpointFactory types.BinaryNetworkEndpointFactory,
	offchainConfigDigester types.OffchainConfigDigester,
	offchainKeyring types.OffchainKeyring,
	onchainKeyring types.OnchainKeyring,
	mercuryPluginFactory ocr3types.MercuryPluginFactory,
) {
	subs := subprocesses.Subprocesses{}
	defer subs.Wait()

	var chTelemetrySend chan<- *serialization.TelemetryWrapper
	{
		chTelemetry := make(chan *serialization.TelemetryWrapper, 100)
		chTelemetrySend = chTelemetry
		subs.Go(func() {
			forwardTelemetry(ctx, logger, monitoringEndpoint, chTelemetry)
		})
	}

	metricsRegistererWrapper := metricshelper.NewPrometheusRegistererWrapper(metricsRegisterer, logger)

	runWithContractConfig(
		ctx,

		configTracker,
		database,
		func(ctx context.Context, logger loghelper.LoggerWithContext, contractConfig types.ContractConfig) (err error, retry bool) {
			skipResourceExhaustionChecks := localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode

			fromAccount, err := contractTransmitter.FromAccount(ctx)
			if err != nil {
				return fmt.Errorf("ManagedMercuryOracle: error getting FromAccount: %w", err), true
			}

			ocr3OnchainKeyring := mercuryshim.NewMercuryOCR3OnchainKeyring(onchainKeyring)

			sharedConfig, oid, err := ocr3config.SharedConfigFromContractConfig[mercuryshim.MercuryReportInfo](
				skipResourceExhaustionChecks,
				contractConfig,
				offchainKeyring,
				ocr3OnchainKeyring,
				netEndpointFactory.PeerID(),
				fromAccount,
			)
			if err != nil {
				return fmt.Errorf("ManagedMercuryOracle: error while decoding ContractConfig: %w", err), false
			}

			// Run with new config
			peerIDs := []string{}
			for _, identity := range sharedConfig.OracleIdentities {
				peerIDs = append(peerIDs, identity.PeerID)
			}

			childLogger := logger.MakeChild(commontypes.LogFields{
				"oid": oid,
			})

			maxDurationInitialization := util.NilCoalesce(sharedConfig.MaxDurationInitialization, localConfig.DefaultMaxDurationInitialization)
			initCtx, initCancel := context.WithTimeout(ctx, maxDurationInitialization)
			defer initCancel()

			ins := loghelper.NewIfNotStopped(
				maxDurationInitialization+common.ReportingPluginTimeoutWarningGracePeriod,
				func() {
					logger.Error("ManagedMercuryOracle: MercuryPluginFactory.NewMercuryPlugin is taking too long", commontypes.LogFields{
						"maxDuration": maxDurationInitialization,
					})
				},
			)

			mercuryPlugin, mercuryPluginInfo, err := mercuryPluginFactory.NewMercuryPlugin(initCtx, ocr3types.MercuryPluginConfig{
				sharedConfig.ConfigDigest,
				oid,
				sharedConfig.N(),
				sharedConfig.F,
				sharedConfig.OnchainConfig,
				sharedConfig.ReportingPluginConfig,
				sharedConfig.MinRoundInterval(),
				sharedConfig.MaxDurationObservation,
			})

			ins.Stop()

			if err != nil {
				return fmt.Errorf("ManagedMercuryOracle: error during NewReportingPlugin(): %w", err), true
			}
			defer loghelper.CloseLogError(
				mercuryPlugin,
				logger,
				"ManagedMercuryOracle: error during reportingPlugin.Close()",
			)

			registerer := prometheus.WrapRegistererWith(
				prometheus.Labels{
					// disambiguate different protocol instances by configDigest
					"config_digest": sharedConfig.ConfigDigest.String(),
					// disambiguate different oracle instances by offchainPublicKey
					"offchain_public_key": fmt.Sprintf("%x", offchainKeyring.OffchainPublicKey()),
				},
				metricsRegistererWrapper,
			)

			if err := validateMercuryPluginLimits(mercuryPluginInfo.Limits); err != nil {
				logger.Error("ManagedMercuryOracle: invalid MercuryPluginInfo", commontypes.LogFields{
					"error":             err,
					"mercuryPluginInfo": mercuryPluginInfo,
				})
				return fmt.Errorf("ManagedMercuryOracle: invalid MercuryPluginInfo"), false
			}

			reportingPluginLimits := mercuryshim.ReportingPluginLimits(mercuryPluginInfo.Limits)

			lims, err := limits.OCR3Limits(sharedConfig.PublicConfig, reportingPluginLimits, ocr3OnchainKeyring.MaxSignatureLength())
			if err != nil {
				logger.Error("ManagedMercuryOracle: error during limits", commontypes.LogFields{
					"error":                 err,
					"publicConfig":          sharedConfig.PublicConfig,
					"reportingPluginLimits": reportingPluginLimits,
					"maxSigLen":             ocr3OnchainKeyring.MaxSignatureLength(),
				})
				return fmt.Errorf("ManagedMercuryOracle: error during limits"), false
			}
			binNetEndpoint, err := netEndpointFactory.NewEndpoint(
				sharedConfig.ConfigDigest,
				peerIDs,
				v2bootstrappers,
				sharedConfig.F,
				lims,
			)
			if err != nil {
				logger.Error("ManagedMercuryOracle: error during NewEndpoint", commontypes.LogFields{
					"error":           err,
					"peerIDs":         peerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return fmt.Errorf("ManagedMercuryOracle: error during NewEndpoint"), true
			}

			// No need to binNetEndpoint.Start/Close since netEndpoint will handle that for us

			netEndpoint := shim.NewOCR3SerializingEndpoint[mercuryshim.MercuryReportInfo](
				chTelemetrySend,
				sharedConfig.ConfigDigest,
				binNetEndpoint,
				ocr3OnchainKeyring.MaxSignatureLength(),
				childLogger,
				registerer,
				reportingPluginLimits,
				sharedConfig.N(),
				sharedConfig.F,
			)
			if err := netEndpoint.Start(); err != nil {
				return fmt.Errorf("ManagedMercuryOracle: error during netEndpoint.Start(): %w", err), true
			}
			defer loghelper.CloseLogError(
				netEndpoint,
				logger,
				"ManagedMercuryOracle: error during netEndpoint.Close()",
			)

			reportingPluginConfig := ocr3types.ReportingPluginConfig{
				sharedConfig.ConfigDigest,
				oid,
				sharedConfig.N(),
				sharedConfig.F,
				sharedConfig.OnchainConfig,
				sharedConfig.ReportingPluginConfig,
				sharedConfig.DeltaRound,
				sharedConfig.MaxDurationQuery,
				sharedConfig.MaxDurationObservation,
				sharedConfig.MaxDurationShouldAcceptAttestedReport,
				sharedConfig.MaxDurationShouldTransmitAcceptedReport,
			}
			reportingPlugin := &mercuryshim.MercuryReportingPlugin{
				reportingPluginConfig,
				logger,
				mercuryPlugin,
				mercuryPluginInfo.Limits,
			}

			protocol.RunOracle[mercuryshim.MercuryReportInfo](
				ctx,
				sharedConfig,
				mercuryshim.NewMercuryOCR3ContractTransmitter(contractTransmitter),
				&shim.SerializingOCR3Database{database},
				oid,
				localConfig,
				childLogger,
				registerer,
				netEndpoint,
				offchainKeyring,
				ocr3OnchainKeyring,
				shim.LimitCheckOCR3ReportingPlugin[mercuryshim.MercuryReportInfo]{reportingPlugin, reportingPluginLimits},
				shim.NewOCR3TelemetrySender(chTelemetrySend, childLogger, localConfig.EnableTransmissionTelemetry),
			)

			return nil, false
		},
		localConfig,
		logger,
		offchainConfigDigester,
		defaultRetryParams(),
	)
}

func validateMercuryPluginLimits(limits ocr3types.MercuryPluginLimits) error {
	var err error
	if !(0 <= limits.MaxObservationLength && limits.MaxObservationLength <= ocr3types.MaxMaxMercuryObservationLength) {
		err = errors.Join(err, fmt.Errorf("MaxObservationLength (%v) out of range. Should be between 0 and %v", limits.MaxObservationLength, ocr3types.MaxMaxMercuryObservationLength))
	}
	if !(0 <= limits.MaxReportLength && limits.MaxReportLength <= ocr3types.MaxMaxMercuryReportLength) {
		err = errors.Join(err, fmt.Errorf("MaxReportLength (%v) out of range. Should be between 0 and %v", limits.MaxReportLength, ocr3types.MaxMaxMercuryReportLength))
	}
	return err
}
