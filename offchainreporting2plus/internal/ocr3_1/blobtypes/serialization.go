package blobtypes

import (
	"crypto/ed25519"
	"encoding"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/mt"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes/serialization"
	"google.golang.org/protobuf/proto"
)

var _ encoding.BinaryMarshaler = LightCertifiedBlob{}
var _ encoding.BinaryAppender = LightCertifiedBlob{}
var _ encoding.BinaryUnmarshaler = &LightCertifiedBlob{}

func (lc LightCertifiedBlob) AppendBinary(b []byte) ([]byte, error) {
	pbSignatures := make([]*serialization.AttributedBlobAvailabilitySignature, 0, len(lc.AttributedBlobAvailabilitySignatures))
	for _, sig := range lc.AttributedBlobAvailabilitySignatures {
		pbSignatures = append(pbSignatures, serialization.NewAttributedBlobAvailabilitySignature(
			sig.Signature,
			uint32(sig.Signer),
		))
	}

	pbLightCertifiedBlob := serialization.NewLightCertifiedBlob(
		lc.ChunkDigestsRoot[:],
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

	var chunkDigestsRoot mt.Digest
	if len(pbLightCertifiedBlob.ChunkDigestsRoot) != len(mt.Digest{}) {
		return fmt.Errorf("invalid chunk digests root length: expected %d bytes, got %d", len(mt.Digest{}), len(pbLightCertifiedBlob.ChunkDigestsRoot))
	}
	copy(chunkDigestsRoot[:], pbLightCertifiedBlob.ChunkDigestsRoot)

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
		chunkDigestsRoot,
		pbLightCertifiedBlob.PayloadLength,
		pbLightCertifiedBlob.ExpirySeqNr,
		commontypes.OracleID(pbLightCertifiedBlob.Submitter),
		signatures,
	}
	return nil
}
