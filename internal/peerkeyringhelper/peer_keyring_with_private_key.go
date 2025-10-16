package peerkeyringhelper

import (
	"crypto/ed25519"
	"fmt"

	ragetypes "github.com/RoSpaceDev/libocr/ragep2p/types"
)

type PeerKeyringWithPrivateKey struct {
	privateKey    ed25519.PrivateKey
	peerPublicKey ragetypes.PeerPublicKey
}

var _ ragetypes.PeerKeyring = &PeerKeyringWithPrivateKey{}

func NewPeerKeyringWithPrivateKey(privateKey ed25519.PrivateKey) (*PeerKeyringWithPrivateKey, error) {
	if err := ed25519SanityCheck(privateKey); err != nil {
		return nil, fmt.Errorf("ed25519 sanity check failed: %w", err)
	}
	peerPublicKey, err := ragetypes.PeerPublicKeyFromGenericPublicKey(privateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("StaticallySizedEd25519PublicKey failed even though sanity check succeeded: %w", err)
	}
	return &PeerKeyringWithPrivateKey{privateKey, peerPublicKey}, nil
}

// PublicKey implements ragetypes.PeerKeyring.
func (s *PeerKeyringWithPrivateKey) PublicKey() ragetypes.PeerPublicKey {
	return s.peerPublicKey
}

// Sign implements ragetypes.PeerKeyring.
func (s *PeerKeyringWithPrivateKey) Sign(msg []byte) (signature []byte, err error) {
	return ed25519.Sign(s.privateKey, msg), nil
}
