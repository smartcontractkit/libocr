package protocol

import (
	"context"
	"math"
	"slices"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/protocol/requestergadget"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/RoSpaceDev/libocr/subprocesses"
	"github.com/google/btree"
)

const (
	// An oracle sends a STATE-SYNC-SUMMARY message every DeltaStateSyncHeartbeat
	DeltaStateSyncHeartbeat time.Duration = 1 * time.Second
)

func RunStateSync[RI any](
	ctx context.Context,

	chNetToStateSync <-chan MessageToStateSyncWithSender[RI],
	chOutcomeGenerationToStateSync <-chan EventToStateSync[RI],
	chReportAttestationToStateSync <-chan EventToStateSync[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	kvDb KeyValueDatabase,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	reportingPlugin ocr3_1types.ReportingPlugin[RI],
) {
	chNotificationToStateBlockReplay := make(chan struct{})
	chNotificationToStateDestroyIfNeeded := make(chan struct{})

	subs := subprocesses.Subprocesses{}
	defer subs.Wait()
	subs.Go(func() {
		RunStateSyncDestroyIfNeeded(ctx, logger, kvDb, chNotificationToStateDestroyIfNeeded)
	})
	subs.Go(func() {
		RunStateSyncReap(ctx, config, logger, database, kvDb)
	})
	subs.Go(func() {
		RunStateSyncBlockReplay(ctx, logger, kvDb, chNotificationToStateBlockReplay)
	})

	newStateSyncState(ctx,
		chNetToStateSync,
		chNotificationToStateBlockReplay,
		chNotificationToStateDestroyIfNeeded,
		chOutcomeGenerationToStateSync,
		chReportAttestationToStateSync,
		config, database, id, kvDb, logger, netSender).run()
}

type syncMode int

const (
	syncModeUnknown syncMode = iota
	syncModeBlock
	syncModeTree
	syncModeFetchSnapshotBlock
)

type stateSyncState[RI any] struct {
	ctx context.Context

	chNetToStateSync                     <-chan MessageToStateSyncWithSender[RI]
	chNotificationToStateBlockReplay     chan<- struct{}
	chNotificationToStateDestroyIfNeeded chan<- struct{}
	chOutcomeGenerationToStateSync       <-chan EventToStateSync[RI]
	chReportAttestationToStateSync       <-chan EventToStateSync[RI]
	config                               ocr3config.SharedConfig
	database                             Database
	id                                   commontypes.OracleID
	kvDb                                 KeyValueDatabase
	logger                               loghelper.LoggerWithContext
	netSender                            NetworkSender[RI]

	highestPersistedStateTransitionBlockSeqNr uint64
	lowestPersistedStateTransitionBlockSeqNr  uint64
	highestCommittedSeqNr                     uint64

	oracles []*syncOracle

	highestHeardSeqNr uint64

	blockSyncState blockSyncState[RI]
	treeSyncState  treeSyncState[RI]

	syncMode syncMode

	tSendSummary <-chan time.Time
}

type syncOracle struct {
	// lowestPersistedSeqNr is the lowest sequence number the oracle still has an attested
	// state transition block for
	lowestPersistedSeqNr uint64
	// highestCommittedSeqNr is the highest sequence number the oracle has committed to
	highestCommittedSeqNr uint64
	lastSummaryReceivedAt time.Time
}

func (stasy *stateSyncState[RI]) run() {
	stasy.refreshStateSyncState()
	stasy.logger.Info("StateSync: running", commontypes.LogFields{
		"highestPersistedStateTransitionBlockSeqNr": stasy.highestPersistedStateTransitionBlockSeqNr,
	})

	for {
		select {
		case msg := <-stasy.chNetToStateSync:
			msg.msg.processStateSync(stasy, msg.sender)
		case ev := <-stasy.chOutcomeGenerationToStateSync:
			ev.processStateSync(stasy)
		case ev := <-stasy.chReportAttestationToStateSync:
			ev.processStateSync(stasy)
		case <-stasy.tSendSummary:
			stasy.eventTSendSummaryTimeout()
		case <-stasy.blockSyncState.blockRequesterGadget.Ticker():
			stasy.blockSyncState.blockRequesterGadget.Tick()
		case <-stasy.treeSyncState.treeChunkRequesterGadget.Ticker():
			stasy.treeSyncState.treeChunkRequesterGadget.Tick()
		case <-stasy.ctx.Done():
		}

		// ensure prompt exit
		select {
		case <-stasy.ctx.Done():
			stasy.logger.Info("StateSync: exiting", nil)
			return
		default:
		}
	}
}

func (stasy *stateSyncState[RI]) pleaseTryToReplayBlock() {
	select {
	case stasy.chNotificationToStateBlockReplay <- struct{}{}:
	default:
	}
}

func (stasy *stateSyncState[RI]) pleaseDestroyStateIfNeeded() {
	select {
	case stasy.chNotificationToStateDestroyIfNeeded <- struct{}{}:
	default:
	}
}

func (stasy *stateSyncState[RI]) eventStateSyncRequest(ev EventStateSyncRequest[RI]) {
	stasy.logger.Debug("received EventStateSyncRequest", commontypes.LogFields{
		"heardSeqNr": ev.SeqNr,
	})

	stasy.pleaseTryToReplayBlock()

	if ev.SeqNr <= stasy.highestHeardSeqNr {
		return
	}

	stasy.logger.Debug("highest heard sequence number increased from EventStateSyncRequest", commontypes.LogFields{
		"old": stasy.highestHeardSeqNr,
		"new": ev.SeqNr,
	})
	stasy.highestHeardSeqNr = ev.SeqNr
	stasy.tryToKickStartSync()
}

func (stasy *stateSyncState[RI]) refreshStateSyncState() (ok bool) {
	kvReadTxn, err := stasy.kvDb.NewReadTransactionUnchecked()
	if err != nil {
		stasy.logger.Warn("failed to create new transaction", commontypes.LogFields{
			"error": err,
		})
		return
	}
	defer kvReadTxn.Discard()

	highestCommittedToKVSeqNr, err := kvReadTxn.ReadHighestCommittedSeqNr()
	if err != nil {
		stasy.logger.Error("failed to get highest committed to kv seq nr during refresh", commontypes.LogFields{
			"error": err,
		})
		return
	}

	stasy.highestCommittedSeqNr = highestCommittedToKVSeqNr

	if highestCommittedToKVSeqNr > stasy.highestPersistedStateTransitionBlockSeqNr {

		stasy.highestPersistedStateTransitionBlockSeqNr = highestCommittedToKVSeqNr
		stasy.reapBlockBuffer()
	}

	lowestPersistedSeqNr, err := kvReadTxn.ReadLowestPersistedSeqNr()
	if err != nil {
		stasy.logger.Warn("failed to read lowest persisted seq nr", commontypes.LogFields{
			"error": err,
		})
		return
	}
	stasy.lowestPersistedStateTransitionBlockSeqNr = lowestPersistedSeqNr

	treeSyncStatus, err := kvReadTxn.ReadTreeSyncStatus()
	if err != nil {
		stasy.logger.Warn("failed to read tree sync status", commontypes.LogFields{
			"error": err,
		})
		return
	}
	stasy.treeSyncState.treeSyncPhase = treeSyncStatus.Phase
	stasy.treeSyncState.targetSeqNr = treeSyncStatus.TargetSeqNr
	stasy.treeSyncState.targetStateRootDigest = treeSyncStatus.TargetStateRootDigest
	stasy.treeSyncState.pendingKeyDigestRanges = treeSyncStatus.PendingKeyDigestRanges
	ok = true
	return
}

func (stasy *stateSyncState[RI]) eventTSendSummaryTimeout() {
	defer func() {
		stasy.tSendSummary = time.After(DeltaStateSyncHeartbeat)
	}()
	if !stasy.refreshStateSyncState() {
		return
	}
	if stasy.treeSyncState.treeSyncPhase != TreeSyncPhaseInactive {

		stasy.netSender.Broadcast(MessageStateSyncSummary[RI]{
			math.MaxUint64,
			0,
		})
		return
	}
	stasy.netSender.Broadcast(MessageStateSyncSummary[RI]{
		stasy.lowestPersistedStateTransitionBlockSeqNr,
		stasy.highestCommittedSeqNr,
	})
}

func (stasy *stateSyncState[RI]) messageStateSyncSummary(msg MessageStateSyncSummary[RI], sender commontypes.OracleID) {
	stasy.logger.Debug("received messageStateSyncSummary", commontypes.LogFields{
		"sender":                   sender,
		"msgLowestPersistedSeqNr":  msg.LowestPersistedSeqNr,
		"msgHighestCommittedSeqNr": msg.HighestCommittedSeqNr,
	})
	stasy.oracles[sender] = &syncOracle{
		msg.LowestPersistedSeqNr,
		msg.HighestCommittedSeqNr,
		time.Now(),
	}
	stasy.updateHighestHeardFromSummaries()

	stasy.tryToKickStartSync()
}

func (stasy *stateSyncState[RI]) summaryFreshnessCutoff() time.Duration {
	return stasy.config.DeltaProgress / 4
}

type honestOraclePruneStatus int

const (
	_ honestOraclePruneStatus = iota
	honestOraclePruneStatusCannotDecideYet
	honestOraclePruneStatusWouldNotPrune
	honestOraclePruneStatusWouldPrune
)

// for a given sequence number, is there guaranteed to be at least one honest oracle that could help us
// sync the committed state at `seqNr`?
// - honestOraclePruneStatusWouldNotPrune: yes
// - honestOraclePruneStatusWouldNotPrune: no
// - honestOraclePruneStatusCannotDecideYet: we don't have enough information to answer the question
func (stasy *stateSyncState[RI]) findSomeHonestOraclePruneStatus(seqNr uint64) honestOraclePruneStatus {
	wouldNotPrune := 0
	wouldPrune := 0

	for i, oracle := range stasy.oracles {
		if time.Since(oracle.lastSummaryReceivedAt) > stasy.summaryFreshnessCutoff() {

			continue
		}

		if commontypes.OracleID(i) == stasy.id {

			continue
		}

		if oracle.lowestPersistedSeqNr <= seqNr {
			wouldNotPrune++
		} else {
			wouldPrune++
		}
	}

	if wouldNotPrune > stasy.config.F {
		return honestOraclePruneStatusWouldNotPrune
	} else if wouldPrune > stasy.config.F {
		return honestOraclePruneStatusWouldPrune
	} else {
		return honestOraclePruneStatusCannotDecideYet
	}
}

func (stasy *stateSyncState[RI]) updateHighestHeardFromSummaries() {
	hiSeqNrs := make([]uint64, len(stasy.oracles))
	for i, oracle := range stasy.oracles {

		hiSeqNrs[i] = oracle.highestCommittedSeqNr
	}

	slices.Sort(hiSeqNrs)
	candidateHighestHeard := hiSeqNrs[len(hiSeqNrs)-1-stasy.config.F]
	if candidateHighestHeard <= stasy.highestHeardSeqNr {
		return
	}

	stasy.logger.Debug("highest heard sequence number increased from MessageStateSyncSummary", commontypes.LogFields{
		"old": stasy.highestHeardSeqNr,
		"new": candidateHighestHeard,
	})
	stasy.highestHeardSeqNr = candidateHighestHeard
}

type blockSyncOrTreeSyncDecision int

const (
	_ blockSyncOrTreeSyncDecision = iota
	blockSyncOrTreeSyncDecisionCannotDecideYet
	blockSyncOrTreeSyncDecisionBlockSync
	blockSyncOrTreeSyncDecisionTreeSync
)

func (stasy *stateSyncState[RI]) decideBlockSyncOrTreeSyncBasedOnSummariesAndHighestHeard() blockSyncOrTreeSyncDecision {

	ourStartingPointSeqNr := max(stasy.highestCommittedSeqNr, stasy.highestPersistedStateTransitionBlockSeqNr) + 1
	switch stasy.findSomeHonestOraclePruneStatus(ourStartingPointSeqNr) {
	case honestOraclePruneStatusWouldNotPrune:
		return blockSyncOrTreeSyncDecisionBlockSync
	case honestOraclePruneStatusWouldPrune:
		return blockSyncOrTreeSyncDecisionTreeSync
	case honestOraclePruneStatusCannotDecideYet:
	}
	return blockSyncOrTreeSyncDecisionCannotDecideYet
}

func (stasy *stateSyncState[RI]) pickSomeTreeSyncTarget() (uint64, bool) {
	if snapshotSeqNr(stasy.highestHeardSeqNr) == stasy.highestHeardSeqNr {
		return stasy.highestHeardSeqNr, true
	} else {
		snapshotIndex := snapshotIndexFromSeqNr(stasy.highestHeardSeqNr)
		if snapshotIndex > 0 {
			return maxSeqNrWithSnapshotIndex(snapshotIndex - 1), true
		} else {
			return 0, false
		}
	}
}

func (stasy *stateSyncState[RI]) needToRetargetTreeSync() bool {
	switch stasy.findSomeHonestOraclePruneStatus(stasy.treeSyncState.targetSeqNr) {
	case honestOraclePruneStatusWouldNotPrune:
		return false
	case honestOraclePruneStatusWouldPrune:

		return true
	case honestOraclePruneStatusCannotDecideYet:
		return false
	}
	return false
}

func (stasy *stateSyncState[RI]) treeSyncCompleted() {
	stasy.syncMode = syncModeUnknown
	stasy.tryToKickStartSync()
}

func (stasy *stateSyncState[RI]) treeSyncNeedsSnapshotBlock() {
	stasy.syncMode = syncModeFetchSnapshotBlock
	stasy.blockSyncState.blockRequesterGadget.PleaseRecheckPendingItems()
}

func (stasy *stateSyncState[RI]) tryToKickStartSync() {

	switch stasy.syncMode {
	case syncModeTree:

		stasy.evolveTreeSyncPhase()
		return
	case syncModeFetchSnapshotBlock:

		if stasy.needToRetargetTreeSync() {
			stasy.logger.Warn("not guaranteed to be able to fetch the tree-sync target block, giving up", commontypes.LogFields{
				"targetSeqNr": stasy.treeSyncState.targetSeqNr,
			})
			stasy.syncMode = syncModeUnknown
		} else {
			return
		}
	case syncModeBlock:
		stasy.tryCompleteBlockSync()
		stasy.blockSyncState.blockRequesterGadget.PleaseRecheckPendingItems()

	case syncModeUnknown:

	}

	if !stasy.refreshStateSyncState() {
		stasy.logger.Warn("cannot kick start sync, failed to refresh stateSyncState", nil)
		return
	}

	decision := stasy.decideBlockSyncOrTreeSyncBasedOnSummariesAndHighestHeard()

	if decision == blockSyncOrTreeSyncDecisionCannotDecideYet {
		stasy.logger.Debug("cannot decide whether to block-sync or tree-sync, yet", nil)
		return
	}

	if stasy.treeSyncState.treeSyncPhase != TreeSyncPhaseInactive || decision == blockSyncOrTreeSyncDecisionTreeSync {
		stasy.logger.Debug("switching to tree-sync mode", nil)
		stasy.syncMode = syncModeTree
		stasy.evolveTreeSyncPhase()
		return
	}

	if decision == blockSyncOrTreeSyncDecisionBlockSync {
		stasy.logger.Debug("switching to block-sync mode", nil)
		stasy.syncMode = syncModeBlock
		// requester gadget will take it from here
		stasy.blockSyncState.blockRequesterGadget.PleaseRecheckPendingItems()
		return
	}
}

func newStateSyncState[RI any](
	ctx context.Context,
	chNetToStateSync <-chan MessageToStateSyncWithSender[RI],
	chNotificationToStateBlockReplay chan<- struct{},
	chNotificationToStateDestroyIfNeeded chan<- struct{},
	chOutcomeGenerationToStateSync <-chan EventToStateSync[RI],
	chReportAttestationToStateSync <-chan EventToStateSync[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	kvDb KeyValueDatabase,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
) *stateSyncState[RI] {
	oracles := make([]*syncOracle, 0)
	for i := 0; i < config.N(); i++ {
		oracles = append(oracles, &syncOracle{
			0,
			0,
			time.Time{},
		})
	}

	stasy := &stateSyncState[RI]{
		ctx,

		chNetToStateSync,
		chNotificationToStateBlockReplay,
		chNotificationToStateDestroyIfNeeded,
		chOutcomeGenerationToStateSync,
		chReportAttestationToStateSync,
		config,
		database,
		id,
		kvDb,
		logger.MakeUpdated(commontypes.LogFields{"proto": "stasy"}),
		netSender,
		0,
		0,
		0,

		oracles,
		0,

		blockSyncState[RI]{
			logger.MakeUpdated(commontypes.LogFields{"proto": "stasy/block"}),
			nil, // defined right below
			btree.NewG(2, func(a AttestedStateTransitionBlock, b AttestedStateTransitionBlock) bool {
				return a.StateTransitionBlock.SeqNr() < b.StateTransitionBlock.SeqNr()
			}),
		},
		treeSyncState[RI]{
			logger.MakeUpdated(commontypes.LogFields{"proto": "stasy/tree"}),
			nil, // defined right below
			TreeSyncPhaseInactive,
			0,
			StateRootDigest{},
			PendingKeyDigestRanges{},
		},
		syncModeUnknown,
		time.After(DeltaStateSyncHeartbeat),
	}

	stasy.blockSyncState.blockRequesterGadget = requestergadget.NewRequesterGadget[seqNrRange](
		config.N(),
		DeltaMinBlockSyncRequest,
		stasy.sendBlockSyncRequest,
		stasy.getPendingBlocksToRequest,
		stasy.getBlockSyncSeeders,
	)
	stasy.treeSyncState.treeChunkRequesterGadget = requestergadget.NewRequesterGadget[treeSyncChunkRequestItem](
		config.N(),
		DeltaMinTreeSyncRequest,
		stasy.sendTreeSyncChunkRequest,
		stasy.getPendingTreeSyncChunksToRequest,
		stasy.getTreeSyncChunkSeeders,
	)
	return stasy
}
