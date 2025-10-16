package types

import (
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
)

type BinaryMessageOutboundPriority byte

const (
	_ BinaryMessageOutboundPriority = iota
	BinaryMessagePriorityLow
	BinaryMessagePriorityDefault
)

type ResponsePolicy interface {
	isResponsePolicy()
}

type SingleUseSizedLimitedResponsePolicy struct {
	MaxSize         int // TODO the name must demonstrate what size is measured in
	ExpiryTimestamp time.Time
}

func (SingleUseSizedLimitedResponsePolicy) isResponsePolicy() {}

type RequestHandle interface {
	MakeResponse(payload []byte) OutboundBinaryMessageResponse
}

type OutboundBinaryMessage interface {
	isOutboundBinaryMessage()
	GetPayload() []byte
}

var _ OutboundBinaryMessage = OutboundBinaryMessagePlain{}
var _ OutboundBinaryMessage = OutboundBinaryMessageRequest{}
var _ OutboundBinaryMessage = OutboundBinaryMessageResponse{}

type OutboundBinaryMessagePlain struct {
	Payload  []byte
	Priority BinaryMessageOutboundPriority
}

func (OutboundBinaryMessagePlain) isOutboundBinaryMessage() {}

func (o OutboundBinaryMessagePlain) GetPayload() []byte {
	return o.Payload
}

type OutboundBinaryMessageRequest struct {
	ResponsePolicy ResponsePolicy
	Payload        []byte
	Priority       BinaryMessageOutboundPriority
}

func (OutboundBinaryMessageRequest) isOutboundBinaryMessage() {}

func (o OutboundBinaryMessageRequest) GetPayload() []byte {
	return o.Payload
}

type OutboundBinaryMessageResponse struct {
	// By making the request handle private, we want to discourage folks from creating
	// this structure directly (unless they're implementing a BinaryNetworkEndpoint).
	// Note that, with a ragep2p backend (in its current version), we need the Response
	// priority to match the Request priority. Otherwise, responses would be dropped.
	// We try to protect a user of the interface from this sharp edge.
	requestHandle RequestHandle
	Payload       []byte
	Priority      BinaryMessageOutboundPriority
}

// Don't use this function unless you're implementing a BinaryNetworkEndpoint!
// The purpose of this function is to enable implementers of a RequestHandle instance to
// generate a OutboundBinaryMessageResponse in RequestHandle.MakeResponse()
func MustMakeOutboundBinaryMessageResponse(requestHandle RequestHandle, payload []byte, priority BinaryMessageOutboundPriority) OutboundBinaryMessageResponse {
	return OutboundBinaryMessageResponse{
		requestHandle,
		payload,
		priority,
	}
}

// Don't use this function unless you're implementing a BinaryNetworkEndpoint!
func MustGetOutboundBinaryMessageResponseRequestHandle(msg OutboundBinaryMessageResponse) RequestHandle {
	return msg.requestHandle
}

func (OutboundBinaryMessageResponse) isOutboundBinaryMessage() {}

func (o OutboundBinaryMessageResponse) GetPayload() []byte {
	return o.Payload
}

type InboundBinaryMessage interface {
	isInboundBinaryMessage()
	GetPayload() []byte
}

var _ InboundBinaryMessage = InboundBinaryMessagePlain{}
var _ InboundBinaryMessage = InboundBinaryMessageRequest{}
var _ InboundBinaryMessage = InboundBinaryMessageResponse{}

type InboundBinaryMessagePlain struct {
	Payload  []byte
	Priority BinaryMessageOutboundPriority // the priority the sender used for transmitting this message
}

func (InboundBinaryMessagePlain) isInboundBinaryMessage() {}

func (i InboundBinaryMessagePlain) GetPayload() []byte {
	return i.Payload
}

type InboundBinaryMessageRequest struct {
	RequestHandle RequestHandle
	Payload       []byte
	Priority      BinaryMessageOutboundPriority // the priority the sender used for transmitting this request
}

func (InboundBinaryMessageRequest) isInboundBinaryMessage() {}

func (i InboundBinaryMessageRequest) GetPayload() []byte {
	return i.Payload
}

type InboundBinaryMessageResponse struct {
	Payload  []byte
	Priority BinaryMessageOutboundPriority // the priority the sender used for transmitting this response
}

func (InboundBinaryMessageResponse) isInboundBinaryMessage() {}

func (i InboundBinaryMessageResponse) GetPayload() []byte {
	return i.Payload
}

type InboundBinaryMessageWithSender struct {
	InboundBinaryMessage
	Sender commontypes.OracleID
}

type BinaryNetworkEndpoint2 interface {
	SendTo(msg OutboundBinaryMessage, to commontypes.OracleID)
	Broadcast(msg OutboundBinaryMessage)
	Receive() <-chan InboundBinaryMessageWithSender
	Close() error
}
