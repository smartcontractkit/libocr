package frame

import (
	"encoding/binary"
	"fmt"

	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/internaltypes"
	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

const (
	streamIDSize  = internaltypes.StreamIDSize
	requestIDSize = internaltypes.RequestIDSize
)

var errMaxMessageSizeExceeded = fmt.Errorf("frame header error: the message size must not exceed %v bytes", types.MaxMessageLength)
var errOpenStreamPayloadSizeExceeded = fmt.Errorf(
	"frame header error: the payload size specified in a 'OpenStream' frame header must be at most %v bytes",
	types.MaxStreamNameLength,
)
var errCloseStreamPayloadSizeInvalid = fmt.Errorf("frame header error: the payload size specified in a 'CloseStream' frame header must be 0")
var errFrameHeaderSizeInvalid = fmt.Errorf("frame header error: invalid header length")
var errInvalidFrameType = fmt.Errorf("frame header error: invalid frame type")

type FrameType byte

const (
	_ FrameType = iota
	FrameTypeOpenStream
	FrameTypeCloseStream
	FrameTypeMessagePlain
	FrameTypeMessageRequest
	FrameTypeMessageResponse
)

//go-sumtype:decl FrameHeader

type FrameHeader interface {
	isFrameHeader()

	GetType() FrameType
	GetPayloadSize() int
	GetStreamID() internaltypes.StreamID

	Encode() []byte
}

const (
	baseFrameHeaderSize            = 1 + 32 + 4
	openStreamFrameHeaderSize      = baseFrameHeaderSize
	closeStreamFrameHeaderSize     = baseFrameHeaderSize
	MaxControlFrameHeaderSize      = baseFrameHeaderSize // maximum size of a control frame (open/close) header
	messagePlainFrameHeaderSize    = baseFrameHeaderSize
	messageRequestFrameHeaderSize  = baseFrameHeaderSize + requestIDSize
	messageResponseFrameHeaderSize = baseFrameHeaderSize + requestIDSize
	MaxFrameHeaderSize             = baseFrameHeaderSize + requestIDSize // whatever value is the largest header size from above
)

var frameHeaderSizes = map[FrameType]int{
	FrameTypeOpenStream:      openStreamFrameHeaderSize,
	FrameTypeCloseStream:     closeStreamFrameHeaderSize,
	FrameTypeMessagePlain:    messagePlainFrameHeaderSize,
	FrameTypeMessageRequest:  messageRequestFrameHeaderSize,
	FrameTypeMessageResponse: messageResponseFrameHeaderSize,
}

// The different frame header types must be wire-compatible with the previously used frame header structure.
// See encodeBaseFrameHeader(...) and decodeBaseFrameHeader().
//
// type frameHeader struct {
//     Type          frameType
//     StreamID      streamID
//     PayloadSize   uint32
// }

type OpenStreamFrameHeader struct {
	StreamID    internaltypes.StreamID
	PayloadSize int
}

type CloseStreamFrameHeader struct {
	StreamID internaltypes.StreamID
}

type MessagePlainFrameHeader struct {
	StreamID    internaltypes.StreamID
	PayloadSize int
}

type MessageRequestFrameHeader struct {
	StreamID    internaltypes.StreamID
	PayloadSize int
	RequestID   internaltypes.RequestID
}

type MessageResponseFrameHeader struct {
	StreamID    internaltypes.StreamID
	PayloadSize int
	RequestID   internaltypes.RequestID
}

func (OpenStreamFrameHeader) isFrameHeader()      {}
func (CloseStreamFrameHeader) isFrameHeader()     {}
func (MessagePlainFrameHeader) isFrameHeader()    {}
func (MessageRequestFrameHeader) isFrameHeader()  {}
func (MessageResponseFrameHeader) isFrameHeader() {}

func (OpenStreamFrameHeader) GetType() FrameType {
	return FrameTypeOpenStream
}
func (CloseStreamFrameHeader) GetType() FrameType {
	return FrameTypeCloseStream
}
func (MessagePlainFrameHeader) GetType() FrameType {
	return FrameTypeMessagePlain
}
func (MessageRequestFrameHeader) GetType() FrameType {
	return FrameTypeMessageRequest
}
func (MessageResponseFrameHeader) GetType() FrameType {
	return FrameTypeMessageResponse
}

func (h OpenStreamFrameHeader) GetPayloadSize() int {
	return h.PayloadSize
}
func (h CloseStreamFrameHeader) GetPayloadSize() int {
	return 0
}
func (h MessagePlainFrameHeader) GetPayloadSize() int {
	return h.PayloadSize
}
func (h MessageRequestFrameHeader) GetPayloadSize() int {
	return h.PayloadSize
}
func (h MessageResponseFrameHeader) GetPayloadSize() int {
	return h.PayloadSize
}

func (h OpenStreamFrameHeader) GetStreamID() internaltypes.StreamID {
	return h.StreamID
}
func (h CloseStreamFrameHeader) GetStreamID() internaltypes.StreamID {
	return h.StreamID
}
func (h MessagePlainFrameHeader) GetStreamID() internaltypes.StreamID {
	return h.StreamID
}
func (h MessageRequestFrameHeader) GetStreamID() internaltypes.StreamID {
	return h.StreamID
}
func (h MessageResponseFrameHeader) GetStreamID() internaltypes.StreamID {
	return h.StreamID
}

func encodeBaseFrameHeader(frameType FrameType, streamID internaltypes.StreamID, payloadSize int, extraBufferCapacity int) []byte {
	buffer := make([]byte, 0, baseFrameHeaderSize+extraBufferCapacity)
	buffer = append(buffer, byte(frameType))
	buffer = append(buffer, streamID[:]...)
	buffer = binary.BigEndian.AppendUint32(buffer, uint32(payloadSize))
	return buffer
}

func decodeBaseFrameHeader(encoded []byte, expectedType FrameType, expectedSize int) (internaltypes.StreamID, int, error) {
	var streamID internaltypes.StreamID

	if len(encoded) != expectedSize {
		return internaltypes.StreamID{}, 0, errFrameHeaderSizeInvalid
	}
	if FrameType(encoded[0]) != expectedType {
		return internaltypes.StreamID{}, 0, errInvalidFrameType
	}

	payloadSize := binary.BigEndian.Uint32(encoded[1+streamIDSize:])
	if payloadSize > types.MaxMessageLength {
		return internaltypes.StreamID{}, 0, errMaxMessageSizeExceeded
	}

	copy(streamID[:], encoded[1:streamIDSize+1])

	return streamID, int(payloadSize), nil
}

func (h OpenStreamFrameHeader) Encode() []byte {
	return encodeBaseFrameHeader(FrameTypeOpenStream, h.StreamID, h.PayloadSize, 0)
}

func (h CloseStreamFrameHeader) Encode() []byte {
	return encodeBaseFrameHeader(FrameTypeCloseStream, h.StreamID, 0, 0)
}

func (h MessagePlainFrameHeader) Encode() []byte {
	return encodeBaseFrameHeader(FrameTypeMessagePlain, h.StreamID, h.PayloadSize, 0)
}

func (h MessageRequestFrameHeader) Encode() []byte {
	buffer := encodeBaseFrameHeader(FrameTypeMessageRequest, h.StreamID, h.PayloadSize, requestIDSize)
	buffer = append(buffer, h.RequestID[:]...)
	return buffer
}

func (h MessageResponseFrameHeader) Encode() []byte {
	buffer := encodeBaseFrameHeader(FrameTypeMessageResponse, h.StreamID, h.PayloadSize, requestIDSize)
	buffer = append(buffer, h.RequestID[:]...)
	return buffer
}

func decodeFrameHeader(encoded []byte) (FrameHeader, error) {
	if len(encoded) == 0 {
		return nil, errFrameHeaderSizeInvalid
	}

	switch FrameType(encoded[0]) {
	case FrameTypeOpenStream:
		return decodeOpenStreamFrameHeader(encoded)
	case FrameTypeCloseStream:
		return decodeCloseStreamFrameHeader(encoded)
	case FrameTypeMessagePlain:
		return decodeMessagePlainFrameHeader(encoded)
	case FrameTypeMessageRequest:
		return decodeMessageRequestFrameHeader(encoded)
	case FrameTypeMessageResponse:
		return decodeMessageResponseFrameHeader(encoded)
	default:
		return nil, errInvalidFrameType
	}
}

func decodeOpenStreamFrameHeader(encoded []byte) (OpenStreamFrameHeader, error) {
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, FrameTypeOpenStream, openStreamFrameHeaderSize)
	if err != nil {
		return OpenStreamFrameHeader{}, err
	}
	if payloadSize > types.MaxStreamNameLength {
		return OpenStreamFrameHeader{}, errOpenStreamPayloadSizeExceeded
	}

	return OpenStreamFrameHeader{streamID, payloadSize}, nil
}

func decodeCloseStreamFrameHeader(encoded []byte) (CloseStreamFrameHeader, error) {
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, FrameTypeCloseStream, closeStreamFrameHeaderSize)
	if err != nil {
		return CloseStreamFrameHeader{}, err
	}
	if payloadSize != 0 {
		return CloseStreamFrameHeader{}, errCloseStreamPayloadSizeInvalid
	}

	return CloseStreamFrameHeader{streamID}, nil
}

func decodeMessagePlainFrameHeader(encoded []byte) (MessagePlainFrameHeader, error) {
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, FrameTypeMessagePlain, messagePlainFrameHeaderSize)
	if err != nil {
		return MessagePlainFrameHeader{}, err
	}
	return MessagePlainFrameHeader{streamID, payloadSize}, err
}

func decodeMessageRequestFrameHeader(encoded []byte) (MessageRequestFrameHeader, error) {
	var requestID internaltypes.RequestID
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, FrameTypeMessageRequest, messageRequestFrameHeaderSize)
	if err != nil {
		return MessageRequestFrameHeader{}, err
	}

	copy(requestID[:], encoded[baseFrameHeaderSize:])
	return MessageRequestFrameHeader{streamID, payloadSize, requestID}, nil
}

func decodeMessageResponseFrameHeader(encoded []byte) (MessageResponseFrameHeader, error) {
	var requestID internaltypes.RequestID
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, FrameTypeMessageResponse, messageResponseFrameHeaderSize)
	if err != nil {
		return MessageResponseFrameHeader{}, err
	}

	copy(requestID[:], encoded[baseFrameHeaderSize:])
	return MessageResponseFrameHeader{streamID, payloadSize, requestID}, nil
}

var ErrReadFrameHeaderReadFailed = fmt.Errorf("failed to read frame header")

type FrameHeaderReader struct {
	readFn func(buf []byte) bool
	buf    []byte
}

// readFn must return false if it wasn't able to completely fill the buffer.
// readFn is assumed to be stateful and remember what was read already.
func MakeFrameHeaderReader(readFn func(buf []byte) bool) FrameHeaderReader {
	return FrameHeaderReader{readFn, make([]byte, MaxFrameHeaderSize)}
}

func (r *FrameHeaderReader) ReadFrameHeader() (header FrameHeader, err error) {
	// Read frame type.
	if !r.readFn(r.buf[:1]) {
		return nil, ErrReadFrameHeaderReadFailed
	}

	// Get the length of the frame header for the given type. Abort if the type is invalid.
	frameType := FrameType(r.buf[0])
	headerSize, ok := frameHeaderSizes[frameType]
	if !ok {
		return nil, fmt.Errorf("invalid frame type: %d", frameType)
	}

	// Read the rest of the frame header.
	if !r.readFn(r.buf[1:headerSize]) {
		return nil, ErrReadFrameHeaderReadFailed
	}

	// Decode the frame header.
	header, err = decodeFrameHeader(r.buf[:headerSize])
	if err != nil {
		return nil, fmt.Errorf("%w caused by raw frame header %x", err, r.buf[:headerSize])
	}

	return header, nil
}
