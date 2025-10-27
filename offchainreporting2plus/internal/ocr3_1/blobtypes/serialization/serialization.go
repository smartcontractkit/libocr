package serialization

import (
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

func NewAttributedBlobAvailabilitySignature(
	signature []byte,
	signer uint32,
) *AttributedBlobAvailabilitySignature {
	return &AttributedBlobAvailabilitySignature{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		signature,
		signer,
	}
}

func NewLightCertifiedBlob(
	chunkDigestsRoot []byte,
	payloadLength uint64,
	expirySeqNr uint64,
	submitter uint32,
	attributedBlobAvailabilitySignatures []*AttributedBlobAvailabilitySignature,
) *LightCertifiedBlob {
	return &LightCertifiedBlob{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		chunkDigestsRoot,
		payloadLength,
		expirySeqNr,
		submitter,
		attributedBlobAvailabilitySignatures,
	}
}
