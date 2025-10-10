package types

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"encoding"
	"fmt"

	"github.com/mr-tron/base58"
)

// Address represents a network address & port such as "192.168.1.2:8080". It
// can also contain special bind addresses such as "0.0.0.0:80".
type Address string

// PeerID represents a unique identifier for another peer.
type PeerID [ed25519.PublicKeySize]byte

var (
	_ fmt.Stringer               = PeerID{}
	_ encoding.TextMarshaler     = PeerID{}
	_ encoding.TextUnmarshaler   = &PeerID{}
	_ encoding.BinaryMarshaler   = PeerID{}
	_ encoding.BinaryUnmarshaler = &PeerID{}
)

func (p PeerID) String() string {
	text, err := p.MarshalText()
	if err != nil {
		return fmt.Sprintf("<PeerID: failed to stringify due to error '%s'>", err)
	}
	return string(text)
}

func (p PeerID) MarshalText() (text []byte, err error) {
	bin, err := p.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return []byte(base58.Encode(bin)), nil
}

func (p *PeerID) UnmarshalText(text []byte) error {
	b58, err := base58.Decode(string(text))
	if err != nil {
		return fmt.Errorf("failed to base58 decode: %w", err)
	}
	return p.UnmarshalBinary(b58)
}

// This magic value comes from libp2p's encoding of peer ids, see https://docs.libp2p.io/concepts/peer-id/
func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}

func (p PeerID) MarshalBinary() (data []byte, err error) {
	return append(libp2pMagic(), p[:]...), nil
}

func (p *PeerID) UnmarshalBinary(data []byte) error {
	expectedMagic := libp2pMagic()
	expectedLength := len(expectedMagic) + len(PeerID{})
	if len(data) != expectedLength {
		return fmt.Errorf("unexpected peerid length (was %d, expected %d)", len(data), expectedLength)
	}
	actualMagic := data[:len(expectedMagic)]
	if !bytes.Equal(actualMagic, expectedMagic) {
		return fmt.Errorf("unexpected peerid magic (was %x, expected %x)", actualMagic, expectedMagic)
	}
	copy(p[:], data[len(expectedMagic):])
	return nil
}

func PeerPublicKeyFromGenericPublicKey(publickey crypto.PublicKey) (PeerPublicKey, error) {
	var result PeerPublicKey

	pkslice, ok := publickey.(ed25519.PublicKey)
	if !ok {
		return result, fmt.Errorf("invalid ed25519 public key")
	}
	if ed25519.PublicKeySize != len(pkslice) {
		return result, fmt.Errorf("invalid key size (expected %d, actual %d)", ed25519.PublicKeySize, len(pkslice))
	}
	copy(result[:], pkslice)
	return result, nil
}

func MustPeerPublicKeyFromGenericPublicKey(publickey crypto.PublicKey) PeerPublicKey {
	result, err := PeerPublicKeyFromGenericPublicKey(publickey)
	if err != nil {
		panic(err)
	}
	return result
}

func PeerIDFromPeerPublicKey(pk PeerPublicKey) PeerID {
	return PeerID(pk)
}

func PeerIDFromPublicKey(pk ed25519.PublicKey) (PeerID, error) {
	peerPublicKey, err := PeerPublicKeyFromGenericPublicKey(pk)
	return PeerIDFromPeerPublicKey(peerPublicKey), err
}

func PeerIDFromPrivateKey(sk ed25519.PrivateKey) (PeerID, error) {
	return PeerIDFromPublicKey(sk.Public().(ed25519.PublicKey))
}

func PeerIDFromKeyring(keyring PeerKeyring) PeerID {
	return PeerIDFromPeerPublicKey(keyring.PublicKey())
}

type PeerInfo struct {
	ID    PeerID
	Addrs []Address
}

// PeerPublicKey is the public key used to cryptographically identify an
// oracle in p2p related communications and peer discovery.
type PeerPublicKey [ed25519.PublicKeySize]byte

func Ed25519PublicKeyFromPeerPublicKey(pk PeerPublicKey) ed25519.PublicKey {
	return pk[:]
}

type PeerKeyring interface {
	// Sign returns an EdDSA-Ed25519 signature on msg produced using the
	// standard library's ed25519.Sign function. Must be fast.
	Sign(msg []byte) (signature []byte, err error)

	// PublicKey returns the public component of the keypair used in Sign.
	PublicKey() PeerPublicKey
}

// TokenBucketParams contains the two parameters for a token bucket rate
// limiter.
type TokenBucketParams struct {
	Rate     float64
	Capacity uint32
}

const MaxStreamsPerPeer = 2_000

const MaxMessageLength = 1024 * 1024 * 1024 // 1 GiB. This must be smaller than INT32_MAX

const MaxStreamNameLength = 256
