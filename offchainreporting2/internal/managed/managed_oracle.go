package managed

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/big"
	"sort"
	"time"

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
// configuration updates and translating from commontypes.BinaryNetworkEndpoint to
// protocol.NetworkEndpoint.
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
	mo := managedOracleState{
		ctx: ctx,

		v2bootstrappers:        v2bootstrappers,
		configTracker:          configTracker,
		contractTransmitter:    contractTransmitter,
		database:               database,
		onchainKeyring:         onchainKeyring,
		localConfig:            localConfig,
		logger:                 logger,
		monitoringEndpoint:     monitoringEndpoint,
		netEndpointFactory:     netEndpointFactory,
		offchainKeyring:        offchainKeyring,
		reportingPluginFactory: reportingPluginFactory,
		configDigester:         prefixCheckConfigDigester{offchainConfigDigester},
	}
	mo.run()
}

type managedOracleState struct {
	ctx context.Context

	v2bootstrappers        []commontypes.BootstrapperLocator
	config                 config.SharedConfig
	configTracker          types.ContractConfigTracker
	contractTransmitter    types.ContractTransmitter
	database               types.Database
	onchainKeyring         types.OnchainKeyring
	localConfig            types.LocalConfig
	logger                 loghelper.LoggerWithContext
	monitoringEndpoint     commontypes.MonitoringEndpoint
	netEndpointFactory     types.BinaryNetworkEndpointFactory
	offchainKeyring        types.OffchainKeyring
	reportingPluginFactory types.ReportingPluginFactory

	chTelemetry        chan<- *serialization.TelemetryWrapper
	configDigester     prefixCheckConfigDigester
	netEndpoint        *shim.SerializingEndpoint
	reportingPlugin    types.ReportingPlugin
	oracleCancel       context.CancelFunc
	oracleSubprocesses subprocesses.Subprocesses
	otherSubprocesses  subprocesses.Subprocesses
}

func (mo *managedOracleState) run() {
	{
		chTelemetry := make(chan *serialization.TelemetryWrapper, 100)
		mo.chTelemetry = chTelemetry
		mo.otherSubprocesses.Go(func() {
			forwardTelemetry(mo.ctx, mo.logger, mo.monitoringEndpoint, chTelemetry)
		})
	}

	mo.otherSubprocesses.Go(func() {
		collectGarbage(mo.ctx, mo.database, mo.localConfig, mo.logger)
	})

	// Restore config from database, so that we can run even if the ethereum node
	// isn't working.
	mo.restoreFromDatabase()

	// Only start tracking config after we attempted to load config from db
	chNewConfig := make(chan types.ContractConfig, 5)
	mo.otherSubprocesses.Go(func() {
		TrackConfig(mo.ctx, mo.configDigester, mo.configTracker, mo.config.ConfigDigest, mo.localConfig, mo.logger, chNewConfig)
	})

	for {
		select {
		case change := <-chNewConfig:
			mo.logger.Info("ManagedOracle: switching between configs", commontypes.LogFields{
				"oldConfigDigest": mo.config.ConfigDigest.Hex(),
				"newConfigDigest": change.ConfigDigest.Hex(),
			})
			mo.configChanged(change)
		case <-mo.ctx.Done():
			mo.logger.Info("ManagedOracle: winding down", nil)
			mo.closeOracle()
			mo.otherSubprocesses.Wait()
			mo.logger.Info("ManagedOracle: exiting", nil)
			return // Exit ManagedOracle event loop altogether
		}
	}
}

func (mo *managedOracleState) restoreFromDatabase() {
	// Restore config from database, so that we can run even if the ethereum node
	// isn't working.
	var cc *types.ContractConfig
	ok := mo.otherSubprocesses.BlockForAtMost(
		mo.ctx,
		mo.localConfig.DatabaseTimeout,
		func(ctx context.Context) {
			cc = loadConfigFromDatabase(ctx, mo.database, mo.logger)
		},
	)
	if !ok {
		mo.logger.Error("ManagedOracle: database timed out while attempting to restore configuration", commontypes.LogFields{
			"timeout": mo.localConfig.DatabaseTimeout,
		})
		return
	}

	if cc == nil {
		mo.logger.Info("ManagedOracle: found no configuration to restore", commontypes.LogFields{})
		return
	}

	if err := mo.configDigester.CheckContractConfig(*cc); err != nil {
		mo.logger.Error("ManagedOracle: error checking ConfigDigest while attempting to restore configuration", commontypes.LogFields{
			"err":            err,
			"contractConfig": *cc,
		})
		return
	}

	mo.configChanged(*cc)
}

func (mo *managedOracleState) closeOracle() {
	if mo.oracleCancel != nil {
		mo.oracleCancel()
		mo.oracleSubprocesses.Wait()
		if err := mo.netEndpoint.Close(); err != nil {
			mo.logger.Error("ManagedOracle: error while closing BinaryNetworkEndpoint", commontypes.LogFields{
				"error": err,
			})
			// nothing to be done about it, let's try to carry on.
		}
		if err := mo.reportingPlugin.Close(); err != nil {
			mo.logger.Error("ManagedOracle: error while closing ReportingPlugin", commontypes.LogFields{
				"error": err,
			})
			// nothing to be done about it, let's try to carry on.
		}
		mo.oracleCancel = nil
		mo.netEndpoint = nil
		mo.reportingPlugin = nil
	}
}

func (mo *managedOracleState) configChanged(contractConfig types.ContractConfig) {
	// Cease any operation from earlier configs
	mo.closeOracle()

	// Decode contractConfig

	skipChainSpecificChecks := true //mo.localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode
	var err error
	var oid commontypes.OracleID
	mo.config, oid, err = config.SharedConfigFromContractConfig(
		big.NewInt(1),
		skipChainSpecificChecks,
		contractConfig,
		mo.offchainKeyring,
		mo.onchainKeyring,
		mo.netEndpointFactory.PeerID(),
		mo.contractTransmitter.FromAccount(),
	)
	if err != nil {
		mo.logger.Error("ManagedOracle: error while updating config", commontypes.LogFields{
			"error": err,
		})
		return
	}

	// Run with new config
	peerIDs := []string{}
	for _, identity := range mo.config.OracleIdentities {
		peerIDs = append(peerIDs, identity.PeerID)
	}

	childLogger := mo.logger.MakeChild(commontypes.LogFields{
		"configDigest": fmt.Sprintf("%x", mo.config.ConfigDigest),
		"oid":          oid,
	})

	reportingPlugin, reportingPluginInfo, err := mo.reportingPluginFactory.NewReportingPlugin(types.ReportingPluginConfig{
		mo.config.ConfigDigest,
		mo.config.N(),
		mo.config.F,
		mo.config.OnchainConfig,
		mo.config.ReportingPluginConfig,
		mo.config.DeltaRound,
		mo.config.MaxDurationQuery,
		mo.config.MaxDurationObservation,
		mo.config.MaxDurationReport,
		mo.config.MaxDurationShouldAcceptFinalizedReport,
		mo.config.MaxDurationShouldTransmitAcceptedReport,
	})
	mo.reportingPlugin = reportingPlugin
	if err != nil {
		mo.logger.Error("ManagedOracle: error during MakeReportingPlugin()", commontypes.LogFields{
			"error":        err,
			"configDigest": mo.config.ConfigDigest,
		})
		return
	}

	lims := limits(mo.config.PublicConfig, reportingPluginInfo, mo.onchainKeyring.MaxSignatureLength())
	binNetEndpoint, err := mo.netEndpointFactory.NewEndpoint(
		mo.config.ConfigDigest,
		peerIDs,
		mo.v2bootstrappers,
		mo.config.F,
		lims,
	)
	if err != nil {
		mo.logger.Error("ManagedOracle: error during NewEndpoint", commontypes.LogFields{
			"error":           err,
			"configDigest":    mo.config.ConfigDigest,
			"peerIDs":         peerIDs,
			"v2bootstrappers": mo.v2bootstrappers,
		})
		return
	}

	netEndpoint := shim.NewSerializingEndpoint(
		mo.chTelemetry,
		mo.config.ConfigDigest,
		binNetEndpoint,
		childLogger,
	)
	mo.netEndpoint = netEndpoint
	if err := netEndpoint.Start(); err != nil {
		mo.logger.Error("ManagedOracle: error during netEndpoint.Start()", commontypes.LogFields{
			"error":        err,
			"configDigest": mo.config.ConfigDigest,
		})
		return
	}

	var reportQuorum int
	if reportingPluginInfo.UniqueReports {
		// We requre greater than (n+f)/2 signatures to reach a byzantine
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
		reportQuorum = (mo.config.N()+mo.config.F)/2 + 1
	} else {
		reportQuorum = mo.config.F + 1
	}

	oracleCtx, oracleCancel := context.WithCancel(mo.ctx)
	mo.oracleCancel = oracleCancel
	mo.oracleSubprocesses.Go(func() {
		defer oracleCancel()
		protocol.RunOracle(
			oracleCtx,
			mo.config,
			mo.contractTransmitter,
			mo.database,
			oid,
			mo.localConfig,
			childLogger,
			mo.netEndpoint,
			mo.offchainKeyring,
			mo.onchainKeyring,
			mo.reportingPlugin,
			reportQuorum,
			shim.MakeTelemetrySender(mo.chTelemetry, childLogger),
		)
	})

	childCtx, childCancel := context.WithTimeout(mo.ctx, mo.localConfig.DatabaseTimeout)
	defer childCancel()
	if err := mo.database.WriteConfig(childCtx, contractConfig); err != nil {
		mo.logger.ErrorIfNotCanceled("ManagedOracle: error writing new config to database", childCtx, commontypes.LogFields{
			"configDigest": mo.config.ConfigDigest,
			"config":       contractConfig,
			"error":        err,
		})
	}
}

func limits(cfg config.PublicConfig, reportingPluginInfo types.ReportingPluginInfo, maxSigLen int) types.BinaryNetworkEndpointLimits {
	const overhead = 256

	maxLenNewEpoch := overhead
	maxLenObserveReq := reportingPluginInfo.MaxQueryLen + overhead
	maxLenObserve := reportingPluginInfo.MaxObservationLen + overhead
	maxLenReportReq := (reportingPluginInfo.MaxObservationLen+ed25519.SignatureSize)*cfg.N() + overhead
	maxLenReport := reportingPluginInfo.MaxReportLen + ed25519.SignatureSize + overhead
	maxLenFinal := reportingPluginInfo.MaxReportLen + maxSigLen*cfg.N() + overhead
	maxLenFinalEcho := maxLenFinal

	maxMessageSize := max(maxLenObserveReq, maxLenObserve, maxLenReportReq, maxLenReport, maxLenFinal, maxLenFinalEcho)

	messagesRate := (1.0*float64(time.Second)/float64(cfg.DeltaResend) +
		1.0*float64(time.Second)/float64(cfg.DeltaProgress) +
		1.0*float64(time.Second)/float64(cfg.DeltaRound) +
		3.0*float64(time.Second)/float64(cfg.DeltaRound) +
		2.0*float64(time.Second)/float64(cfg.DeltaRound)) * 2.0

	messagesCapacity := (2 + 6) * 2

	bytesRate := float64(time.Second)/float64(cfg.DeltaResend)*float64(maxLenNewEpoch) +
		float64(time.Second)/float64(cfg.DeltaProgress)*float64(maxLenNewEpoch) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenObserveReq) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenObserve) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenReportReq) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenReport) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenFinal) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenFinalEcho)

	bytesCapacity := (maxLenNewEpoch + maxLenObserveReq + maxLenObserve + maxLenReportReq + maxLenReport + maxLenFinal + maxLenFinalEcho) * 2

	return types.BinaryNetworkEndpointLimits{
		maxMessageSize,
		messagesRate,
		messagesCapacity,
		bytesRate,
		bytesCapacity,
	}
}

func max(x int, xs ...int) int {
	sort.Ints(xs)
	if len(xs) == 0 || xs[len(xs)-1] < x {
		return x
	} else {
		return xs[len(xs)-1]
	}
}
