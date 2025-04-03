package protocol

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/scheduler"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol/queue"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

func RunStatePersistence[RI any](
	ctx context.Context,

	chNetToStatePersistence <-chan MessageToStatePersistenceWithSender[RI],
	chOutcomeGenerationToStatePersistence <-chan EventToStatePersistence[RI],
	chReportAttestationToStatePersistence <-chan EventToStatePersistence[RI],
	chStatePersistenceToOutcomeGeneration chan<- EventToOutcomeGeneration[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	kvStore ocr3_1types.KeyValueStore,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	restoredState StatePersistenceState,
	restoredHighestCommittedToKVSeqNr uint64,
) {
	sched := scheduler.NewScheduler[EventToStatePersistence[RI]]()
	defer sched.Close()

	newStatePersistenceState(ctx, chNetToStatePersistence,
		chOutcomeGenerationToStatePersistence,
		chReportAttestationToStatePersistence,
		chStatePersistenceToOutcomeGeneration,
		config, database, id, kvStore, logger, netSender, sched).run(restoredState, restoredHighestCommittedToKVSeqNr)
}

const maxPersistedAttestedStateTransitionBlocks int = math.MaxInt

type statePersistenceState[RI any] struct {
	ctx context.Context

	chNetToStatePersistence               <-chan MessageToStatePersistenceWithSender[RI]
	chOutcomeGenerationToStatePersistence <-chan EventToStatePersistence[RI]
	chReportAttestationToStatePersistence <-chan EventToStatePersistence[RI]
	chStatePersistenceToOutcomeGeneration chan<- EventToOutcomeGeneration[RI]
	config                                ocr3config.SharedConfig
	database                              Database
	id                                    commontypes.OracleID
	kvStore                               ocr3_1types.KeyValueStore
	logger                                loghelper.LoggerWithContext
	netSender                             NetworkSender[RI]

	kvStoreTxn                                      *kvStoreTxn
	highestComputedStateTransitionSeqNr             uint64
	highestPersistedStateTransitionBlockSeqNr       uint64
	highestCommittedToKVSeqNr                       uint64
	highestAcknowledgedCommittedToKVSeqNr           uint64
	highestHeardSeqNr                               uint64
	readyToSendBlockSyncReq                         bool
	retrievedStateTransitionBlockSeqNrQueue         *queue.Queue[uint64]
	eventProduceStateTransition                     EventProduceStateTransition[RI]
	notifyOutcomeGenerationOfOpenedKVTransaction    bool
	notifyOutcomeGenerationOfReceivedKVTransaction  bool
	notifyOutcomeGenerationOfCommittedKVTransaction bool

	blockSyncState blockSyncState[RI]
	treeSyncState  treeSyncState
}

type kvStoreTxn struct {
	transaction ocr3_1types.KeyValueStoreTransaction
	replayed    bool
}

func (state *statePersistenceState[RI]) run(restoredState StatePersistenceState, restoredCommittedToKVStoreSeqNr uint64) {
	state.highestPersistedStateTransitionBlockSeqNr = restoredState.HighestPersistedStateTransitionBlockSeqNr
	state.highestCommittedToKVSeqNr = restoredCommittedToKVStoreSeqNr
	state.highestAcknowledgedCommittedToKVSeqNr = restoredCommittedToKVStoreSeqNr
	state.highestComputedStateTransitionSeqNr = restoredCommittedToKVStoreSeqNr
	state.logger.Info("StatePersistence: running", commontypes.LogFields{
		"restoredHighestPersistedStateTransitionBlockSeqNr": restoredState.HighestPersistedStateTransitionBlockSeqNr,
		"restoredCommittedToKVStoreSeqNr":                   restoredCommittedToKVStoreSeqNr,
	})

	for {
		var nilOrEventProduceStateTransition chan<- EventToOutcomeGeneration[RI]
		if state.notifyOutcomeGenerationOfOpenedKVTransaction {
			nilOrEventProduceStateTransition = state.chStatePersistenceToOutcomeGeneration
		} else {
			nilOrEventProduceStateTransition = nil
		}

		var nilOrEventReplayStateTransition chan<- EventToOutcomeGeneration[RI]
		nextBlock, ok := state.nextStateTransitionBlockToReplay()
		if ok {
			nilOrEventReplayStateTransition = state.chStatePersistenceToOutcomeGeneration
		} else {
			nilOrEventReplayStateTransition = nil
		}

		var nilOrEventReceivedKVTransaction chan<- EventToOutcomeGeneration[RI]
		if state.notifyOutcomeGenerationOfReceivedKVTransaction {
			nilOrEventReceivedKVTransaction = state.chStatePersistenceToOutcomeGeneration
		} else {
			nilOrEventReceivedKVTransaction = nil
		}

		var nilOrEventCommittedKVStateTransaction chan<- EventToOutcomeGeneration[RI]
		if state.notifyOutcomeGenerationOfCommittedKVTransaction {
			nilOrEventCommittedKVStateTransaction = state.chStatePersistenceToOutcomeGeneration
		} else {
			nilOrEventCommittedKVStateTransaction = nil
		}

		select {
		case nilOrEventProduceStateTransition <- EventProduceStateTransition[RI]{
			state.eventProduceStateTransition.RoundCtx,
			state.eventProduceStateTransition.Txn,
			state.eventProduceStateTransition.Query,
			state.eventProduceStateTransition.Asos,
			state.eventProduceStateTransition.Prepared,
			state.eventProduceStateTransition.StateTransitionOutputDigest,
			state.eventProduceStateTransition.ReportsPlusPrecursorDigest,
			state.eventProduceStateTransition.CommitQC,
		}:
			state.notifyOutcomeGenerationOfOpenedKVTransaction = false
		case nilOrEventReplayStateTransition <- EventReplayVerifiedStateTransition[RI]{nextBlock}:
		case nilOrEventReceivedKVTransaction <- EventAcknowledgedComputedStateTransition[RI]{state.highestComputedStateTransitionSeqNr}:
			state.notifyOutcomeGenerationOfReceivedKVTransaction = false
		case nilOrEventCommittedKVStateTransaction <- EventCommittedKVTransaction[RI]{state.highestCommittedToKVSeqNr}:
			state.notifyOutcomeGenerationOfCommittedKVTransaction = false
		case msg := <-state.chNetToStatePersistence:
			msg.msg.processStatePersistence(state, msg.sender)
		case ev := <-state.chOutcomeGenerationToStatePersistence:
			ev.processStatePersistence(state)
		case ev := <-state.chReportAttestationToStatePersistence:
			ev.processStatePersistence(state)
		case ev := <-state.blockSyncState.scheduler.Scheduled():
			ev.processStatePersistence(state)
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

func (state *statePersistenceState[RI]) eventKVTransactionRequest(ev EventKVTransactionRequest[RI]) {
	state.logger.Debug("StatePersistence: receivedEventKVTransactionRequest", commontypes.LogFields{
		"seqNr": ev.RoundCtx.SeqNr,
	})
	if ev.RoundCtx.SeqNr <= state.highestCommittedToKVSeqNr {
		state.logger.Debug("StatePersistence: dropping KV transaction request, we have already committed this transaction", commontypes.LogFields{
			"seqNr":              ev.RoundCtx.SeqNr,
			"committedToKVSeqNr": state.highestCommittedToKVSeqNr,
		})
		return
	}

	txn, err := state.kvStore.NewTransaction(ev.RoundCtx.SeqNr)
	if err != nil && errors.Is(err, ocr3_1types.ErrDuplicateTransaction) {
		state.logger.Debug("could not open new kv transaction", commontypes.LogFields{
			"seqNr": ev.RoundCtx.SeqNr,
			"err":   err,
		})
		return
	} else if err != nil {
		state.logger.Error("could not open new kv transaction", commontypes.LogFields{
			"seqNr": ev.RoundCtx.SeqNr,
			"err":   err,
		})
		return
	}

	// if ev.CommitQC != nil we are dealing with a replay
	if ev.CommitQC != nil {
		state.kvStoreTxn = &kvStoreTxn{
			txn,
			true,
		}
	} else {
		state.kvStoreTxn = &kvStoreTxn{
			txn,
			false,
		}
	}

	state.logger.Debug("new transaction", commontypes.LogFields{
		"txnSeqNr":           txn.SeqNr(),
		"committedToKVSeqNr": state.highestCommittedToKVSeqNr,
	})
	state.eventProduceStateTransition = EventProduceStateTransition[RI]{
		ev.RoundCtx,
		*state.kvStoreTxn,
		ev.Query,
		ev.Asos,
		ev.Prepared,
		ev.StateTransitionOutputDigest,
		ev.ReportsPlusPrecursorDigest,
		ev.CommitQC,
	}
	state.notifyOutcomeGenerationOfOpenedKVTransaction = true
}

func (state *statePersistenceState[RI]) eventComputedStateTransition(ev EventComputedStateTransition[RI]) {
	state.logger.Debug("received EventComputedStateTransition", commontypes.LogFields{
		"txnSeqNr": ev.SeqNr,
	})
	if state.kvStoreTxn != nil && state.kvStoreTxn.transaction != nil && ev.SeqNr != state.kvStoreTxn.transaction.SeqNr() {
		state.logger.Debug("StatePersistence: discarding EventComputedStateTransition, invalid seqNr", commontypes.LogFields{
			"evSeqNr":  ev.SeqNr,
			"txnSeqNr": state.kvStoreTxn.transaction.SeqNr(),
		})
		return
	}
	state.highestComputedStateTransitionSeqNr = ev.SeqNr
	state.notifyOutcomeGenerationOfReceivedKVTransaction = true

}

func (state *statePersistenceState[RI]) eventAttestedStateTransitionBlock(ev EventAttestedStateTransitionBlock[RI]) {
	seqNr := ev.AttestedStateTransitionBlock.StateTransitionBlock.SeqNr()
	state.logger.Debug("received EventAttestedStateTransitionBlock", commontypes.LogFields{
		"evSeqNr": seqNr,
	})

	if ev.AttestedStateTransitionBlock.StateTransitionBlock.SeqNr() != state.highestComputedStateTransitionSeqNr {
		state.logger.Critical("we received an attestedStateTransitionBlock out of order", commontypes.LogFields{
			"attestedStateTransitionBlockSeqNr":   ev.AttestedStateTransitionBlock.StateTransitionBlock.SeqNr(),
			"highestComputedStateTransitionSeqNr": state.highestComputedStateTransitionSeqNr,
		})
	}

	defer state.heardSeqNr(seqNr)

	if seqNr == state.highestPersistedStateTransitionBlockSeqNr+1 {
		if err := state.persist(ev.AttestedStateTransitionBlock); err != nil {
			state.logger.Error("failed to persist attested state transition", commontypes.LogFields{
				"error": err,
				"seqNr": seqNr,
			})
		}

		err := state.database.WriteStatePersistenceState(
			state.ctx, state.config.ConfigDigest,
			StatePersistenceState{
				state.highestPersistedStateTransitionBlockSeqNr,
			})
		if err != nil {
			state.logger.Error("Failed to persist the key-value sequence number", commontypes.LogFields{
				"SeqNr": seqNr,
				"err":   err,
			})
			return
		}
	}

	if state.kvStoreTxn == nil || state.kvStoreTxn.transaction == nil {
		state.logger.Error("State persistence protocol has not created a kv transaction yet", commontypes.LogFields{
			"highestComputedStateTransitionSeqNr": state.highestComputedStateTransitionSeqNr,
		})
		return
	}

	if state.kvStoreTxn.transaction.SeqNr() != seqNr {
		state.logger.Error("Cannot commit to the key value store, committed sequence number is different from the sequence number of the state persistence key value store transaction",
			commontypes.LogFields{
				"committedSeqNr": seqNr,
				"txnSeqNr":       state.kvStoreTxn.transaction.SeqNr(),
			})
		return
	}

	err := state.kvStoreTxn.transaction.Commit()
	if err != nil {
		state.logger.Warn("Failed to commit to the key value store", commontypes.LogFields{
			"committedSeqNr": seqNr,
			"err":            err,
		})
		state.kvStoreTxn.transaction.Discard()
		state.kvStoreTxn = nil
		return
	}
	state.kvStoreTxn = nil

	if first, isNotEmpty := state.retrievedStateTransitionBlockSeqNrQueue.Peek(); isNotEmpty {
		if *first == seqNr {
			if popped, ok := state.retrievedStateTransitionBlockSeqNrQueue.Pop(); ok {
				state.logger.Debug("popped from queue committed to KV seqNr", commontypes.LogFields{
					"committedSeqNr": seqNr,
					"poppedSeqNr":    popped,
				})
			}
		}
	}

	state.highestCommittedToKVSeqNr = seqNr
	state.logger.Debug("persisted transaction to KV store", commontypes.LogFields{
		"highestCommittedToKVSeqNr": state.highestCommittedToKVSeqNr,
	})
	state.clearStaleBlockSyncRequests()
	state.notifyOutcomeGenerationOfCommittedKVTransaction = true
}

func (state *statePersistenceState[RI]) eventAcknowledgedCommittedKVTransaction(ev EventAcknowledgedCommittedKVTransaction[RI]) {
	state.logger.Debug("received EventAcknowledgedCommittedKVTransaction", commontypes.LogFields{
		"txnSeqNr": ev.SeqNr,
	})
	if ev.SeqNr != state.highestCommittedToKVSeqNr {
		state.logger.Critical("we received an  AcknowledgedCommittedKVTransaction event out of order", commontypes.LogFields{
			"eventSeqNr":                ev.SeqNr,
			"highestCommittedToKVSeqNr": state.highestCommittedToKVSeqNr,
		})
	}
	state.highestAcknowledgedCommittedToKVSeqNr = ev.SeqNr
}

// This event is triggered upon outgen enters a new epoch, and it's purpose is to discard a transaction
// that did not get committed in case the epoch had timed out. We should only discard transactions
// from the common path, i.e. transactions that are not being replayed.
func (state *statePersistenceState[RI]) eventDiscardKVTransaction(ev EventDiscardKVTransaction[RI]) {
	state.logger.Debug("received EventDiscardKVTransaction", commontypes.LogFields{
		"evSeqNr": ev.SeqNr,
	})
	if state.kvStoreTxn != nil && state.kvStoreTxn.transaction != nil &&
		ev.SeqNr == state.kvStoreTxn.transaction.SeqNr() {
		state.kvStoreTxn.transaction.Discard()
		state.kvStoreTxn = nil
	}
}

func (state *statePersistenceState[RI]) eventStateSyncRequest(ev EventStateSyncRequest[RI]) {
	state.logger.Debug("received EventStateSyncRequest", commontypes.LogFields{
		"heardSeqNr": ev.SeqNr,
	})
	state.heardSeqNr(ev.SeqNr)
}

func (state *statePersistenceState[RI]) heardSeqNr(seqNr uint64) {
	if seqNr > state.highestHeardSeqNr {
		state.highestHeardSeqNr = seqNr
		state.logger.Debug("increased highestHeardSeqNr seqNr", commontypes.LogFields{
			"heardSeqNr": seqNr,
		})
		state.highestHeardIncreased()
	}
}

func (state *statePersistenceState[RI]) persist(astb AttestedStateTransitionBlock) error {
	if astb.StateTransitionBlock.SeqNr() != state.highestPersistedStateTransitionBlockSeqNr+1 {
		return fmt.Errorf("cannot persist out of order state transition block: expected %d, got %d",
			state.highestPersistedStateTransitionBlockSeqNr+1,
			astb.StateTransitionBlock.SeqNr(),
		)
	}
	err := state.database.WriteAttestedStateTransitionBlock(state.ctx,
		state.config.ConfigDigest,
		astb.StateTransitionBlock.SeqNr(),
		astb,
	)
	if err != nil {
		return err
	}
	err = state.database.WriteStatePersistenceState(
		state.ctx, state.config.ConfigDigest,
		StatePersistenceState{
			astb.StateTransitionBlock.SeqNr(),
		})
	if err != nil {
		return err
	}
	state.highestPersistedStateTransitionBlockSeqNr = astb.StateTransitionBlock.SeqNr()
	state.logger.Debug("persisted state transition block", commontypes.LogFields{
		"highestPersisted": state.highestPersistedStateTransitionBlockSeqNr,
	})
	return nil
}

func (state *statePersistenceState[RI]) nextStateTransitionBlockToReplay() (AttestedStateTransitionBlock, bool) {
	if next, isNotEmpty := state.retrievedStateTransitionBlockSeqNrQueue.Peek(); isNotEmpty {
		if *next <= state.highestAcknowledgedCommittedToKVSeqNr {
			state.logger.Debug("no need to replay state transition block", commontypes.LogFields{
				"nextSeqNr":                             *next,
				"highestCommittedToKVSeqNr":             state.highestCommittedToKVSeqNr,
				"highestAcknowledgedCommittedToKVSeqNr": state.highestAcknowledgedCommittedToKVSeqNr,
			})
			state.retrievedStateTransitionBlockSeqNrQueue.Pop()
			return AttestedStateTransitionBlock{}, false
		}
		if *next > state.highestAcknowledgedCommittedToKVSeqNr+1 {
			state.logger.Error("cannot replay state transition block, missing state transitions in between", commontypes.LogFields{
				"nextBlockSeqNr":                        *next,
				"highestCommittedToKVSeqNr":             state.highestCommittedToKVSeqNr,
				"highestAcknowledgedCommittedToKVSeqNr": state.highestAcknowledgedCommittedToKVSeqNr,
			})
			return AttestedStateTransitionBlock{}, false
		}
		// else nextSeqNr = state.highestAcknowledgedCommittedToKVSeqNr + 1
		if state.kvStoreTxn != nil {
			state.logger.Trace("no need to replay state transition block yet, we have an open transaction in progress", commontypes.LogFields{
				"nextSeqNr":                 *next,
				"highestCommittedToKVSeqNr": state.highestCommittedToKVSeqNr,
			})
			return AttestedStateTransitionBlock{}, false
		}
		nextBlock, err := state.database.ReadAttestedStateTransitionBlock(state.ctx, state.config.ConfigDigest, *next)
		if err != nil {
			state.logger.Error("failed to read attested state transition block from database", commontypes.LogFields{
				"nextSeqNr": *next,
				"error":     err,
			})
			return AttestedStateTransitionBlock{}, false
		}
		if *next != nextBlock.StateTransitionBlock.SeqNr() {
			state.logger.Error("read from database block with unexpected sequence number", commontypes.LogFields{
				"nextSeqNr":                 *next,
				"stateTransitionBlockSeqNr": nextBlock.StateTransitionBlock.SeqNr(),
			})
			return AttestedStateTransitionBlock{}, false
		}
		state.logger.Debug("next state transition block to replay", commontypes.LogFields{
			"nextSeqNr": *next,
		})
		return nextBlock, true
	}
	return AttestedStateTransitionBlock{}, false
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
	chOutcomeGenerationToStatePersistence <-chan EventToStatePersistence[RI],
	chReportAttestationToStatePersistence <-chan EventToStatePersistence[RI],
	chStatePersistenceToOutcomeGeneration chan<- EventToOutcomeGeneration[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	kvStore ocr3_1types.KeyValueStore,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
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
	return &statePersistenceState[RI]{
		ctx,

		chNetToStatePersistence,
		chOutcomeGenerationToStatePersistence,
		chReportAttestationToStatePersistence,
		chStatePersistenceToOutcomeGeneration,
		config,
		database,
		id,
		kvStore,
		logger.MakeUpdated(commontypes.LogFields{"proto": "state"}),
		netSender,
		nil,
		0,
		0,
		0,
		0,
		0,
		true,
		queue.NewQueue[uint64](),

		EventProduceStateTransition[RI]{},
		false,
		false,
		false,
		blockSyncState[RI]{
			logger.MakeUpdated(commontypes.LogFields{"proto": "stateBlockSync"}),
			oracles,
			scheduler,
		},
		treeSyncState{},
	}
}
