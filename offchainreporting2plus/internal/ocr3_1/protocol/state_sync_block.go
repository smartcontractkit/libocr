package protocol

import (
	"fmt"
	"time"

	"github.com/google/btree"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol/requestergadget"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const (
	maxCumulativeWriteSetBytes    = ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes
	maxMaxCumulativeWriteSetBytes = 2 * maxCumulativeWriteSetBytes
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
	cfgMaxBlocksPerBlockSyncResponse := uint64(stasy.config.GetMaxBlocksPerBlockSyncResponse())
	cfgBlockSyncLookahead := stasy.config.GetMaxParallelRequestedBlocks()
	stasy.blockSyncState.sortedBlockBuffer.Ascend(func(astb AttestedStateTransitionBlock) bool {
		seqNr := astb.StateTransitionBlock.SeqNr()
		if lastSeqNr+1 < seqNr {
			// [lastSeqNr+1..seqNr) (exclusive) is a gap to fill

			for rangeStartSeqNr := lastSeqNr + 1; rangeStartSeqNr < seqNr; rangeStartSeqNr += cfgMaxBlocksPerBlockSyncResponse {
				rangeEndExclSeqNr := min(rangeStartSeqNr+cfgMaxBlocksPerBlockSyncResponse, seqNr)
				pending = append(pending, seqNrRange{rangeStartSeqNr, rangeEndExclSeqNr})
			}
		}
		lastSeqNr = seqNr
		return true
	})

	for rangeStartSeqNr := lastSeqNr + 1; rangeStartSeqNr <= stasy.highestPersistedStateTransitionBlockSeqNr+cfgBlockSyncLookahead && rangeStartSeqNr <= stasy.highestHeardSeqNr; rangeStartSeqNr += cfgMaxBlocksPerBlockSyncResponse {
		maxRangeEndExclSeqNr := stasy.highestPersistedStateTransitionBlockSeqNr + cfgBlockSyncLookahead + 1
		rangeEndExclSeqNr := min(rangeStartSeqNr+cfgMaxBlocksPerBlockSyncResponse, maxRangeEndExclSeqNr)
		// If this is the last range that is pending, we're fine with asking for
		// blocks beyond highestHeardSeqNr, as long as we're not asking for more
		// than cfgMaxBlocksPerBlockSyncResponse blocks.
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

	requestInfo := &types.RequestInfo{
		time.Now().Add(stasy.config.GetDeltaBlockSyncResponseTimeout()),
	}
	msg := MessageBlockSyncRequest[RI]{
		types.EmptyRequestHandleForOutboundRequest,
		requestInfo,
		seqNrRange.StartSeqNr,
		seqNrRange.EndExclSeqNr,
		maxCumulativeWriteSetBytes,
	}
	stasy.netSender.SendTo(msg, target)
	return requestInfo, true
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
		cfgMaxBlocksPerBlockSyncResponse := uint64(stasy.config.GetMaxBlocksPerBlockSyncResponse())
		maxBlocksInResponseU64 := msg.EndExclSeqNr - msg.StartSeqNr
		if maxBlocksInResponseU64 > cfgMaxBlocksPerBlockSyncResponse {
			maxBlocksInResponseU64 = cfgMaxBlocksPerBlockSyncResponse
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

	var astbs []AttestedStateTransitionBlock
	var gstb *GenesisStateTransitionBlock

	// Special case: If someone requests the genesis block, they're looking for
	// the genesis state transition block from the prev instance.
	if msg.StartSeqNr == stasy.genesisSeqNr && msg.EndExclSeqNr == stasy.genesisSeqNr+1 {
		prevInstanceGenesisStateTransitionBlock, err := tx.ReadPrevInstanceGenesisStateTransitionBlock()
		if err != nil {
			stasy.blockSyncState.logger.Error("error reading prev instance genesis state transition block", commontypes.LogFields{
				"error": err,
			})
			return
		}
		gstb = prevInstanceGenesisStateTransitionBlock
	} else {
		var err error
		astbs, _, err = tx.ReadAttestedStateTransitionBlocks(msg.StartSeqNr, maxBlocksInResponse)
		if err != nil {
			stasy.blockSyncState.logger.Error("failed to read attested state transition blocks", commontypes.LogFields{
				"error": err,
			})
			return
		}
	}

	maxCumulativeWriteSetBytes := min(msg.MaxCumulativeWriteSetBytes, maxMaxCumulativeWriteSetBytes)

	// trim suffix of astbs to fit in maxCumulativeWriteSetBytes and avoid seq nr gaps
	cumulativeWriteSetBytes := 0
	validPrefixCount := 0
	for i, astb := range astbs {
		// check if block has the expected sequence number, if not break
		seqNr := astb.StateTransitionBlock.SeqNr()
		expectedSeqNr := msg.StartSeqNr + uint64(i)
		if seqNr != expectedSeqNr {
			// response shouldn't contain gaps
			break
		}

		// check if block's write set size would cause response to be too big
		thisBlockWriteSetBytes := writeSetBytes(&astb)
		if cumulativeWriteSetBytes+thisBlockWriteSetBytes < cumulativeWriteSetBytes {
			// overflow, shouldn't happen but we are careful
			break
		}
		if cumulativeWriteSetBytes+thisBlockWriteSetBytes > maxCumulativeWriteSetBytes {
			// block would cause response to be too big
			break
		}
		cumulativeWriteSetBytes += thisBlockWriteSetBytes
		validPrefixCount++
	}
	astbs = astbs[:validPrefixCount]

	goAway := len(astbs) == 0 && gstb == nil
	stasy.blockSyncState.logger.Debug("sending MessageBlockSyncResponse", commontypes.LogFields{
		"highestPersisted":           stasy.highestPersistedStateTransitionBlockSeqNr,
		"lowestPersisted":            stasy.lowestPersistedStateTransitionBlockSeqNr,
		"requestStartSeqNr":          msg.StartSeqNr,
		"requestEndExclSeqNr":        msg.EndExclSeqNr,
		"blocks":                     len(astbs),
		"hasGenesisBlock":            gstb != nil,
		"cumulativeWriteSetBytes":    cumulativeWriteSetBytes,
		"maxCumulativeWriteSetBytes": maxCumulativeWriteSetBytes,
		"goAway":                     goAway,
		"to":                         sender,
	})
	stasy.netSender.SendTo(MessageBlockSyncResponse[RI]{
		msg.RequestHandle,
		msg.StartSeqNr,
		msg.EndExclSeqNr,
		astbs,
		gstb,
	}, sender)
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

	requestedGenesis := msg.RequestStartSeqNr == stasy.genesisSeqNr && msg.RequestEndExclSeqNr == stasy.genesisSeqNr+1

	// Validate response structure based on what was requested
	if requestedGenesis {
		// Genesis request: must have GenesisStateTransitionBlock, must NOT have AttestedStateTransitionBlocks
		if len(msg.AttestedStateTransitionBlocks) > 0 {
			stasy.blockSyncState.logger.Warn("dropping MessageBlockSyncResponse for genesis with unexpected attested blocks", commontypes.LogFields{
				"sender":            sender,
				"requestSeqNrRange": requestSeqNrRange,
				"numBlocks":         len(msg.AttestedStateTransitionBlocks),
			})
			stasy.blockSyncState.blockRequesterGadget.MarkBadResponse(requestSeqNrRange, sender)
			return
		}
		if msg.GenesisStateTransitionBlock == nil {
			stasy.blockSyncState.logger.Debug("dropping MessageBlockSyncResponse for genesis, go-away", commontypes.LogFields{
				"sender":            sender,
				"requestSeqNrRange": requestSeqNrRange,
			})
			stasy.blockSyncState.blockRequesterGadget.MarkGoAwayResponse(requestSeqNrRange, sender)
			return
		}
	} else {
		// Non-genesis request: must NOT have GenesisStateTransitionBlock
		if msg.GenesisStateTransitionBlock != nil {
			stasy.blockSyncState.logger.Warn("dropping MessageBlockSyncResponse for non-genesis with unexpected genesis block", commontypes.LogFields{
				"sender":            sender,
				"requestSeqNrRange": requestSeqNrRange,
			})
			stasy.blockSyncState.blockRequesterGadget.MarkBadResponse(requestSeqNrRange, sender)
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
	}

	stasy.blockSyncState.logger.Debug("received MessageBlockSyncResponse", commontypes.LogFields{
		"sender":          sender,
		"blocks":          len(msg.AttestedStateTransitionBlocks),
		"hasGenesisBlock": msg.GenesisStateTransitionBlock != nil,
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

	// Verify genesis state transition block if present
	if msg.GenesisStateTransitionBlock != nil {
		if err := stasy.verifyGenesisStateTransitionBlock(*msg.GenesisStateTransitionBlock); err != nil {
			stasy.blockSyncState.logger.Warn("dropping MessageBlockSyncResponse with genesis block that does not verify", commontypes.LogFields{
				"sender":              sender,
				"requestStartSeqNr":   msg.RequestStartSeqNr,
				"requestEndExclSeqNr": msg.RequestEndExclSeqNr,
				"genesisBlockSeqNr":   msg.GenesisStateTransitionBlock.SeqNr,
				"error":               err,
			})
			stasy.blockSyncState.blockRequesterGadget.MarkBadResponse(requestSeqNrRange, sender)
			return
		}
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

		if err := stasy.verifyAttestedStateTransitionBlock(astb); err != nil {
			stasy.blockSyncState.logger.Warn("dropping MessageBlockSyncResponse with block that does not verify", commontypes.LogFields{
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
		var block AttestedOrGenesisStateTransitionBlock
		if msg.GenesisStateTransitionBlock != nil {
			block = msg.GenesisStateTransitionBlock
		} else if len(msg.AttestedStateTransitionBlocks) > 0 {
			block = &msg.AttestedStateTransitionBlocks[0]
		}
		if block == nil {
			stasy.blockSyncState.logger.Critical("unexpected: no block in response after validation passed", nil)
		} else if err := stasy.acceptTreeSyncTargetBlockFromBlockSync(block); err != nil {
			stasy.blockSyncState.logger.Error("error accepting tree-sync target block, will try again", commontypes.LogFields{
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

func (stasy *stateSyncState[RI]) verifyAttestedStateTransitionBlock(astb AttestedStateTransitionBlock) error {
	return astb.Verify(stasy.config.PublicConfig)
}

func (stasy *stateSyncState[RI]) verifyGenesisStateTransitionBlock(gstb GenesisStateTransitionBlock) error {
	if gstb.SeqNr != stasy.genesisSeqNr {
		return fmt.Errorf("genesis state transition block seqNr %d does not match expected genesis seqNr %d", gstb.SeqNr, stasy.genesisSeqNr)
	}
	return VerifyGenesisStateTransitionBlockFromPrevInstance(stasy.config.PublicConfig, gstb)
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

func writeSetBytes(astb *AttestedStateTransitionBlock) int {
	runningCumulativeWriteSetBytes := 0
	for _, kvPair := range astb.StateTransitionBlock.StateWriteSet.Entries {
		runningCumulativeWriteSetBytes += len(kvPair.Key) + len(kvPair.Value)
	}
	return runningCumulativeWriteSetBytes
}
