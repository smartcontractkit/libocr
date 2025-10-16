package demuxer

import (
	"fmt"
	"math"
	"sync"

	"github.com/RoSpaceDev/libocr/internal/ringbuffer"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/internaltypes"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/ratelimit"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/responselimit"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/stream2types"
	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

type ShouldPushResult int

const (
	_ ShouldPushResult = iota
	ShouldPushResultYes
	ShouldPushResultMessageTooBig
	ShouldPushResultMessagesLimitExceeded
	ShouldPushResultBytesLimitExceeded
	ShouldPushResultUnknownStream
	ShouldPushResultResponseRejected
)

type PushResult int

const (
	_ PushResult = iota
	PushResultSuccess
	PushResultDropped
	PushResultUnknownStream
)

type PopResult int

const (
	_ PopResult = iota
	PopResultSuccess
	PopResultEmpty
	PopResultUnknownStream
)

type demuxerStream struct {
	buffer          *ringbuffer.RingBuffer[stream2types.InboundBinaryMessage]
	chSignal        chan struct{}
	maxMessageSize  int
	messagesLimiter ratelimit.TokenBucket
	bytesLimiter    ratelimit.TokenBucket
}

// Demuxer helps, on the receiving side of a connection, to demux the arriving messages
// into the correct streams. In the process, it checks rate limits.
type Demuxer struct {
	mutex           sync.Mutex
	streams         map[internaltypes.StreamID]*demuxerStream
	responseChecker *responselimit.ResponseChecker
}

func NewDemuxer() *Demuxer {
	return &Demuxer{
		sync.Mutex{},
		map[internaltypes.StreamID]*demuxerStream{},
		responselimit.NewResponseChecker(),
	}
}

func makeTokenBucket(params types.TokenBucketParams) ratelimit.TokenBucket {
	tb := ratelimit.TokenBucket{}
	tb.SetRate(ratelimit.MillitokensPerSecond(math.Ceil(params.Rate * 1000)))
	tb.SetCapacity(params.Capacity)
	tb.AddTokens(params.Capacity)
	return tb
}

func updateTokenBucket(bucket *ratelimit.TokenBucket, params types.TokenBucketParams) {
	bucket.SetRate(ratelimit.MillitokensPerSecond(math.Ceil(params.Rate * 1000)))
	bucket.SetCapacity(params.Capacity)
}

// panics if incomingBufferSize is 0 or less
func (d *Demuxer) AddStream(
	sid internaltypes.StreamID,
	maxIncomingBufferedMessages int,
	maxMessageSize int,
	messagesLimit types.TokenBucketParams,
	bytesLimit types.TokenBucketParams,
) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, ok := d.streams[sid]; ok {
		return false
	}

	d.streams[sid] = &demuxerStream{
		ringbuffer.NewRingBuffer[stream2types.InboundBinaryMessage](maxIncomingBufferedMessages),
		make(chan struct{}, 1),
		maxMessageSize,
		makeTokenBucket(messagesLimit),
		makeTokenBucket(bytesLimit),
	}
	return true
}

func (d *Demuxer) UpdateStream(
	sid internaltypes.StreamID,
	maxIncomingBufferedMessages int,
	maxMessageSize int,
	messagesLimit types.TokenBucketParams,
	bytesLimit types.TokenBucketParams,
) bool {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return false
	}

	if maxIncomingBufferedMessages <= 0 {
		return false
	}

	s.buffer.SetCap(maxIncomingBufferedMessages)
	s.maxMessageSize = maxMessageSize
	updateTokenBucket(&s.messagesLimiter, messagesLimit)
	updateTokenBucket(&s.bytesLimiter, bytesLimit)

	return true
}

func (d *Demuxer) RemoveStream(sid internaltypes.StreamID) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.streams, sid)
	d.responseChecker.ClearPoliciesForStream(sid)
}

func (d *Demuxer) SetPolicy(sid internaltypes.StreamID, rid internaltypes.RequestID, policy responselimit.ResponsePolicy) {
	// reponseChecker.SetPolicy is threadsafe, no need to acquire d.mutex
	d.responseChecker.SetPolicy(sid, rid, policy)
}

func (d *Demuxer) ShouldPush(sid internaltypes.StreamID, size int) ShouldPushResult {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return ShouldPushResultUnknownStream
	}

	if size > s.maxMessageSize {
		return ShouldPushResultMessageTooBig
	}

	messagesLimiterAllow := s.messagesLimiter.RemoveTokens(1)
	bytesLimiterAllow := s.bytesLimiter.RemoveTokens(uint32(size))

	if !messagesLimiterAllow {
		return ShouldPushResultMessagesLimitExceeded
	}

	if !bytesLimiterAllow {
		return ShouldPushResultBytesLimitExceeded
	}

	return ShouldPushResultYes
}

func (d *Demuxer) ShouldPushResponse(sid internaltypes.StreamID, rid internaltypes.RequestID, size int) ShouldPushResult {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.streams[sid]
	if !ok {
		return ShouldPushResultUnknownStream
	}

	checkResult := d.responseChecker.CheckResponse(sid, rid, size)
	switch checkResult {
	case responselimit.ResponseCheckResultReject:
		return ShouldPushResultResponseRejected
	case responselimit.ResponseCheckResultAllow:
		return ShouldPushResultYes
	}

	// The above switch should be exhaustive. If it is not the linter is expected to catch this.
	panic(fmt.Sprintf("unexpected ragep2p.responseCheckResult: %#v", checkResult))
}

func (d *Demuxer) PushMessage(sid internaltypes.StreamID, msg stream2types.InboundBinaryMessage) PushResult {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return PushResultUnknownStream
	}

	result := PushResultSuccess
	if _, evicted := s.buffer.PushEvict(msg); evicted {
		result = PushResultDropped
	}

	select {
	case s.chSignal <- struct{}{}:
	default:
	}

	return result
}

// Pops a message from the underlying stream's buffer.
// Returns a non-nil value iff popResult == popResultSuccess.
func (d *Demuxer) PopMessage(sid internaltypes.StreamID) (stream2types.InboundBinaryMessage, PopResult) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return nil, PopResultUnknownStream
	}

	result, ok := s.buffer.Pop()
	if !ok {
		return nil, PopResultEmpty
	}

	if !s.buffer.IsEmpty() {
		select {
		case s.chSignal <- struct{}{}:
		default:
		}
	}

	return result, PopResultSuccess
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
func (d *Demuxer) SignalMaybePending(sid internaltypes.StreamID) <-chan struct{} {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s, ok := d.streams[sid]
	if !ok {
		return nil
	}

	return s.chSignal
}
