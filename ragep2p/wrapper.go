package ragep2p

import (
	"github.com/RoSpaceDev/libocr/networking/ragep2pwrapper"
	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

func Wrapped(host *Host) ragep2pwrapper.Host {
	return &hostWrapper{host}
}

var _ ragep2pwrapper.Host = &hostWrapper{}

type hostWrapper struct {
	host *Host
}

var _ ragep2pwrapper.Stream = &streamWrapper{}

type streamWrapper struct {
	stream *Stream
}

func (h *hostWrapper) Start() error {
	return h.host.Start()
}

func (h *hostWrapper) Close() error {
	return h.host.Close()
}

func (h *hostWrapper) ID() types.PeerID {
	return h.host.ID()
}

func (h *hostWrapper) NewStream(
	other types.PeerID,
	streamName string,
	outgoingBufferSize int,
	incomingBufferSize int,
	maxMessageLength int,
	messagesLimit types.TokenBucketParams,
	bytesLimit types.TokenBucketParams,
) (ragep2pwrapper.Stream, error) {
	stream, err := h.host.NewStream(other, streamName, outgoingBufferSize, incomingBufferSize, maxMessageLength, messagesLimit, bytesLimit)
	if err != nil {
		return nil, err
	}
	return &streamWrapper{stream}, nil
}

func (h *hostWrapper) RawWrappee() any {
	return h.host
}

func (s *streamWrapper) Other() types.PeerID {
	return s.stream.Other()
}

func (s *streamWrapper) Name() string {
	return s.stream.Name()
}

func (s *streamWrapper) SendMessage(data []byte) {
	s.stream.SendMessage(data)
}

func (s *streamWrapper) ReceiveMessages() <-chan []byte {
	return s.stream.ReceiveMessages()
}

func (s *streamWrapper) Close() error {
	return s.stream.Close()
}
