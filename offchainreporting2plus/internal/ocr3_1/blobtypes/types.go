package blobtypes

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding"
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

// Returns a byte slice whose first four bytes are the string "ocr3" and the rest
// of which is the sum returned by h. Used for domain separation vs ocr2, where
// we just directly sign sha256 hashes.
//
// Any signatures made with the OffchainKeyring should use ocr3_1DomainSeparatedSum!
func ocr3_1DomainSeparatedSum(h hash.Hash) []byte {
	result := make([]byte, 0, 6+32)
	result = append(result, []byte("ocr3.1")...)
	return h.Sum(result)
}

type BlobChunkDigest [32]byte

var _ fmt.Stringer = BlobChunkDigest{}

func (bcd BlobChunkDigest) String() string {
	return fmt.Sprintf("%x", bcd[:])
}

func MakeBlobChunkDigest(chunk []byte) BlobChunkDigest {
	h := sha256.New()
	h.Write(chunk)
	var result BlobChunkDigest
	h.Sum(result[:0])
	return result
}

type BlobDigest [32]byte

var _ fmt.Stringer = BlobDigest{}

func (bd BlobDigest) String() string {
	return fmt.Sprintf("%x", bd[:])
}

func MakeBlobDigest(
	configDigest types.ConfigDigest,
	chunkDigests []BlobChunkDigest,
	payloadLength uint64,
	expirySeqNr uint64,
	submitter commontypes.OracleID,
) BlobDigest {
	h := sha256.New()

	_, _ = h.Write(configDigest[:])

	_ = binary.Write(h, binary.BigEndian, uint64(len(chunkDigests)))
	for _, chunkDigest := range chunkDigests {

		_, _ = h.Write(chunkDigest[:])
	}

	_ = binary.Write(h, binary.BigEndian, payloadLength)

	_ = binary.Write(h, binary.BigEndian, expirySeqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(submitter))

	var result BlobDigest
	h.Sum(result[:0])
	return result
}

const blobAvailabilitySignatureDomainSeparator = "ocr3.1/BlobAvailabilitySignature/"

type BlobAvailabilitySignature []byte

func MakeBlobAvailabilitySignature(
	blobDigest BlobDigest,
	signer func(msg []byte) ([]byte, error),
) (BlobAvailabilitySignature, error) {
	return signer(blobAvailabilitySignatureMsg(blobDigest))
}

func (sig BlobAvailabilitySignature) Verify(
	blobDigest BlobDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, blobAvailabilitySignatureMsg(blobDigest), sig)
	if !ok {
		return fmt.Errorf("BlobAvailabilitySignature failed to verify")
	}

	return nil
}

func blobAvailabilitySignatureMsg(
	blobDigest BlobDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(blobAvailabilitySignatureDomainSeparator))

	_, _ = h.Write(blobDigest[:])

	return ocr3_1DomainSeparatedSum(h)
}

type AttributedBlobAvailabilitySignature struct {
	Signature BlobAvailabilitySignature
	Signer    commontypes.OracleID
}

type LightCertifiedBlob struct {
	ChunkDigests  []BlobChunkDigest
	PayloadLength uint64
	ExpirySeqNr   uint64
	Submitter     commontypes.OracleID

	AttributedBlobAvailabilitySignatures []AttributedBlobAvailabilitySignature
}

func (lc *LightCertifiedBlob) Verify(
	configDigest types.ConfigDigest,
	oracleIdentities []config.OracleIdentity,
	fPlusOneSize int,
	byzQuorumSize int,
) error {
	if !(fPlusOneSize <= len(lc.AttributedBlobAvailabilitySignatures) && len(lc.AttributedBlobAvailabilitySignatures) <= byzQuorumSize) {
		return fmt.Errorf("wrong number of signatures, expected in range [%d, %d] for quorum but got %d", fPlusOneSize, byzQuorumSize, len(lc.AttributedBlobAvailabilitySignatures))
	}

	blobDigest := MakeBlobDigest(
		configDigest,
		lc.ChunkDigests,
		lc.PayloadLength,
		lc.ExpirySeqNr,
		lc.Submitter,
	)

	seen := make(map[commontypes.OracleID]bool)
	for i, abs := range lc.AttributedBlobAvailabilitySignatures {
		if seen[abs.Signer] {
			return fmt.Errorf("duplicate signature by %v", abs.Signer)
		}
		seen[abs.Signer] = true
		if !(0 <= int(abs.Signer) && int(abs.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", abs.Signer)
		}
		if err := abs.Signature.Verify(blobDigest, oracleIdentities[abs.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, abs.Signer, oracleIdentities[abs.Signer].OffchainPublicKey, err)
		}
	}

	return nil
}

var _ BlobHandleSumType = &LightCertifiedBlob{}

func (lc *LightCertifiedBlob) isBlobHandleSumType() {}

// go-sumtype:decl BlobHandleSumType

type BlobHandleSumType interface {
	isBlobHandleSumType()
	encoding.BinaryAppender
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type BlobHandle struct {
	_ [0]func()

	blobHandleSumType BlobHandleSumType
}

func ExtractBlobHandleSumType(h BlobHandle) BlobHandleSumType {
	return h.blobHandleSumType
}

func MakeBlobHandle(b BlobHandleSumType) BlobHandle {
	return BlobHandle{
		[0]func(){},
		b,
	}
}

var _ encoding.BinaryAppender = BlobHandle{}
var _ encoding.BinaryMarshaler = BlobHandle{}
var _ encoding.BinaryUnmarshaler = &BlobHandle{}

type blobHandleSumTypeVariant byte

const (
	_ blobHandleSumTypeVariant = iota
	blobHandleSumTypeVariantLightCertifiedBlob
)

func (h BlobHandle) AppendBinary(b []byte) ([]byte, error) {
	var variant blobHandleSumTypeVariant
	switch h.blobHandleSumType.(type) {
	case *LightCertifiedBlob:
		variant = blobHandleSumTypeVariantLightCertifiedBlob
	}

	prefix := append(b, byte(variant))
	final, err := h.blobHandleSumType.AppendBinary(prefix)
	if err != nil {
		return nil, err
	}

	return final, nil
}

func (h BlobHandle) MarshalBinary() ([]byte, error) {
	return h.AppendBinary(nil)
}

func (h *BlobHandle) UnmarshalBinary(data []byte) error {
	if len(data) < 1 {
		return fmt.Errorf("data too short to unmarshal BlobHandle: %d bytes", len(data))
	}
	variant, rest := blobHandleSumTypeVariant(data[0]), data[1:]
	switch variant {
	case blobHandleSumTypeVariantLightCertifiedBlob:
		var lc LightCertifiedBlob
		if err := lc.UnmarshalBinary(rest); err != nil {
			return fmt.Errorf("failed to unmarshal LightCertifiedBlob: %w", err)
		}
		h.blobHandleSumType = &lc
	default:
		return fmt.Errorf("unknown BlobHandle version: %d", data[0])
	}
	return nil
}
