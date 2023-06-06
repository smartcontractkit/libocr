package managed

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/mercuryshim"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/shim"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
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

	runWithContractConfig(
		ctx,

		configTracker,
		database,
		func(ctx context.Context, contractConfig types.ContractConfig, logger loghelper.LoggerWithContext) {
			skipResourceExhaustionChecks := localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode

			fromAccount, err := contractTransmitter.FromAccount()
			if err != nil {
				logger.Error("ManagedMercuryOracle: error getting FromAccount", commontypes.LogFields{
					"error": err,
				})
				return
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
				logger.Error("ManagedMercuryOracle: error while updating config", commontypes.LogFields{
					"error": err,
				})
				return
			}

			// Run with new config
			peerIDs := []string{}
			for _, identity := range sharedConfig.OracleIdentities {
				peerIDs = append(peerIDs, identity.PeerID)
			}

			childLogger := logger.MakeChild(commontypes.LogFields{
				"oid": oid,
			})

			mercuryPlugin, mercuryPluginInfo, err := mercuryPluginFactory.NewMercuryPlugin(ocr3types.MercuryPluginConfig{
				sharedConfig.ConfigDigest,
				oid,
				sharedConfig.N(),
				sharedConfig.F,
				sharedConfig.OnchainConfig,
				sharedConfig.ReportingPluginConfig,
				sharedConfig.DeltaRound,
				sharedConfig.MaxDurationObservation,
			})
			if err != nil {
				logger.Error("ManagedMercuryOracle: error during NewReportingPlugin()", commontypes.LogFields{
					"error": err,
				})
				return
			}
			defer loghelper.CloseLogError(
				mercuryPlugin,
				logger,
				"ManagedMercuryOracle: error during reportingPlugin.Close()",
			)

			// if err := validateReportingPluginLimits(merucuryPluginInfo.Limits); err != nil {
			// 	logger.Error("ManagedMercuryOracle: invalid ReportingPluginInfo", commontypes.LogFields{
			// 		"error":               err,
			// 		"reportingPluginInfo": reportingPluginInfo,
			// 	})
			// 	return
			// }

			lims, err := todoLimits()
			if err != nil {
				logger.Error("ManagedMercuryOracle: error during limits", commontypes.LogFields{
					"error":             err,
					"publicConfig":      sharedConfig.PublicConfig,
					"mercuryPluginInfo": mercuryPluginInfo,
					"maxSigLen":         onchainKeyring.MaxSignatureLength(),
				})
				return
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
				return
			}

			// No need to binNetEndpoint.Start/Close since netEndpoint will handle that for us

			netEndpoint := shim.NewOCR3SerializingEndpoint[mercuryshim.MercuryReportInfo](
				chTelemetrySend,
				sharedConfig.ConfigDigest,
				binNetEndpoint,
				childLogger,
			)
			if err := netEndpoint.Start(); err != nil {
				logger.Error("ManagedMercuryOracle: error during netEndpoint.Start()", commontypes.LogFields{
					"error":        err,
					"configDigest": sharedConfig.ConfigDigest,
				})
				return
			}
			defer loghelper.CloseLogError(
				netEndpoint,
				logger,
				"ManagedMercuryOracle: error during netEndpoint.Close()",
			)

			ocr3PluginConfig := ocr3types.OCR3PluginConfig{
				sharedConfig.ConfigDigest,
				oid,
				sharedConfig.N(),
				sharedConfig.F,
				sharedConfig.OnchainConfig,
				sharedConfig.ReportingPluginConfig,
				sharedConfig.DeltaRound,
				sharedConfig.MaxDurationQuery,
				sharedConfig.MaxDurationObservation,
				sharedConfig.MaxDurationShouldAcceptFinalizedReport,
				sharedConfig.MaxDurationShouldTransmitAcceptedReport,
			}
			ocr3Plugin := &mercuryshim.MercuryOCR3Plugin{
				ocr3PluginConfig,
				mercuryPlugin,
			}

			protocol.RunOracle[mercuryshim.MercuryReportInfo](
				ctx,
				sharedConfig,
				mercuryshim.NewMercuryOCR3ContractTransmitter(contractTransmitter),
				&shim.SerializingOCR3Database{database},
				oid,
				localConfig,
				childLogger,
				netEndpoint,
				offchainKeyring,
				mercuryshim.NewMercuryOCR3OnchainKeyring(onchainKeyring),

				// shim.LimitCheckReportingPlugin{reportingPlugin, reportingPluginInfo.Limits},
				ocr3Plugin,
				shim.MakeOCR3TelemetrySender(chTelemetrySend, childLogger),
			)
		},
		localConfig,
		logger,
		offchainConfigDigester,
	)
}

// func validateReportingPluginLimits(limits types.ReportingPluginLimits) error {
// 	var err error
// 	if !(0 <= limits.MaxQueryLength && limits.MaxQueryLength <= types.MaxMaxQueryLength) {
// 		err = multierr.Append(err, fmt.Errorf("MaxQueryLength (%v) out of range. Should be between 0 and %v", limits.MaxQueryLength, types.MaxMaxQueryLength))
// 	}
// 	if !(0 <= limits.MaxObservationLength && limits.MaxObservationLength <= types.MaxMaxObservationLength) {
// 		err = multierr.Append(err, fmt.Errorf("MaxObservationLength (%v) out of range. Should be between 0 and %v", limits.MaxObservationLength, types.MaxMaxObservationLength))
// 	}
// 	if !(0 <= limits.MaxReportLength && limits.MaxReportLength <= types.MaxMaxReportLength) {
// 		err = multierr.Append(err, fmt.Errorf("MaxReportLength (%v) out of range. Should be between 0 and %v", limits.MaxReportLength, types.MaxMaxReportLength))
// 	}
// 	return err
// }
