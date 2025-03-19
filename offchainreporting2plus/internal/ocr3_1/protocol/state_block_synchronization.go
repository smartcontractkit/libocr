package protocol

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/scheduler"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const (
	MaxBlocksSent int = 10

	// Maximum delay between a BLOCK-SYNC-REQ and a BLOCK-SYNC response. We'll try
	// with another oracle if we don't get a response in this time.
	DeltaMaxBlockSyncRequest time.Duration = 1 * time.Second

	// Minimum delay between two consecutive BLOCK-SYNC-REQ requests
	DeltaMinBlockSyncRequest = 10 * time.Millisecond

	// An oracle sends a BLOCK-SYNC-SUMMARY message every DeltaBlockSyncHeartbeat
	DeltaBlockSyncHeartbeat time.Duration = 5 * time.Second

	maxBlockSyncSize int = 50000000
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
	for _, oracle := range state.blockSyncState.oracles {
		req := oracle.inFlightRequest
		if req == nil {
			continue
		}
		if req.message.HighestCommittedSeqNr < state.highestCommittedToKVSeqNr {
			// this is a stale request
			state.blockSyncState.logger.Debug("removing stale BlockSyncRequest", commontypes.LogFields{
				"requestedFrom":             req.requestedFrom,
				"requestedFromSeqNr":        req.message.HighestCommittedSeqNr,
				"highestCommittedToKVSeqNr": state.highestCommittedToKVSeqNr,
			})
			oracle.inFlightRequest = nil
		}
	}
}

func (state *statePersistenceState[RI]) trySendNextRequest() {
	if !state.readyToSendBlockSyncReq {
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

	if state.highestHeardSeqNr > state.highestCommittedToKVSeqNr {
		state.sendBlockSyncReq(state.highestCommittedToKVSeqNr)
	}
}

func (state *statePersistenceState[RI]) tryComplete() {
	state.clearStaleBlockSyncRequests()
	state.trySendNextRequest()
}

func (state *statePersistenceState[RI]) processBlockSyncSummaryHeartbeat() {
	defer state.blockSyncState.scheduler.ScheduleDelay(EventBlockSyncSummaryHeartbeat[RI]{}, DeltaBlockSyncHeartbeat)
	lowestPersistedSeqNr := 0
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
	if len(candidates) <= state.config.F {

		state.blockSyncState.logger.Debug("not enough candidate oracles for MessageBlockSyncRequest, restarting from scratch", nil)
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
	msg := MessageBlockSyncRequest[RI]{seqNr, nonce.Uint64()}
	state.netSender.SendTo(MessageBlockSyncRequestWrapper[RI]{
		msg,
		nil,
		types.SingleUseSizedLimitedResponsePolicy{
			maxBlockSyncSize,
			time.Now().Add(DeltaMaxBlockSyncRequest),
		},
	}, target)
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

func (state *statePersistenceState[RI]) messageBlockSyncReq(msgWrapper RequestMessage[RI], sender commontypes.OracleID) {
	if _, ok := msgWrapper.GetSerializableRequestMessage().(MessageBlockSyncRequest[RI]); !ok {
		state.blockSyncState.logger.Error("message type assertion failed", commontypes.LogFields{
			"sender": sender,
		})
	}
	msg := msgWrapper.GetSerializableRequestMessage().(MessageBlockSyncRequest[RI])
	state.blockSyncState.logger.Debug("received MessageBlockSyncRequest", commontypes.LogFields{
		"sender":                   sender,
		"msgHighestCommittedSeqNr": msg.HighestCommittedSeqNr,
	})
	loSeqNr := msg.HighestCommittedSeqNr + 1

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
		state.netSender.SendTo(MessageBlockSyncWrapper[RI]{
			MessageBlockSync[RI]{
				astbs,
				msg.Nonce,
			},
			msgWrapper.GetRequestHandle(),
		}, sender)
	} else {
		state.blockSyncState.logger.Debug("no blocks to send, not responding to MessageBlockSyncRequest", commontypes.LogFields{
			"highestPersisted": state.highestPersistedStateTransitionBlockSeqNr,
			"loSeqNr":          loSeqNr,
			"to":               sender,
		})
	}
}

func (state *statePersistenceState[RI]) messageBlockSync(msgWrapper ResponseMessage[RI], sender commontypes.OracleID) {
	if _, ok := msgWrapper.GetSerializableResponseMessage().(MessageBlockSync[RI]); !ok {
		state.blockSyncState.logger.Error("message type assertion failed", commontypes.LogFields{
			"sender": sender,
		})
	}
	msg := msgWrapper.GetSerializableResponseMessage().(MessageBlockSync[RI])
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
		state.blockSyncState.logger.Debug("retrieved state transition block", commontypes.LogFields{
			"seqNr": astb.StateTransitionBlock.SeqNr(),
		})
		if astb.StateTransitionBlock.SeqNr() <= state.highestCommittedToKVSeqNr {
			state.blockSyncState.logger.Debug("no need to replay this state transition", commontypes.LogFields{
				"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
				"highestPersistedToKvSeqNr": state.highestCommittedToKVSeqNr,
			})
			continue
		}
		if astb.StateTransitionBlock.SeqNr() > state.highestPersistedStateTransitionBlockSeqNr+1 {

			state.blockSyncState.logger.Warn("dropping MessageBlockSync which creates gaps in persisted blocks", commontypes.LogFields{
				"stateTransitionBlockSeqNr":                 astb.StateTransitionBlock.SeqNr(),
				"highestPersistedStateTransitionBlockSeqNr": state.highestPersistedStateTransitionBlockSeqNr,
				"sender": sender,
			})
			return
		} else if astb.StateTransitionBlock.SeqNr() <= state.highestPersistedStateTransitionBlockSeqNr {
			state.blockSyncState.logger.Debug("no need to persist this block, we have done so already", commontypes.LogFields{
				"stateTransitionBlockSeqNr":            astb.StateTransitionBlock.SeqNr(),
				"highestPersistedStateTransitionBlock": state.highestPersistedStateTransitionBlockSeqNr,
			})
		} else if astb.StateTransitionBlock.SeqNr() == state.highestPersistedStateTransitionBlockSeqNr+1 {
			err := state.persist(astb)
			if err != nil {
				state.blockSyncState.logger.Error("error persisting state transition block", commontypes.LogFields{
					"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
					"error":                     err,
				})
				return
			}
		}
		// we might have already persisted the block but we might have not committed it to KV yet
		if last, isNotEmpty := state.retrievedStateTransitionBlockSeqNrQueue.PeekLast(); !isNotEmpty || last < state.highestPersistedStateTransitionBlockSeqNr {
			state.blockSyncState.logger.Debug("pushing to retrievedStateTransitionBlockSeqNrQueue", commontypes.LogFields{
				"stateTransitionBlockSeqNr": astb.StateTransitionBlock.SeqNr(),
				"lastSeqNr":                 last,
			})
			state.retrievedStateTransitionBlockSeqNrQueue.Push(astb.StateTransitionBlock.SeqNr())
		}
	}
}
