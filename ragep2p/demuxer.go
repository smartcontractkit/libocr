package ragep2p

import (
	"fmt"
	"math"
	"sync"

	"github.com/smartcontractkit/libocr/internal/ringbuffer"
	"github.com/smartcontractkit/libocr/ragep2p/internal/ratelimit"
	"github.com/smartcontractkit/libocr/ragep2p/internal/responselimit"
)

type shouldPushResult int

const (
	_ shouldPushResult = iota
	shouldPushResultYes
	shouldPushResultMessageTooBig
	shouldPushResultMessagesLimitExceeded
	shouldPushResultBytesLimitExceeded
	shouldPushResultUnknownStream
	shouldPushResultResponseRejected
)

type pushResult int

const (
	_ pushResult = iota
	pushResultSuccess
	pushResultDropped
	pushResultUnknownStream
)

type popResult int

const (
	_ popResult = iota
	popResultSuccess
	popResultEmpty
	popResultUnknownStream
)

type demuxerStream struct {
	buffer          *ringbuffer.RingBuffer[InboundBinaryMessage]
	chSignal        chan struct{}
	maxMessageSize  int
	messagesLimiter ratelimit.TokenBucket
	bytesLimiter    ratelimit.TokenBucket
}

type demuxer struct {
	mutex           sync.Mutex
	streams         map[streamID]*demuxerStream
	responseChecker *responselimit.ResponseChecker
}

func newDemuxer() *demuxer {
	return &demuxer{
		sync.Mutex{},
		map[streamID]*demuxerStream{},
		responselimit.NewResponseChecker(),
	}
}

func makeRateLimiter(params TokenBucketParams) ratelimit.TokenBucket {
	tb := ratelimit.TokenBucket{}
	tb.SetRate(ratelimit.MillitokensPerSecond(math.Ceil(params.Rate * 1000)))
	tb.SetCapacity(params.Capacity)
	tb.AddTokens(params.Capacity)
	return tb
}

func (d *demuxer) AddStream(
	sid streamID,
	incomingBufferSize int,
	maxMessageSize int,
	messagesLimit TokenBucketParams,
	bytesLimit TokenBucketParams,
) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, ok := d.streams[sid]; ok {
		return false
	}

	d.streams[sid] = &demuxerStream{
		ringbuffer.NewRingBuffer[InboundBinaryMessage](incomingBufferSize),
		make(chan struct{}, 1),
		maxMessageSize,
		makeRateLimiter(messagesLimit),
		makeRateLimiter(bytesLimit),
	}
	return true
}

func (d *demuxer) RemoveStream(sid streamID) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.streams, sid)
	d.responseChecker.ClearPoliciesForStream(sid)
}

func (d *demuxer) ShouldPush(sid streamID, size int) shouldPushResult {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return shouldPushResultUnknownStream
	}

	if size > s.maxMessageSize {
		return shouldPushResultMessageTooBig
	}

	messagesLimiterAllow := s.messagesLimiter.RemoveTokens(1)
	bytesLimiterAllow := s.bytesLimiter.RemoveTokens(uint32(size))

	if !messagesLimiterAllow {
		return shouldPushResultMessagesLimitExceeded
	}

	if !bytesLimiterAllow {
		return shouldPushResultBytesLimitExceeded
	}

	return shouldPushResultYes
}

func (d *demuxer) ShouldPushResponse(sid streamID, rid requestID, size int) shouldPushResult {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.streams[sid]
	if !ok {
		return shouldPushResultUnknownStream
	}

	checkResult := d.responseChecker.CheckResponse(sid, rid, size)
	switch checkResult {
	case responselimit.ResponseCheckResultReject:
		return shouldPushResultResponseRejected
	case responselimit.ResponseCheckResultAllow:
		return shouldPushResultYes
	}

	// The above switch should be exhaustive. If it is not the linter is expected to catch this.
	panic(fmt.Sprintf("unexpected ragep2p.responseCheckResult: %#v", checkResult))
}

func (d *demuxer) PushMessage(sid streamID, msg InboundBinaryMessage) pushResult {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return pushResultUnknownStream
	}

	result := pushResultSuccess
	if _, evicted := s.buffer.PushEvict(msg); evicted {
		result = pushResultDropped
	}

	select {
	case s.chSignal <- struct{}{}:
	default:
	}

	return result
}

// Pops a message from the underlying stream's buffer.
// Returns a non-nil value iff popResult == popResultSuccess.
func (d *demuxer) PopMessage(sid streamID) (InboundBinaryMessage, popResult) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return nil, popResultUnknownStream
	}

	result, ok := s.buffer.Pop()
	if !ok {
		return nil, popResultEmpty
	}

	if !s.buffer.IsEmpty() {
		select {
		case s.chSignal <- struct{}{}:
		default:
		}
	}

	return result, popResultSuccess
}

// The signals received via the returned channel are NOT a reliable indicator that the buffer is NON-empty. Depending on
// the exact interleaving of the goroutines (in particular, authenticatedConnectionReadLoop, and receiveLoop), a call
// to PopMessage() - after receiving a signal through the channel - may return (nil, popResultEmpty).
//
// Example execution timeline for a buffer size of 1:
//
// | authenticatedConnectionReadLoop   buffer   receiveLoop
// |                                   []
// | demux.PushMessage(m1)
// |                                   [m1]
// | send signal to s.chSignal
// |                                            signal received (case <-chSignalMaybePending triggers)
// | demux.PushMessage(m2), buffer
// | overflows and m1 is dropped
// |                                   [m2]
// |                                            demux.PopMessage() returns (m2, popResultSuccess)
// |                                   []
// | send signal to s.chSignal
// |                                            signal received (case <-chSignalMaybePending triggers)
// |                                            demux.PopMessage() returns (nil, popResultEmpty)
// ▼
// time
func (d *demuxer) SignalMaybePending(sid streamID) <-chan struct{} {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return nil
	}

	return s.chSignal
}
