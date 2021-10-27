package managed

import (
	"context"
	"sort"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/shim"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// RunManagedOracle runs a "managed" version of protocol.RunOracle. It handles
// setting up telemetry, garbage collection, configuration updates, translating
// from commontypes.BinaryNetworkEndpoint to protocol.NetworkEndpoint, and
// creation/teardown of reporting plugins.
func RunManagedOracle(
	ctx context.Context,

	v2bootstrappers []commontypes.BootstrapperLocator,
	configTracker types.ContractConfigTracker,
	contractTransmitter types.ContractTransmitter,
	database types.Database,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
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

	subs.Go(func() {
		collectGarbage(ctx, database, localConfig, logger)
	})

	runWithContractConfig(
		ctx,

		configTracker,
		database,
		func(ctx context.Context, contractConfig types.ContractConfig, logger loghelper.LoggerWithContext) {
			skipResourceExhaustionChecks := localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode
			sharedConfig, oid, err := config.SharedConfigFromContractConfig(
				skipResourceExhaustionChecks,
				contractConfig,
				offchainKeyring,
				onchainKeyring,
				netEndpointFactory.PeerID(),
				contractTransmitter.FromAccount(),
			)
			if err != nil {
				logger.Error("ManagedOracle: error while updating config", commontypes.LogFields{
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

			reportingPlugin, reportingPluginInfo, err := reportingPluginFactory.NewReportingPlugin(types.ReportingPluginConfig{
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
			if err != nil {
				logger.Error("ManagedOracle: error during NewReportingPlugin()", commontypes.LogFields{
					"error": err,
				})
				return
			}
			if err := reportingPlugin.Start(); err != nil {
				logger.Error("ManagedOracle: error during ReportingPlugin.Start()", commontypes.LogFields{
					"error":               err,
					"reportingPluginInfo": reportingPluginInfo,
				})
				return
			}
			defer loghelper.CloseLogError(
				reportingPlugin,
				logger,
				"ManagedOracle: error during reportingPlugin.Close()",
			)

			lims := limits(sharedConfig.PublicConfig, reportingPluginInfo, onchainKeyring.MaxSignatureLength())
			binNetEndpoint, err := netEndpointFactory.NewEndpoint(
				sharedConfig.ConfigDigest,
				peerIDs,
				v2bootstrappers,
				sharedConfig.F,
				lims,
			)
			if err != nil {
				logger.Error("ManagedOracle: error during NewEndpoint", commontypes.LogFields{
					"error":           err,
					"peerIDs":         peerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return
			}

			// No need to binNetEndpoint.Start/Close since netEndpoint will handle that for us

			netEndpoint := shim.NewSerializingEndpoint(
				chTelemetrySend,
				sharedConfig.ConfigDigest,
				binNetEndpoint,
				childLogger,
			)
			if err := netEndpoint.Start(); err != nil {
				logger.Error("ManagedOracle: error during netEndpoint.Start()", commontypes.LogFields{
					"error":        err,
					"configDigest": sharedConfig.ConfigDigest,
				})
				return
			}
			defer loghelper.CloseLogError(
				reportingPlugin,
				logger,
				"ManagedOracle: error during netEndpoint.Close()",
			)

			var reportQuorum int
			if reportingPluginInfo.UniqueReports {
				// We require greater than (n+f)/2 signatures to reach a byzantine
				// quroum. This ensures unique reports since each honest node will sign
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
				netEndpoint,
				offchainKeyring,
				onchainKeyring,
				reportingPlugin,
				reportQuorum,
				shim.MakeTelemetrySender(chTelemetrySend, childLogger),
			)
		},
		localConfig,
		logger,
		offchainConfigDigester,
	)
}

func max(x int, xs ...int) int {
	sort.Ints(xs)
	if len(xs) == 0 || xs[len(xs)-1] < x {
		return x
	} else {
		return xs[len(xs)-1]
	}
}
