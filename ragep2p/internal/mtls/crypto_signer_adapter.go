package mtls

import (
	"crypto"
	"crypto/ed25519"
	"fmt"
	"io"

	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

type peerKeyringCryptoSignerAdapter struct {
	keyring types.PeerKeyring
}

// Public implements crypto.Signer.
func (p *peerKeyringCryptoSignerAdapter) Public() crypto.PublicKey {
	pk := p.keyring.PublicKey()
	return ed25519.PublicKey(pk[:])
}

// Sign implements crypto.Signer.
func (p *peerKeyringCryptoSignerAdapter) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	// We can safely ignore the io.Reader providing randomness, since we use
	// deterministic Ed25519 here.
	if opts != crypto.Hash(0) {
		return nil, fmt.Errorf("unexpected SignerOpts for Ed25519: %v", opts)
	}
	return p.keyring.Sign(digest)
}

var _ crypto.Signer = &peerKeyringCryptoSignerAdapter{}
