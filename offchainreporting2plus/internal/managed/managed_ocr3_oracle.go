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
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3/protocol"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/shim"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
	"github.com/prometheus/client_golang/prometheus"
)

// RunManagedOCR3Oracle runs a "managed" version of protocol.RunOracle. It handles
// setting up telemetry, garbage collection, configuration updates, translating
// from commontypes.BinaryNetworkEndpoint to protocol.NetworkEndpoint, and
// creation/teardown of reporting plugins.
func RunManagedOCR3Oracle[RI any](
	ctx context.Context,

	v2bootstrappers []commontypes.BootstrapperLocator,
	configTracker types.ContractConfigTracker,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	database ocr3types.Database,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	monitoringEndpoint commontypes.MonitoringEndpoint,
	netEndpointFactory types.BinaryNetworkEndpointFactory,
	offchainConfigDigester types.OffchainConfigDigester,
	offchainKeyring types.OffchainKeyring,
	onchainKeyring ocr3types.OnchainKeyring[RI],
	reportingPluginFactory ocr3types.ReportingPluginFactory[RI],
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
				return fmt.Errorf("ManagedOCR3Oracle: error getting FromAccount: %w", err), true
			}

			sharedConfig, oid, err := ocr3config.SharedConfigFromContractConfig(
				skipResourceExhaustionChecks,
				contractConfig,
				offchainKeyring,
				onchainKeyring,
				netEndpointFactory.PeerID(),
				fromAccount,
			)
			if err != nil {
				return fmt.Errorf("ManagedOCR3Oracle: error while decoding ContractConfig: %w", err), false
			}

			registerer := prometheus.WrapRegistererWith(
				prometheus.Labels{
					// disambiguate different protocol instances by configDigest
					"config_digest": sharedConfig.ConfigDigest.String(),
					// disambiguate different oracle instances by offchainPublicKey
					"offchain_public_key": fmt.Sprintf("%x", offchainKeyring.OffchainPublicKey()),
				},
				metricsRegistererWrapper,
			)

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
					logger.Error("ManagedOCR3Oracle: ReportingPluginFactory.NewReportingPlugin is taking too long", commontypes.LogFields{
						"maxDuration": maxDurationInitialization,
					})
				},
			)

			reportingPlugin, reportingPluginInfo, err := reportingPluginFactory.NewReportingPlugin(initCtx, ocr3types.ReportingPluginConfig{
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
			})

			ins.Stop()

			if err != nil {
				return fmt.Errorf("ManagedOCR3Oracle: error during NewReportingPlugin(): %w", err), true
			}
			defer loghelper.CloseLogError(
				reportingPlugin,
				logger,
				"ManagedOCR3Oracle: error during reportingPlugin.Close()",
			)

			if err := validateOCR3ReportingPluginLimits(reportingPluginInfo.Limits); err != nil {
				logger.Error("ManagedOCR3Oracle: invalid ReportingPluginInfo", commontypes.LogFields{
					"error":               err,
					"reportingPluginInfo": reportingPluginInfo,
				})
				return fmt.Errorf("ManagedOCR3Oracle: invalid MercuryPluginInfo"), false
			}

			maxSigLen := onchainKeyring.MaxSignatureLength()
			lims, err := limits.OCR3Limits(sharedConfig.PublicConfig, reportingPluginInfo.Limits, maxSigLen)
			if err != nil {
				logger.Error("ManagedOCR3Oracle: error during limits", commontypes.LogFields{
					"error":                 err,
					"publicConfig":          sharedConfig.PublicConfig,
					"reportingPluginLimits": reportingPluginInfo.Limits,
					"maxSigLen":             maxSigLen,
				})
				return fmt.Errorf("ManagedOCR3Oracle: error during limits"), false
			}
			binNetEndpoint, err := netEndpointFactory.NewEndpoint(
				sharedConfig.ConfigDigest,
				peerIDs,
				v2bootstrappers,
				sharedConfig.F,
				lims,
			)
			if err != nil {
				logger.Error("ManagedOCR3Oracle: error during NewEndpoint", commontypes.LogFields{
					"error":           err,
					"peerIDs":         peerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return fmt.Errorf("ManagedOCR3Oracle: error during NewEndpoint"), true
			}

			// No need to binNetEndpoint.Start/Close since netEndpoint will handle that for us

			netEndpoint := shim.NewOCR3SerializingEndpoint[RI](
				chTelemetrySend,
				sharedConfig.ConfigDigest,
				binNetEndpoint,
				maxSigLen,
				childLogger,
				registerer,
				reportingPluginInfo.Limits,
				sharedConfig.N(),
				sharedConfig.F,
			)
			if err := netEndpoint.Start(); err != nil {
				return fmt.Errorf("ManagedOCR3Oracle: error during netEndpoint.Start(): %w", err), true
			}
			defer loghelper.CloseLogError(
				netEndpoint,
				logger,
				"ManagedOCR3Oracle: error during netEndpoint.Close()",
			)

			protocol.RunOracle[RI](
				ctx,
				sharedConfig,
				contractTransmitter,
				&shim.SerializingOCR3Database{database},
				oid,
				localConfig,
				childLogger,
				registerer,
				netEndpoint,
				offchainKeyring,
				onchainKeyring,
				shim.LimitCheckOCR3ReportingPlugin[RI]{reportingPlugin, reportingPluginInfo.Limits},
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

func validateOCR3ReportingPluginLimits(limits ocr3types.ReportingPluginLimits) error {
	var err error
	if !(0 <= limits.MaxQueryLength && limits.MaxQueryLength <= ocr3types.MaxMaxQueryLength) {
		err = errors.Join(err, fmt.Errorf("MaxQueryLength (%v) out of range. Should be between 0 and %v", limits.MaxQueryLength, ocr3types.MaxMaxQueryLength))
	}
	if !(0 <= limits.MaxObservationLength && limits.MaxObservationLength <= ocr3types.MaxMaxObservationLength) {
		err = errors.Join(err, fmt.Errorf("MaxObservationLength (%v) out of range. Should be between 0 and %v", limits.MaxObservationLength, ocr3types.MaxMaxObservationLength))
	}
	if !(0 <= limits.MaxOutcomeLength && limits.MaxOutcomeLength <= ocr3types.MaxMaxOutcomeLength) {
		err = errors.Join(err, fmt.Errorf("MaxOutcomeLength (%v) out of range. Should be between 0 and %v", limits.MaxOutcomeLength, ocr3types.MaxMaxOutcomeLength))
	}
	if !(0 <= limits.MaxReportLength && limits.MaxReportLength <= ocr3types.MaxMaxReportLength) {
		err = errors.Join(err, fmt.Errorf("MaxReportLength (%v) out of range. Should be between 0 and %v", limits.MaxReportLength, ocr3types.MaxMaxReportLength))
	}
	if !(0 <= limits.MaxReportCount && limits.MaxReportCount <= ocr3types.MaxMaxReportCount) {
		err = errors.Join(err, fmt.Errorf("MaxReportCount (%v) out of range. Should be between 0 and %v", limits.MaxReportCount, ocr3types.MaxMaxReportCount))
	}
	return err
}
