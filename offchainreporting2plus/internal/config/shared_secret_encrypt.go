package config

import (
	"crypto/aes"
	"fmt"
	"io"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/curve25519"
)

// EncryptSharedSecretDeterministic constructs a SharedSecretEncryptions from
// a set of ConfigEncryptionPublicKeys, the sharedSecret, and an
// ephemeral secret key
func EncryptSharedSecretDeterministic(
	publicKeys []types.ConfigEncryptionPublicKey,
	sharedSecret *[SharedSecretSize]byte,
	ephemeralSk *[curve25519.ScalarSize]byte,
) (SharedSecretEncryptions, error) {
	ephemeralPk, err := curve25519.X25519(ephemeralSk[:], curve25519.Basepoint)
	if err != nil {
		return SharedSecretEncryptions{}, fmt.Errorf("while computing ephemeral pk: %w", err)
	}

	var ephemeralPkArray [curve25519.PointSize]byte
	copy(ephemeralPkArray[:], ephemeralPk)

	encryptedSharedSecrets := []EncryptedSharedSecret{}
	for _, pk := range publicKeys { // encrypt sharedSecret with each pk
		pkBytes := [curve25519.PointSize]byte(pk)
		dhPoint, err := curve25519.X25519(ephemeralSk[:], pkBytes[:])
		if err != nil {
			return SharedSecretEncryptions{}, fmt.Errorf("while computing dhPoint for %x: %w", pkBytes, err)
		}

		key := crypto.Keccak256(dhPoint)[:16]

		encryptedSharedSecret := EncryptedSharedSecret(aesEncryptBlock(key, sharedSecret[:]))
		encryptedSharedSecrets = append(encryptedSharedSecrets, encryptedSharedSecret)
	}

	return SharedSecretEncryptions{
		ephemeralPkArray,
		common.BytesToHash(crypto.Keccak256(sharedSecret[:])),
		encryptedSharedSecrets,
	}, nil
}

// EncryptSharedSecret constructs a SharedSecretEncryptions from
// a set of SharedSecretEncryptionPublicKeys, the sharedSecret, and a cryptographic
// randomness source
func EncryptSharedSecret(
	keys []types.ConfigEncryptionPublicKey,
	sharedSecret *[SharedSecretSize]byte,
	rand io.Reader,
) (SharedSecretEncryptions, error) {
	var sk [curve25519.ScalarSize]byte
	_, err := io.ReadFull(rand, sk[:])
	if err != nil {
		return SharedSecretEncryptions{}, fmt.Errorf("could not read enough randomness for encryption: %w", err)
	}
	return EncryptSharedSecretDeterministic(keys, sharedSecret, &sk)
}

// Encrypt one block with AES-128
func aesEncryptBlock(key, plaintext []byte) [16]byte {
	if len(key) != 16 {
		panic("key has wrong length")
	}
	if len(plaintext) != 16 {
		panic("ciphertext has wrong length")
	}

	cipher, err := aes.NewCipher(key)
	if err != nil {
		// assertion
		panic(fmt.Sprintf("Unexpected error during aes.NewCipher: %v", err))
	}

	var ciphertext [16]byte
	cipher.Encrypt(ciphertext[:], plaintext)
	return ciphertext
}
