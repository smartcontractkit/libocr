// Temporary wrapper package to allow switching between ragep2p and ragep2pnew.
// Eventually, we will remove this package and use ragep2pnew directly.
package ragep2pwrapper

import (
	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

type Host interface {
	Start() error
	Close() error
	ID() types.PeerID
	NewStream(
		other types.PeerID,
		streamName string,
		outgoingBufferSize int,
		incomingBufferSize int,
		maxMessageLength int,
		messagesLimit types.TokenBucketParams,
		bytesLimit types.TokenBucketParams,
	) (Stream, error)

	RawWrappee() any
}

type Stream interface {
	Other() types.PeerID
	Name() string
	SendMessage(data []byte)
	ReceiveMessages() <-chan []byte
	Close() error
}
