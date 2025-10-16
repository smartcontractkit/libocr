package managed

import (
	"context"
	"errors"
	"fmt"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/internal/metricshelper"
	"github.com/RoSpaceDev/libocr/internal/util"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/managed/limits"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr2/protocol"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr2/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/shim"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
	"github.com/prometheus/client_golang/prometheus"
)

// RunManagedOCR2Oracle runs a "managed" version of protocol.RunOracle. It handles
// setting up telemetry, garbage collection, configuration updates, translating
// from commontypes.BinaryNetworkEndpoint to protocol.NetworkEndpoint, and
// creation/teardown of reporting plugins.
func RunManagedOCR2Oracle(
	ctx context.Context,

	v2bootstrappers []commontypes.BootstrapperLocator,
	configTracker types.ContractConfigTracker,
	contractTransmitter types.ContractTransmitter,
	database types.Database,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	monitoringEndpoint commontypes.MonitoringEndpoint,
	netEndpointFactory types.BinaryNetworkEndpointFactory,
	offchainConfigDigester types.OffchainConfigDigester,
	offchainKeyring types.OffchainKeyring,
	onchainKeyring types.OnchainKeyring,
	reportingPluginFactory types.ReportingPluginFactory,
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

	subs.Go(func() {
		collectGarbage(ctx, database, localConfig, logger)
	})

	runWithContractConfig(
		ctx,

		configTracker,
		database,
		func(ctx context.Context, logger loghelper.LoggerWithContext, contractConfig types.ContractConfig) (err error, retry bool) {
			skipResourceExhaustionChecks := localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode

			fromAccount, err := contractTransmitter.FromAccount(ctx)
			if err != nil {
				logger.Error("ManagedOCR2Oracle: error getting FromAccount", commontypes.LogFields{
					"error": err,
				})
				return
			}

			sharedConfig, oid, err := ocr2config.SharedConfigFromContractConfig(
				skipResourceExhaustionChecks,
				contractConfig,
				offchainKeyring,
				onchainKeyring,
				netEndpointFactory.PeerID(),
				fromAccount,
			)
			if err != nil {
				logger.Error("ManagedOCR2Oracle: error while updating config", commontypes.LogFields{
					"error": err,
				})
				return
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
				maxDurationInitialization+protocol.ReportingPluginTimeoutWarningGracePeriod,
				func() {
					logger.Error("ManagedOCR2Oracle: ReportingPluginFactory.NewReportingPlugin is taking too long", commontypes.LogFields{
						"maxDuration": maxDurationInitialization,
					})
				},
			)

			reportingPlugin, reportingPluginInfo, err := reportingPluginFactory.NewReportingPlugin(initCtx, types.ReportingPluginConfig{
				sharedConfig.ConfigDigest,
				oid,
				sharedConfig.N(),
				sharedConfig.F,
				sharedConfig.OnchainConfig,
				sharedConfig.ReportingPluginConfig,
				sharedConfig.DeltaRound,
				sharedConfig.MaxDurationQuery,
				sharedConfig.MaxDurationObservation,
				sharedConfig.MaxDurationReport,
				sharedConfig.MaxDurationShouldAcceptFinalizedReport,
				sharedConfig.MaxDurationShouldTransmitAcceptedReport,
			})

			ins.Stop()

			if err != nil {
				logger.Error("ManagedOCR2Oracle: error during NewReportingPlugin()", commontypes.LogFields{
					"error": err,
				})
				return
			}
			defer loghelper.CloseLogError(
				reportingPlugin,
				logger,
				"ManagedOCR2Oracle: error during reportingPlugin.Close()",
			)
			if err := validateReportingPluginLimits(reportingPluginInfo.Limits); err != nil {
				logger.Error("ManagedOCR2Oracle: invalid ReportingPluginInfo", commontypes.LogFields{
					"error":               err,
					"reportingPluginInfo": reportingPluginInfo,
				})
				return fmt.Errorf("ManagedOCR2Oracle: invalid ReportingPluginInfo"), false
			}

			maxSigLen := onchainKeyring.MaxSignatureLength()
			lims, err := limits.OCR2Limits(sharedConfig.PublicConfig, reportingPluginInfo.Limits, maxSigLen)
			if err != nil {
				logger.Error("ManagedOCR2Oracle: error during limits", commontypes.LogFields{
					"error":               err,
					"publicConfig":        sharedConfig.PublicConfig,
					"reportingPluginInfo": reportingPluginInfo,
					"maxSigLen":           maxSigLen,
				})
				return fmt.Errorf("ManagedOCR2Oracle: error during limits"), false
			}
			binNetEndpoint, err := netEndpointFactory.NewEndpoint(
				sharedConfig.ConfigDigest,
				peerIDs,
				v2bootstrappers,
				sharedConfig.F,
				lims,
			)
			if err != nil {
				logger.Error("ManagedOCR2Oracle: error during NewEndpoint", commontypes.LogFields{
					"error":           err,
					"peerIDs":         peerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return fmt.Errorf("ManagedOCR2Oracle: error during NewEndpoint"), true
			}

			// No need to binNetEndpoint.Start/Close since netEndpoint will handle that for us

			netEndpoint := shim.NewOCR2SerializingEndpoint(
				chTelemetrySend,
				sharedConfig.ConfigDigest,
				binNetEndpoint,
				childLogger,
				registerer,
				reportingPluginInfo.Limits,
			)
			if err := netEndpoint.Start(); err != nil {
				return fmt.Errorf("ManagedOCR2Oracle: error during netEndpoint.Start(): %w", err), true
			}
			defer loghelper.CloseLogError(
				netEndpoint,
				logger,
				"ManagedOCR2Oracle: error during netEndpoint.Close()",
			)

			var reportQuorum int
			if reportingPluginInfo.UniqueReports {
				// We require greater than (n+f)/2 signatures to reach a byzantine
				// quorum. This ensures unique reports since each honest node will sign
				// at most one report for any given (epoch, round).
				//
				// Argument:
				//
				// (n+f)/2 = ((n-f)+f+f)/2 = (n-f)/2 + f
				//
				// There are (n-f) honest nodes, so to get two reports for an (epoch,
				// round) to reach  quorum, we'd need an honest node to sign two reports
				// which contradicts the assumption that an honest node will sign at
				// most one report for any given (epoch, round).
				reportQuorum = (sharedConfig.N()+sharedConfig.F)/2 + 1
			} else {
				reportQuorum = sharedConfig.F + 1
			}

			protocol.RunOracle(
				ctx,
				sharedConfig,
				contractTransmitter,
				database,
				oid,
				localConfig,
				childLogger,
				registerer,
				netEndpoint,
				offchainKeyring,
				onchainKeyring,
				shim.LimitCheckReportingPlugin{reportingPlugin, reportingPluginInfo.Limits},
				reportQuorum,
				shim.NewOCR2TelemetrySender(chTelemetrySend, childLogger, localConfig.EnableTransmissionTelemetry),
			)

			return nil, false
		},
		localConfig,
		logger,
		offchainConfigDigester,
		defaultRetryParams(),
	)
}

func validateReportingPluginLimits(limits types.ReportingPluginLimits) error {
	var err error
	if !(0 <= limits.MaxQueryLength && limits.MaxQueryLength <= types.MaxMaxQueryLength) {
		err = errors.Join(err, fmt.Errorf("MaxQueryLength (%v) out of range. Should be between 0 and %v", limits.MaxQueryLength, types.MaxMaxQueryLength))
	}
	if !(0 <= limits.MaxObservationLength && limits.MaxObservationLength <= types.MaxMaxObservationLength) {
		err = errors.Join(err, fmt.Errorf("MaxObservationLength (%v) out of range. Should be between 0 and %v", limits.MaxObservationLength, types.MaxMaxObservationLength))
	}
	if !(0 <= limits.MaxReportLength && limits.MaxReportLength <= types.MaxMaxReportLength) {
		err = errors.Join(err, fmt.Errorf("MaxReportLength (%v) out of range. Should be between 0 and %v", limits.MaxReportLength, types.MaxMaxReportLength))
	}
	return err
}
