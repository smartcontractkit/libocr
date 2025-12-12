package serialization

import (
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/jmt"
	"github.com/smartcontractkit/libocr/internal/mt"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"google.golang.org/protobuf/proto"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

// Serialize encodes a protocol.Message into a binary payload
func Serialize[RI any](m protocol.Message[RI]) ([]byte, *MessageWrapper, error) {
	tpm := toProtoMessage[RI]{}
	pb, err := tpm.messageWrapper(m)
	if err != nil {
		return nil, nil, err
	}
	b, err := proto.Marshal(pb)
	if err != nil {
		return nil, nil, err
	}
	return b, pb, nil
}

func SerializeCertifiedPrepareOrCommit(cpoc protocol.CertifiedPrepareOrCommit) ([]byte, error) {
	if cpoc == nil {
		return nil, fmt.Errorf("cannot serialize nil CertifiedPrepareOrCommit")
	}

	tpm := toProtoMessage[struct{}]{}

	return proto.Marshal(tpm.certifiedPrepareOrCommit(cpoc))
}

func SerializeTreeSyncStatus(m protocol.TreeSyncStatus) ([]byte, error) {
	tpm := toProtoMessage[struct{}]{}
	pb := TreeSyncStatus{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		TreeSyncPhase(m.Phase),
		m.TargetSeqNr,
		m.TargetStateRootDigest[:],
		tpm.pendingKeyDigestRanges(m.PendingKeyDigestRanges),
	}
	return proto.Marshal(&pb)
}

func SerializePacemakerState(m protocol.PacemakerState) ([]byte, error) {
	pb := PacemakerState{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		m.Epoch,
		m.HighestSentNewEpochWish,
	}

	return proto.Marshal(&pb)
}

func SerializeBlobMeta(m protocol.BlobMeta) ([]byte, error) {
	tpm := toProtoMessage[struct{}]{}
	pb := BlobMeta{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		m.PayloadLength,
		m.ChunkHaves,
		tpm.chunkDigests(m.ChunkDigests),
		m.ExpirySeqNr,
		uint32(m.Submitter),
	}
	return proto.Marshal(&pb)
}

func SerializeBlobQuotaStats(bqs protocol.BlobQuotaStats) ([]byte, error) {
	pb := BlobQuotaStats{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		bqs.Count,
		bqs.CumulativePayloadLength,
	}
	return proto.Marshal(&pb)
}

func SerializeAttestedStateTransitionBlock(astb protocol.AttestedStateTransitionBlock) ([]byte, error) {
	tpm := toProtoMessage[struct{}]{}

	return proto.Marshal(tpm.attestedStateTransitionBlock(astb))
}

func SerializeStateTransitionBlock(stb protocol.StateTransitionBlock) ([]byte, error) {
	tpm := toProtoMessage[struct{}]{}

	return proto.Marshal(tpm.stateTransitionBlock(stb))
}

func SerializeGenesisStateTransitionBlock(gstb protocol.GenesisStateTransitionBlock) ([]byte, error) {
	tpm := toProtoMessage[struct{}]{}

	return proto.Marshal(tpm.genesisStateTransitionBlock(&gstb))
}

// Deserialize decodes a binary payload into a protocol.Message
func Deserialize[RI any](n int, b []byte, requestHandle types.RequestHandle) (protocol.Message[RI], *MessageWrapper, error) {
	pb := &MessageWrapper{}
	if err := proto.Unmarshal(b, pb); err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal protobuf: %w", err)
	}

	fpm := fromProtoMessage[RI]{n, requestHandle}
	m, err := fpm.messageWrapper(pb)
	if err != nil {
		return nil, nil, fmt.Errorf("could not translate protobuf to protocol.Message: %w", err)
	}
	return m, pb, nil
}

func DeserializeTrustedPrepareOrCommit(b []byte) (protocol.CertifiedPrepareOrCommit, error) {
	pb := CertifiedPrepareOrCommit{}
	if err := proto.Unmarshal(b, &pb); err != nil {
		return nil, err
	}

	// We trust the PrepareOrCommit we deserialize here, so we can simply use the maximum number
	// of oracles for n.
	n := types.MaxOracles
	// we intentionally leave this as nil since prepare and commits don't carry a request handle
	var requestHandle types.RequestHandle
	fpm := fromProtoMessage[struct{}]{n, requestHandle}
	return fpm.certifiedPrepareOrCommit(&pb)
}

func DeserializeTreeSyncStatus(b []byte) (protocol.TreeSyncStatus, error) {
	pb := TreeSyncStatus{}
	if err := proto.Unmarshal(b, &pb); err != nil {
		return protocol.TreeSyncStatus{}, err
	}
	var stateRootDigest protocol.StateRootDigest
	if len(stateRootDigest) < len(pb.TargetStateRootDigest) {
		return protocol.TreeSyncStatus{}, fmt.Errorf("invalid state root digest length expected at most %d, got %d", len(stateRootDigest), len(pb.TargetStateRootDigest))
	}
	copy(stateRootDigest[:], pb.TargetStateRootDigest)
	fpm := fromProtoMessage[struct{}]{}
	pkdr, err := fpm.pendingKeyDigestRanges(pb.PendingKeyDigestRanges)
	if err != nil {
		return protocol.TreeSyncStatus{}, err
	}
	return protocol.TreeSyncStatus{
		protocol.TreeSyncPhase(pb.Phase),
		pb.TargetSeqNr,
		stateRootDigest,
		pkdr,
	}, nil
}

func DeserializePacemakerState(b []byte) (protocol.PacemakerState, error) {
	pb := PacemakerState{}
	if err := proto.Unmarshal(b, &pb); err != nil {
		return protocol.PacemakerState{}, err
	}

	return protocol.PacemakerState{
		pb.Epoch,
		pb.HighestSentNewEpochWish,
	}, nil
}

func DeserializeBlobMeta(b []byte) (protocol.BlobMeta, error) {
	pb := BlobMeta{}
	if err := proto.Unmarshal(b, &pb); err != nil {
		return protocol.BlobMeta{}, err
	}
	fpm := fromProtoMessage[struct{}]{}
	chunkDigests, err := fpm.chunkDigests(pb.ChunkDigests)
	if err != nil {
		return protocol.BlobMeta{}, err
	}
	return protocol.BlobMeta{
		pb.PayloadLength,
		pb.ChunkHaves,
		chunkDigests,
		pb.ExpirySeqNr,
		commontypes.OracleID(pb.Submitter),
	}, nil
}

func DeserializeBlobQuotaStats(b []byte) (protocol.BlobQuotaStats, error) {
	pb := BlobQuotaStats{}
	if err := proto.Unmarshal(b, &pb); err != nil {
		return protocol.BlobQuotaStats{}, err
	}
	return protocol.BlobQuotaStats{
		pb.Count,
		pb.CumulativePayloadLength,
	}, nil
}

func DeserializeAttestedStateTransitionBlock(b []byte) (protocol.AttestedStateTransitionBlock, error) {
	pb := AttestedStateTransitionBlock{}
	if err := proto.Unmarshal(b, &pb); err != nil {
		return protocol.AttestedStateTransitionBlock{}, err
	}

	n := types.MaxOracles
	// we intentionally leave this as nil since attested state transition blocks don't carry a request handle
	var requestHandle types.RequestHandle
	fpm := fromProtoMessage[struct{}]{n, requestHandle}
	return fpm.attestedStateTransitionBlock(&pb)
}

func DeserializeStateTransitionBlock(b []byte) (protocol.StateTransitionBlock, error) {
	pb := StateTransitionBlock{}
	if err := proto.Unmarshal(b, &pb); err != nil {
		return protocol.StateTransitionBlock{}, err
	}

	fpm := fromProtoMessage[struct{}]{}
	return fpm.stateTransitionBlock(&pb)
}

func DeserializeGenesisStateTransitionBlock(b []byte) (protocol.GenesisStateTransitionBlock, error) {
	pb := GenesisStateTransitionBlock{}
	if err := proto.Unmarshal(b, &pb); err != nil {
		return protocol.GenesisStateTransitionBlock{}, err
	}

	fpm := fromProtoMessage[struct{}]{}
	gstb, err := fpm.genesisStateTransitionBlock(&pb)
	if err != nil {
		return protocol.GenesisStateTransitionBlock{}, err
	}
	if gstb == nil {
		return protocol.GenesisStateTransitionBlock{}, fmt.Errorf("unexpectedly got nil genesis state transition block")
	}
	return *gstb, nil
}

//
// *toProtoMessage
//

type toProtoMessage[RI any] struct{}

func (tpm *toProtoMessage[RI]) messageWrapper(m protocol.Message[RI]) (*MessageWrapper, error) {
	msgWrapper := MessageWrapper{}
	switch v := m.(type) {
	case protocol.MessageNewEpochWish[RI]:
		pm := &MessageNewEpochWish{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
		}
		msgWrapper.Msg = &MessageWrapper_MessageNewEpochWish{pm}
	case protocol.MessageEpochStartRequest[RI]:
		pm := &MessageEpochStartRequest{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			tpm.certifiedPrepareOrCommit(v.HighestCertified),
			tpm.signedHighestCertifiedTimestamp(v.SignedHighestCertifiedTimestamp),
		}
		msgWrapper.Msg = &MessageWrapper_MessageEpochStartRequest{pm}
	case protocol.MessageEpochStart[RI]:
		pm := &MessageEpochStart{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			tpm.epochStartProof(v.EpochStartProof),
			v.Abdicate,
		}
		msgWrapper.Msg = &MessageWrapper_MessageEpochStart{pm}
	case protocol.MessageRoundStart[RI]:
		pm := &MessageRoundStart{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			v.Query,
		}
		msgWrapper.Msg = &MessageWrapper_MessageRoundStart{pm}
	case protocol.MessageObservation[RI]:
		pm := &MessageObservation{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			tpm.signedObservation(v.SignedObservation),
		}
		msgWrapper.Msg = &MessageWrapper_MessageObservation{pm}

	case protocol.MessageProposal[RI]:
		pbasos := make([]*AttributedSignedObservation, 0, len(v.AttributedSignedObservations))
		for _, aso := range v.AttributedSignedObservations {
			pbasos = append(pbasos, tpm.attributedSignedObservation(aso))
		}
		pm := &MessageProposal{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			pbasos,
		}
		msgWrapper.Msg = &MessageWrapper_MessageProposal{pm}
	case protocol.MessagePrepare[RI]:
		pm := &MessagePrepare{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			v.Signature,
		}
		msgWrapper.Msg = &MessageWrapper_MessagePrepare{pm}
	case protocol.MessageCommit[RI]:
		pm := &MessageCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			v.Signature,
		}
		msgWrapper.Msg = &MessageWrapper_MessageCommit{pm}
	case protocol.MessageReportSignatures[RI]:
		pm := &MessageReportSignatures{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.SeqNr,
			v.ReportSignatures,
			v.ReportsPlusPrecursorDigest[:],
		}
		msgWrapper.Msg = &MessageWrapper_MessageReportSignatures{pm}
	case protocol.MessageReportsPlusPrecursorRequest[RI]:
		pm := &MessageReportsPlusPrecursorRequest{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.SeqNr,
		}
		msgWrapper.Msg = &MessageWrapper_MessageReportsPlusPrecursorRequest{pm}
	case protocol.MessageReportsPlusPrecursor[RI]:
		pm := &MessageReportsPlusPrecursor{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.SeqNr,
			v.ReportsPlusPrecursor,
		}
		msgWrapper.Msg = &MessageWrapper_MessageReportsPlusPrecursor{pm}
	case protocol.MessageBlockSyncRequest[RI]:
		pm := &MessageBlockSyncRequest{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.StartSeqNr,
			v.EndExclSeqNr,
			uint32(v.MaxCumulativeWriteSetBytes),
		}
		msgWrapper.Msg = &MessageWrapper_MessageBlockSyncRequest{pm}
	case protocol.MessageBlockSyncResponse[RI]:
		astbs := make([]*AttestedStateTransitionBlock, 0, len(v.AttestedStateTransitionBlocks))
		for _, astb := range v.AttestedStateTransitionBlocks {
			astbs = append(astbs, tpm.attestedStateTransitionBlock(astb))
		}
		pm := &MessageBlockSyncResponse{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.RequestStartSeqNr,
			v.RequestEndExclSeqNr,
			astbs,
			tpm.genesisStateTransitionBlock(v.GenesisStateTransitionBlock),
		}
		msgWrapper.Msg = &MessageWrapper_MessageBlockSyncResponse{pm}
	case protocol.MessageStateSyncSummary[RI]:
		pm := &MessageStateSyncSummary{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.LowestPersistedSeqNr,
			v.HighestCommittedSeqNr,
		}
		msgWrapper.Msg = &MessageWrapper_MessageStateSyncSummary{pm}
	case protocol.MessageTreeSyncChunkRequest[RI]:
		pm := &MessageTreeSyncChunkRequest{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.ToSeqNr,
			v.StartIndex[:],
			v.EndInclIndex[:],
			uint32(v.MaxCumulativeKeysPlusValuesBytes),
		}
		msgWrapper.Msg = &MessageWrapper_MessageTreeSyncChunkRequest{pm}
	case protocol.MessageTreeSyncChunkResponse[RI]:
		pm := &MessageTreeSyncChunkResponse{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.ToSeqNr,
			v.StartIndex[:],
			v.RequestEndInclIndex[:],
			v.GoAway,
			v.EndInclIndex[:],
			tpm.treeSyncChunkKeyValues(v.KeyValues),
			tpm.treeSyncChunkBoundingLeaves(v.BoundingLeaves),
		}
		msgWrapper.Msg = &MessageWrapper_MessageTreeSyncChunkResponse{pm}
	case protocol.MessageBlobOffer[RI]:
		pm := &MessageBlobOffer{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.ChunkDigestsRoot[:],
			v.PayloadLength,
			v.ExpirySeqNr,
		}
		msgWrapper.Msg = &MessageWrapper_MessageBlobOffer{pm}
	case protocol.MessageBlobOfferResponse[RI]:
		pm := &MessageBlobOfferResponse{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.BlobDigest[:],
			v.RejectOffer,
			v.Signature[:],
		}
		msgWrapper.Msg = &MessageWrapper_MessageBlobOfferResponse{pm}
	case protocol.MessageBlobChunkRequest[RI]:
		pm := &MessageBlobChunkRequest{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.BlobDigest[:],
			v.ChunkIndex,
		}
		msgWrapper.Msg = &MessageWrapper_MessageBlobChunkRequest{pm}
	case protocol.MessageBlobChunkResponse[RI]:
		pbProof := make([][]byte, len(v.Proof))
		for i, p := range v.Proof {
			pbProof[i] = p[:]
		}
		pm := &MessageBlobChunkResponse{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.BlobDigest[:],
			v.ChunkIndex,
			v.GoAway,
			v.Chunk,
			pbProof,
		}
		msgWrapper.Msg = &MessageWrapper_MessageBlobChunkResponse{pm}
	default:
		return nil, fmt.Errorf("unable to serialize message of type %T", m)

	}
	return &msgWrapper, nil
}

func (tpm *toProtoMessage[RI]) chunkDigests(chunkDigests []protocol.BlobChunkDigest) [][]byte {
	cds := make([][]byte, 0, len(chunkDigests))
	for _, chunkDigest := range chunkDigests {
		cds = append(cds, chunkDigest[:])
	}
	return cds
}

func (tpm *toProtoMessage[RI]) certifiedPrepareOrCommit(cpoc protocol.CertifiedPrepareOrCommit) *CertifiedPrepareOrCommit {
	switch v := cpoc.(type) {
	case *protocol.CertifiedPrepare:
		prepareQuorumCertificate := make([]*AttributedPrepareSignature, 0, len(v.PrepareQuorumCertificate))
		for _, aps := range v.PrepareQuorumCertificate {
			prepareQuorumCertificate = append(prepareQuorumCertificate, &AttributedPrepareSignature{
				// zero-initialize protobuf built-ins
				protoimpl.MessageState{},
				0,
				nil,
				// fields
				aps.Signature,
				uint32(aps.Signer),
			})
		}
		return &CertifiedPrepareOrCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			&CertifiedPrepareOrCommit_Prepare{&CertifiedPrepare{
				// zero-initialize protobuf built-ins
				protoimpl.MessageState{},
				0,
				nil,
				// fields
				v.PrevHistoryDigest[:],
				v.Epoch(),
				v.SeqNr(),
				v.StateTransitionInputsDigest[:],
				v.StateWriteSetDigest[:],
				v.StateRootDigest[:],
				v.ReportsPlusPrecursorDigest[:],
				prepareQuorumCertificate,
			}},
		}
	case *protocol.CertifiedCommit:
		return &CertifiedPrepareOrCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			&CertifiedPrepareOrCommit_Commit{tpm.CertifiedCommit(*v)},
		}
	case *protocol.GenesisFromScratch:
		return &CertifiedPrepareOrCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			&CertifiedPrepareOrCommit_GenesisFromScratch{&GenesisFromScratch{}},
		}
	case *protocol.GenesisFromPrevInstance:
		return &CertifiedPrepareOrCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			&CertifiedPrepareOrCommit_GenesisFromPrevInstance{
				&GenesisFromPrevInstance{
					// zero-initialize protobuf built-ins
					protoimpl.MessageState{},
					0,
					nil,
					// fields
					v.PrevHistoryDigest[:],
					v.PrevSeqNr,
				}},
		}
	default:
		// It's safe to crash here since the "protocol.*" versions of these values
		// come from the trusted, local environment.
		panic("unrecognized protocol.CertifiedPrepareOrCommit implementation")
	}
}

func (tpm *toProtoMessage[RI]) CertifiedCommit(cpocc protocol.CertifiedCommit) *CertifiedCommit {
	commitQuorumCertificate := make([]*AttributedCommitSignature, 0, len(cpocc.CommitQuorumCertificate))
	for _, aps := range cpocc.CommitQuorumCertificate {
		commitQuorumCertificate = append(commitQuorumCertificate, &AttributedCommitSignature{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			aps.Signature,
			uint32(aps.Signer),
		})
	}
	return &CertifiedCommit{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		cpocc.PrevHistoryDigest[:],
		cpocc.Epoch(),
		cpocc.SeqNr(),
		cpocc.StateTransitionInputsDigest[:],
		cpocc.StateWriteSetDigest[:],
		cpocc.StateRootDigest[:],
		cpocc.ReportsPlusPrecursorDigest[:],
		commitQuorumCertificate,
	}
}

func (tpm *toProtoMessage[RI]) attributedSignedHighestCertifiedTimestamp(ashct protocol.AttributedSignedHighestCertifiedTimestamp) *AttributedSignedHighestCertifiedTimestamp {
	return &AttributedSignedHighestCertifiedTimestamp{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		tpm.signedHighestCertifiedTimestamp(ashct.SignedHighestCertifiedTimestamp),
		uint32(ashct.Signer),
	}
}

func (tpm *toProtoMessage[RI]) signedHighestCertifiedTimestamp(shct protocol.SignedHighestCertifiedTimestamp) *SignedHighestCertifiedTimestamp {
	return &SignedHighestCertifiedTimestamp{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		tpm.highestCertifiedTimestamp(shct.HighestCertifiedTimestamp),
		shct.Signature,
	}
}

func (tpm *toProtoMessage[RI]) highestCertifiedTimestamp(hct protocol.HighestCertifiedTimestamp) *HighestCertifiedTimestamp {
	return &HighestCertifiedTimestamp{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		hct.SeqNr,
		hct.CommittedElsePrepared,
		hct.Epoch,
	}
}

func (tpm *toProtoMessage[RI]) epochStartProof(srqc protocol.EpochStartProof) *EpochStartProof {
	highestCertifiedProof := make([]*AttributedSignedHighestCertifiedTimestamp, 0, len(srqc.HighestCertifiedProof))
	for _, ashct := range srqc.HighestCertifiedProof {
		highestCertifiedProof = append(highestCertifiedProof, tpm.attributedSignedHighestCertifiedTimestamp(ashct))
	}
	return &EpochStartProof{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		tpm.certifiedPrepareOrCommit(srqc.HighestCertified),
		highestCertifiedProof,
	}
}

func (tpm *toProtoMessage[RI]) signedObservation(o protocol.SignedObservation) *SignedObservation {
	return &SignedObservation{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		o.Observation,
		o.Signature,
	}
}

func (tpm *toProtoMessage[RI]) attributedSignedObservation(aso protocol.AttributedSignedObservation) *AttributedSignedObservation {
	return &AttributedSignedObservation{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		tpm.signedObservation(aso.SignedObservation),
		uint32(aso.Observer),
	}
}

func (tpm *toProtoMessage[RI]) attestedStateTransitionBlock(astb protocol.AttestedStateTransitionBlock) *AttestedStateTransitionBlock {
	attributedSignatures := make([]*AttributedCommitSignature, 0, len(astb.AttributedCommitSignatures))
	for _, as := range astb.AttributedCommitSignatures {
		attributedSignatures = append(attributedSignatures, &AttributedCommitSignature{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			as.Signature,
			uint32(as.Signer),
		})
	}
	return &AttestedStateTransitionBlock{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		tpm.stateTransitionBlock(astb.StateTransitionBlock),
		attributedSignatures,
	}
}

func (tpm *toProtoMessage[RI]) stateTransitionBlock(stb protocol.StateTransitionBlock) *StateTransitionBlock {
	return &StateTransitionBlock{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		stb.PrevHistoryDigest[:],
		stb.Epoch,
		stb.SeqNr(),
		stb.StateTransitionInputsDigest[:],
		tpm.stateWriteSet(stb.StateWriteSet),
		stb.StateRootDigest[:],
		stb.ReportsPlusPrecursorDigest[:],
	}
}

func (tpm *toProtoMessage[RI]) genesisStateTransitionBlock(gstb *protocol.GenesisStateTransitionBlock) *GenesisStateTransitionBlock {
	if gstb == nil {
		return nil
	}
	return &GenesisStateTransitionBlock{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		gstb.PrevHistoryDigest[:],
		gstb.SeqNr,
		gstb.StateTransitionInputsDigest[:],
		gstb.StateWriteSetDigest[:],
		gstb.StateRootDigest[:],
		gstb.ReportsPlusPrecursorDigest[:],
	}
}

func (tpm *toProtoMessage[RI]) stateWriteSet(sto protocol.StateWriteSet) *StateWriteSet {
	pbWriteSet := make([]*KeyValueModification, 0, len(sto.Entries))
	for _, kvmod := range sto.Entries {
		pbWriteSet = append(pbWriteSet, &KeyValueModification{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			kvmod.Key,
			kvmod.Value,
			kvmod.Deleted,
		})
	}
	return &StateWriteSet{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		pbWriteSet,
	}
}

func (tpm *toProtoMessage[RI]) treeSyncChunkProof(proof []jmt.Digest) [][]byte {
	pns := make([][]byte, 0, len(proof))
	for _, pn := range proof {
		pns = append(pns, pn[:])
	}
	return pns
}

func (tpm *toProtoMessage[RI]) treeSyncChunkBoundingLeaves(boundingLeaves []jmt.BoundingLeaf) []*BoundingLeaf {
	pbbls := make([]*BoundingLeaf, 0, len(boundingLeaves))
	for _, bl := range boundingLeaves {
		pbbls = append(pbbls, &BoundingLeaf{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			tpm.leafKeyAndValueDigests(bl.Leaf),
			tpm.treeSyncChunkProof(bl.Siblings),
		})
	}
	return pbbls
}

func (tpm *toProtoMessage[RI]) leafKeyAndValueDigests(leafKeyAndValueDigests jmt.LeafKeyAndValueDigests) *LeafKeyAndValueDigests {
	return &LeafKeyAndValueDigests{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		leafKeyAndValueDigests.KeyDigest[:],
		leafKeyAndValueDigests.ValueDigest[:],
	}
}

func (tpm *toProtoMessage[RI]) treeSyncChunkKeyValues(kvps []protocol.KeyValuePair) []*KeyValuePair {
	pbkvps := make([]*KeyValuePair, 0, len(kvps))
	for _, kvp := range kvps {
		pbkvps = append(pbkvps, &KeyValuePair{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			kvp.Key,
			kvp.Value,
		})
	}
	return pbkvps
}

func (tpm *toProtoMessage[RI]) pendingKeyDigestRanges(pkdr protocol.PendingKeyDigestRanges) []*KeyDigestRange {
	allRanges := pkdr.All()
	pbRanges := make([]*KeyDigestRange, 0, len(allRanges))
	for _, r := range allRanges {
		pbRanges = append(pbRanges, &KeyDigestRange{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			r.StartIndex[:],
			r.EndInclIndex[:],
		})
	}
	return pbRanges
}

//
// *fromProtoMessage
//

type fromProtoMessage[RI any] struct {
	n             int
	requestHandle types.RequestHandle
}

func (fpm *fromProtoMessage[RI]) messageWrapper(wrapper *MessageWrapper) (protocol.Message[RI], error) {
	switch msg := wrapper.Msg.(type) {
	case *MessageWrapper_MessageNewEpochWish:
		return fpm.messageNewEpochWish(wrapper.GetMessageNewEpochWish())
	case *MessageWrapper_MessageEpochStartRequest:
		return fpm.messageEpochStartRequest(wrapper.GetMessageEpochStartRequest())
	case *MessageWrapper_MessageEpochStart:
		return fpm.messageEpochStart(wrapper.GetMessageEpochStart())
	case *MessageWrapper_MessageRoundStart:
		return fpm.messageRoundStart(wrapper.GetMessageRoundStart())
	case *MessageWrapper_MessageObservation:
		return fpm.messageObservation(wrapper.GetMessageObservation())
	case *MessageWrapper_MessageProposal:
		return fpm.messageProposal(wrapper.GetMessageProposal())
	case *MessageWrapper_MessagePrepare:
		return fpm.messagePrepare(wrapper.GetMessagePrepare())
	case *MessageWrapper_MessageCommit:
		return fpm.messageCommit(wrapper.GetMessageCommit())
	case *MessageWrapper_MessageReportSignatures:
		return fpm.messageReportSignatures(wrapper.GetMessageReportSignatures())
	case *MessageWrapper_MessageReportsPlusPrecursorRequest:
		return fpm.messageReportsPlusPrecursorRequest(wrapper.GetMessageReportsPlusPrecursorRequest())
	case *MessageWrapper_MessageReportsPlusPrecursor:
		return fpm.messageReportsPlusPrecursor(wrapper.GetMessageReportsPlusPrecursor())
	case *MessageWrapper_MessageBlockSyncRequest:
		return fpm.messageBlockSyncRequest(wrapper.GetMessageBlockSyncRequest())
	case *MessageWrapper_MessageBlockSyncResponse:
		return fpm.messageBlockSyncResponse(wrapper.GetMessageBlockSyncResponse())
	case *MessageWrapper_MessageStateSyncSummary:
		return fpm.messageStateSyncSummary(wrapper.GetMessageStateSyncSummary())
	case *MessageWrapper_MessageTreeSyncChunkRequest:
		return fpm.messageTreeSyncChunkRequest(wrapper.GetMessageTreeSyncChunkRequest())
	case *MessageWrapper_MessageTreeSyncChunkResponse:
		return fpm.messageTreeSyncChunkResponse(wrapper.GetMessageTreeSyncChunkResponse())
	case *MessageWrapper_MessageBlobOffer:
		return fpm.messageBlobOffer(wrapper.GetMessageBlobOffer())
	case *MessageWrapper_MessageBlobOfferResponse:
		return fpm.messageBlobOfferResponse(wrapper.GetMessageBlobOfferResponse())
	case *MessageWrapper_MessageBlobChunkRequest:
		return fpm.messageBlobChunkRequest(wrapper.GetMessageBlobChunkRequest())
	case *MessageWrapper_MessageBlobChunkResponse:
		return fpm.messageBlobChunkResponse(wrapper.GetMessageBlobChunkResponse())

	default:
		return nil, fmt.Errorf("unrecognized Msg type %T", msg)
	}
}

func (fpm *fromProtoMessage[RI]) messageNewEpochWish(m *MessageNewEpochWish) (protocol.MessageNewEpochWish[RI], error) {
	if m == nil {
		return protocol.MessageNewEpochWish[RI]{}, fmt.Errorf("unable to extract a MessageNewEpochWish value")
	}
	return protocol.MessageNewEpochWish[RI]{
		m.Epoch,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageEpochStartRequest(m *MessageEpochStartRequest) (protocol.MessageEpochStartRequest[RI], error) {
	if m == nil {
		return protocol.MessageEpochStartRequest[RI]{}, fmt.Errorf("unable to extract a MessageEpochStartRequest value")
	}
	hc, err := fpm.certifiedPrepareOrCommit(m.HighestCertified)
	if err != nil {
		return protocol.MessageEpochStartRequest[RI]{}, err
	}
	shct, err := fpm.signedHighestCertifiedTimestamp(m.SignedHighestCertifiedTimestamp)
	if err != nil {
		return protocol.MessageEpochStartRequest[RI]{}, err
	}
	return protocol.MessageEpochStartRequest[RI]{
		m.Epoch,
		hc,
		shct,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageEpochStart(m *MessageEpochStart) (protocol.MessageEpochStart[RI], error) {
	if m == nil {
		return protocol.MessageEpochStart[RI]{}, fmt.Errorf("unable to extract a MessageEpochStart value")
	}
	srqc, err := fpm.epochStartProof(m.EpochStartProof)
	if err != nil {
		return protocol.MessageEpochStart[RI]{}, err
	}
	return protocol.MessageEpochStart[RI]{
		m.Epoch,
		srqc,
		m.Abdicate,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageProposal(m *MessageProposal) (protocol.MessageProposal[RI], error) {
	if m == nil {
		return protocol.MessageProposal[RI]{}, fmt.Errorf("unable to extract a MessageProposal value")
	}
	asos, err := fpm.attributedSignedObservations(m.AttributedSignedObservations)
	if err != nil {
		return protocol.MessageProposal[RI]{}, err
	}
	return protocol.MessageProposal[RI]{
		m.Epoch,
		m.SeqNr,
		asos,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messagePrepare(m *MessagePrepare) (protocol.MessagePrepare[RI], error) {
	if m == nil {
		return protocol.MessagePrepare[RI]{}, fmt.Errorf("unable to extract a MessagePrepare value")
	}
	return protocol.MessagePrepare[RI]{
		m.Epoch,
		m.SeqNr,
		m.Signature,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageCommit(m *MessageCommit) (protocol.MessageCommit[RI], error) {
	if m == nil {
		return protocol.MessageCommit[RI]{}, fmt.Errorf("unable to extract a MessageCommit value")
	}
	return protocol.MessageCommit[RI]{
		m.Epoch,
		m.SeqNr,
		m.Signature,
	}, nil
}

func (fpm *fromProtoMessage[RI]) certifiedPrepareOrCommit(m *CertifiedPrepareOrCommit) (protocol.CertifiedPrepareOrCommit, error) {
	if m == nil {
		return nil, fmt.Errorf("unable to extract a CertifiedPrepareOrCommit value")
	}
	switch poc := m.PrepareOrCommit.(type) {
	case *CertifiedPrepareOrCommit_Prepare:
		cpocp, err := fpm.certifiedPrepare(poc.Prepare)
		if err != nil {
			return nil, err
		}
		return &cpocp, nil
	case *CertifiedPrepareOrCommit_Commit:
		cpocc, err := fpm.certifiedCommit(poc.Commit)
		if err != nil {
			return nil, err
		}
		return &cpocc, nil
	case *CertifiedPrepareOrCommit_GenesisFromScratch:
		return &protocol.GenesisFromScratch{}, nil
	case *CertifiedPrepareOrCommit_GenesisFromPrevInstance:
		gfpi, err := fpm.genesisFromPrevInstance(poc.GenesisFromPrevInstance)
		if err != nil {
			return nil, err
		}
		return &gfpi, nil
	default:
		return nil, fmt.Errorf("unknown case of CertifiedPrepareOrCommit")
	}
}

func (fpm *fromProtoMessage[RI]) genesisFromPrevInstance(m *GenesisFromPrevInstance) (protocol.GenesisFromPrevInstance, error) {
	if m == nil {
		return protocol.GenesisFromPrevInstance{}, fmt.Errorf("unable to extract a GenesisFromPrevInstance value")
	}
	var prevHistoryDigest protocol.HistoryDigest
	copy(prevHistoryDigest[:], m.PrevHistoryDigest)
	return protocol.GenesisFromPrevInstance{
		prevHistoryDigest,
		m.PrevSeqNr,
	}, nil
}

func (fpm *fromProtoMessage[RI]) certifiedPrepare(m *CertifiedPrepare) (protocol.CertifiedPrepare, error) {
	if m == nil {
		return protocol.CertifiedPrepare{}, fmt.Errorf("unable to extract a CertifiedPrepare value")
	}
	var prevHistoryDigest protocol.HistoryDigest
	copy(prevHistoryDigest[:], m.PrevHistoryDigest)
	var inputsDigest protocol.StateTransitionInputsDigest
	copy(inputsDigest[:], m.StateTransitionInputsDigest)
	var writeSetDigest protocol.StateWriteSetDigest
	copy(writeSetDigest[:], m.StateWriteSetDigest)
	var stateRootDigest protocol.StateRootDigest
	copy(stateRootDigest[:], m.StateRootDigest)
	var reportsPlusPrecursorDigest protocol.ReportsPlusPrecursorDigest
	copy(reportsPlusPrecursorDigest[:], m.ReportsPlusPrecursorDigest)

	prepareQuorumCertificate := make([]protocol.AttributedPrepareSignature, 0, len(m.PrepareQuorumCertificate))
	for _, aps := range m.PrepareQuorumCertificate {
		signer, err := fpm.oracleID(aps.GetSigner())
		if err != nil {
			return protocol.CertifiedPrepare{}, err
		}
		prepareQuorumCertificate = append(prepareQuorumCertificate, protocol.AttributedPrepareSignature{
			aps.GetSignature(),
			signer,
		})
	}
	return protocol.CertifiedPrepare{
		prevHistoryDigest,
		m.Epoch,
		m.SeqNr,
		inputsDigest,
		writeSetDigest,
		stateRootDigest,
		reportsPlusPrecursorDigest,
		prepareQuorumCertificate,
	}, nil

}

func (fpm *fromProtoMessage[RI]) certifiedCommit(m *CertifiedCommit) (protocol.CertifiedCommit, error) {
	if m == nil {
		return protocol.CertifiedCommit{}, fmt.Errorf("unable to extract a CertifiedCommit value")
	}
	var prevHistoryDigest protocol.HistoryDigest
	copy(prevHistoryDigest[:], m.PrevHistoryDigest)
	var inputsDigest protocol.StateTransitionInputsDigest
	copy(inputsDigest[:], m.StateTransitionInputsDigest)
	var writeSetDigest protocol.StateWriteSetDigest
	copy(writeSetDigest[:], m.StateWriteSetDigest)
	var stateRootDigest protocol.StateRootDigest
	copy(stateRootDigest[:], m.StateRootDigest)
	var reportsPlusPrecursorDigest protocol.ReportsPlusPrecursorDigest
	copy(reportsPlusPrecursorDigest[:], m.ReportsPlusPrecursorDigest)

	commitQuorumCertificate := make([]protocol.AttributedCommitSignature, 0, len(m.CommitQuorumCertificate))
	for _, aps := range m.CommitQuorumCertificate {
		signer, err := fpm.oracleID(aps.GetSigner())
		if err != nil {
			return protocol.CertifiedCommit{}, err
		}
		commitQuorumCertificate = append(commitQuorumCertificate, protocol.AttributedCommitSignature{
			aps.GetSignature(),
			signer,
		})
	}
	return protocol.CertifiedCommit{
		prevHistoryDigest,
		m.Epoch,
		m.SeqNr,
		inputsDigest,
		writeSetDigest,
		stateRootDigest,
		reportsPlusPrecursorDigest,
		commitQuorumCertificate,
	}, nil
}

func (fpm *fromProtoMessage[RI]) signedHighestCertifiedTimestamp(m *SignedHighestCertifiedTimestamp) (protocol.SignedHighestCertifiedTimestamp, error) {
	if m == nil {
		return protocol.SignedHighestCertifiedTimestamp{}, fmt.Errorf("unable to extract a SignedHighestCertifiedTimestamp value")
	}
	hct, err := fpm.highestCertifiedTimestamp(m.HighestCertifiedTimestamp)
	if err != nil {
		return protocol.SignedHighestCertifiedTimestamp{}, err
	}
	return protocol.SignedHighestCertifiedTimestamp{
		hct,
		m.Signature,
	}, nil
}

func (fpm *fromProtoMessage[RI]) highestCertifiedTimestamp(m *HighestCertifiedTimestamp) (protocol.HighestCertifiedTimestamp, error) {
	if m == nil {
		return protocol.HighestCertifiedTimestamp{}, fmt.Errorf("unable to extract a HighestCertifiedTimestamp value")
	}
	return protocol.HighestCertifiedTimestamp{
		m.SeqNr,
		m.CommittedElsePrepared,
		m.Epoch,
	}, nil
}

func (fpm *fromProtoMessage[RI]) epochStartProof(m *EpochStartProof) (protocol.EpochStartProof, error) {
	if m == nil {
		return protocol.EpochStartProof{}, fmt.Errorf("unable to extract a EpochStartProof value")
	}
	hc, err := fpm.certifiedPrepareOrCommit(m.HighestCertified)
	if err != nil {
		return protocol.EpochStartProof{}, err
	}
	hctqc := make([]protocol.AttributedSignedHighestCertifiedTimestamp, 0, len(m.HighestCertifiedProof))
	for _, ashct := range m.HighestCertifiedProof {
		signer, err := fpm.oracleID(ashct.GetSigner())
		if err != nil {
			return protocol.EpochStartProof{}, err
		}
		hctqc = append(hctqc, protocol.AttributedSignedHighestCertifiedTimestamp{
			protocol.SignedHighestCertifiedTimestamp{
				protocol.HighestCertifiedTimestamp{
					ashct.GetSignedHighestCertifiedTimestamp().GetHighestCertifiedTimestamp().GetSeqNr(),
					ashct.GetSignedHighestCertifiedTimestamp().GetHighestCertifiedTimestamp().GetCommittedElsePrepared(),
					ashct.GetSignedHighestCertifiedTimestamp().GetHighestCertifiedTimestamp().GetEpoch(),
				},
				ashct.GetSignedHighestCertifiedTimestamp().GetSignature(),
			},
			signer,
		})
	}

	return protocol.EpochStartProof{
		hc,
		hctqc,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageRoundStart(m *MessageRoundStart) (protocol.MessageRoundStart[RI], error) {
	if m == nil {
		return protocol.MessageRoundStart[RI]{}, fmt.Errorf("unable to extract a MessageRoundStart value")
	}
	return protocol.MessageRoundStart[RI]{
		fpm.requestHandle,
		m.Epoch,
		m.SeqNr,
		m.Query,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageObservation(m *MessageObservation) (protocol.MessageObservation[RI], error) {
	if m == nil {
		return protocol.MessageObservation[RI]{}, fmt.Errorf("unable to extract a MessageObservation value")
	}
	so, err := fpm.signedObservation(m.SignedObservation)
	if err != nil {
		return protocol.MessageObservation[RI]{}, err
	}
	return protocol.MessageObservation[RI]{
		types.EmptyRequestHandleForInboundResponse,
		m.Epoch,
		m.SeqNr,
		so,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageReportSignatures(m *MessageReportSignatures) (protocol.MessageReportSignatures[RI], error) {
	if m == nil {
		return protocol.MessageReportSignatures[RI]{}, fmt.Errorf("unable to extract a MessageReportSignatures value")
	}

	var reportsPlusPrecursorDigest protocol.ReportsPlusPrecursorDigest
	copy(reportsPlusPrecursorDigest[:], m.ReportsPlusPrecursorDigest)

	return protocol.MessageReportSignatures[RI]{
		m.SeqNr,
		m.ReportSignatures,
		reportsPlusPrecursorDigest,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageReportsPlusPrecursorRequest(m *MessageReportsPlusPrecursorRequest) (protocol.MessageReportsPlusPrecursorRequest[RI], error) {
	if m == nil {
		return protocol.MessageReportsPlusPrecursorRequest[RI]{}, fmt.Errorf("unable to extract a MessageReportsPlusPrecursorRequest value")
	}
	return protocol.MessageReportsPlusPrecursorRequest[RI]{
		fpm.requestHandle,
		m.SeqNr,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageReportsPlusPrecursor(m *MessageReportsPlusPrecursor) (protocol.MessageReportsPlusPrecursor[RI], error) {
	if m == nil {
		return protocol.MessageReportsPlusPrecursor[RI]{}, fmt.Errorf("unable to extract a MessageReportsPlusPrecursor value")
	}
	return protocol.MessageReportsPlusPrecursor[RI]{
		fpm.requestHandle,
		m.SeqNr,
		m.ReportsPlusPrecursor,
	}, nil
}

func (fpm *fromProtoMessage[RI]) attributedSignedObservations(pbasos []*AttributedSignedObservation) ([]protocol.AttributedSignedObservation, error) {
	asos := make([]protocol.AttributedSignedObservation, 0, len(pbasos))
	for _, pbaso := range pbasos {
		aso, err := fpm.attributedSignedObservation(pbaso)
		if err != nil {
			return nil, err
		}
		asos = append(asos, aso)
	}
	return asos, nil
}

func (fpm *fromProtoMessage[RI]) attributedSignedObservation(m *AttributedSignedObservation) (protocol.AttributedSignedObservation, error) {
	if m == nil {
		return protocol.AttributedSignedObservation{}, fmt.Errorf("unable to extract an AttributedSignedObservation value")
	}

	signedObservation, err := fpm.signedObservation(m.SignedObservation)
	if err != nil {
		return protocol.AttributedSignedObservation{}, err
	}

	observer, err := fpm.oracleID(m.Observer)
	if err != nil {
		return protocol.AttributedSignedObservation{}, err
	}

	return protocol.AttributedSignedObservation{
		signedObservation,
		observer,
	}, nil
}

func (fpm *fromProtoMessage[RI]) signedObservation(m *SignedObservation) (protocol.SignedObservation, error) {
	if m == nil {
		return protocol.SignedObservation{}, fmt.Errorf("unable to extract a SignedObservation value")
	}

	return protocol.SignedObservation{
		m.Observation,
		m.Signature,
	}, nil
}

func (fpm *fromProtoMessage[RI]) oracleID(m uint32) (commontypes.OracleID, error) {
	oid := commontypes.OracleID(m)
	if int(oid) >= fpm.n {
		return 0, fmt.Errorf("invalid OracleID: %d", m)
	}
	return oid, nil
}

func (fpm *fromProtoMessage[RI]) messageBlockSyncRequest(m *MessageBlockSyncRequest) (protocol.MessageBlockSyncRequest[RI], error) {
	if m == nil {
		return protocol.MessageBlockSyncRequest[RI]{}, fmt.Errorf("unable to extract a MessageBlockSyncRequest value")
	}
	return protocol.MessageBlockSyncRequest[RI]{
		fpm.requestHandle,
		types.EmptyRequestInfoForInboundRequest,
		m.StartSeqNr,
		m.EndExclSeqNr,
		int(m.MaxCumulativeWriteSetBytes),
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageBlockSyncResponse(m *MessageBlockSyncResponse) (protocol.MessageBlockSyncResponse[RI], error) {
	if m == nil {
		return protocol.MessageBlockSyncResponse[RI]{}, fmt.Errorf("unable to extract a MessageBlockSyncResponse value")
	}
	astbs, err := fpm.attestedStateTransitionBlocks(m.AttestedStateTransitionBlocks)
	if err != nil {
		return protocol.MessageBlockSyncResponse[RI]{}, err
	}
	gstb, err := fpm.genesisStateTransitionBlock(m.GenesisStateTransitionBlock)
	if err != nil {
		return protocol.MessageBlockSyncResponse[RI]{}, err
	}
	return protocol.MessageBlockSyncResponse[RI]{
		types.EmptyRequestHandleForInboundResponse,
		m.RequestStartSeqNr,
		m.RequestEndExclSeqNr,
		astbs,
		gstb,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageStateSyncSummary(m *MessageStateSyncSummary) (protocol.MessageStateSyncSummary[RI], error) {
	if m == nil {
		return protocol.MessageStateSyncSummary[RI]{}, fmt.Errorf("unable to extract a MessageStateSyncSummary value")
	}
	return protocol.MessageStateSyncSummary[RI]{
		m.LowestPersistedSeqNr,
		m.HighestCommittedSeqNr,
	}, nil
}
func (fpm *fromProtoMessage[RI]) attestedStateTransitionBlocks(pbastbs []*AttestedStateTransitionBlock) ([]protocol.AttestedStateTransitionBlock, error) {
	astbs := make([]protocol.AttestedStateTransitionBlock, 0, len(pbastbs))
	for _, pbastb := range pbastbs {
		astb, err := fpm.attestedStateTransitionBlock(pbastb)
		if err != nil {
			return nil, err
		}
		astbs = append(astbs, astb)
	}
	return astbs, nil
}

func (fpm *fromProtoMessage[RI]) attestedStateTransitionBlock(m *AttestedStateTransitionBlock) (protocol.AttestedStateTransitionBlock, error) {
	if m == nil {
		return protocol.AttestedStateTransitionBlock{}, fmt.Errorf("unable to extract a AttestedStateTransitionBlock value")
	}
	stb, err := fpm.stateTransitionBlock(m.StateTransitionBlock)
	if err != nil {
		return protocol.AttestedStateTransitionBlock{}, err
	}
	asigs, err := fpm.attributedCommitSignatures(m.AttributedSignatures)
	if err != nil {
		return protocol.AttestedStateTransitionBlock{}, err
	}
	return protocol.AttestedStateTransitionBlock{
		stb,
		asigs,
	}, nil
}

func (fpm *fromProtoMessage[RI]) stateTransitionBlock(m *StateTransitionBlock) (protocol.StateTransitionBlock, error) {
	if m == nil {
		return protocol.StateTransitionBlock{}, fmt.Errorf("unable to extract a StateTransitionBlock value")
	}
	var prevHistoryDigest protocol.HistoryDigest
	copy(prevHistoryDigest[:], m.PrevHistoryDigest)
	var inputsDigest protocol.StateTransitionInputsDigest
	copy(inputsDigest[:], m.StateTransitionInputsDigest)
	var stateRootDigest protocol.StateRootDigest
	copy(stateRootDigest[:], m.StateRootDigest)
	writeSet, err := fpm.stateWriteSet(m.StateWriteSet)
	if err != nil {
		return protocol.StateTransitionBlock{}, err
	}
	var reportsPlusPrecursorDigest protocol.ReportsPlusPrecursorDigest
	copy(reportsPlusPrecursorDigest[:], m.ReportsPlusPrecursorDigest)
	return protocol.StateTransitionBlock{
		prevHistoryDigest,
		m.Epoch,
		m.SeqNr,
		inputsDigest,
		writeSet,
		stateRootDigest,
		reportsPlusPrecursorDigest,
	}, nil
}

func (fpm *fromProtoMessage[RI]) genesisStateTransitionBlock(m *GenesisStateTransitionBlock) (*protocol.GenesisStateTransitionBlock, error) {
	if m == nil {
		return nil, nil
	}
	var prevHistoryDigest protocol.HistoryDigest
	copy(prevHistoryDigest[:], m.PrevHistoryDigest)
	var inputsDigest protocol.StateTransitionInputsDigest
	copy(inputsDigest[:], m.StateTransitionInputsDigest)
	var writeSetDigest protocol.StateWriteSetDigest
	copy(writeSetDigest[:], m.StateWriteSetDigest)
	var stateRootDigest protocol.StateRootDigest
	copy(stateRootDigest[:], m.StateRootDigest)
	var reportsPlusPrecursorDigest protocol.ReportsPlusPrecursorDigest
	copy(reportsPlusPrecursorDigest[:], m.ReportsPlusPrecursorDigest)
	return &protocol.GenesisStateTransitionBlock{
		prevHistoryDigest,
		m.SeqNr,
		inputsDigest,
		writeSetDigest,
		stateRootDigest,
		reportsPlusPrecursorDigest,
	}, nil
}

func (fpm *fromProtoMessage[RI]) attributedCommitSignatures(pbasigs []*AttributedCommitSignature) ([]protocol.AttributedCommitSignature, error) {
	asigs := make([]protocol.AttributedCommitSignature, 0, len(pbasigs))
	for _, pbasig := range pbasigs {
		asig, err := fpm.attributedCommitSignature(pbasig)
		if err != nil {
			return nil, err
		}
		asigs = append(asigs, asig)
	}
	return asigs, nil
}

func (fpm *fromProtoMessage[RI]) attributedCommitSignature(m *AttributedCommitSignature) (protocol.AttributedCommitSignature, error) {
	if m == nil {
		return protocol.AttributedCommitSignature{}, fmt.Errorf("unable to extract an AttributedCommitSignature value")
	}
	signer, err := fpm.oracleID(m.GetSigner())
	if err != nil {
		return protocol.AttributedCommitSignature{}, err
	}
	return protocol.AttributedCommitSignature{
		m.Signature,
		signer,
	}, nil
}

func (fpm *fromProtoMessage[RI]) keyValuePairWithDeletions(pbkvmod *KeyValueModification) (protocol.KeyValuePairWithDeletions, error) {
	if pbkvmod == nil {
		return protocol.KeyValuePairWithDeletions{}, fmt.Errorf("unable to extract a KeyValuePairWithDeletions value")
	}

	return protocol.KeyValuePairWithDeletions{
		pbkvmod.Key,
		pbkvmod.Value,
		pbkvmod.Deleted,
	}, nil
}

func (fpm *fromProtoMessage[RI]) stateWriteSet(m *StateWriteSet) (protocol.StateWriteSet, error) {
	if m == nil {
		return protocol.StateWriteSet{}, fmt.Errorf("unable to extract an StateWriteSet value")
	}

	writeSet := make([]protocol.KeyValuePairWithDeletions, 0, len(m.Entries))
	for _, pbkvmod := range m.Entries {
		kvwd, err := fpm.keyValuePairWithDeletions(pbkvmod)
		if err != nil {
			return protocol.StateWriteSet{}, err
		}
		writeSet = append(writeSet, kvwd)
	}

	return protocol.StateWriteSet{writeSet}, nil
}

func (fpm *fromProtoMessage[RI]) messageBlobOffer(m *MessageBlobOffer) (protocol.MessageBlobOffer[RI], error) {
	if m == nil {
		return protocol.MessageBlobOffer[RI]{}, fmt.Errorf("unable to extract a MessageBlobOffer value")
	}
	var chunkDigestsRoot mt.Digest
	copy(chunkDigestsRoot[:], m.ChunkDigestsRoot)
	return protocol.MessageBlobOffer[RI]{
		fpm.requestHandle,
		types.EmptyRequestInfoForInboundRequest,
		chunkDigestsRoot,
		m.PayloadLength,
		m.ExpirySeqNr,
	}, nil
}

func (fpm *fromProtoMessage[RI]) chunkDigests(pbcds [][]byte) ([]protocol.BlobChunkDigest, error) {
	cds := make([]protocol.BlobChunkDigest, 0, len(pbcds))
	for _, pbcd := range pbcds {
		var blockChunkDigest protocol.BlobChunkDigest
		copy(blockChunkDigest[:], pbcd)

		cds = append(cds, blockChunkDigest)
	}
	return cds, nil
}

func (fpm *fromProtoMessage[RI]) messageBlobOfferResponse(m *MessageBlobOfferResponse) (protocol.MessageBlobOfferResponse[RI], error) {
	if m == nil {
		return protocol.MessageBlobOfferResponse[RI]{}, fmt.Errorf("unable to extract a MessageBlobOfferResponse value")
	}

	var blobDigest protocol.BlobDigest
	copy(blobDigest[:], m.BlobDigest)

	return protocol.MessageBlobOfferResponse[RI]{
		types.EmptyRequestHandleForInboundResponse,
		blobDigest,
		m.RejectOffer,
		m.Signature,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageBlobChunkRequest(m *MessageBlobChunkRequest) (protocol.MessageBlobChunkRequest[RI], error) {
	if m == nil {
		return protocol.MessageBlobChunkRequest[RI]{}, fmt.Errorf("unable to extract a MessageBlobChunkRequest value")
	}

	var blobDigest protocol.BlobDigest
	copy(blobDigest[:], m.BlobDigest)

	return protocol.MessageBlobChunkRequest[RI]{
		fpm.requestHandle,
		types.EmptyRequestInfoForInboundRequest,
		blobDigest,
		m.ChunkIndex,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageBlobChunkResponse(m *MessageBlobChunkResponse) (protocol.MessageBlobChunkResponse[RI], error) {
	if m == nil {
		return protocol.MessageBlobChunkResponse[RI]{}, fmt.Errorf("unable to extract a MessageBlobChunkResponse value")
	}

	var blobDigest protocol.BlobDigest
	copy(blobDigest[:], m.BlobDigest)

	proof := make([]mt.Digest, 0, len(m.Proof))
	for _, p := range m.Proof {
		var pmt mt.Digest
		copy(pmt[:], p)
		proof = append(proof, pmt)
	}

	return protocol.MessageBlobChunkResponse[RI]{
		types.EmptyRequestHandleForInboundResponse,
		blobDigest,
		m.ChunkIndex,
		m.GoAway,
		m.Chunk,
		proof,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageTreeSyncChunkRequest(m *MessageTreeSyncChunkRequest) (protocol.MessageTreeSyncChunkRequest[RI], error) {
	if m == nil {
		return protocol.MessageTreeSyncChunkRequest[RI]{}, fmt.Errorf("unable to extract an MessageTreeSyncRequest value")
	}
	startIndex, err := fpm.digest(m.StartIndex)
	if err != nil {
		return protocol.MessageTreeSyncChunkRequest[RI]{}, err
	}
	endInclIndex, err := fpm.digest(m.EndInclIndex)
	if err != nil {
		return protocol.MessageTreeSyncChunkRequest[RI]{}, err
	}
	return protocol.MessageTreeSyncChunkRequest[RI]{
		fpm.requestHandle,
		types.EmptyRequestInfoForInboundRequest,
		m.ToSeqNr,
		startIndex,
		endInclIndex,
		int(m.MaxCumulativeKeysPlusValuesBytes),
	}, nil
}

func (fpm *fromProtoMessage[RI]) treeSyncChunkProof(pbpns [][]byte) ([]jmt.Digest, error) {
	proof := make([]jmt.Digest, 0, len(pbpns))
	for _, pbpn := range pbpns {
		var pn jmt.Digest
		if len(pbpn) != len(pn) {
			return proof, fmt.Errorf("invalid proof node length, expected %d, got %d ", len(pn), len(pbpn))
		}
		copy(pn[:], pbpn)
		proof = append(proof, pn)
	}
	return proof, nil
}

func (fpm *fromProtoMessage[RI]) keyValuePair(pbkvp *KeyValuePair) (protocol.KeyValuePair, error) {
	if pbkvp == nil {
		return protocol.KeyValuePair{}, fmt.Errorf("unable to extract a KeyValuePair value")
	}

	return protocol.KeyValuePair{
		pbkvp.Key,
		pbkvp.Value,
	}, nil
}

func (fpm *fromProtoMessage[RI]) treeSyncChunkKeyValues(pbkvps []*KeyValuePair) ([]protocol.KeyValuePair, error) {
	kvps := make([]protocol.KeyValuePair, 0, len(pbkvps))
	for _, pbkvp := range pbkvps {
		kv, err := fpm.keyValuePair(pbkvp)
		if err != nil {
			return nil, err
		}
		kvps = append(kvps, kv)
	}
	return kvps, nil
}

func (fpm *fromProtoMessage[RI]) digest(pbdigest []byte) (jmt.Digest, error) {
	if len(pbdigest) != len(jmt.Digest{}) {
		return jmt.Digest{}, fmt.Errorf("digest must be %d bytes, got %d", len(jmt.Digest{}), len(pbdigest))
	}
	var digest jmt.Digest
	copy(digest[:], pbdigest)
	return digest, nil
}

func (fpm *fromProtoMessage[RI]) leafKeyAndValueDigests(pblkd *LeafKeyAndValueDigests) (jmt.LeafKeyAndValueDigests, error) {
	if pblkd == nil {
		return jmt.LeafKeyAndValueDigests{}, fmt.Errorf("unable to extract a LeafKeyAndValueDigests value")
	}
	keyDigest, err := fpm.digest(pblkd.KeyDigest)
	if err != nil {
		return jmt.LeafKeyAndValueDigests{}, err
	}
	valueDigest, err := fpm.digest(pblkd.ValueDigest)
	if err != nil {
		return jmt.LeafKeyAndValueDigests{}, err
	}
	return jmt.LeafKeyAndValueDigests{
		keyDigest,
		valueDigest,
	}, nil
}

func (fpm *fromProtoMessage[RI]) boundingLeaf(pbbl *BoundingLeaf) (jmt.BoundingLeaf, error) {
	if pbbl == nil {
		return jmt.BoundingLeaf{}, fmt.Errorf("unable to extract a BoundingLeaf value")
	}
	leafKeyAndValueDigests, err := fpm.leafKeyAndValueDigests(pbbl.Leaf)
	if err != nil {
		return jmt.BoundingLeaf{}, err
	}
	siblings, err := fpm.treeSyncChunkProof(pbbl.Siblings)
	if err != nil {
		return jmt.BoundingLeaf{}, err
	}
	return jmt.BoundingLeaf{
		leafKeyAndValueDigests,
		siblings,
	}, nil
}

func (fpm *fromProtoMessage[RI]) treeSyncChunkBoundingLeaves(pbbls []*BoundingLeaf) ([]jmt.BoundingLeaf, error) {
	boundingLeaves := make([]jmt.BoundingLeaf, 0, len(pbbls))
	for _, pbbl := range pbbls {
		bl, err := fpm.boundingLeaf(pbbl)
		if err != nil {
			return nil, err
		}
		boundingLeaves = append(boundingLeaves, bl)
	}
	return boundingLeaves, nil
}

func (fpm *fromProtoMessage[RI]) messageTreeSyncChunkResponse(m *MessageTreeSyncChunkResponse) (protocol.MessageTreeSyncChunkResponse[RI], error) {
	if m == nil {
		return protocol.MessageTreeSyncChunkResponse[RI]{}, fmt.Errorf("unable to extract an MessageTreeSyncChunk value")
	}
	startIndex, err := fpm.digest(m.StartIndex)
	if err != nil {
		return protocol.MessageTreeSyncChunkResponse[RI]{}, err
	}
	requestEndInclIndex, err := fpm.digest(m.RequestEndInclIndex)
	if err != nil {
		return protocol.MessageTreeSyncChunkResponse[RI]{}, err
	}
	encInclIndex, err := fpm.digest(m.EndInclIndex)
	if err != nil {
		return protocol.MessageTreeSyncChunkResponse[RI]{}, err
	}
	keyValues, err := fpm.treeSyncChunkKeyValues(m.KeyValues)
	if err != nil {
		return protocol.MessageTreeSyncChunkResponse[RI]{}, err
	}
	boundingLeaves, err := fpm.treeSyncChunkBoundingLeaves(m.BoundingLeaves)
	if err != nil {
		return protocol.MessageTreeSyncChunkResponse[RI]{}, err
	}
	return protocol.MessageTreeSyncChunkResponse[RI]{
		types.EmptyRequestHandleForInboundResponse,
		m.ToSeqNr,
		startIndex,
		requestEndInclIndex,
		m.GoAway,
		encInclIndex,
		keyValues,
		boundingLeaves,
	}, nil
}

func (fpm *fromProtoMessage[RI]) keyDigestRange(pbr *KeyDigestRange) (protocol.KeyDigestRange, error) {
	if pbr == nil {
		return protocol.KeyDigestRange{}, fmt.Errorf("unable to extract a KeyDigestRange value")
	}
	startIndex, err := fpm.digest(pbr.StartIndex)
	if err != nil {
		return protocol.KeyDigestRange{}, err
	}
	endInclIndex, err := fpm.digest(pbr.EndInclIndex)
	if err != nil {
		return protocol.KeyDigestRange{}, err
	}
	return protocol.KeyDigestRange{startIndex, endInclIndex}, nil
}

func (fpm *fromProtoMessage[RI]) pendingKeyDigestRanges(pbRanges []*KeyDigestRange) (protocol.PendingKeyDigestRanges, error) {
	ranges := make([]protocol.KeyDigestRange, 0, len(pbRanges))
	for _, pbr := range pbRanges {
		kr, err := fpm.keyDigestRange(pbr)
		if err != nil {
			return protocol.PendingKeyDigestRanges{}, err
		}
		ranges = append(ranges, kr)
	}
	return protocol.NewPendingKeyDigestRanges(ranges), nil
}
