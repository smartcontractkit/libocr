package peeridhelper

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"

	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

// context string acting as domain separator according to RFC8446
const context = "ragep2p arbitrary message signing"

// We hash messages longer than shortMessageMax to avoid copying large messages (and allocating lots of memory in the process)
const shortMessageMax = 1024

// MakePeerIDSignatureDomainSeparatedPayload provides a safe way to generate signature payloads for
// signing by the keypair associated with a PeerID. The payloads are guaranteed to not be confused
// with messages signed within ragep2p.
//
// Be sure to use domain separators that are unique to your application and do not
// collide with other applications!
//
// The output of this function must be directly used as the payload for the underlying ed25519.Sign function.
func MakePeerIDSignatureDomainSeparatedPayload(domainSeparator string, message []byte) []byte {
	// We use the same domain separation scheme as TLS 1.3, but with a different context string. Details below.
	// RFC8446 specifies the following domain separation scheme for CertificateVerify messages:
	// The digital signature is then computed over the concatenation of:
	// -  A string that consists of octet 32 (0x20) repeated 64 times
	// -  The context string
	// -  A single 0 byte which serves as the separator
	// -  The content to be signed
	// We use the same domain separation scheme, but use a different context string
	// to ensure domain separation with TLS.

	capacity := 64 + // 64 times 0x20 from RFC8446
		len(context) + // context from RFC8446
		1 + // 0x0 byte from RFC8446
		8 + // length of domainSeparator
		len(domainSeparator) // domainSeparator
	if len(message) <= shortMessageMax {
		capacity += 1 + // 0x0 to indicate message is "short"
			8 + // length of message
			len(message) // message
	} else {
		capacity += 1 + // 0x1 to indicate message is "long"
			sha256.Size // sha256 hash of message
	}

	// Follow RFC8446 ...
	payload := make([]byte, 0, capacity)
	for range 64 {
		payload = append(payload, 0x20)
	}
	payload = append(payload, []byte(context)...)
	payload = append(payload, 0x0)

	// ... and now that we have domain separation from all things TLS,
	// we perform further domain separation between different users of
	// MakePeerIDSignaturePayload ...
	payload = binary.BigEndian.AppendUint64(payload, uint64(len(domainSeparator)))
	payload = append(payload, []byte(domainSeparator)...)
	// ... and finally append the message (or a hash of it, if the message is too long)
	if len(message) <= shortMessageMax {
		payload = append(payload, 0x0)

		payload = binary.BigEndian.AppendUint64(payload, uint64(len(message)))
		payload = append(payload, message...)
	} else {
		payload = append(payload, 0x1)

		h := sha256.New()
		h.Write(message)
		payload = h.Sum(payload)
	}

	return payload
}

// DomainSeparatedPeerKeyring is a wrapper around a PeerKeyring that ensures
// messages signed with it are domain separated from messages signed within ragep2p.
//
// Be sure to use domain separators that are unique to your application and do not
// collide with other applications!
//
// We intentionally do not implement the PeerKeyring interface, because we don't
// want to encourage nesting of DomainSeparatedPeerKeyring.
type DomainSeparatedPeerKeyring struct {
	domainSeparator string
	peerKeyring     types.PeerKeyring
}

func NewDomainSeparatedPeerKeyring(domainSeparator string, peerKeyring types.PeerKeyring) *DomainSeparatedPeerKeyring {
	return &DomainSeparatedPeerKeyring{
		domainSeparator,
		peerKeyring,
	}
}

func (dk *DomainSeparatedPeerKeyring) Sign(message []byte) (signature []byte, err error) {
	return dk.peerKeyring.Sign(MakePeerIDSignatureDomainSeparatedPayload(dk.domainSeparator, message))
}

func (dk *DomainSeparatedPeerKeyring) Verify(publicKey types.PeerPublicKey, message []byte, signature []byte) bool {
	// not needed strictly speaking, added for defense in depth
	if len(publicKey) != ed25519.PublicKeySize {
		return false
	}

	return ed25519.Verify(ed25519.PublicKey(publicKey[:]), MakePeerIDSignatureDomainSeparatedPayload(dk.domainSeparator, message), signature)
}

func (dk *DomainSeparatedPeerKeyring) DomainSeparatorAndPublicKey() (string, types.PeerPublicKey) {
	return dk.domainSeparator, dk.peerKeyring.PublicKey()
}
