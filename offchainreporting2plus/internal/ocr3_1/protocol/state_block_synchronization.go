package protocol

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/scheduler"
)

const (
	MaxBlocksSent int = 10

	// Maximum delay between a BLOCK-SYNC-REQ and a BLOCK-SYNC response. We'll try
	// with another oracle if we don't get a response in this time.

	DeltaMaxBlockSyncRequest time.Duration = 1 * time.Second

	// Minimum delay between two consecutive BLOCK-SYNC-REQ requests
	DeltaMinBlockSyncRequest = 10 * time.Millisecond

	// An oracle sends a BLOCK-SYNC-SUMMARY message every DeltaBlockSyncHeartbeat
	DeltaBlockSyncHeartbeat time.Duration = time.Duration(math.MaxInt64)

	MaxBlockSyncSize int = 50000000
)

type blockSyncState[RI any] struct {
	logger    commontypes.Logger
	oracles   []*blockSyncTargetOracle[RI]
	scheduler *scheduler.Scheduler[EventToStatePersistence[RI]]
}

func (state *statePersistenceState[RI]) numInflightRequests() int {
	count := 0
	for _, oracle := range state.blockSyncState.oracles {
		if oracle.inFlightRequest != nil {
			count++
		}
	}
	return count
}

type blockSyncTargetOracle[RI any] struct {
	// lowestPersistedSeqNr is the lowest sequence number the oracle still has an attested
	// state transition block for

	lowestPersistedSeqNr  uint64
	lastSummaryReceivedAt time.Time
	// whether it is viable to send the next block sync request to this oracle
	candidate bool
	// the current inflight request to this oracle, nil otherwise
	inFlightRequest *inFlightRequest[RI]
}

type inFlightRequest[RI any] struct {
	message       MessageBlockSyncRequest[RI]
	requestedFrom commontypes.OracleID
}

func (state *statePersistenceState[RI]) highestHeardIncreased() {
	state.trySendNextRequest()
}

func (state *statePersistenceState[RI]) clearStaleBlockSyncRequests() {
	state.refreshHighestPersistedStateTransitionBlockSeqNr()
	nowPersistedSeqNr := state.highestPersistedStateTransitionBlockSeqNr
	for _, oracle := range state.blockSyncState.oracles {
		req := oracle.inFlightRequest
		if req == nil {
			continue
		}
		thenPersistedSeqNr := req.message.HighestCommittedSeqNr
		if thenPersistedSeqNr < nowPersistedSeqNr {
			// this is a stale request
			state.blockSyncState.logger.Debug("removing stale BlockSyncRequest", commontypes.LogFields{
				"requestedFrom":      req.requestedFrom,
				"thenPersistedSeqNr": thenPersistedSeqNr,
				"nowPersistedSeqNr":  nowPersistedSeqNr,
			})
			oracle.inFlightRequest = nil
		}
	}
}

func (state *statePersistenceState[RI]) trySendNextRequest() {
	if !state.readyToSendBlockSyncReq {
		state.blockSyncState.logger.Trace("trySendNextRequest: not marked as ready to send BlockSyncRequest, dropping", nil)
		return
	}
	if state.numInflightRequests() != 0 {
		// if numInflightRequests > 0, we are already waiting for a response which
		// we'll either receive or timeout, but regardless it will carry us over
		// until state.highestHeard is retrieved
		state.blockSyncState.logger.Debug("we are already fetching blocks", commontypes.LogFields{
			"numInflightRequests": state.numInflightRequests(),
		})
		return
	}

	state.refreshHighestPersistedStateTransitionBlockSeqNr()
	reqSeqNr := state.highestPersistedStateTransitionBlockSeqNr
	if state.highestHeardSeqNr > reqSeqNr {
		state.blockSyncState.logger.Trace("trySendNextRequest: highestHeardSeqNr > highestPersistedStateTransitionBlockSeqNr, sending BlockSyncRequest", commontypes.LogFields{
			"highestHeardSeqNr":     state.highestHeardSeqNr,
			"highestPersistedSeqNr": state.highestPersistedStateTransitionBlockSeqNr,
		})
		state.sendBlockSyncReq(reqSeqNr)
	}
}

func (state *statePersistenceState[RI]) tryComplete() {
	state.clearStaleBlockSyncRequests()
	state.trySendNextRequest()
}

func (state *statePersistenceState[RI]) processBlockSyncSummaryHeartbeat() {
	defer state.blockSyncState.scheduler.ScheduleDelay(EventBlockSyncSummaryHeartbeat[RI]{}, DeltaBlockSyncHeartbeat)
	lowestPersistedSeqNr := 0
	state.refreshHighestPersistedStateTransitionBlockSeqNr()
	if state.highestPersistedStateTransitionBlockSeqNr >= 1 {
		lowestPersistedSeqNr = 1
	}
	state.netSender.Broadcast(MessageBlockSyncSummary[RI]{
		uint64(lowestPersistedSeqNr),
	})
}

func (state *statePersistenceState[RI]) messageBlockSyncSummary(msg MessageBlockSyncSummary[RI], sender commontypes.OracleID) {
	state.blockSyncState.logger.Debug("received messageBlockSyncSummary", commontypes.LogFields{
		"sender":                  sender,
		"msgLowestPersistedSeqNr": msg.LowestPersistedSeqNr,
	})
	oracle := state.blockSyncState.oracles[sender]
	oracle.lowestPersistedSeqNr = msg.LowestPersistedSeqNr
	oracle.lastSummaryReceivedAt = time.Now()
}

func (state *statePersistenceState[RI]) processExpiredBlockSyncRequest(requestedFrom commontypes.OracleID, nonce uint64) {
	oracle := state.blockSyncState.oracles[requestedFrom]
	if oracle.inFlightRequest == nil {
		return
	}
	if oracle.inFlightRequest.message.Nonce == nonce {
		oracle.inFlightRequest = nil
		oracle.candidate = false
	}
	state.tryComplete()
}

func (state *statePersistenceState[RI]) sendBlockSyncReq(seqNr uint64) {
	candidates := make([]commontypes.OracleID, 0, state.config.N())
	for oracleID, oracle := range state.blockSyncState.oracles {
		if commontypes.OracleID(oracleID) == state.id {
			continue
		}
		if oracle.candidate {

			candidates = append(candidates, commontypes.OracleID(oracleID))
		}
	}

	if len(candidates) == 0 {

		state.blockSyncState.logger.Debug("not candidate oracles for MessageBlockSyncRequest, restarting from scratch", nil)
		candidates = make([]commontypes.OracleID, 0, state.config.N())
		for oracleID, oracle := range state.blockSyncState.oracles {
			oracle.candidate = true
			candidates = append(candidates, commontypes.OracleID(oracleID))
		}
	}
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates))))
	if err != nil {
		state.blockSyncState.logger.Critical("unexpected error returned by rand.Int", commontypes.LogFields{
			"error": err,
		})
		return
	}
	target := candidates[int(randomIndex.Int64())]
	nonce, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 64))
	if err != nil {
		state.blockSyncState.logger.Critical("unexpected error returned by rand.Int", commontypes.LogFields{
			"error": err,
		})
		return
	}
	state.blockSyncState.logger.Debug("sending MessageBlockSyncRequest", commontypes.LogFields{
		"highestCommittedSeqNr": seqNr,
		"target":                target,
	})
	msg := MessageBlockSyncRequest[RI]{
		nil, // TODO: consider using a sentinel value here, e.g. "EmptyRequestHandleForInboundResponse"
		seqNr,
		nonce.Uint64(),
	}
	state.netSender.SendTo(msg, target)
	if !(0 <= int(target) && int(target) < len(state.blockSyncState.oracles)) {
		state.blockSyncState.logger.Critical("target oracle out of bounds", commontypes.LogFields{
			"target": target,
			"N":      state.config.N(),
		})
		return
	}
	state.blockSyncState.oracles[target].inFlightRequest = &inFlightRequest[RI]{msg, target}
	state.blockSyncState.scheduler.ScheduleDelay(EventExpiredBlockSyncRequest[RI]{target, nonce.Uint64()}, DeltaMaxBlockSyncRequest)
	state.readyToSendBlockSyncReq = false
	state.blockSyncState.scheduler.ScheduleDelay(EventReadyToSendNextBlockSyncRequest[RI]{}, DeltaMinBlockSyncRequest)
}

func (state *statePersistenceState[RI]) messageBlockSyncReq(msg MessageBlockSyncRequest[RI], sender commontypes.OracleID) {
	state.blockSyncState.logger.Debug("received MessageBlockSyncRequest", commontypes.LogFields{
		"sender":                   sender,
		"msgHighestCommittedSeqNr": msg.HighestCommittedSeqNr,
	})
	loSeqNr := msg.HighestCommittedSeqNr + 1

	state.refreshHighestPersistedStateTransitionBlockSeqNr()

	var (
		astbs   []AttestedStateTransitionBlock
		hiSeqNr uint64
	)
	for seqNr := loSeqNr; len(astbs) < MaxBlocksSent && seqNr <= state.highestPersistedStateTransitionBlockSeqNr; seqNr++ {

		astb, err := state.database.ReadAttestedStateTransitionBlock(state.ctx, state.config.ConfigDigest, seqNr)
		if err != nil {
			state.blockSyncState.logger.Error("Database.ReadAttestedStateTransitionBlock failed while producing MessageBlockSync", commontypes.LogFields{
				"seqNr": seqNr,
				"error": err,
			})
			break // Stopping to not produce a gap.
		}

		if astb.StateTransitionBlock.SeqNr() != seqNr {
			break // Stopping to not produce a gap.
		}
		astbs = append(astbs, astb)
		hiSeqNr = seqNr
	}

	if len(astbs) > 0 {
		state.blockSyncState.logger.Debug("sending MessageBlockSync", commontypes.LogFields{
			"highestPersisted": state.highestPersistedStateTransitionBlockSeqNr,
			"loSeqNr":          loSeqNr,
			"hiSeqNr":          hiSeqNr,
			"to":               sender,
		})
		state.netSender.SendTo(MessageBlockSync[RI]{
			msg.RequestHandle,
			astbs,
			msg.Nonce,
		}, sender)
	} else {
		state.blockSyncState.logger.Debug("no blocks to send, not responding to MessageBlockSyncRequest", commontypes.LogFields{
			"highestPersisted": state.highestPersistedStateTransitionBlockSeqNr,
			"loSeqNr":          loSeqNr,
			"to":               sender,
		})
	}
}

func (state *statePersistenceState[RI]) messageBlockSync(msg MessageBlockSync[RI], sender commontypes.OracleID) {
	state.blockSyncState.logger.Debug("received MessageBlockSync", commontypes.LogFields{
		"sender": sender,
	})
	req := state.blockSyncState.oracles[sender].inFlightRequest
	if req == nil {
		state.blockSyncState.logger.Warn("dropping unexpected MessageBlockSync", commontypes.LogFields{
			"nonce":  msg.Nonce,
			"sender": sender,
		})
		return
	}

	if msg.Nonce != req.message.Nonce {
		state.blockSyncState.logger.Warn("dropping MessageBlockSync with unexpected nonce", commontypes.LogFields{
			"expectedNonce": req.message.Nonce,
			"actualNonce":   msg.Nonce,
			"sender":        sender,
		})
		return
	}

	// so that any future response with the same nonce will become invalid
	state.blockSyncState.oracles[sender].inFlightRequest = nil

	// at this point we know we've received a response from the correct oracle

	// 1. if any of the following logic errors out, we will immediately notice
	// and start re-requesting from where we left off, even if we partially
	// persist the blocks in this response
	// 2. if the logic succeeds, we'll move to requesting for the next sequence
	// number, until we reach highestHeardSeqNr
	defer state.tryComplete()
	if len(msg.AttestedStateTransitionBlocks) > MaxBlocksSent {
		state.blockSyncState.logger.Warn("dropping MessageBlockSync with more blocks than the maximum allowed number", commontypes.LogFields{
			"blockNum":         len(msg.AttestedStateTransitionBlocks),
			"expectedBlockNum": MaxBlocksSent,
			"sender":           sender,
		})
		return
	}
	for i, astb := range msg.AttestedStateTransitionBlocks {
		if astb.StateTransitionBlock.SeqNr() != req.message.HighestCommittedSeqNr+uint64(i)+1 {
			state.blockSyncState.logger.Warn("dropping MessageBlockSync with out of order state transition blocks", commontypes.LogFields{
				"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
				"sender":                    sender,
			})
			return
		}
	}
	for _, astb := range msg.AttestedStateTransitionBlocks {
		if err := astb.Verify(state.config.ConfigDigest, state.config.OracleIdentities, state.config.ByzQuorumSize()); err != nil {
			state.blockSyncState.logger.Warn("dropping MessageBlockSync with invalid attestation", commontypes.LogFields{
				"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
				"sender":                    sender,
				"error":                     err,
			})
			return
		}
	}

	for _, astb := range msg.AttestedStateTransitionBlocks {
		state.refreshHighestPersistedStateTransitionBlockSeqNr()
		expectedSeqNr := state.highestPersistedStateTransitionBlockSeqNr + 1
		seqNr := astb.StateTransitionBlock.SeqNr()

		state.blockSyncState.logger.Debug("retrieved state transition block", commontypes.LogFields{
			"stateTransitionBlockSeqNr": seqNr,
		})

		if seqNr > expectedSeqNr {

			state.blockSyncState.logger.Warn("dropping MessageBlockSync which creates gaps in persisted blocks", commontypes.LogFields{
				"stateTransitionBlockSeqNr":                 astb.StateTransitionBlock.SeqNr(),
				"highestPersistedStateTransitionBlockSeqNr": state.highestPersistedStateTransitionBlockSeqNr,
				"sender": sender,
			})
			return
		} else if seqNr < expectedSeqNr {
			state.blockSyncState.logger.Debug("no need to persist this block, we have done so already", commontypes.LogFields{
				"stateTransitionBlockSeqNr":            astb.StateTransitionBlock.SeqNr(),
				"highestPersistedStateTransitionBlock": state.highestPersistedStateTransitionBlockSeqNr,
			})
		} else {
			werr := state.persist(astb)
			if werr != nil {

				{
					rastb, rerr := state.database.ReadAttestedStateTransitionBlock(state.ctx, state.config.ConfigDigest, seqNr)
					if rerr == nil && rastb.StateTransitionBlock.SeqNr() == seqNr {

						continue
					}
				}

				state.blockSyncState.logger.Error("error persisting state transition block", commontypes.LogFields{
					"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
					"error":                     werr,
				})
				return
			}
		}
	}
}
