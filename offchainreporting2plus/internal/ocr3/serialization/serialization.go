package serialization

import (
	"fmt"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3/protocol"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"

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

// Deserialize decodes a binary payload into a protocol.Message
func Deserialize[RI any](n int, b []byte) (protocol.Message[RI], *MessageWrapper, error) {
	pb := &MessageWrapper{}
	if err := proto.Unmarshal(b, pb); err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal protobuf: %w", err)
	}

	fpm := fromProtoMessage[RI]{n}
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
	fpm := fromProtoMessage[struct{}]{n}
	return fpm.certifiedPrepareOrCommit(&pb)
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
			v.Signature31,
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
			v.Signature31,
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
			v.Signature31,
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
		}
		msgWrapper.Msg = &MessageWrapper_MessageReportSignatures{pm}
	case protocol.MessageCertifiedCommitRequest[RI]:
		pm := &MessageCertifiedCommitRequest{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.SeqNr,
		}
		msgWrapper.Msg = &MessageWrapper_MessageCertifiedCommitRequest{pm}
	case protocol.MessageCertifiedCommit[RI]:
		pm := &MessageCertifiedCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			tpm.CertifiedCommit(v.CertifiedCommit),
		}
		msgWrapper.Msg = &MessageWrapper_MessageCertifiedCommit{pm}

	default:
		return nil, fmt.Errorf("unable to serialize message of type %T", m)

	}
	return &msgWrapper, nil
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
				uint64(v.PrepareEpoch),
				v.SeqNr,
				v.OutcomeInputsDigest[:],
				v.Outcome,
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
		uint64(cpocc.CommitEpoch),
		cpocc.SeqNr,
		cpocc.Outcome,
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
		shct.Signature31,
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

//
// *fromProtoMessage
//

type fromProtoMessage[RI any] struct {
	n int
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
	case *MessageWrapper_MessageCertifiedCommitRequest:
		return fpm.messageCertifiedCommitRequest(wrapper.GetMessageCertifiedCommitRequest())
	case *MessageWrapper_MessageCertifiedCommit:
		return fpm.messageCertifiedCommit(wrapper.GetMessageCertifiedCommit())
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
		m.Signature31,
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
		m.Signature31,
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
	default:
		return nil, fmt.Errorf("unknown case of CertifiedPrepareOrCommit")
	}
}

func (fpm *fromProtoMessage[RI]) certifiedPrepare(m *CertifiedPrepare) (protocol.CertifiedPrepare, error) {
	if m == nil {
		return protocol.CertifiedPrepare{}, fmt.Errorf("unable to extract a CertifiedPrepare value")
	}
	var outcomeInputsDigest protocol.OutcomeInputsDigest
	copy(outcomeInputsDigest[:], m.OutcomeInputsDigest)
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
		m.PrepareEpoch,
		m.SeqNr,
		outcomeInputsDigest,
		m.Outcome,
		prepareQuorumCertificate,
	}, nil

}

func (fpm *fromProtoMessage[RI]) certifiedCommit(m *CertifiedCommit) (protocol.CertifiedCommit, error) {
	if m == nil {
		return protocol.CertifiedCommit{}, fmt.Errorf("unable to extract a CertifiedCommit value")
	}
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
		m.CommitEpoch,
		m.SeqNr,
		m.Outcome,
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
		m.Signature31,
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
				ashct.GetSignedHighestCertifiedTimestamp().GetSignature31(),
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
		m.Epoch,
		m.SeqNr,
		m.Query,
		m.Signature31,
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
		m.Epoch,
		m.SeqNr,
		so,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageReportSignatures(m *MessageReportSignatures) (protocol.MessageReportSignatures[RI], error) {
	if m == nil {
		return protocol.MessageReportSignatures[RI]{}, fmt.Errorf("unable to extract a MessageReportSignatures value")
	}
	return protocol.MessageReportSignatures[RI]{
		m.SeqNr,
		m.ReportSignatures,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageCertifiedCommitRequest(m *MessageCertifiedCommitRequest) (protocol.MessageCertifiedCommitRequest[RI], error) {
	if m == nil {
		return protocol.MessageCertifiedCommitRequest[RI]{}, fmt.Errorf("unable to extract a MessageCertifiedCommitRequest value")
	}
	return protocol.MessageCertifiedCommitRequest[RI]{
		m.SeqNr,
	}, nil
}

func (fpm *fromProtoMessage[RI]) messageCertifiedCommit(m *MessageCertifiedCommit) (protocol.MessageCertifiedCommit[RI], error) {
	if m == nil {
		return protocol.MessageCertifiedCommit[RI]{}, fmt.Errorf("unable to extract a MessageCertifiedCommit value")
	}
	cpocc, err := fpm.certifiedCommit(m.CertifiedCommit)
	if err != nil {
		return protocol.MessageCertifiedCommit[RI]{}, err
	}
	return protocol.MessageCertifiedCommit[RI]{
		cpocc,
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
