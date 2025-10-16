package ratelimitaggregator

import (
	"sync"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

// Aggregator aggregates the rate limits of all streams of a peer for o11y
// and sanity checking purposes.
type Aggregator struct {
	logger loghelper.LoggerWithContext

	mutex                        sync.Mutex
	messagesTokenBucketAggregate TokenBucketAggregate
	bytesTokenBucketAggregate    TokenBucketAggregate
}

type TokenBucketAggregate struct {
	Rate     float64
	Capacity float64
}

func NewAggregator(logger loghelper.LoggerWithContext) *Aggregator {
	return &Aggregator{
		logger,

		sync.Mutex{},
		TokenBucketAggregate{},
		TokenBucketAggregate{},
	}
}

// These are so large that nobody sane should hit them in practice.
const (
	warnThresholdMessagesRate     = 10_000
	warnThresholdMessagesCapacity = 5 * warnThresholdMessagesRate
	warnThresholdBytesRate        = 1_024 * 1_024 * 1_024
	warnThresholdBytesCapacity    = 5 * warnThresholdBytesRate
)

func (crl *Aggregator) AddStream(messagesLimit types.TokenBucketParams, bytesLimit types.TokenBucketParams) {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	crl.messagesTokenBucketAggregate.Rate += messagesLimit.Rate
	crl.messagesTokenBucketAggregate.Capacity += float64(messagesLimit.Capacity)
	crl.bytesTokenBucketAggregate.Rate += bytesLimit.Rate
	crl.bytesTokenBucketAggregate.Capacity += float64(bytesLimit.Capacity)

	if crl.messagesTokenBucketAggregate.Rate > warnThresholdMessagesRate {
		crl.logger.Warn("aggregate messages rate exceeds warning threshold, this likely points to buggy configuration of some ragep2p Stream", commontypes.LogFields{
			"messagesRate": crl.messagesTokenBucketAggregate.Rate,
			"threshold":    warnThresholdMessagesRate,
		})
	}
	if crl.messagesTokenBucketAggregate.Capacity > warnThresholdMessagesCapacity {
		crl.logger.Warn("aggregate messages capacity exceeds warning threshold, this likely points to buggy configuration of some ragep2p Stream", commontypes.LogFields{
			"messagesCapacity": crl.messagesTokenBucketAggregate.Capacity,
			"threshold":        warnThresholdMessagesCapacity,
		})
	}
	if crl.bytesTokenBucketAggregate.Rate > warnThresholdBytesRate {
		crl.logger.Warn("aggregate bytes rate exceeds warning threshold, this likely points to buggy configuration of some ragep2p Stream", commontypes.LogFields{
			"bytesRate": crl.bytesTokenBucketAggregate.Rate,
			"threshold": warnThresholdBytesRate,
		})
	}
	if crl.bytesTokenBucketAggregate.Capacity > warnThresholdBytesCapacity {
		crl.logger.Warn("aggregate bytes capacity exceeds warning threshold, this likely points to buggy configuration of some ragep2p Stream", commontypes.LogFields{
			"bytesCapacity": crl.bytesTokenBucketAggregate.Capacity,
			"threshold":     warnThresholdBytesCapacity,
		})
	}
}

func (crl *Aggregator) RemoveStream(messagesLimit types.TokenBucketParams, bytesLimit types.TokenBucketParams) {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	crl.messagesTokenBucketAggregate.Rate -= messagesLimit.Rate
	crl.messagesTokenBucketAggregate.Capacity -= float64(messagesLimit.Capacity)
	crl.bytesTokenBucketAggregate.Rate -= bytesLimit.Rate
	crl.bytesTokenBucketAggregate.Capacity -= float64(bytesLimit.Capacity)
}

func (crl *Aggregator) Aggregates() (messages TokenBucketAggregate, bytes TokenBucketAggregate) {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	return crl.messagesTokenBucketAggregate, crl.bytesTokenBucketAggregate
}
