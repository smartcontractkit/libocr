package ragep2p

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	internal_types "github.com/smartcontractkit/libocr/ragep2p/internal/types"
)

const (
	streamIDSize  = internal_types.StreamIDSize
	requestIDSize = internal_types.RequestIDSize
)

var errMaxMessageSizeExceeded = fmt.Errorf("the message size must not exceed %v bytes", MaxMessageLength)
var errOpenStreamPayloadSizeExceeded = fmt.Errorf(
	"the payload size specified in a 'OpenStream' frame header must be at most %v bytes",
	MaxStreamNameLength,
)
var errCloseStreamPayloadSizeInvalid = fmt.Errorf("the payload size specified in a 'CloseStream' frame header must be 0")
var errFrameHeaderSizeInvalid = fmt.Errorf("frame decoding error, invalid header length")
var errInvalidFrameType = fmt.Errorf("frame decoding error, invalid frame type")

type frameType byte

const (
	_ frameType = iota
	frameTypeOpenStream
	frameTypeCloseStream
	frameTypeMessage
	frameTypeRequest
	frameTypeResponse
)

type frameHeader interface {
	Type() frameType
	Encode() []byte
	PayloadSize() int
	StreamID() streamID
}

const (
	baseFrameHeaderSize        = 1 + 32 + 4
	openStreamFrameHeaderSize  = baseFrameHeaderSize
	closeStreamFrameHeaderSize = baseFrameHeaderSize
	messageFrameHeaderSize     = baseFrameHeaderSize
	requestFrameHeaderSize     = baseFrameHeaderSize + requestIDSize
	responseFrameHeaderSize    = baseFrameHeaderSize + requestIDSize
	maxFrameHeaderSize         = baseFrameHeaderSize + requestIDSize // whatever value is the largest header size from above
)

var frameHeaderSizes = map[frameType]int{
	frameTypeOpenStream:  openStreamFrameHeaderSize,
	frameTypeCloseStream: closeStreamFrameHeaderSize,
	frameTypeMessage:     messageFrameHeaderSize,
	frameTypeRequest:     requestFrameHeaderSize,
	frameTypeResponse:    responseFrameHeaderSize,
}

// The different frame header types must be wire-compatible with the previously used frame header structure.
// See encodeBaseFrameHeader(...) and decodeBaseFrameHeader().
//
// type frameHeader struct {
//     Type          frameType
//     StreamID      streamID
//     PayloadSize   uint32
// }

type openStreamFrameHeader struct {
	streamID    streamID
	payloadSize int
}

type closeStreamFrameHeader struct {
	streamID streamID
}

type messageFrameHeader struct {
	streamID    streamID
	payloadSize int
}

type requestFrameHeader struct {
	streamID    streamID
	payloadSize int
	requestID   requestID
}

type responseFrameHeader struct {
	streamID    streamID
	payloadSize int
	requestID   requestID
}

func (_ openStreamFrameHeader) Type() frameType {
	return frameTypeOpenStream
}
func (_ closeStreamFrameHeader) Type() frameType {
	return frameTypeCloseStream
}
func (_ messageFrameHeader) Type() frameType {
	return frameTypeMessage
}
func (_ requestFrameHeader) Type() frameType {
	return frameTypeRequest
}
func (_ responseFrameHeader) Type() frameType {
	return frameTypeResponse
}

func (h openStreamFrameHeader) PayloadSize() int {
	return h.payloadSize
}
func (h closeStreamFrameHeader) PayloadSize() int {
	return 0
}
func (h messageFrameHeader) PayloadSize() int {
	return h.payloadSize
}
func (h requestFrameHeader) PayloadSize() int {
	return h.payloadSize
}
func (h responseFrameHeader) PayloadSize() int {
	return h.payloadSize
}

func (h openStreamFrameHeader) StreamID() streamID {
	return h.streamID
}
func (h closeStreamFrameHeader) StreamID() streamID {
	return h.streamID
}
func (h messageFrameHeader) StreamID() streamID {
	return h.streamID
}
func (h requestFrameHeader) StreamID() streamID {
	return h.streamID
}
func (h responseFrameHeader) StreamID() streamID {
	return h.streamID
}

func encodeBaseFrameHeader(frameType frameType, streamID streamID, payloadSize int, extraBufferCapacity int) []byte {
	buffer := make([]byte, 0, baseFrameHeaderSize+extraBufferCapacity)
	buffer = append(buffer, byte(frameType))
	buffer = append(buffer, streamID[:]...)
	buffer = binary.BigEndian.AppendUint32(buffer, uint32(payloadSize))
	return buffer
}

func decodeBaseFrameHeader(encoded []byte, expectedType frameType, expectedSize int) (streamID, int, error) {
	var streamID streamID

	if len(encoded) != expectedSize {
		return streamID, 0, errFrameHeaderSizeInvalid
	}
	if frameType(encoded[0]) != expectedType {
		return streamID, 0, errInvalidFrameType
	}

	payloadSize := binary.BigEndian.Uint32(encoded[1+streamIDSize:])
	if payloadSize > MaxMessageLength {
		return streamID, 0, errMaxMessageSizeExceeded
	}

	copy(streamID[:], encoded[1:streamIDSize+1])

	return streamID, int(payloadSize), nil
}

func (h openStreamFrameHeader) Encode() []byte {
	return encodeBaseFrameHeader(frameTypeOpenStream, h.streamID, h.payloadSize, 0)
}

func (h closeStreamFrameHeader) Encode() []byte {
	return encodeBaseFrameHeader(frameTypeCloseStream, h.streamID, 0, 0)
}

func (h messageFrameHeader) Encode() []byte {
	return encodeBaseFrameHeader(frameTypeMessage, h.streamID, h.payloadSize, 0)
}

func (h requestFrameHeader) Encode() []byte {
	buffer := encodeBaseFrameHeader(frameTypeRequest, h.streamID, h.payloadSize, requestIDSize)
	buffer = append(buffer, h.requestID[:]...)
	return buffer
}

func (h responseFrameHeader) Encode() []byte {
	buffer := encodeBaseFrameHeader(frameTypeResponse, h.streamID, h.payloadSize, requestIDSize)
	buffer = append(buffer, h.requestID[:]...)
	return buffer
}

func decodeFrameHeader(encoded []byte) (frameHeader, error) {
	if len(encoded) == 0 {
		return nil, errFrameHeaderSizeInvalid
	}

	switch frameType(encoded[0]) {
	case frameTypeOpenStream:
		return decodeOpenStreamFrameHeader(encoded)
	case frameTypeCloseStream:
		return decodeCloseStreamFrameHeader(encoded)
	case frameTypeMessage:
		return decodeMessageFrameHeader(encoded)
	case frameTypeRequest:
		return decodeRequestFrameHeader(encoded)
	case frameTypeResponse:
		return decodeResponseFrameHeader(encoded)
	default:
		return nil, errInvalidFrameType
	}
}

func decodeOpenStreamFrameHeader(encoded []byte) (openStreamFrameHeader, error) {
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, frameTypeOpenStream, openStreamFrameHeaderSize)
	if err != nil {
		return openStreamFrameHeader{}, err
	}
	if payloadSize > MaxStreamNameLength {
		return openStreamFrameHeader{}, errOpenStreamPayloadSizeExceeded
	}

	return openStreamFrameHeader{streamID, payloadSize}, nil
}

func decodeCloseStreamFrameHeader(encoded []byte) (closeStreamFrameHeader, error) {
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, frameTypeCloseStream, closeStreamFrameHeaderSize)
	if err != nil {
		return closeStreamFrameHeader{}, err
	}
	if payloadSize != 0 {
		return closeStreamFrameHeader{}, errCloseStreamPayloadSizeInvalid
	}

	return closeStreamFrameHeader{streamID}, nil
}

func decodeMessageFrameHeader(encoded []byte) (messageFrameHeader, error) {
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, frameTypeMessage, messageFrameHeaderSize)
	if err != nil {
		return messageFrameHeader{}, err
	}
	return messageFrameHeader{streamID, payloadSize}, err
}

func decodeRequestFrameHeader(encoded []byte) (requestFrameHeader, error) {
	var requestID requestID
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, frameTypeRequest, requestFrameHeaderSize)
	if err != nil {
		return requestFrameHeader{}, err
	}

	copy(requestID[:], encoded[baseFrameHeaderSize:])
	return requestFrameHeader{streamID, payloadSize, requestID}, nil
}

func decodeResponseFrameHeader(encoded []byte) (responseFrameHeader, error) {
	var requestID requestID
	streamID, payloadSize, err := decodeBaseFrameHeader(encoded, frameTypeResponse, responseFrameHeaderSize)
	if err != nil {
		return responseFrameHeader{}, err
	}

	copy(requestID[:], encoded[baseFrameHeaderSize:])
	return responseFrameHeader{streamID, payloadSize, requestID}, nil
}

func getStreamID(self types.PeerID, other types.PeerID, name string) streamID {
	if bytes.Compare(self[:], other[:]) < 0 {
		return getStreamID(other, self, name)
	}

	h := sha256.New()
	h.Write(self[:])
	h.Write(other[:])
	// this is fine because self and other are of constant length. if more than
	// one variable length item is ever added here, we should also hash lengths
	// to prevent collisions.
	h.Write([]byte(name))

	var result streamID
	copy(result[:], h.Sum(nil))
	return result
}

func getRandomRequestID() (requestID, error) {
	requestID := requestID{}
	_, err := rand.Read(requestID[:])
	return requestID, err
}

func getRandomStreamID() (streamID, error) {
	streamID := streamID{}
	_, err := rand.Read(streamID[:])
	return streamID, err
}
