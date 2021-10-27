package tracing

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	offchainreporting "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/managed"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/protocol"
	ti "github.com/smartcontractkit/libocr/offchainreporting2/testimplementations"
	"github.com/smartcontractkit/libocr/offchainreporting2/testsetup"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

func SetUpTestForTracing(
	testArgs testsetup.TestSetupArgs,
	oracleArgs []offchainreporting.OracleArgs,
	identities []config.OracleIdentity,
	tracer *Tracer,
	network *Network,
	dbFactory *ti.InMemoryDatabaseFactory,
) (
	newOracleArgs map[OracleID]offchainreporting.OracleArgs,
	twinOracleArgs map[OracleID]offchainreporting.OracleArgs,
) {
	// wrap contract transmitters with mutexes so that twins can share them.
	// If the oracle doesn't have a twin, the locked transmitter makes no difference but makes this code simpler.
	lockedContractTransmitters := make([]types.ContractTransmitter, len(oracleArgs))
	for i, originalArgs := range oracleArgs {
		lockedContractTransmitters[i] = MakeLockedContractTransmitter(originalArgs.ContractTransmitter)
	}
	// Arguments for the healthy oracles. Most are borrowed from the output of testsetup.MakeTestSetup.
	// However, some are wrapped with tracing instrumentation: network, contract interfaces, database, config digester and the plugin.
	newOracleArgs = map[OracleID]offchainreporting.OracleArgs{}
	for i, originalArgs := range oracleArgs {
		peerID := identities[i].PeerID
		oracleID := OracleID{commontypes.OracleID(i), false}
		args := offchainreporting.OracleArgs{
			network.NewEndpointFactory(oracleID, peerID),                                      //BinaryNetworkEndpointFactory
			originalArgs.V2Bootstrappers,                                                      //V2Bootstrappers
			MakeContractConfigTracker(tracer, oracleID, originalArgs.ContractConfigTracker),   //ContractConfigTracker
			MakeContractTransmitter(tracer, oracleID, lockedContractTransmitters[i]),          //ContractTransmitter
			MakeDatabase(tracer, oracleID, dbFactory.MakeDatabase(oracleID.Int())),            //Database
			originalArgs.LocalConfig,                                                          //LocalConfig
			originalArgs.Logger,                                                               //Logger
			originalArgs.MonitoringEndpoint,                                                   //MonitoringEndpoint
			MakeOffchainConfigDigester(tracer, oracleID, originalArgs.OffchainConfigDigester), //OffchainConfigDigester
			originalArgs.OffchainKeyring,                                                      //OffchainKeyring
			originalArgs.OnchainKeyring,                                                       //OnchainKeyring
			MakeReportingPluginFactory(tracer, oracleID, originalArgs.ReportingPluginFactory), //ReportingPluginFactory
		}
		newOracleArgs[oracleID] = args
	}
	// Arguments for the twins.
	// Twins share the same backend dependencies except for the networking, database, logger and reporting plugin
	twinOracleArgs = map[OracleID]offchainreporting.OracleArgs{}
	for _, twinID := range testArgs.Twins {
		originalOracleID := FromInt(-twinID)
		twinOracleID := FromInt(twinID)
		originalArgs := oracleArgs[int(originalOracleID.OracleID)]
		originalPeerID := identities[int(originalOracleID.OracleID)].PeerID
		twinLogger := loghelper.MakeRootLoggerWithContext(originalArgs.Logger).MakeChild(commontypes.LogFields{
			"id": twinID,
		})
		twinArgs := offchainreporting.OracleArgs{
			network.NewEndpointFactory(twinOracleID, originalPeerID),                                                  // BinaryNetworkEndpointFactory:
			originalArgs.V2Bootstrappers,                                                                              // V2Bootstrappers
			MakeContractConfigTracker(tracer, twinOracleID, originalArgs.ContractConfigTracker),                       // ContractConfigTracker:
			MakeContractTransmitter(tracer, twinOracleID, lockedContractTransmitters[int(originalOracleID.OracleID)]), // ContractTransmitter:
			MakeDatabase(tracer, twinOracleID, dbFactory.MakeDatabase(twinOracleID.Int())),                            // Database:
			originalArgs.LocalConfig,        // LocalConfig:
			twinLogger,                      // Logger:
			originalArgs.MonitoringEndpoint, // MonitoringEndpoint:
			MakeOffchainConfigDigester(tracer, twinOracleID, originalArgs.OffchainConfigDigester), // OffchainConfigDigester:
			originalArgs.OffchainKeyring, // OffchainKeyring
			originalArgs.OnchainKeyring,  // OnchainKeyring
			MakeReportingPluginFactory(tracer, twinOracleID, originalArgs.ReportingPluginFactory), // ReportingPluginFactory
		}
		twinOracleArgs[twinOracleID] = twinArgs
	}
	return newOracleArgs, twinOracleArgs
}

func ExtractArgs(argsMappings ...map[OracleID]offchainreporting.OracleArgs) []offchainreporting.OracleArgs {
	out := []offchainreporting.OracleArgs{}
	for _, argsMapping := range argsMappings {
		for _, args := range argsMapping {
			out = append(out, args)
		}
	}
	return out
}

func RunSimulation(ctx context.Context, oracleArgs []offchainreporting.OracleArgs) {
	var wg sync.WaitGroup
	for _, params := range oracleArgs {
		wg.Add(1)
		go func(ctx context.Context, params offchainreporting.OracleArgs) {
			defer wg.Done()
			managed.RunManagedOracle(
				ctx,
				params.V2Bootstrappers,
				params.ContractConfigTracker,
				params.ContractTransmitter,
				params.Database,
				params.LocalConfig,
				loghelper.MakeRootLoggerWithContext(params.Logger),
				params.MonitoringEndpoint,
				params.BinaryNetworkEndpointFactory,
				params.OffchainConfigDigester,
				params.OffchainKeyring,
				params.OnchainKeyring,
				params.ReportingPluginFactory,
			)
		}(ctx, params)
	}
	wg.Wait()
}

// SetLeaderForOracles will update the DBs of all oracles that are registered with
// the epoch that corresponds to the desired leader.
// Note that the test Database implementation ignores the configDigest and will always serve the same configuration.
// Note that twins don't knows of themselves as twins so to set a pair of twins (x, -x) as leaders it is sufficient
// to set newLeader=x!
func SetLeaderForOracles(
	ctx context.Context,
	dbFactory *ti.InMemoryDatabaseFactory,
	params testsetup.TestSetupArgs,
	newLeader int,
	leaderSelectionKey [16]byte,
) {
	oracleIDs, twinOracleIDs := testsetup.GetOracleIDs(params.N), params.Twins
	epoch := getEpochForLeader(commontypes.OracleID(newLeader), len(oracleIDs), leaderSelectionKey)
	digest, _ := types.BytesToConfigDigest(make([]byte, 16))
	newState := types.PersistentState{
		Epoch:                epoch,
		HighestSentEpoch:     epoch,
		HighestReceivedEpoch: repeat(epoch, len(oracleIDs)),
	}
	for _, oracleID := range append(oracleIDs, twinOracleIDs...) {
		slice := dbFactory.GetDatabase(oracleID)
		if slice == nil {
			continue
		}
		_ = slice.WriteState(ctx, digest, newState)
	}
}

// Helpers

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func getEpochForLeader(newLeader commontypes.OracleID, n int, key [16]byte) uint32 {
	for epoch := uint32(0); epoch < 100; epoch++ {
		if protocol.Leader(epoch, n, key) == newLeader {
			return epoch
		}
	}
	panic(fmt.Sprintf("unable to find an epoch for %d to be leader", newLeader))
}

func repeat(val uint32, length int) []uint32 {
	out := make([]uint32, length)
	for i := 0; i < length; i++ {
		out[i] = val
	}
	return out
}

func FindHonestOracleIDs(args testsetup.TestSetupArgs) (honest, dishonest []commontypes.OracleID) {
	honest, dishonest = []commontypes.OracleID{}, []commontypes.OracleID{}
	twinsSet := map[commontypes.OracleID]struct{}{}
	for _, oid := range args.Twins {
		twinsSet[commontypes.OracleID(abs(int(oid)))] = struct{}{}
	}
	for i := 0; i < args.N; i++ {
		if _, isTwin := twinsSet[commontypes.OracleID(i)]; isTwin {
			dishonest = append(dishonest, commontypes.OracleID(i))
		} else {
			honest = append(honest, commontypes.OracleID(i))
		}
	}
	return honest, dishonest
}

func RemoveTracesFromOracles(traces []Trace, oids []commontypes.OracleID) []Trace {
	oidSet := map[commontypes.OracleID]struct{}{}
	for _, oid := range oids {
		oidSet[oid] = struct{}{}
	}
	out := []Trace{}
	for _, raw := range traces {
		var oracleID commontypes.OracleID
		switch trace := raw.(type) {
		case *SendTo:
			oracleID = trace.Common.Originator.OracleID
		case *Broadcast:
			oracleID = trace.Common.Originator.OracleID
		case *Receive:
			oracleID = trace.Common.Originator.OracleID
		case *Drop:
			oracleID = trace.Common.Originator.OracleID
		case *EndpointStart:
			oracleID = trace.Common.Originator.OracleID
		case *EndpointClose:
			oracleID = trace.Common.Originator.OracleID

		case *ReadState:
			oracleID = trace.Common.Originator.OracleID
		case *WriteState:
			oracleID = trace.Common.Originator.OracleID
		case *ReadConfig:
			oracleID = trace.Common.Originator.OracleID
		case *WriteConfig:
			oracleID = trace.Common.Originator.OracleID
		case *StorePendingTransmission:
			oracleID = trace.Common.Originator.OracleID
		case *PendingTransmissionsWithConfigDigest:
			oracleID = trace.Common.Originator.OracleID
		case *DeletePendingTransmission:
			oracleID = trace.Common.Originator.OracleID
		case *DeletePendingTransmissionsOlderThan:
			oracleID = trace.Common.Originator.OracleID

		case *Notify:
			oracleID = trace.Common.Originator.OracleID
		case *LatestConfigDetails:
			oracleID = trace.Common.Originator.OracleID
		case *LatestConfig:
			oracleID = trace.Common.Originator.OracleID
		case *LatestBlockHeight:
			oracleID = trace.Common.Originator.OracleID

		case *Transmit:
			oracleID = trace.Common.Originator.OracleID
		case *LatestConfigDigestAndEpoch:
			oracleID = trace.Common.Originator.OracleID
		case *FromAccount:
			oracleID = trace.Common.Originator.OracleID

		case *Query:
			oracleID = trace.Common.Originator.OracleID
		case *Observation:
			oracleID = trace.Common.Originator.OracleID
		case *Report:
			oracleID = trace.Common.Originator.OracleID
		case *ShouldAcceptFinalizedReport:
			oracleID = trace.Common.Originator.OracleID
		case *ShouldTransmitAcceptedReport:
			oracleID = trace.Common.Originator.OracleID
		case *PluginStart:
			oracleID = trace.Common.Originator.OracleID
		case *PluginClose:
			oracleID = trace.Common.Originator.OracleID

		case *ConfigDigest:
			oracleID = trace.Common.Originator.OracleID
		case *ConfigDigestPrefix:
			oracleID = trace.Common.Originator.OracleID

		default:
			panic(fmt.Sprintf("unexpected trace %#v", trace))
		}
		if _, isRejectedOracleID := oidSet[oracleID]; !isRejectedOracleID {
			out = append(out, raw)
		}
	}
	return out
}
