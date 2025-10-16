package protocol

import (
	"bytes"
	"fmt"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/jmt"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/protocol/requestergadget"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
)

const (
	// Maximum delay between a TREE-SYNC-REQ and TREE-SYNC-CHUNK response. We'll try
	// with another oracle if we don't get a response in this time.
	DeltaMaxTreeSyncRequest time.Duration = 1 * time.Second
	// Minimum delay between two consecutive BLOCK-SYNC-REQ requests
	DeltaMinTreeSyncRequest = 10 * time.Millisecond

	// The maximum number of key-value pairs that an oracle will send in a single tree-sync chunk
	MaxTreeSyncChunkKeys = 128
	// Maximum number of bytes in of the combined keys and values length in a chunk.

	MaxTreeSyncChunkKeysPlusValuesLength = 2 * (ocr3_1types.MaxMaxKeyValueKeyLength + ocr3_1types.MaxMaxKeyValueValueLength)

	MaxMaxParallelTreeSyncChunkFetches = 8
)

func (stasy *stateSyncState[RI]) maxParallelTreeSyncChunkFetches() int {
	return max(1, min(MaxMaxParallelTreeSyncChunkFetches, stasy.config.N()-1))
}

func (stasy *stateSyncState[RI]) newPendingKeyDigestRanges() PendingKeyDigestRanges {
	numSplits := stasy.maxParallelTreeSyncChunkFetches()
	split, err := splitKeyDigestSpaceEvenly(numSplits)
	if err != nil {
		stasy.logger.Error("failed to create even key digest range split, reverting to a single range", commontypes.LogFields{
			"numSplits": numSplits,
			"error":     err,
		})
		return NewPendingKeyDigestRanges([]KeyDigestRange{{jmt.MinDigest, jmt.MaxDigest}})
	}
	return NewPendingKeyDigestRanges(split)
}

type treeSyncChunkRequestItem struct {
	targetSeqNr    uint64
	keyDigestRange KeyDigestRange
}

type treeSyncState[RI any] struct {
	logger                   commontypes.Logger
	treeChunkRequesterGadget *requestergadget.RequesterGadget[treeSyncChunkRequestItem]

	treeSyncPhase         TreeSyncPhase
	targetSeqNr           uint64
	targetStateRootDigest StateRootDigest

	pendingKeyDigestRanges PendingKeyDigestRanges
}

func (stasy *stateSyncState[RI]) sendTreeSyncChunkRequest(item treeSyncChunkRequestItem, target commontypes.OracleID) (*requestergadget.RequestInfo, bool) {
	if stasy.syncMode != syncModeTree || stasy.treeSyncState.treeSyncPhase != TreeSyncPhaseActive || stasy.treeSyncState.targetStateRootDigest == (StateRootDigest{}) {
		return nil, false
	}
	stasy.treeSyncState.logger.Debug("sending MessageTreeSyncChunkRequest", commontypes.LogFields{
		"targetSeqNr":    item.targetSeqNr,
		"keyDigestRange": item.keyDigestRange,
		"target":         target,
	})
	msg := MessageTreeSyncChunkRequest[RI]{
		nil,
		item.targetSeqNr,
		item.keyDigestRange.StartIndex,
		item.keyDigestRange.EndInclIndex,
	}
	stasy.netSender.SendTo(msg, target)
	return &requestergadget.RequestInfo{
		time.Now().Add(DeltaMaxTreeSyncRequest),
	}, true
}

func (stasy *stateSyncState[RI]) getPendingTreeSyncChunksToRequest() []treeSyncChunkRequestItem {
	if stasy.syncMode != syncModeTree || stasy.treeSyncState.treeSyncPhase != TreeSyncPhaseActive || stasy.treeSyncState.targetStateRootDigest == (StateRootDigest{}) {
		return nil
	}
	var pending []treeSyncChunkRequestItem
	for _, keyDigestRange := range stasy.treeSyncState.pendingKeyDigestRanges.All() {
		pending = append(pending, treeSyncChunkRequestItem{stasy.treeSyncState.targetSeqNr, keyDigestRange})
	}
	return pending
}

func (stasy *stateSyncState[RI]) getTreeSyncChunkSeeders(_ treeSyncChunkRequestItem) map[commontypes.OracleID]struct{} {
	seeders := make(map[commontypes.OracleID]struct{})
	for oid := range stasy.oracles {

		if commontypes.OracleID(oid) == stasy.id {
			continue
		}
		seeders[commontypes.OracleID(oid)] = struct{}{}
	}
	return seeders
}

func (stasy *stateSyncState[RI]) evolveTreeSyncPhase() {
	if stasy.syncMode != syncModeTree {
		return
	}

	if !stasy.refreshStateSyncState() {
		return
	}

	stasy.treeSyncState.logger.Debug("trying to evolve tree-sync phase", commontypes.LogFields{
		"phase":       stasy.treeSyncState.treeSyncPhase,
		"targetSeqNr": stasy.treeSyncState.targetSeqNr,
	})

	switch stasy.treeSyncState.treeSyncPhase {
	case TreeSyncPhaseInactive:
		newTargetSeqNr, found := stasy.pickSomeTreeSyncTarget()
		if !found {
			return
		}
		stasy.treeSyncState.logger.Debug("initializing new tree-sync", commontypes.LogFields{
			"newTargetSeqNr": newTargetSeqNr,
		})

		newTreeSyncStatus := TreeSyncStatus{
			TreeSyncPhaseWaiting,
			newTargetSeqNr,
			StateRootDigest{},
			stasy.newPendingKeyDigestRanges(),
		}

		kvReadWriteTxn, err := stasy.kvDb.NewSerializedReadWriteTransactionUnchecked()
		if err != nil {
			return
		}
		defer kvReadWriteTxn.Discard()

		if err := kvReadWriteTxn.WriteTreeSyncStatus(newTreeSyncStatus); err != nil {
			stasy.treeSyncState.logger.Error("failed to write tree-sync status", commontypes.LogFields{
				"err": err,
			})
			return
		}
		if err := kvReadWriteTxn.Commit(); err != nil {
			stasy.treeSyncState.logger.Error("failed to commit", commontypes.LogFields{
				"err": err,
			})
			return
		}
		stasy.refreshStateSyncState()
		return
	case TreeSyncPhaseWaiting:
		stasy.treeSyncState.logger.Debug("tree-sync waiting for key-value store cleanup 🧹", nil)
		stasy.pleaseDestroyStateIfNeeded()
		return
	case TreeSyncPhaseActive:
		if stasy.needToRetargetTreeSync() {
			stasy.treeSyncState.logger.Debug("not enough oracles to help us tree-sync to current target, we must re-target", commontypes.LogFields{
				"targetSeqNr": stasy.treeSyncState.targetSeqNr,
			})

			newTargetSeqNr, found := stasy.pickSomeTreeSyncTarget()
			if !found {
				return
			}

			stasy.treeSyncState.logger.Debug("tree-sync needed to re-target, and we found a new target", commontypes.LogFields{
				"targetSeqNr":    stasy.treeSyncState.targetSeqNr,
				"newTargetSeqNr": newTargetSeqNr,
			})

			newTreeSyncStatus := TreeSyncStatus{
				TreeSyncPhaseWaiting,
				newTargetSeqNr,
				StateRootDigest{},
				stasy.newPendingKeyDigestRanges(),
			}

			kvReadWriteTxn, err := stasy.kvDb.NewSerializedReadWriteTransactionUnchecked()
			if err != nil {
				return
			}
			defer kvReadWriteTxn.Discard()

			if err := kvReadWriteTxn.WriteTreeSyncStatus(newTreeSyncStatus); err != nil {
				stasy.treeSyncState.logger.Error("failed to write tree-sync status", commontypes.LogFields{
					"err": err,
				})
				return
			}
			if err := kvReadWriteTxn.Commit(); err != nil {
				stasy.treeSyncState.logger.Error("failed to commit", commontypes.LogFields{
					"err": err,
				})
				return
			}
			stasy.refreshStateSyncState()
			return
		} else {
			// our target seq nr is fine, and we are active
			// yield to block sync to fetch the target state root digest if necessary
			if stasy.treeSyncState.targetStateRootDigest == (StateRootDigest{}) {
				stasy.treeSyncState.logger.Debug("tree-sync yielding to block-sync to fetch the target state root digest", commontypes.LogFields{
					"targetSeqNr": stasy.treeSyncState.targetSeqNr,
				})
				stasy.treeSyncNeedsSnapshotBlock()
				return
			}

		}
	}
}

func (stasy *stateSyncState[RI]) acceptTreeSyncTargetBlockFromBlockSync(block AttestedStateTransitionBlock) error {
	if stasy.syncMode != syncModeFetchSnapshotBlock {
		return fmt.Errorf("not in fetch snapshot block mode")
	}

	if !stasy.refreshStateSyncState() {
		return fmt.Errorf("not accepting block without refreshed state")
	}

	if stasy.treeSyncState.treeSyncPhase != TreeSyncPhaseActive {
		return fmt.Errorf("not accepting block in unexpected tree-sync phase %v", stasy.treeSyncState.treeSyncPhase)
	}

	seqNr := block.StateTransitionBlock.SeqNr()
	if seqNr != stasy.treeSyncState.targetSeqNr {
		return fmt.Errorf("tree-sync target block sequence number does not match expected target sequence number")
	}

	stateRootDigest := block.StateTransitionBlock.StateRootDigest

	kvReadWriteTxn, err := stasy.kvDb.NewSerializedReadWriteTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("failed to create kv read/write transaction: %w", err)
	}
	defer kvReadWriteTxn.Discard()
	if err = kvReadWriteTxn.WriteTreeSyncStatus(TreeSyncStatus{
		stasy.treeSyncState.treeSyncPhase,
		stasy.treeSyncState.targetSeqNr,
		stateRootDigest,
		stasy.treeSyncState.pendingKeyDigestRanges,
	}); err != nil {
		return fmt.Errorf("failed to write tree-sync status: %w", err)
	}

	if err = kvReadWriteTxn.WriteAttestedStateTransitionBlock(seqNr, block); err != nil {
		return fmt.Errorf("failed to write attested state transition block: %w", err)
	}

	if err = kvReadWriteTxn.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	stasy.treeSyncState.logger.Debug("tree-sync accepted verified state root digest", commontypes.LogFields{
		"targetSeqNr": seqNr,
		"rootDigest":  stateRootDigest,
	})

	stasy.syncMode = syncModeTree
	stasy.refreshStateSyncState()
	stasy.treeSyncState.treeChunkRequesterGadget.PleaseRecheckPendingItems()
	return nil
}

func (stasy *stateSyncState[RI]) messageTreeSyncChunkRequest(msg MessageTreeSyncChunkRequest[RI], sender commontypes.OracleID) {
	stasy.treeSyncState.logger.Debug("received MessageTreeSyncChunkRequest", commontypes.LogFields{
		"sender":     sender,
		"toSeqNr":    msg.ToSeqNr,
		"startIndex": msg.StartIndex,
	})

	if !mustTakeSnapshot(msg.ToSeqNr) {
		stasy.treeSyncState.logger.Warn("dropping MessageTreeSyncChunkRequest with invalid SeqNr", commontypes.LogFields{
			"toSeqNr": msg.ToSeqNr,
		})
		return
	}

	kvReadTxn, err := stasy.kvDb.NewReadTransactionUnchecked()
	if err != nil {
		stasy.logger.Warn("failed to create new transaction", commontypes.LogFields{
			"error": err,
		})
		return
	}
	defer kvReadTxn.Discard()

	highestCommittedSeqNr, err := kvReadTxn.ReadHighestCommittedSeqNr()
	if err != nil {
		stasy.logger.Warn("failed to read highest committed seq nr", commontypes.LogFields{
			"error": err,
		})
		return
	}

	lowestPersistedSeqNr, err := kvReadTxn.ReadLowestPersistedSeqNr()
	if err != nil {
		stasy.logger.Warn("failed to read lowest persisted seq nr", commontypes.LogFields{
			"error": err,
		})
		return
	}

	treeSyncStatus, err := kvReadTxn.ReadTreeSyncStatus()
	if err != nil {
		stasy.logger.Warn("failed to read tree sync status", commontypes.LogFields{
			"error": err,
		})
		return
	}

	if treeSyncStatus.Phase != TreeSyncPhaseInactive || !(lowestPersistedSeqNr <= msg.ToSeqNr && msg.ToSeqNr <= highestCommittedSeqNr) {
		kvReadTxn.Discard()
		stasy.treeSyncState.logger.Debug("sending MessageTreeSyncChunkResponse to go-away", commontypes.LogFields{
			"sender":                sender,
			"toSeqNr":               msg.ToSeqNr,
			"lowestPersistedSeqNr":  lowestPersistedSeqNr,
			"highestCommittedSeqNr": highestCommittedSeqNr,
			"treeSyncPhase":         treeSyncStatus.Phase,
		})
		stasy.netSender.SendTo(MessageTreeSyncChunkResponse[RI]{
			msg.RequestHandle,
			msg.ToSeqNr,
			msg.StartIndex,
			jmt.Digest{},
			true,
			jmt.Digest{},
			nil,
			nil,
		}, sender)
		return
	}

	endInclIndex, boundingLeaves, keyValues, err := kvReadTxn.ReadTreeSyncChunk(
		msg.ToSeqNr,
		msg.StartIndex,
		msg.EndInclIndex,
	)
	if err != nil {
		stasy.treeSyncState.logger.Warn("failed to read chunk", commontypes.LogFields{
			"sender":     sender,
			"ToSeqNr":    msg.ToSeqNr,
			"startIndex": msg.StartIndex,
			"err":        err,
		})
		return
	}

	chunk := MessageTreeSyncChunkResponse[RI]{
		msg.RequestHandle,
		msg.ToSeqNr,
		msg.StartIndex,
		msg.EndInclIndex,
		false,
		endInclIndex,
		keyValues,
		boundingLeaves,
	}

	stasy.treeSyncState.logger.Debug("sent MessageTreeSyncChunkResponse", commontypes.LogFields{
		"target":         sender,
		"toSeqNr":        msg.ToSeqNr,
		"startIndex":     fmt.Sprintf("%x", msg.StartIndex),
		"endInclIndex":   fmt.Sprintf("%x", endInclIndex),
		"proofLen":       proofLen(boundingLeaves),
		"keyValuesCount": len(keyValues),
	})

	stasy.netSender.SendTo(chunk, sender)
}

func proofLen(boundingLeaves []jmt.BoundingLeaf) int {
	proofLen := 0
	for _, bl := range boundingLeaves {
		proofLen += len(bl.Siblings)
	}
	return proofLen
}

func (stasy *stateSyncState[RI]) messageTreeSyncChunkResponse(msg MessageTreeSyncChunkResponse[RI], sender commontypes.OracleID) {
	msgSeqNr := msg.ToSeqNr
	requestedKeyDigestRange := KeyDigestRange{msg.StartIndex, msg.RequestEndInclIndex}
	item := treeSyncChunkRequestItem{
		msgSeqNr,
		requestedKeyDigestRange,
	}
	if !stasy.treeSyncState.treeChunkRequesterGadget.CheckAndMarkResponse(item, sender) {
		stasy.treeSyncState.logger.Warn("dropping MessageTreeSyncChunkResponse: check and mark response failed", commontypes.LogFields{
			"sender":   sender,
			"msgSeqNr": msgSeqNr,
		})
		return
	}

	if msg.GoAway {
		stasy.treeSyncState.logger.Warn("dropping MessageTreeSyncChunkResponse: go-away", commontypes.LogFields{
			"sender":   sender,
			"msgSeqNr": msgSeqNr,
		})
		stasy.treeSyncState.treeChunkRequesterGadget.MarkGoAwayResponse(item, sender)
		return
	}

	if !(bytes.Compare(msg.EndInclIndex[:], msg.RequestEndInclIndex[:]) <= 0) {
		stasy.treeSyncState.logger.Warn("dropping MessageTreeSyncChunkResponse: end incl index is out of bounds", commontypes.LogFields{
			"sender":              sender,
			"msgSeqNr":            msgSeqNr,
			"requestEndInclIndex": msg.RequestEndInclIndex,
			"endInclIndex":        msg.EndInclIndex,
		})
		stasy.treeSyncState.treeChunkRequesterGadget.MarkBadResponse(item, sender)
		return
	}

	receivedKeyDigestRange := KeyDigestRange{msg.StartIndex, msg.EndInclIndex}

	stasy.treeSyncState.logger.Debug("received MessageTreeSyncChunkResponse", commontypes.LogFields{
		"sender":         sender,
		"startIndex":     fmt.Sprintf("%x", msg.StartIndex),
		"endInclIndex":   fmt.Sprintf("%x", msg.EndInclIndex),
		"proofLen":       proofLen(msg.BoundingLeaves),
		"keyValuesCount": len(msg.KeyValues),
	})

	if !stasy.refreshStateSyncState() {
		stasy.treeSyncState.logger.Warn("dropping MessageTreeSyncChunkResponse: could not refresh state", commontypes.LogFields{
			"sender":   sender,
			"msgSeqNr": msgSeqNr,
		})
		return
	}

	// Make sure that we already have the target state root digest
	if stasy.treeSyncState.targetStateRootDigest == (StateRootDigest{}) {
		stasy.treeSyncState.logger.Warn("dropping MessageTreeSyncChunkResponse: we do not have target state root digest, yet", commontypes.LogFields{
			"sender":   sender,
			"msgSeqNr": msgSeqNr,
		})
		return
	}

	// We have not sent a tree sync request yet or the response arrived way too late and we have already synced
	if stasy.treeSyncState.treeSyncPhase != TreeSyncPhaseActive {
		stasy.treeSyncState.logger.Warn("dropping MessageTreeSyncChunkResponse: tree-sync is not active", commontypes.LogFields{
			"sender":   sender,
			"msgSeqNr": msgSeqNr,
			"phase":    stasy.treeSyncState.treeSyncPhase,
			"mode":     stasy.syncMode,
		})
		return
	}

	// Check that target seqNrs match
	if stasy.treeSyncState.targetSeqNr != msgSeqNr {
		stasy.treeSyncState.logger.Warn("dropping MessageTreeSyncChunkResponse: message SeqNr does not match expected target", commontypes.LogFields{
			"sender":        sender,
			"msgSeqNr":      msgSeqNr,
			"expectedSeqNr": stasy.treeSyncState.targetSeqNr,
		})
		return
	}

	kvReadWriteTxn, err := stasy.kvDb.NewSerializedReadWriteTransactionUnchecked()
	if err != nil {
		stasy.treeSyncState.logger.Warn("could not create kv read/write transaction", commontypes.LogFields{
			"err": err,
		})
	}
	defer kvReadWriteTxn.Discard()

	// Verify and write chunk
	verifyAndWriteTreeSyncChunkResult, err := kvReadWriteTxn.VerifyAndWriteTreeSyncChunk(
		stasy.treeSyncState.targetStateRootDigest,
		stasy.treeSyncState.targetSeqNr,
		msg.StartIndex,
		msg.EndInclIndex,
		msg.BoundingLeaves,
		msg.KeyValues,
	)

	switch verifyAndWriteTreeSyncChunkResult {
	case VerifyAndWriteTreeSyncChunkResultUnrelatedError:
		stasy.treeSyncState.logger.Warn("failed to apply chunk", commontypes.LogFields{
			"sender": sender,
			"err":    err,
		})
		return
	case VerifyAndWriteTreeSyncChunkResultByzantine:
		stasy.treeSyncState.logger.Warn("byzantine chunk, marking as bad response", commontypes.LogFields{
			"sender":     sender,
			"startIndex": fmt.Sprintf("%x", msg.StartIndex),
			"err":        err,
		})
		stasy.treeSyncState.treeChunkRequesterGadget.MarkBadResponse(item, sender)
		return
	case VerifyAndWriteTreeSyncChunkResultOkComplete:
		stasy.treeSyncState.treeChunkRequesterGadget.MarkGoodResponse(item, sender)

		if err = kvReadWriteTxn.WriteTreeSyncStatus(TreeSyncStatus{
			TreeSyncPhaseInactive,
			0,
			(StateRootDigest{}),
			PendingKeyDigestRanges{},
		}); err != nil {
			stasy.treeSyncState.logger.Error("failed to write tree-sync status", commontypes.LogFields{
				"sender": sender,
				"err":    err,
			})
			return
		}
		if err = kvReadWriteTxn.WriteLowestPersistedSeqNr(stasy.treeSyncState.targetSeqNr); err != nil {
			stasy.treeSyncState.logger.Error("failed to write lowest persisted seq nr", commontypes.LogFields{
				"sender": sender,
				"err":    err,
			})
			return
		}
		if err = kvReadWriteTxn.WriteHighestCommittedSeqNr(stasy.treeSyncState.targetSeqNr); err != nil {
			stasy.treeSyncState.logger.Error("failed to write highest committed sequence number", commontypes.LogFields{
				"sender": sender,
				"err":    err,
			})
			return
		}
		if err = kvReadWriteTxn.Commit(); err != nil {
			stasy.treeSyncState.logger.Error("failed to commit", commontypes.LogFields{
				"sender": sender,
				"err":    err,
			})
			return
		}
		stasy.treeSyncState.logger.Info("tree synchronization to snapshot completed 🌲", commontypes.LogFields{
			"sender":      sender,
			"targetSeqNr": stasy.treeSyncState.targetSeqNr,
			"rootDigest":  fmt.Sprintf("%x", stasy.treeSyncState.targetStateRootDigest),
		})
		stasy.treeSyncCompleted()
		return
	case VerifyAndWriteTreeSyncChunkResultOkNeedMore:
		stasy.treeSyncState.treeChunkRequesterGadget.MarkGoodResponse(item, sender)

		updatedPendingKeyDigestRanges := stasy.treeSyncState.pendingKeyDigestRanges.WithReceivedRange(receivedKeyDigestRange)

		if err = kvReadWriteTxn.WriteTreeSyncStatus(TreeSyncStatus{
			stasy.treeSyncState.treeSyncPhase,
			stasy.treeSyncState.targetSeqNr,
			stasy.treeSyncState.targetStateRootDigest,
			updatedPendingKeyDigestRanges,
		}); err != nil {
			stasy.treeSyncState.logger.Error("failed to write tree-sync status", commontypes.LogFields{
				"sender": sender,
				"err":    err,
			})
			return
		}
		if err = kvReadWriteTxn.Commit(); err != nil {
			stasy.treeSyncState.logger.Error("failed to commit", commontypes.LogFields{
				"sender": sender,
				"err":    err,
			})
			return
		}

		stasy.treeSyncState.logger.Debug("applied chunk 🍃", commontypes.LogFields{
			"sender":                 sender,
			"startIndex":             fmt.Sprintf("%x", msg.StartIndex),
			"endInclIndex":           fmt.Sprintf("%x", msg.EndInclIndex),
			"keyValuesCount":         len(msg.KeyValues),
			"pendingKeyDigestRanges": fmt.Sprintf("%x", updatedPendingKeyDigestRanges),
		})

		stasy.treeSyncState.pendingKeyDigestRanges = updatedPendingKeyDigestRanges
		stasy.treeSyncState.treeChunkRequesterGadget.PleaseRecheckPendingItems()
		return
	}
	panic("unreachable")
}

func mustTakeSnapshot(seqNr uint64) bool {
	return seqNr%SnapshotInterval == 0
}
