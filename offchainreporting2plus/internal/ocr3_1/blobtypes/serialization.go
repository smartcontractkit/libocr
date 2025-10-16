package blobtypes

import (
	"crypto/ed25519"
	"encoding"
	"fmt"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes/serialization"
	"google.golang.org/protobuf/proto"
)

var _ encoding.BinaryMarshaler = LightCertifiedBlob{}
var _ encoding.BinaryAppender = LightCertifiedBlob{}
var _ encoding.BinaryUnmarshaler = &LightCertifiedBlob{}

func (lc LightCertifiedBlob) AppendBinary(b []byte) ([]byte, error) {
	pbChunkDigests := make([][]byte, 0, len(lc.ChunkDigests))
	for _, digest := range lc.ChunkDigests {
		pbChunkDigests = append(pbChunkDigests, digest[:])
	}

	pbSignatures := make([]*serialization.AttributedBlobAvailabilitySignature, 0, len(lc.AttributedBlobAvailabilitySignatures))
	for _, sig := range lc.AttributedBlobAvailabilitySignatures {
		pbSignatures = append(pbSignatures, serialization.NewAttributedBlobAvailabilitySignature(
			sig.Signature,
			uint32(sig.Signer),
		))
	}

	pbLightCertifiedBlob := serialization.NewLightCertifiedBlob(
		pbChunkDigests,
		lc.PayloadLength,
		lc.ExpirySeqNr,
		uint32(lc.Submitter),
		pbSignatures,
	)

	opts := proto.MarshalOptions{}
	ret, err := opts.MarshalAppend(b, pbLightCertifiedBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to MarshalAppend LightCertifiedBlob protobuf: %w", err)
	}
	return ret, nil
}

func (lc LightCertifiedBlob) MarshalBinary() ([]byte, error) {
	return lc.AppendBinary(nil)
}

func (lc *LightCertifiedBlob) UnmarshalBinary(data []byte) error {
	pbLightCertifiedBlob := serialization.LightCertifiedBlob{}
	if err := proto.Unmarshal(data, &pbLightCertifiedBlob); err != nil {
		return fmt.Errorf("failed to unmarshal LightCertifiedBlob protobuf: %w", err)
	}

	chunkDigests := make([]BlobChunkDigest, 0, len(pbLightCertifiedBlob.ChunkDigests))
	for _, digest := range pbLightCertifiedBlob.ChunkDigests {
		if len(digest) != len(BlobChunkDigest{}) {
			return fmt.Errorf("invalid chunk digest length: expected %d bytes, got %d", len(BlobChunkDigest{}), len(digest))
		}
		var chunkDigest BlobChunkDigest
		copy(chunkDigest[:], digest)
		chunkDigests = append(chunkDigests, chunkDigest)
	}

	signatures := make([]AttributedBlobAvailabilitySignature, 0, len(pbLightCertifiedBlob.AttributedBlobAvailabilitySignatures))
	for _, sig := range pbLightCertifiedBlob.AttributedBlobAvailabilitySignatures {
		if len(sig.Signature) != ed25519.SignatureSize {
			return fmt.Errorf("invalid signature length: expected %d bytes, got %d", ed25519.SignatureSize, len(sig.Signature))
		}
		signatures = append(signatures, AttributedBlobAvailabilitySignature{
			sig.Signature,
			commontypes.OracleID(sig.Signer),
		})
	}

	*lc = LightCertifiedBlob{
		chunkDigests,
		pbLightCertifiedBlob.PayloadLength,
		pbLightCertifiedBlob.ExpirySeqNr,
		commontypes.OracleID(pbLightCertifiedBlob.Submitter),
		signatures,
	}
	return nil
}
