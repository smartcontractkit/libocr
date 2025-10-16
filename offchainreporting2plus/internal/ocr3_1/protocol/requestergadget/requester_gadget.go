package requestergadget

import (
	"cmp"
	"maps"
	"slices"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
)

type RequestInfo struct {
	// If no response is received by this time, the request is considered timed out.
	ExpiryTimestamp time.Time
}

func NewRequesterGadget[Item comparable](
	n int,
	requestInterval time.Duration, // Wait interval between requests to the same seeder
	sendRequestFn func(Item, commontypes.OracleID) (*RequestInfo, bool), // Invoked by the RequesterGadget to send a request for the given item to the given seeder.
	getPendingItemsFn func() []Item, // Invoked by the RequesterGadget to get the list of items that should be requested. RequesterGadget will attempt to request items earlier in the list first.
	getSeedersFn func(Item) map[commontypes.OracleID]struct{}, // Invoked by the RequesterGadget to get the list of seeders that can serve the given item.
) *RequesterGadget[Item] {
	oracles := make(map[commontypes.OracleID]*oracleState, n)
	for i := range n {
		oracles[commontypes.OracleID(i)] = &oracleState{
			time.Time{},
			0,
		}
	}
	return &RequesterGadget[Item]{
		oracles,
		requestInterval,
		make(map[Item]*pendingItemState),
		time.After(0),
		sendRequestFn,
		getPendingItemsFn,
		getSeedersFn,
	}
}

// PleaseRecheckPendingItems must be called by the protocol when the output of
// getPendingItemsFn or getSeedersFn has changed.
func (rg *RequesterGadget[Item]) PleaseRecheckPendingItems() {
	rg.chTick = time.After(0)
}

// CheckAndMarkResponse must be called by the protocol when a response is
// received, to ensure that the response matches a request that the gadget has
// sent. It will return true even if the request has technically timed out in
// some cases.
func (rg *RequesterGadget[Item]) CheckAndMarkResponse(item Item, sender commontypes.OracleID) bool {
	rg.PleaseRecheckPendingItems() // overly sensitive, but easier to reason about
	if pendingItem, ok := rg.ourPendingItems[item]; ok {
		if pendingItem.pendingRequestOrNil == nil {
			return false
		}
		pendingRequest := pendingItem.pendingRequestOrNil
		if pendingRequest.seeder == sender {
			pendingItem.pendingRequestOrNil = nil
			return true
		}
	}
	return false
}

// We temporarily exclude a responder for an item when they time out, send a go
// away, or send a bad response. We clear the exclusion list when we've excluded
// them all but still haven't received the item.
func (rg *RequesterGadget[Item]) temporaryExcludeResponderForItem(item Item, sender commontypes.OracleID) {
	if pendingItem, ok := rg.ourPendingItems[item]; ok {
		pendingItem.temporarilyExcludedSeeders[sender] = struct{}{}
	}
}

func (rg *RequesterGadget[Item]) MarkGoAwayResponse(item Item, sender commontypes.OracleID) {
	rg.temporaryExcludeResponderForItem(item, sender)
}

func (rg *RequesterGadget[Item]) MarkGoodResponder(sender commontypes.OracleID) {
	rg.oracles[sender].score++
}

func (rg *RequesterGadget[Item]) MarkGoodResponse(_ Item, sender commontypes.OracleID) {
	rg.MarkGoodResponder(sender)
}

func (rg *RequesterGadget[Item]) MarkBadResponder(sender commontypes.OracleID) {
	rg.oracles[sender].score /= 2
}

func (rg *RequesterGadget[Item]) MarkBadResponse(item Item, sender commontypes.OracleID) {
	rg.temporaryExcludeResponderForItem(item, sender)
	rg.MarkBadResponder(sender)
}

// Only called by the requester gadget itself. The protocol using this gadget
// has no way of knowing a request was sent or timed out.
func (rg *RequesterGadget[Item]) markTimedOutResponse(item Item, sender commontypes.OracleID) {
	rg.MarkBadResponse(item, sender)
}

func (rg *RequesterGadget[Item]) rankedSeeders(seeders map[commontypes.OracleID]struct{}, excluded map[commontypes.OracleID]struct{}) []commontypes.OracleID {
	type scoredSeeder struct {
		seeder commontypes.OracleID
		score  uint64
	}
	scoredSeeders := make([]scoredSeeder, 0, len(seeders))
	for seeder := range seeders {
		if _, ok := excluded[seeder]; ok {
			continue
		}
		scoredSeeders = append(scoredSeeders, scoredSeeder{
			seeder,
			rg.oracles[seeder].score,
		})
	}
	slices.SortFunc(scoredSeeders, func(a, b scoredSeeder) int {
		// higher score goes first
		return cmp.Compare(b.score, a.score)
	})

	ranks := make([]commontypes.OracleID, 0, len(scoredSeeders))
	for _, scoredSeeder := range scoredSeeders {
		ranks = append(ranks, scoredSeeder.seeder)
	}
	return shuffle(ranks)
}

func (rg *RequesterGadget[Item]) Ticker() <-chan time.Time {
	return rg.chTick
}

const maxNextTickInterval = 15 * time.Second

func (rg *RequesterGadget[Item]) Tick() {

	now := time.Now()

	pendingItems := rg.getPendingItemsFn()
	// Discard any pending requests for no longer needed items.
	maps.DeleteFunc(rg.ourPendingItems, func(item Item, _ *pendingItemState) bool {
		return !slices.Contains(pendingItems, item)
	})

	nextTick := now.Add(maxNextTickInterval)

	for _, item := range pendingItems {
		// Add state for this item if we didn't have it before.
		if _, ok := rg.ourPendingItems[item]; !ok {
			rg.ourPendingItems[item] = &pendingItemState{
				nil,
				make(map[commontypes.OracleID]struct{}),
			}
		}
		pendingItemState := rg.ourPendingItems[item]
		pendingRequestOrNil := pendingItemState.pendingRequestOrNil

		var shouldRequestNow bool
		if pendingRequestOrNil != nil {
			pendingRequest := pendingRequestOrNil
			if pendingRequest.expiryTimestamp.Before(now) {
				// Previous request timed out.
				rg.markTimedOutResponse(item, pendingRequest.seeder)
				shouldRequestNow = true
			}
		} else {
			shouldRequestNow = true
		}

		if !shouldRequestNow {
			continue
		}

		seeders := rg.getSeedersFn(item)
		rankedNonExcludedSeeders := rg.rankedSeeders(seeders, pendingItemState.temporarilyExcludedSeeders)
		// If we have no remaining seeders because we have excluded all of them,
		// clear the exclusion list. We still need to make progress fetching the
		// thing, and we could have excluded the oracles due to a transient
		// issue on our end even.
		if len(rankedNonExcludedSeeders) == 0 && len(seeders) != 0 {
			clear(pendingItemState.temporarilyExcludedSeeders)
			rankedNonExcludedSeeders = rg.rankedSeeders(seeders, pendingItemState.temporarilyExcludedSeeders)
		}

		for _, seeder := range rankedNonExcludedSeeders {
			if rg.oracles[seeder].nextPossibleSendTimestamp.After(now) {
				continue
			}

			// try sending to this oracle
			requestInfo, ok := rg.sendRequestFn(item, seeder)
			if !ok {
				continue
			}

			rg.oracles[seeder].nextPossibleSendTimestamp = now.Add(rg.requestInterval)
			pendingItemState.pendingRequestOrNil = &pendingRequest{
				seeder,
				requestInfo.ExpiryTimestamp,
			}
			break
		}

		var nextTickForThisRequest time.Time
		if pendingItemState.pendingRequestOrNil != nil {
			// We sent a request in this tick, but want to recheck in case of
			// timeout.
			nextTickForThisRequest = pendingItemState.pendingRequestOrNil.expiryTimestamp
		} else if len(rankedNonExcludedSeeders) > 0 {
			// We didn't manage to send a request in this tick, so want to send
			// a request in the next tick, preferably to the best ranked seeder.
			// No guarantee we'll be able to get to them first though, prior
			// pending requests in the list will have priority.
			nextTickForThisRequest = rg.oracles[rankedNonExcludedSeeders[0]].nextPossibleSendTimestamp // <= now.Add(rg.requestInterval)
		}

		if nextTickForThisRequest.Before(nextTick) {
			nextTick = nextTickForThisRequest
		}
	}

	rg.chTick = time.After(time.Until(nextTick))
}

// A RequesterGadget helps us track and send requests for some data
// from a set of *seeders* that may or may not be able to serve the requests.
// Seeders may be byzantine, crashed, or just slow.
//
// Not thread-safe. RequesterGadget is expected to be integrated into a single subprotocol
// event loop via selecting on Ticker() and then calling Tick().
//
// Response processing is not handled by RequesterGadget. It is the responsibility
// of the subprotocol integrating RequesterGadget. After a response is received, the subprotocol must call
// CheckAndMarkResponse. It should also call one of MarkGoAwayResponse, MarkGoodResponse,
// MarkBadResponse, MarkGoodResponder, MarkBadResponder once the response has been processed.
type RequesterGadget[Item comparable] struct {
	oracles         map[commontypes.OracleID]*oracleState
	requestInterval time.Duration
	ourPendingItems map[Item]*pendingItemState
	chTick          <-chan time.Time

	sendRequestFn     func(Item, commontypes.OracleID) (*RequestInfo, bool)
	getPendingItemsFn func() []Item
	getSeedersFn      func(Item) map[commontypes.OracleID]struct{}
}

type pendingItemState struct {
	pendingRequestOrNil        *pendingRequest
	temporarilyExcludedSeeders map[commontypes.OracleID]struct{}
}

type pendingRequest struct {
	seeder          commontypes.OracleID
	expiryTimestamp time.Time
}

type oracleState struct {
	nextPossibleSendTimestamp time.Time
	score                     uint64
}
