package muxer

import (
	"fmt"
	"sync"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/internal/randmap"
	"github.com/RoSpaceDev/libocr/internal/ringbuffer"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/internaltypes"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/stream2types"
)

const invertPrioritiesEvery = 8

type streamRecord struct {
	streamName    string
	priority      stream2types.StreamPriority
	enabled       bool
	messageBuffer *ringbuffer.RingBuffer[stream2types.OutboundBinaryMessage]
}

type Muxer struct {
	logger loghelper.LoggerWithContext

	chSignal chan struct{}

	mutex                                    sync.Mutex
	streamRecords                            map[internaltypes.StreamID]*streamRecord
	defaultPriorityStreamsWithPendingMessage *randmap.Map[internaltypes.StreamID, *ringbuffer.RingBuffer[stream2types.OutboundBinaryMessage]]
	lowPriorityStreamsWithPendingMessage     *randmap.Map[internaltypes.StreamID, *ringbuffer.RingBuffer[stream2types.OutboundBinaryMessage]]
	popCount                                 uint
}

func NewMuxer(logger loghelper.LoggerWithContext) *Muxer {
	return &Muxer{
		logger,

		make(chan struct{}, 1),
		sync.Mutex{},
		map[internaltypes.StreamID]*streamRecord{},
		randmap.NewMap[internaltypes.StreamID, *ringbuffer.RingBuffer[stream2types.OutboundBinaryMessage]](),
		randmap.NewMap[internaltypes.StreamID, *ringbuffer.RingBuffer[stream2types.OutboundBinaryMessage]](),

		0,
	}
}

// Adds a stream to the Muxer. The stream is initially disabled.
//
// panics if maxOutgoingBufferedMessages is 0 or less
func (mux *Muxer) AddStream(
	sid internaltypes.StreamID,
	streamName string,
	priority stream2types.StreamPriority,
	maxOutgoingBufferedMessages int,
) bool {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()

	if _, exists := mux.streamRecords[sid]; exists {
		return false
	}

	mux.streamRecords[sid] = &streamRecord{
		streamName,
		priority,
		false,
		ringbuffer.NewRingBuffer[stream2types.OutboundBinaryMessage](maxOutgoingBufferedMessages),
	}

	mux.logger.Debug("Muxer: stream added", commontypes.LogFields{
		"sid":                         sid,
		"streamName":                  streamName,
		"priority":                    priority,
		"maxOutgoingBufferedMessages": maxOutgoingBufferedMessages,
	})
	return true
}

// panics if maxOutgoingBufferedMessages is 0 or less
func (mux *Muxer) UpdateStream(sid internaltypes.StreamID, maxOutgoingBufferedMessages int) bool {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()

	streamRecord, ok := mux.streamRecords[sid]
	if !ok {
		return false
	}

	streamRecord.messageBuffer.SetCap(maxOutgoingBufferedMessages)
	return true
}

func (mux *Muxer) RemoveStream(sid internaltypes.StreamID) bool {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()

	streamRecord, ok := mux.streamRecords[sid]
	if !ok {
		return false
	}

	streamsWithPendingMessage := mux.streamsWithPendingMessageForPriority(streamRecord.priority)
	streamsWithPendingMessage.Delete(sid)

	delete(mux.streamRecords, sid)

	mux.logger.Debug("Muxer: stream removed", commontypes.LogFields{
		"sid": sid,
	})
	return true
}

// Enables a stream to emit messages via Pop().
func (mux *Muxer) EnableStream(sid internaltypes.StreamID) bool {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()

	streamRecord, ok := mux.streamRecords[sid]
	if !ok {
		return false
	}

	streamRecord.enabled = true

	if streamRecord.messageBuffer.IsEmpty() {
		return true
	}

	streamsWithPendingMessage := mux.streamsWithPendingMessageForPriority(streamRecord.priority)
	streamsWithPendingMessage.Set(sid, streamRecord.messageBuffer)

	select {
	case mux.chSignal <- struct{}{}:
	default:
	}

	mux.logger.Debug("Muxer: stream enabled", commontypes.LogFields{
		"sid": sid,
	})
	return true
}

// Disables a stream from emitting messages via Pop(). Messages can
// still be pushed to the stream's buffer via PushEvict() and will
// be emitted once the stream is enabled again.
func (mux *Muxer) DisableStream(sid internaltypes.StreamID) bool {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()

	streamRecord, ok := mux.streamRecords[sid]
	if !ok {
		return false
	}

	streamRecord.enabled = false

	streamsWithPendingMessage := mux.streamsWithPendingMessageForPriority(streamRecord.priority)
	streamsWithPendingMessage.Delete(sid)

	mux.logger.Debug("Muxer: stream disabled", commontypes.LogFields{
		"sid": sid,
	})
	return true
}

// Pushes a message to the stream's ring buffer of messages and evicts the oldest message if the buffer is full.
func (mux *Muxer) PushEvict(sid internaltypes.StreamID, m stream2types.OutboundBinaryMessage) bool {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()

	streamRecord, ok := mux.streamRecords[sid]
	if !ok {
		return false
	}

	streamRecord.messageBuffer.PushEvict(m)

	if streamRecord.enabled {
		streamsWithPendingMessage := mux.streamsWithPendingMessageForPriority(streamRecord.priority)
		streamsWithPendingMessage.Set(sid, streamRecord.messageBuffer)
	}

	select {
	case mux.chSignal <- struct{}{}:
	default:
	}

	return true
}

// Signals that there are (maybe) messages available via Pop(). Note that unlike
// its counterpart in the demuxer, Pop() will not raise the flag again. A consumer
// of this signal is assumed to keep Pop()-ing until there are no more messages.
func (mux *Muxer) SignalMaybePending() <-chan struct{} {
	return mux.chSignal
}

func (mux *Muxer) streamsWithPendingMessageForPriority(priority stream2types.StreamPriority) *randmap.Map[internaltypes.StreamID, *ringbuffer.RingBuffer[stream2types.OutboundBinaryMessage]] {
	switch priority {
	case stream2types.StreamPriorityDefault:
		return mux.defaultPriorityStreamsWithPendingMessage
	case stream2types.StreamPriorityLow:
		return mux.lowPriorityStreamsWithPendingMessage
	}
	panic(fmt.Sprintf("unexpected stream priority: %#v", priority))
}

func (mux *Muxer) popForPriority(priority stream2types.StreamPriority) (stream2types.OutboundBinaryMessage, internaltypes.StreamID) {
	streamsWithPendingMessage := mux.streamsWithPendingMessageForPriority(priority)

	if streamsWithPendingMessage.Size() == 0 {
		return nil, internaltypes.StreamID{}
	}

	entry, ok := streamsWithPendingMessage.GetRandom()
	if !ok {
		return nil, internaltypes.StreamID{}
	}
	sid, messageBuffer := entry.Key, entry.Value

	msg, ok := messageBuffer.Pop()
	if !ok {
		mux.logger.Error("Muxer: message buffer is unexpectedly empty, this points to a bug", commontypes.LogFields{
			"sid": sid,
		})
		return nil, internaltypes.StreamID{}
	}

	if messageBuffer.IsEmpty() {
		streamsWithPendingMessage.Delete(sid)
	}

	return msg, sid
}

// If any stream is enabled and has messages buffered, one of these will be returned.
// The policy for choosing which message to return is an implementation detail and may change.
// For now, we implement a simple policy that prefers DefaultPriority messages most
// of the time, and LowPriority messages once in a while. If there are multiple enabled streams
// of equal priority with buffered messages, we choose one uniformly at random.
func (mux *Muxer) Pop() (stream2types.OutboundBinaryMessage, internaltypes.StreamID) {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()

	// Overflow is harmless
	mux.popCount++

	var highPriority, lowPriority stream2types.StreamPriority
	if mux.popCount%invertPrioritiesEvery == 0 {
		highPriority = stream2types.StreamPriorityLow
		lowPriority = stream2types.StreamPriorityDefault
	} else {
		highPriority = stream2types.StreamPriorityDefault
		lowPriority = stream2types.StreamPriorityLow
	}

	msg, sid := mux.popForPriority(highPriority)
	if msg != nil {
		return msg, sid
	}

	msg, sid = mux.popForPriority(lowPriority)
	if msg != nil {
		return msg, sid
	}

	mux.logger.Info("Muxer: nothing to pop", commontypes.LogFields{})
	return nil, internaltypes.StreamID{}
}
