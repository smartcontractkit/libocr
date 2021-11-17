package testimplementations

import (
	"io"
	"log"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
)

type OffchainKeyring struct {
	sigPrivKey ed25519.PrivateKey
	dhPrivKey  *[curve25519.ScalarSize]byte
}

var _ types.OffchainKeyring = OffchainKeyring{}

func NewOffchainKeyring(rand io.Reader) *OffchainKeyring {
	_, sigPrivKey, err := ed25519.GenerateKey(rand)
	if err != nil {
		panic(err)
	}
	var dhPrivKey [curve25519.ScalarSize]byte
	_, err = rand.Read(dhPrivKey[:])
	if err != nil {
		panic(err)
	}
	return &OffchainKeyring{
		sigPrivKey,
		&dhPrivKey,
	}
}

func (ok OffchainKeyring) OffchainSign(msg []byte) (signature []byte, err error) {
	sig := ed25519.Sign(ok.sigPrivKey, msg)
	return sig, nil
}

func (ok OffchainKeyring) ConfigDiffieHellman(
	point [curve25519.PointSize]byte,
) (
	sharedPoint [curve25519.PointSize]byte,
	err error,
) {
	p, err := curve25519.X25519(ok.dhPrivKey[:], point[:])
	if err != nil {
		return [curve25519.PointSize]byte{}, err
	}
	copy(sharedPoint[:], p)
	return sharedPoint, nil
}

func (ok OffchainKeyring) OffchainPublicKey() types.OffchainPublicKey {
	return types.OffchainPublicKey(ok.sigPrivKey.Public().(ed25519.PublicKey))
}

func (ok OffchainKeyring) ConfigEncryptionPublicKey() types.ConfigEncryptionPublicKey {
	rv, err := curve25519.X25519(ok.dhPrivKey[:], curve25519.Basepoint)
	if err != nil {
		log.Println("failure while computing public key: " + err.Error())
	}
	var rvFixed [curve25519.PointSize]byte
	copy(rvFixed[:], rv)
	return rvFixed
}
