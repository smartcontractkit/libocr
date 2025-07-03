package protocol

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/scheduler"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

func RunStatePersistence[RI any](
	ctx context.Context,

	chNetToStatePersistence <-chan MessageToStatePersistenceWithSender[RI],
	chReportAttestationToStatePersistence <-chan EventToStatePersistence[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	kvStore KeyValueStore,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	reportingPlugin ocr3_1types.ReportingPlugin[RI],
	restoredState StatePersistenceState,
	restoredHighestCommittedToKVSeqNr uint64,
) {
	sched := scheduler.NewScheduler[EventToStatePersistence[RI]]()
	defer sched.Close()

	newStatePersistenceState(ctx, chNetToStatePersistence,
		chReportAttestationToStatePersistence,
		config, database, id, kvStore, logger, netSender, reportingPlugin, sched).run(restoredState)
}

const maxPersistedAttestedStateTransitionBlocks int = math.MaxInt

type statePersistenceState[RI any] struct {
	ctx context.Context

	chNetToStatePersistence               <-chan MessageToStatePersistenceWithSender[RI]
	chReportAttestationToStatePersistence <-chan EventToStatePersistence[RI]
	tTryReplay                            <-chan time.Time
	config                                ocr3config.SharedConfig
	database                              Database
	id                                    commontypes.OracleID
	kvStore                               KeyValueStore
	logger                                loghelper.LoggerWithContext
	netSender                             NetworkSender[RI]
	reportingPlugin                       ocr3_1types.ReportingPlugin[RI]

	highestPersistedStateTransitionBlockSeqNr uint64
	highestHeardSeqNr                         uint64
	readyToSendBlockSyncReq                   bool

	blockSyncState blockSyncState[RI]
	treeSyncState  treeSyncState
}

func (state *statePersistenceState[RI]) run(restoredState StatePersistenceState) {
	state.highestPersistedStateTransitionBlockSeqNr = restoredState.HighestPersistedStateTransitionBlockSeqNr
	state.logger.Info("StatePersistence: running", commontypes.LogFields{
		"restoredHighestPersistedStateTransitionBlockSeqNr": restoredState.HighestPersistedStateTransitionBlockSeqNr,
	})

	for {
		select {
		case msg := <-state.chNetToStatePersistence:
			msg.msg.processStatePersistence(state, msg.sender)
		case ev := <-state.chReportAttestationToStatePersistence:
			ev.processStatePersistence(state)
		case ev := <-state.blockSyncState.scheduler.Scheduled():
			ev.processStatePersistence(state)
		case <-state.tTryReplay:
			state.eventTTryReplay()
		case <-state.ctx.Done():
		}

		// ensure prompt exit
		select {
		case <-state.ctx.Done():
			state.logger.Info("StatePersistence: exiting", nil)
			// state.scheduler.Close()
			return
		default:
		}
	}
}

func (state *statePersistenceState[RI]) eventTTryReplay() {
	state.logger.Trace("TTryReplay fired", nil)

	progressed := state.tryReplay()

	_, haveNext, retry := state.nextBlockToReplay()
	if progressed || haveNext || retry {

		state.tTryReplay = time.After(0)
	}
}

func (state *statePersistenceState[RI]) tryReplay() bool {
	block, ok, _ := state.nextBlockToReplay()
	if !ok {
		return false
	}
	return state.replayVerifiedBlock(block)
}

func (state *statePersistenceState[RI]) replayVerifiedBlock(stb StateTransitionBlock) (success bool) {
	writeSet := stb.StateTransitionOutputs.WriteSet

	seqNr := stb.SeqNr()
	logger := state.logger.MakeChild(commontypes.LogFields{
		"replay": "YES",
		"seqNr":  seqNr,
	})

	logger.Trace("replaying state transition block", nil)
	kvReadWriteTxn, err := state.kvStore.NewReadWriteTransaction(seqNr)
	if err != nil {
		logger.Error("could not open new kv transaction", commontypes.LogFields{
			"err": err,
		})
		return
	}
	defer kvReadWriteTxn.Discard()

	for _, m := range writeSet {

		var err error
		if m.Deleted {
			err = kvReadWriteTxn.Delete(m.Key)
		} else {
			err = kvReadWriteTxn.Write(m.Key, m.Value)
		}
		if err != nil {
			logger.Error("failed to write write-set modification", commontypes.LogFields{
				"error": err,
				"seqNr": seqNr,
			})
			return
		}
	}

	werr := kvReadWriteTxn.Commit()
	kvReadWriteTxn.Discard()

	if werr != nil {
		kvSeqNr, rerr := state.highestCommittedToKVSeqNr()
		if rerr != nil {
			logger.Error("failed to commit kv transaction, and then failed to read highest committed to kv seq nr", commontypes.LogFields{
				"werror": werr,
				"rerror": rerr,
				"seqNr":  seqNr,
			})
			return
		}
		if kvSeqNr < seqNr {
			logger.Error("failed to commit kv transaction, but not due to conflict without outcome generation", commontypes.LogFields{
				"seqNr":   seqNr,
				"kvSeqNr": kvSeqNr,
				"error":   werr,
			})

			return
		} else {

			return
		}
	}
	success = true
	return
}

func (state *statePersistenceState[RI]) eventStateSyncRequest(ev EventStateSyncRequest[RI]) {
	state.logger.Debug("received EventStateSyncRequest", commontypes.LogFields{
		"heardSeqNr": ev.SeqNr,
	})
	state.heardSeqNr(ev.SeqNr)
}

func (state *statePersistenceState[RI]) heardSeqNr(seqNr uint64) {
	if seqNr > state.highestHeardSeqNr {
		state.logger.Debug("highest heard sequence number increased", commontypes.LogFields{
			"old": state.highestHeardSeqNr,
			"new": seqNr,
		})
		state.highestHeardSeqNr = seqNr
		state.highestHeardIncreased()
	}
}

func (state *statePersistenceState[RI]) refreshHighestPersistedStateTransitionBlockSeqNr() {
	highestCommittedToKVSeqNr, err := state.highestCommittedToKVSeqNr()
	if err != nil {
		state.logger.Error("failed to get highest committed to kv seq nr during refresh", commontypes.LogFields{
			"error": err,
		})
		return
	}
	if highestCommittedToKVSeqNr > state.highestPersistedStateTransitionBlockSeqNr {

		state.highestPersistedStateTransitionBlockSeqNr = highestCommittedToKVSeqNr
	}
}

func (state *statePersistenceState[RI]) persist(verifiedAstb AttestedStateTransitionBlock) error {
	state.refreshHighestPersistedStateTransitionBlockSeqNr()
	expectedSeqNr := state.highestPersistedStateTransitionBlockSeqNr + 1
	seqNr := verifiedAstb.StateTransitionBlock.SeqNr()

	if seqNr != expectedSeqNr {

		return fmt.Errorf("cannot persist out of order state transition block: expected %d, got %d",
			expectedSeqNr,
			seqNr,
		)
	}

	err := state.database.WriteAttestedStateTransitionBlock(
		state.ctx,
		state.config.ConfigDigest,
		seqNr,
		verifiedAstb,
	)
	if err != nil {
		return fmt.Errorf("failed to write attested state transition block %d: %w", seqNr, err)
	}

	state.highestPersistedStateTransitionBlockSeqNr = seqNr
	state.logger.Trace("persisted block", commontypes.LogFields{
		"seqNr": seqNr,
	})
	state.tTryReplay = time.After(0)

	err = state.database.WriteStatePersistenceState(
		state.ctx, state.config.ConfigDigest,
		StatePersistenceState{
			seqNr,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write state persistence state %d: %w", seqNr, err)
	}
	return nil
}

func (state *statePersistenceState[RI]) highestCommittedToKVSeqNr() (uint64, error) {
	return state.kvStore.HighestCommittedSeqNr()
}

func (state *statePersistenceState[RI]) nextBlockToReplay() (block StateTransitionBlock, found bool, retry bool) {
	committedToKVSeqNr, err := state.highestCommittedToKVSeqNr()
	if err != nil {
		state.logger.Error("failed to get highest committed to kv seq nr", commontypes.LogFields{
			"error": err,
		})
		retry = true
		return
	}
	nextSeqNr := committedToKVSeqNr + 1

	astb, err := state.database.ReadAttestedStateTransitionBlock(state.ctx, state.config.ConfigDigest, nextSeqNr)
	if err != nil {
		state.logger.Error("failed to read attested state transition block from database", commontypes.LogFields{
			"nextSeqNr": nextSeqNr,
			"error":     err,
		})
		retry = true
		return
	}
	seqNr := astb.StateTransitionBlock.SeqNr()
	if seqNr == 0 {
		// block not found
		state.logger.Trace("wanted next block to replay not found", commontypes.LogFields{
			"nextSeqNr": nextSeqNr,
		})

		return
	} else if seqNr == nextSeqNr {
		state.logger.Debug("next state transition block to replay", commontypes.LogFields{
			"nextSeqNr": nextSeqNr,
		})

		block = astb.StateTransitionBlock
		found = true
		return
	} else {
		state.logger.Critical("assumption violation, block in database has inconsistent seq nr", commontypes.LogFields{
			"expectedSeqNr": nextSeqNr,
			"actualSeqNr":   seqNr,
			"block":         astb,
		})
		panic("")
	}
}

func (state *statePersistenceState[RI]) eventEventBlockSyncSummaryHeartbeat(ev EventBlockSyncSummaryHeartbeat[RI]) {
	state.processBlockSyncSummaryHeartbeat()
}

func (state *statePersistenceState[RI]) eventExpiredBlockSyncRequest(ev EventExpiredBlockSyncRequest[RI]) {
	state.blockSyncState.logger.Debug("received eventExpiredBlockSyncRequest", commontypes.LogFields{
		"requestedFrom": ev.RequestedFrom,
	})
	state.processExpiredBlockSyncRequest(ev.RequestedFrom, ev.Nonce)
}

func (state *statePersistenceState[RI]) eventReadyToSendNextBlockSyncRequest(ev EventReadyToSendNextBlockSyncRequest[RI]) {
	state.logger.Debug("received eventReadyToSendNextBlockSyncRequest", commontypes.LogFields{})
	state.readyToSendBlockSyncReq = true
	state.trySendNextRequest()
}

func newStatePersistenceState[RI any](
	ctx context.Context,
	chNetToStatePersistence <-chan MessageToStatePersistenceWithSender[RI],

	chReportAttestationToStatePersistence <-chan EventToStatePersistence[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	kvStore KeyValueStore,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	reportingPlugin ocr3_1types.ReportingPlugin[RI],
	scheduler *scheduler.Scheduler[EventToStatePersistence[RI]],
) *statePersistenceState[RI] {
	oracles := make([]*blockSyncTargetOracle[RI], 0)
	for i := 0; i < config.N(); i++ {
		oracles = append(oracles, &blockSyncTargetOracle[RI]{
			0,
			time.Time{},
			true,
			nil,
		})
	}

	scheduler.ScheduleDelay(EventBlockSyncSummaryHeartbeat[RI]{}, DeltaBlockSyncHeartbeat)

	tTryReplay := time.After(0)

	return &statePersistenceState[RI]{
		ctx,

		chNetToStatePersistence,
		chReportAttestationToStatePersistence,
		tTryReplay,
		config,
		database,
		id,
		kvStore,
		logger.MakeUpdated(commontypes.LogFields{"proto": "state"}),
		netSender,
		reportingPlugin,
		0,
		0,
		true,

		blockSyncState[RI]{
			logger.MakeUpdated(commontypes.LogFields{"proto": "stateBlockSync"}),
			oracles,
			scheduler,
		},
		treeSyncState{},
	}
}
