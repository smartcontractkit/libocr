package protocol

import (
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/protocol/requestergadget"
	"github.com/google/btree"
)

const (
	MaxBlocksPerBlockSyncResponse int = 10

	// Minimum delay between two consecutive BLOCK-SYNC-REQ requests
	DeltaMinBlockSyncRequest = 10 * time.Millisecond
	// Maximum delay between a BLOCK-SYNC-REQ and a BLOCK-SYNC response. We'll try
	// with another oracle if we don't get a response in this time.
	DeltaMaxBlockSyncRequest time.Duration = 1 * time.Second

	// We are looking to pipeline fetches of a range of at most
	// BlockSyncLookahead blocks at any given time
	BlockSyncLookahead = 10_000
)

// Half-open range, i.e. [StartSeqNr, EndExclSeqNr)
type seqNrRange struct {
	StartSeqNr   uint64
	EndExclSeqNr uint64
}

type blockSyncState[RI any] struct {
	logger               commontypes.Logger
	blockRequesterGadget *requestergadget.RequesterGadget[seqNrRange]

	sortedBlockBuffer *btree.BTreeG[AttestedStateTransitionBlock]
}

func (stasy *stateSyncState[RI]) bufferBlock(block AttestedStateTransitionBlock) {
	stasy.blockSyncState.sortedBlockBuffer.ReplaceOrInsert(block)
}

func (stasy *stateSyncState[RI]) reapBlockBuffer() {
	for {
		minBlock, ok := stasy.blockSyncState.sortedBlockBuffer.Min()
		if !ok {
			return
		}

		minBlockSeqNr := minBlock.StateTransitionBlock.SeqNr()

		if minBlockSeqNr <= stasy.highestPersistedStateTransitionBlockSeqNr {
			stasy.blockSyncState.sortedBlockBuffer.DeleteMin()
		} else {
			break
		}
	}
}

func (stasy *stateSyncState[RI]) getPendingBlocksToRequest() []seqNrRange {
	switch stasy.syncMode {
	case syncModeTree:
		return nil
	case syncModeUnknown:
		return nil
	case syncModeFetchSnapshotBlock:
		if stasy.treeSyncState.targetSeqNr == 0 {
			stasy.blockSyncState.logger.Critical("assumption violation: tree-sync target sequence number is 0", nil)
			return nil
		}
		return []seqNrRange{{stasy.treeSyncState.targetSeqNr, stasy.treeSyncState.targetSeqNr + 1}}
	case syncModeBlock:
	}

	if stasy.highestHeardSeqNr <= stasy.highestPersistedStateTransitionBlockSeqNr {
		// We are already synced.
		return nil
	}

	var pending []seqNrRange

	stasy.reapBlockBuffer()

	lastSeqNr := stasy.highestPersistedStateTransitionBlockSeqNr
	stasy.blockSyncState.sortedBlockBuffer.Ascend(func(astb AttestedStateTransitionBlock) bool {
		seqNr := astb.StateTransitionBlock.SeqNr()
		if lastSeqNr+1 < seqNr {
			// [lastSeqNr+1..seqNr) (exclusive) is a gap to fill

			for rangeStartSeqNr := lastSeqNr + 1; rangeStartSeqNr < seqNr; rangeStartSeqNr += uint64(MaxBlocksPerBlockSyncResponse) {
				rangeEndExclSeqNr := rangeStartSeqNr + uint64(MaxBlocksPerBlockSyncResponse)
				if rangeEndExclSeqNr > seqNr {
					rangeEndExclSeqNr = seqNr
				}
				pending = append(pending, seqNrRange{rangeStartSeqNr, rangeEndExclSeqNr})
			}
		}
		lastSeqNr = seqNr
		return true
	})

	for rangeStartSeqNr := lastSeqNr + 1; rangeStartSeqNr <= stasy.highestPersistedStateTransitionBlockSeqNr+BlockSyncLookahead && rangeStartSeqNr <= stasy.highestHeardSeqNr; rangeStartSeqNr += uint64(MaxBlocksPerBlockSyncResponse) {
		rangeEndExclSeqNr := rangeStartSeqNr + uint64(MaxBlocksPerBlockSyncResponse)
		if rangeEndExclSeqNr > stasy.highestPersistedStateTransitionBlockSeqNr+BlockSyncLookahead {
			rangeEndExclSeqNr = stasy.highestPersistedStateTransitionBlockSeqNr + BlockSyncLookahead
		}
		// no check for rangeEndExclSeqNr > stasy.highestHeardSeqNr, because there is no harm in asking for more than exists
		pending = append(pending, seqNrRange{rangeStartSeqNr, rangeEndExclSeqNr})
	}

	return pending
}

func (stasy *stateSyncState[RI]) getBlockSyncSeeders(_ seqNrRange) map[commontypes.OracleID]struct{} {
	seeders := make(map[commontypes.OracleID]struct{})
	for oid := range stasy.oracles {

		if commontypes.OracleID(oid) == stasy.id {
			continue
		}
		seeders[commontypes.OracleID(oid)] = struct{}{}
	}
	return seeders
}

func (stasy *stateSyncState[RI]) sendBlockSyncRequest(seqNrRange seqNrRange, target commontypes.OracleID) (*requestergadget.RequestInfo, bool) {
	stasy.blockSyncState.logger.Debug("sending MessageBlockSyncRequest", commontypes.LogFields{
		"seqNrRange": seqNrRange,
		"target":     target,
	})
	msg := MessageBlockSyncRequest[RI]{
		nil, // TODO: consider using a sentinel value here, e.g. "EmptyRequestHandleForInboundResponse"
		seqNrRange.StartSeqNr,
		seqNrRange.EndExclSeqNr,
	}
	stasy.netSender.SendTo(msg, target)
	return &requestergadget.RequestInfo{
		time.Now().Add(DeltaMaxBlockSyncRequest),
	}, true
}

func (stasy *stateSyncState[RI]) messageBlockSyncRequest(msg MessageBlockSyncRequest[RI], sender commontypes.OracleID) {
	if !(msg.StartSeqNr < msg.EndExclSeqNr) {
		stasy.blockSyncState.logger.Warn("dropping MessageBlockSyncRequest with invalid loSeqNr and hiSeqNr", commontypes.LogFields{
			"sender":              sender,
			"requestStartSeqNr":   msg.StartSeqNr,
			"requestEndExclSeqNr": msg.EndExclSeqNr,
		})
		return
	}

	stasy.blockSyncState.logger.Debug("received MessageBlockSyncRequest", commontypes.LogFields{
		"sender":              sender,
		"requestStartSeqNr":   msg.StartSeqNr,
		"requestEndExclSeqNr": msg.EndExclSeqNr,
	})

	var maxBlocksInResponse int
	{
		maxBlocksInResponseU64 := msg.EndExclSeqNr - msg.StartSeqNr
		if maxBlocksInResponseU64 > uint64(MaxBlocksPerBlockSyncResponse) {
			maxBlocksInResponseU64 = uint64(MaxBlocksPerBlockSyncResponse)
		}
		// now we are sure that maxBlocksInResponseU64 will fit an int
		maxBlocksInResponse = int(maxBlocksInResponseU64)
	}

	tx, err := stasy.kvDb.NewReadTransactionUnchecked()
	if err != nil {
		stasy.blockSyncState.logger.Error("failed to create read transaction", commontypes.LogFields{
			"error": err,
		})
		return
	}
	defer tx.Discard()

	astbs, _, err := tx.ReadAttestedStateTransitionBlocks(msg.StartSeqNr, maxBlocksInResponse)
	if err != nil {
		stasy.blockSyncState.logger.Error("failed to read attested state transition blocks", commontypes.LogFields{
			"error": err,
		})
		return
	}

	for i, astb := range astbs {
		seqNr := astb.StateTransitionBlock.SeqNr()
		var expectedSeqNr uint64
		if i == 0 {
			expectedSeqNr = msg.StartSeqNr
		} else {
			expectedSeqNr = astbs[i-1].StateTransitionBlock.SeqNr() + 1
		}
		if seqNr != expectedSeqNr {
			astbs = nil
			break // do not produce gap
		}
	}

	if len(astbs) > 0 {
		stasy.blockSyncState.logger.Debug("sending MessageBlockSyncResponse", commontypes.LogFields{
			"highestPersisted":     stasy.highestPersistedStateTransitionBlockSeqNr,
			"lowestPersisted":      stasy.lowestPersistedStateTransitionBlockSeqNr,
			"requestStartSeqNr":    msg.StartSeqNr,
			"requestEndExclSeqNr":  msg.EndExclSeqNr,
			"responseStartSeqNr":   astbs[0].StateTransitionBlock.SeqNr(),
			"responseEndExclSeqNr": astbs[len(astbs)-1].StateTransitionBlock.SeqNr() + 1,
			"to":                   sender,
		})
		stasy.netSender.SendTo(MessageBlockSyncResponse[RI]{
			msg.RequestHandle,
			msg.StartSeqNr,
			msg.EndExclSeqNr,
			astbs,
		}, sender)
	} else {
		stasy.blockSyncState.logger.Debug("no blocks to send, sending an empty MessageBlockSyncResponse to indicate go-away", commontypes.LogFields{
			"highestPersisted":    stasy.highestPersistedStateTransitionBlockSeqNr,
			"lowestPersisted":     stasy.lowestPersistedStateTransitionBlockSeqNr,
			"requestStartSeqNr":   msg.StartSeqNr,
			"requestEndExclSeqNr": msg.EndExclSeqNr,
			"to":                  sender,
		})
		stasy.netSender.SendTo(MessageBlockSyncResponse[RI]{
			msg.RequestHandle,
			msg.StartSeqNr,
			msg.EndExclSeqNr,
			astbs,
		}, sender)
	}
}

func (stasy *stateSyncState[RI]) messageBlockSyncResponse(msg MessageBlockSyncResponse[RI], sender commontypes.OracleID) {
	requestSeqNrRange := seqNrRange{msg.RequestStartSeqNr, msg.RequestEndExclSeqNr}

	if !stasy.blockSyncState.blockRequesterGadget.CheckAndMarkResponse(requestSeqNrRange, sender) {
		stasy.blockSyncState.logger.Warn("dropping MessageBlockSyncResponse, not allowed", commontypes.LogFields{
			"sender":            sender,
			"requestSeqNrRange": requestSeqNrRange,
		})
		return
	}

	if len(msg.AttestedStateTransitionBlocks) == 0 {
		stasy.blockSyncState.logger.Debug("dropping MessageBlockSyncResponse, go-away", commontypes.LogFields{
			"sender":            sender,
			"requestSeqNrRange": requestSeqNrRange,
		})
		stasy.blockSyncState.blockRequesterGadget.MarkGoAwayResponse(requestSeqNrRange, sender)
		return
	}

	stasy.blockSyncState.logger.Debug("received MessageBlockSyncResponse", commontypes.LogFields{
		"sender": sender,
	})

	switch stasy.syncMode {
	case syncModeFetchSnapshotBlock:
		break
	case syncModeBlock:
		break
	case syncModeTree:
		return
	case syncModeUnknown:
		return
	}

	for i, astb := range msg.AttestedStateTransitionBlocks {
		if astb.StateTransitionBlock.SeqNr() != msg.RequestStartSeqNr+uint64(i) {
			stasy.blockSyncState.logger.Warn("dropping MessageBlockSyncResponse with out of order state transition blocks", commontypes.LogFields{
				"sender":                    sender,
				"requestStartSeqNr":         msg.RequestStartSeqNr,
				"requestEndExclSeqNr":       msg.RequestEndExclSeqNr,
				"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
			})
			stasy.blockSyncState.blockRequesterGadget.MarkBadResponse(requestSeqNrRange, sender)
			return
		}

		if !(astb.StateTransitionBlock.SeqNr() < msg.RequestEndExclSeqNr) {
			stasy.blockSyncState.logger.Warn("dropping MessageBlockSyncResponse with state transition block seqNr that is too large", commontypes.LogFields{
				"sender":                    sender,
				"requestStartSeqNr":         msg.RequestStartSeqNr,
				"requestEndExclSeqNr":       msg.RequestEndExclSeqNr,
				"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
			})
			stasy.blockSyncState.blockRequesterGadget.MarkBadResponse(requestSeqNrRange, sender)
			return
		}

		if err := astb.Verify(stasy.config.ConfigDigest, stasy.config.OracleIdentities, stasy.config.ByzQuorumSize()); err != nil {
			stasy.blockSyncState.logger.Warn("dropping MessageBlockSyncResponse with invalid attestation", commontypes.LogFields{
				"sender":                    sender,
				"requestStartSeqNr":         msg.RequestStartSeqNr,
				"requestEndExclSeqNr":       msg.RequestEndExclSeqNr,
				"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
				"error":                     err,
			})
			stasy.blockSyncState.blockRequesterGadget.MarkBadResponse(requestSeqNrRange, sender)
			return
		}
	}

	stasy.blockSyncState.blockRequesterGadget.MarkGoodResponse(requestSeqNrRange, sender)

	if stasy.syncMode == syncModeFetchSnapshotBlock {
		err := stasy.acceptTreeSyncTargetBlockFromBlockSync(msg.AttestedStateTransitionBlocks[0])
		if err != nil {
			stasy.blockSyncState.logger.Error("error accepting tree-sync target block from block sync, will try again", commontypes.LogFields{
				"error": err,
			})
		}
		return
	}

	for _, astb := range msg.AttestedStateTransitionBlocks {
		stasy.blockSyncState.logger.Debug("buffering state transition block", commontypes.LogFields{
			"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
		})
		stasy.bufferBlock(astb)
	}

	stasy.tryCompleteBlockSync()
}

func (stasy *stateSyncState[RI]) tryCompleteBlockSync() {

	stasy.refreshStateSyncState()

	stasy.reapBlockBuffer()

	minBlock, ok := stasy.blockSyncState.sortedBlockBuffer.Min()
	if !ok {
		return
	}
	if minBlock.StateTransitionBlock.SeqNr() != stasy.highestPersistedStateTransitionBlockSeqNr+1 {
		return
	}

	tx, err := stasy.kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		stasy.blockSyncState.logger.Error("failed to create read transaction", commontypes.LogFields{
			"error": err,
		})
		return
	}
	defer tx.Discard()

	lastSeqNr := stasy.highestPersistedStateTransitionBlockSeqNr
	for {
		astb, ok := stasy.blockSyncState.sortedBlockBuffer.Min()
		if !ok {
			break
		}
		seqNr := astb.StateTransitionBlock.SeqNr()

		if seqNr != lastSeqNr+1 {
			break
		}

		stasy.blockSyncState.logger.Debug("writing state transition block", commontypes.LogFields{
			"stateTransitionBlockSeqNr": seqNr,
		})

		err := tx.WriteAttestedStateTransitionBlock(seqNr, astb)
		if err != nil {
			stasy.blockSyncState.logger.Error("error writing state transition block", commontypes.LogFields{
				"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
				"error":                     err,
			})
			return
		}

		lastSeqNr = seqNr
		stasy.blockSyncState.sortedBlockBuffer.DeleteMin()
	}

	err = tx.Commit()
	if err != nil {
		stasy.blockSyncState.logger.Error("error committing transaction", commontypes.LogFields{
			"error": err,
		})
		return
	}
	stasy.highestPersistedStateTransitionBlockSeqNr = lastSeqNr
	stasy.pleaseTryToReplayBlock()
}
