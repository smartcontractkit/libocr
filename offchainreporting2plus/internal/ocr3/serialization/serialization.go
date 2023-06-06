package serialization

import (
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol"

	"google.golang.org/protobuf/proto"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

// Serialize encodes a protocol.Message into a binary payload
func Serialize[RI any](m protocol.Message[RI]) (b []byte, pbm *MessageWrapper, err error) {
	pbm, err = toProtoMessage(m)
	if err != nil {
		return nil, nil, err
	}
	b, err = proto.Marshal(pbm)
	if err != nil {
		return nil, nil, err
	}
	return b, pbm, nil
}

// Deserialize decodes a binary payload into a protocol.Message
func Deserialize[RI any](b []byte) (protocol.Message[RI], *MessageWrapper, error) {
	pbm := &MessageWrapper{}
	err := proto.Unmarshal(b, pbm)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal protobuf: %w", err)
	}
	m, err := messageWrapperFromProtoMessage[RI](pbm)
	if err != nil {
		return nil, nil, fmt.Errorf("could not translate protobuf to protocol.Message: %w", err)
	}
	return m, pbm, nil
}

//
// *toProtoMessage
//

func toProtoMessage[RI any](m protocol.Message[RI]) (*MessageWrapper, error) {
	msgWrapper := MessageWrapper{}
	switch v := m.(type) {
	case protocol.MessageNewEpoch[RI]:
		pm := &MessageNewEpoch{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
		}
		msgWrapper.Msg = &MessageWrapper_MessageNewEpoch{pm}
	case protocol.MessageReconcile[RI]:
		pm := &MessageReconcile{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			CertifiedPrepareOrCommitToProtoMessage(v.HighestCertified),
			signedHighestCertifiedTimestampToProtoMessage(v.SignedHighestCertifiedTimestamp),
		}
		msgWrapper.Msg = &MessageWrapper_MessageReconcile{pm}
	case protocol.MessageStartEpoch[RI]:
		pm := &MessageStartEpoch{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			startEpochProofToProtoMessage(v.StartEpochProof),
		}
		msgWrapper.Msg = &MessageWrapper_MessageStartEpoch{pm}
	case protocol.MessageStartRound[RI]:
		pm := &MessageStartRound{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			v.Query,
		}
		msgWrapper.Msg = &MessageWrapper_MessageStartRound{pm}
	case protocol.MessageObserve[RI]:
		pm := &MessageObserve{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			signedObservationToProtoMessage(v.SignedObservation),
		}
		msgWrapper.Msg = &MessageWrapper_MessageObserve{pm}

	case protocol.MessagePropose[RI]:
		pbasos := make([]*AttributedSignedObservation, 0, len(v.AttributedSignedObservations))
		for _, aso := range v.AttributedSignedObservations {
			pbasos = append(pbasos, attributedSignedObservationToProtoMessage(aso))
		}
		pm := &MessagePropose{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			uint64(v.Epoch),
			v.SeqNr,
			pbasos,
		}
		msgWrapper.Msg = &MessageWrapper_MessagePropose{pm}
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
	case protocol.MessageFinal[RI]:
		pm := &MessageFinal{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.SeqNr,
			v.ReportSignatures,
		}
		msgWrapper.Msg = &MessageWrapper_MessageFinal{pm}
	case protocol.MessageRequestCertifiedCommit[RI]:
		pm := &MessageRequestCertifiedCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			v.SeqNr,
		}
		msgWrapper.Msg = &MessageWrapper_MessageRequestCertifiedCommit{pm}
	case protocol.MessageSupplyCertifiedCommit[RI]:
		pm := &MessageSupplyCertifiedCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			CertifiedPrepareOrCommitCommitToProtoMessage(v.CertifiedCommit),
		}
		msgWrapper.Msg = &MessageWrapper_MessageSupplyCertifiedCommit{pm}

	default:
		return nil, fmt.Errorf("unable to serialize message of type %T", m)

	}
	return &msgWrapper, nil
}

func CertifiedPrepareOrCommitToProtoMessage(cpoc protocol.CertifiedPrepareOrCommit) *CertifiedPrepareOrCommit {
	switch v := cpoc.(type) {
	case *protocol.CertifiedPrepareOrCommitPrepare:
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
			&CertifiedPrepareOrCommit_Prepare{&CertifiedPrepareOrCommitPrepare{
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
	case *protocol.CertifiedPrepareOrCommitCommit:
		return &CertifiedPrepareOrCommit{
			// zero-initialize protobuf built-ins
			protoimpl.MessageState{},
			0,
			nil,
			// fields
			&CertifiedPrepareOrCommit_Commit{CertifiedPrepareOrCommitCommitToProtoMessage(*v)},
		}
	default:

		panic("unrecognized")
	}
}

func CertifiedPrepareOrCommitCommitToProtoMessage(cpocc protocol.CertifiedPrepareOrCommitCommit) *CertifiedPrepareOrCommitCommit {
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
	return &CertifiedPrepareOrCommitCommit{
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

func attributedSignedHighestCertifiedTimestampToProtoMessage(ashct protocol.AttributedSignedHighestCertifiedTimestamp) *AttributedSignedHighestCertifiedTimestamp {
	return &AttributedSignedHighestCertifiedTimestamp{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		signedHighestCertifiedTimestampToProtoMessage(ashct.SignedHighestCertifiedTimestamp),
		uint32(ashct.Signer),
	}
}

func signedHighestCertifiedTimestampToProtoMessage(shct protocol.SignedHighestCertifiedTimestamp) *SignedHighestCertifiedTimestamp {
	return &SignedHighestCertifiedTimestamp{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		highestCertifiedTimestampToProtoMessage(shct.HighestCertifiedTimestamp),
		shct.Signature,
	}
}

func highestCertifiedTimestampToProtoMessage(hct protocol.HighestCertifiedTimestamp) *HighestCertifiedTimestamp {
	return &HighestCertifiedTimestamp{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		hct.SeqNr,
		hct.CommittedElsePrepared,
	}
}

func startEpochProofToProtoMessage(srqc protocol.StartEpochProof) *StartEpochProof {
	highestCertifiedProof := make([]*AttributedSignedHighestCertifiedTimestamp, 0, len(srqc.HighestCertifiedProof))
	for _, ashct := range srqc.HighestCertifiedProof {
		highestCertifiedProof = append(highestCertifiedProof, attributedSignedHighestCertifiedTimestampToProtoMessage(ashct))
	}
	return &StartEpochProof{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		CertifiedPrepareOrCommitToProtoMessage(srqc.HighestCertified),
		highestCertifiedProof,
	}
}

func signedObservationToProtoMessage(o protocol.SignedObservation) *SignedObservation {
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

func attributedSignedObservationToProtoMessage(aso protocol.AttributedSignedObservation) *AttributedSignedObservation {
	return &AttributedSignedObservation{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		signedObservationToProtoMessage(aso.SignedObservation),
		uint32(aso.Observer),
	}
}

//
// *fromProtoMessage
//

func messageWrapperFromProtoMessage[RI any](wrapper *MessageWrapper) (protocol.Message[RI], error) {
	switch msg := wrapper.Msg.(type) {
	case *MessageWrapper_MessageNewEpoch:
		return messageNewEpochFromProtoMessage[RI](wrapper.GetMessageNewEpoch())
	case *MessageWrapper_MessageReconcile:
		return messageReconcileFromProtoMessage[RI](wrapper.GetMessageReconcile())
	case *MessageWrapper_MessageStartEpoch:
		return messageStartEpochFromProtoMessage[RI](wrapper.GetMessageStartEpoch())
	case *MessageWrapper_MessageStartRound:
		return messageStartRoundFromProtoMessage[RI](wrapper.GetMessageStartRound())
	case *MessageWrapper_MessageObserve:
		return messageObserveFromProtoMessage[RI](wrapper.GetMessageObserve())
	case *MessageWrapper_MessagePropose:
		return messageProposeFromProtoMessage[RI](wrapper.GetMessagePropose())
	case *MessageWrapper_MessagePrepare:
		return messagePrepareFromProtoMessage[RI](wrapper.GetMessagePrepare())
	case *MessageWrapper_MessageCommit:
		return messageCommitFromProtoMessage[RI](wrapper.GetMessageCommit())
	case *MessageWrapper_MessageFinal:
		return messageFinalFromProtoMessage[RI](wrapper.GetMessageFinal())
	case *MessageWrapper_MessageRequestCertifiedCommit:
		return messageRequestCertifiedCommitFromProtoMessage[RI](wrapper.GetMessageRequestCertifiedCommit())
	case *MessageWrapper_MessageSupplyCertifiedCommit:
		return messageSupplyCertifiedCommitFromProtoMessage[RI](wrapper.GetMessageSupplyCertifiedCommit())
	default:
		return nil, fmt.Errorf("unrecognized Msg type %T", msg)
	}
}

func messageNewEpochFromProtoMessage[RI any](m *MessageNewEpoch) (protocol.MessageNewEpoch[RI], error) {
	if m == nil {
		return protocol.MessageNewEpoch[RI]{}, fmt.Errorf("unable to extract a MessageNewEpoch value")
	}
	return protocol.MessageNewEpoch[RI]{
		m.Epoch,
	}, nil
}

func messageReconcileFromProtoMessage[RI any](m *MessageReconcile) (protocol.MessageReconcile[RI], error) {
	if m == nil {
		return protocol.MessageReconcile[RI]{}, fmt.Errorf("unable to extract a MessageReconcile value")
	}
	hc, err := CertifiedPrepareOrCommitFromProtoMessage(m.HighestCertified)
	if err != nil {
		return protocol.MessageReconcile[RI]{}, err
	}
	shct, err := signedHighestCertifiedTimestampFromProtoMessage(m.SignedHighestCertifiedTimestamp)
	if err != nil {
		return protocol.MessageReconcile[RI]{}, err
	}
	return protocol.MessageReconcile[RI]{
		m.Epoch,
		hc,
		shct,
	}, nil
}

func messageStartEpochFromProtoMessage[RI any](m *MessageStartEpoch) (protocol.MessageStartEpoch[RI], error) {
	if m == nil {
		return protocol.MessageStartEpoch[RI]{}, fmt.Errorf("unable to extract a MessageStartEpoch value")
	}
	srqc, err := startEpochProofFromProtoMessage(m.StartEpochProof)
	if err != nil {
		return protocol.MessageStartEpoch[RI]{}, err
	}
	return protocol.MessageStartEpoch[RI]{
		m.Epoch,
		srqc,
	}, nil
}

func messageProposeFromProtoMessage[RI any](m *MessagePropose) (protocol.MessagePropose[RI], error) {
	if m == nil {
		return protocol.MessagePropose[RI]{}, fmt.Errorf("unable to extract a MessagePropose value")
	}
	asos, err := attributedSignedObservationsFromProtoMessage(m.AttributedSignedObservations)
	if err != nil {
		return protocol.MessagePropose[RI]{}, err
	}
	return protocol.MessagePropose[RI]{
		m.Epoch,
		m.SeqNr,
		asos,
	}, nil
}

func messagePrepareFromProtoMessage[RI any](m *MessagePrepare) (protocol.MessagePrepare[RI], error) {
	if m == nil {
		return protocol.MessagePrepare[RI]{}, fmt.Errorf("unable to extract a MessagePrepare value")
	}
	return protocol.MessagePrepare[RI]{
		m.Epoch,
		m.SeqNr,
		m.Signature,
	}, nil
}

func messageCommitFromProtoMessage[RI any](m *MessageCommit) (protocol.MessageCommit[RI], error) {
	if m == nil {
		return protocol.MessageCommit[RI]{}, fmt.Errorf("unable to extract a MessageCommit value")
	}
	return protocol.MessageCommit[RI]{
		m.Epoch,
		m.SeqNr,
		m.Signature,
	}, nil
}

func CertifiedPrepareOrCommitFromProtoMessage(m *CertifiedPrepareOrCommit) (protocol.CertifiedPrepareOrCommit, error) {
	if m == nil {
		return nil, fmt.Errorf("unable to extract a CertifiedPrepareOrCommit value")
	}
	switch poc := m.PrepareOrCommit.(type) {
	case *CertifiedPrepareOrCommit_Prepare:
		prepareQuorumCertificate := make([]protocol.AttributedPrepareSignature, 0, len(poc.Prepare.GetPrepareQuorumCertificate()))
		for _, aps := range poc.Prepare.GetPrepareQuorumCertificate() {
			prepareQuorumCertificate = append(prepareQuorumCertificate, protocol.AttributedPrepareSignature{
				aps.GetSignature(),
				commontypes.OracleID(aps.GetSigner()),
			})
		}

		var outcomeInputsDigest protocol.OutcomeInputsDigest
		copy(outcomeInputsDigest[:], poc.Prepare.OutcomeInputsDigest)

		return &protocol.CertifiedPrepareOrCommitPrepare{
			poc.Prepare.PrepareEpoch,
			poc.Prepare.SeqNr,
			outcomeInputsDigest,
			poc.Prepare.Outcome,
			prepareQuorumCertificate,
		}, nil
	case *CertifiedPrepareOrCommit_Commit:
		cpocc, err := certifiedPrepareOrCommitCommitFromProtoMessage(poc.Commit)
		if err != nil {
			return nil, err
		}
		return &cpocc, nil
	default:
		return nil, fmt.Errorf("unknown case of CertifiedPrepareOrCommit")
	}
}

func certifiedPrepareOrCommitCommitFromProtoMessage(m *CertifiedPrepareOrCommitCommit) (protocol.CertifiedPrepareOrCommitCommit, error) {
	if m == nil {
		return protocol.CertifiedPrepareOrCommitCommit{}, fmt.Errorf("unable to extract a CertifiedPrepareOrCommitCommit value")
	}
	commitQuorumCertificate := make([]protocol.AttributedCommitSignature, 0, len(m.CommitQuorumCertificate))
	for _, aps := range m.CommitQuorumCertificate {
		commitQuorumCertificate = append(commitQuorumCertificate, protocol.AttributedCommitSignature{
			aps.GetSignature(),
			commontypes.OracleID(aps.GetSigner()),
		})
	}
	return protocol.CertifiedPrepareOrCommitCommit{
		m.CommitEpoch,
		m.SeqNr,
		m.Outcome,
		commitQuorumCertificate,
	}, nil
}

func signedHighestCertifiedTimestampFromProtoMessage(m *SignedHighestCertifiedTimestamp) (protocol.SignedHighestCertifiedTimestamp, error) {
	return protocol.SignedHighestCertifiedTimestamp{
		protocol.HighestCertifiedTimestamp{
			m.GetHighestCertifiedTimestamp().GetSeqNr(),
			m.HighestCertifiedTimestamp.GetCommittedElsePrepared(),
		},
		m.GetSignature(),
	}, nil
}

func startEpochProofFromProtoMessage(m *StartEpochProof) (protocol.StartEpochProof, error) {
	if m == nil {
		return protocol.StartEpochProof{}, fmt.Errorf("unable to extract a StartEpochProof value")
	}
	hc, err := CertifiedPrepareOrCommitFromProtoMessage(m.HighestCertified)
	if err != nil {
		return protocol.StartEpochProof{}, err
	}
	hctqc := make([]protocol.AttributedSignedHighestCertifiedTimestamp, 0, len(m.HighestCertifiedProof))
	for _, ashct := range m.HighestCertifiedProof {
		hctqc = append(hctqc, protocol.AttributedSignedHighestCertifiedTimestamp{
			protocol.SignedHighestCertifiedTimestamp{
				protocol.HighestCertifiedTimestamp{
					ashct.GetSignedHighestCertifiedTimestamp().GetHighestCertifiedTimestamp().GetSeqNr(),
					ashct.GetSignedHighestCertifiedTimestamp().GetHighestCertifiedTimestamp().GetCommittedElsePrepared(),
				},
				ashct.GetSignedHighestCertifiedTimestamp().GetSignature(),
			},
			commontypes.OracleID(ashct.GetSigner()),
		})
	}

	return protocol.StartEpochProof{
		hc,
		hctqc,
	}, nil
}

func messageStartRoundFromProtoMessage[RI any](m *MessageStartRound) (protocol.MessageStartRound[RI], error) {
	if m == nil {
		return protocol.MessageStartRound[RI]{}, fmt.Errorf("unable to extract a MessageStartRound value")
	}
	return protocol.MessageStartRound[RI]{
		m.Epoch,
		m.SeqNr,
		m.Query,
	}, nil
}

func messageObserveFromProtoMessage[RI any](m *MessageObserve) (protocol.MessageObserve[RI], error) {
	if m == nil {
		return protocol.MessageObserve[RI]{}, fmt.Errorf("unable to extract a MessageObserve value")
	}
	so, err := signedObservationFromProtoMessage(m.SignedObservation)
	if err != nil {
		return protocol.MessageObserve[RI]{}, err
	}
	return protocol.MessageObserve[RI]{
		m.Epoch,
		m.SeqNr,
		so,
	}, nil
}

func messageFinalFromProtoMessage[RI any](m *MessageFinal) (protocol.MessageFinal[RI], error) {
	if m == nil {
		return protocol.MessageFinal[RI]{}, fmt.Errorf("unable to extract a MessageFinal value")
	}
	return protocol.MessageFinal[RI]{
		m.SeqNr,
		m.ReportSignatures,
	}, nil
}

func messageRequestCertifiedCommitFromProtoMessage[RI any](m *MessageRequestCertifiedCommit) (protocol.MessageRequestCertifiedCommit[RI], error) {
	if m == nil {
		return protocol.MessageRequestCertifiedCommit[RI]{}, fmt.Errorf("unable to extract a MessageRequestCertifiedCommit value")
	}
	return protocol.MessageRequestCertifiedCommit[RI]{
		m.SeqNr,
	}, nil
}

func messageSupplyCertifiedCommitFromProtoMessage[RI any](m *MessageSupplyCertifiedCommit) (protocol.MessageSupplyCertifiedCommit[RI], error) {
	if m == nil {
		return protocol.MessageSupplyCertifiedCommit[RI]{}, fmt.Errorf("unable to extract a MessageSupplyCertifiedCommit value")
	}
	cpocc, err := certifiedPrepareOrCommitCommitFromProtoMessage(m.CertifiedCommit)
	if err != nil {
		return protocol.MessageSupplyCertifiedCommit[RI]{}, err
	}
	return protocol.MessageSupplyCertifiedCommit[RI]{
		cpocc,
	}, nil
}

func attributedSignedObservationsFromProtoMessage(pbasos []*AttributedSignedObservation) ([]protocol.AttributedSignedObservation, error) {
	asos := make([]protocol.AttributedSignedObservation, 0, len(pbasos))
	for _, pbaso := range pbasos {
		aso, err := attributedSignedObservationFromProtoMessage(pbaso)
		if err != nil {
			return nil, err
		}
		asos = append(asos, aso)
	}
	return asos, nil
}

func attributedSignedObservationFromProtoMessage(m *AttributedSignedObservation) (protocol.AttributedSignedObservation, error) {
	if m == nil {
		return protocol.AttributedSignedObservation{}, fmt.Errorf("unable to extract an AttributedSignedObservation value")
	}

	signedObservation, err := signedObservationFromProtoMessage(m.SignedObservation)
	if err != nil {
		return protocol.AttributedSignedObservation{}, err
	}

	return protocol.AttributedSignedObservation{
		signedObservation,
		commontypes.OracleID(m.Observer),
	}, nil
}

func signedObservationFromProtoMessage(m *SignedObservation) (protocol.SignedObservation, error) {
	if m == nil {
		return protocol.SignedObservation{}, fmt.Errorf("unable to extract an SignedObservation value")
	}

	return protocol.SignedObservation{
		m.Observation,
		m.Signature,
	}, nil
}
