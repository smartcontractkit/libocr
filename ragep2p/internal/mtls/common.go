package mtls

import (
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/big"

	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

// Generates a minimal certificate (that wouldn't be considered valid outside this telemetry networking protocol)
// from a PeerKeyring.
func NewMinimalX509CertFromKeyring(keyring types.PeerKeyring) (tls.Certificate, error) {
	template := x509.Certificate{
		SerialNumber: big.NewInt(0), // serial number must be set, so we set it to 0
	}

	// x509 args are of type any, even though crypto.Signer is required
	var signer crypto.Signer = &peerKeyringCryptoSignerAdapter{keyring}

	encodedCert, err := x509.CreateCertificate(rand.Reader, &template, &template, signer.Public(), signer)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("x509.CreateCertificate: %w", err)
	}

	// Uncomment this if you want to get an encoded cert you can feed into openssl x509 etc...
	//
	// var buf bytes.Buffer
	// if err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: encodedCert}); err != nil {
	// 	log.Fatalf("Failed to encode cert into pem format: %v", err)
	// }
	// fmt.Printf("pubkey: %x\nencodedCert: %v\n", signer.Public(), buf.String())

	return tls.Certificate{
		Certificate: [][]byte{encodedCert},

		PrivateKey:                   signer,
		SupportedSignatureAlgorithms: []tls.SignatureScheme{tls.Ed25519},
	}, nil
}

func VerifyCertMatchesPubKey(publicKey [32]byte) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if len(rawCerts) != 1 {
			return fmt.Errorf("required exactly one client certificate")
		}
		cert, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return err
		}
		pk, err := PubKeyFromCert(cert)
		if err != nil {
			return err
		}

		if pk != publicKey {
			return fmt.Errorf("unknown public key on cert: %x doesn't match expected public key %x", pk, publicKey)
		}

		return nil
	}
}

func PubKeyFromCert(cert *x509.Certificate) (pk types.PeerPublicKey, err error) {
	if cert.PublicKeyAlgorithm != x509.Ed25519 {
		return pk, fmt.Errorf("require ed25519 public key")
	}
	return types.PeerPublicKeyFromGenericPublicKey(cert.PublicKey)
}
