package stream2types

import (
	"fmt"

	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/internaltypes"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/responselimit"
	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

type StreamPriority byte

const (
	_ StreamPriority = iota
	StreamPriorityLow
	StreamPriorityDefault
)

type Stream2Limits struct {
	MaxOutgoingBufferedMessages int // number of messages that fit in the outgoing buffer
	MaxIncomingBufferedMessages int // number of messages that fit in the incoming buffer
	MaxMessageLength            int
	MessagesLimit               types.TokenBucketParams // rate limit for (the number of) incoming messages
	BytesLimit                  types.TokenBucketParams // rate limit for (the accumulated size in bytes of) incoming messages
}

func (limits Stream2Limits) Validate() (ValidatedStream2Limits, error) {
	if !(limits.MaxOutgoingBufferedMessages > 0) {
		return ValidatedStream2Limits{}, fmt.Errorf("maxOutgoingBufferedMessages %v is not positive", limits.MaxOutgoingBufferedMessages)
	}
	if !(limits.MaxIncomingBufferedMessages > 0) {
		return ValidatedStream2Limits{}, fmt.Errorf("maxIncomingBufferedMessages %v is not positive", limits.MaxIncomingBufferedMessages)
	}
	if !(limits.MaxMessageLength <= types.MaxMessageLength) {
		return ValidatedStream2Limits{}, fmt.Errorf("maxMessageLength %v is not less than or equal to global MaxMessageLength %v", limits.MaxMessageLength, types.MaxMessageLength)
	}
	if !(0 <= limits.MessagesLimit.Rate) {
		return ValidatedStream2Limits{}, fmt.Errorf("messagesLimit.Rate %v is not non-negative", limits.MessagesLimit.Rate)
	}
	//lint:ignore SA4003
	if !(0 <= limits.MessagesLimit.Capacity) { //nolint:staticcheck
		return ValidatedStream2Limits{}, fmt.Errorf("messagesLimit.Capacity %v is not non-negative", limits.MessagesLimit.Capacity)
	}
	if !(0 <= limits.BytesLimit.Rate) {
		return ValidatedStream2Limits{}, fmt.Errorf("bytesLimit.Rate %v is not non-negative", limits.BytesLimit.Rate)
	}
	//lint:ignore SA4003
	if !(0 <= limits.BytesLimit.Capacity) { //nolint:staticcheck
		return ValidatedStream2Limits{}, fmt.Errorf("bytesLimit.Capacity %v is not non-negative", limits.BytesLimit.Capacity)
	}
	return ValidatedStream2Limits{limits, struct{}{}}, nil
}

type ValidatedStream2Limits struct {
	Stream2Limits
	private struct{}
}

type RequestHandle internaltypes.RequestID

func (r *RequestHandle) MakeResponse(payload []byte) OutboundBinaryMessageResponse {
	return OutboundBinaryMessageResponse{
		internaltypes.RequestID(*r),
		payload,
	}
}

//go-sumtype:decl InboundBinaryMessage

type InboundBinaryMessage interface {
	isInboundBinaryMessage()
}

var _ InboundBinaryMessage = InboundBinaryMessagePlain{}
var _ InboundBinaryMessage = InboundBinaryMessageRequest{}
var _ InboundBinaryMessage = InboundBinaryMessageResponse{}

type InboundBinaryMessagePlain struct {
	Payload []byte
}

func (InboundBinaryMessagePlain) isInboundBinaryMessage() {}

type InboundBinaryMessageRequest struct {
	RequestHandle RequestHandle
	Payload       []byte
}

func (InboundBinaryMessageRequest) isInboundBinaryMessage() {}

type InboundBinaryMessageResponse struct {
	Payload []byte
}

func (InboundBinaryMessageResponse) isInboundBinaryMessage() {}

//go-sumtype:decl OutboundBinaryMessage

type OutboundBinaryMessage interface {
	isOutboundBinaryMessage()
}

var _ OutboundBinaryMessage = OutboundBinaryMessagePlain{}
var _ OutboundBinaryMessage = OutboundBinaryMessageRequest{}
var _ OutboundBinaryMessage = OutboundBinaryMessageResponse{}

type OutboundBinaryMessagePlain struct {
	Payload []byte
}

func (OutboundBinaryMessagePlain) isOutboundBinaryMessage() {}

type OutboundBinaryMessageRequest struct {
	ResponsePolicy responselimit.ResponsePolicy
	Payload        []byte
}

func (OutboundBinaryMessageRequest) isOutboundBinaryMessage() {}

type OutboundBinaryMessageResponse struct {
	requestID internaltypes.RequestID
	Payload   []byte
}

func (OutboundBinaryMessageResponse) isOutboundBinaryMessage() {}

func RequestIDOfOutboundBinaryMessageResponse(m OutboundBinaryMessageResponse) internaltypes.RequestID {
	return m.requestID
}
