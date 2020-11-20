package serialization

import (
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)


func Serialize(m protocol.Message) (b []byte, err error) {
	pmsg, err := toProtoMessage(m)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(pmsg)
}


func Deserialize(b []byte) (protocol.Message, error) {
	msgWrapper := &MessageWrapper{}
	err := proto.Unmarshal(b, msgWrapper)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal protobuf")
	}
	return msgWrapper.fromProtoMessage()
}

func toProtoMessage(m protocol.Message) (proto.Message, error) {
	msgWrapper := MessageWrapper{}
	switch v := m.(type) {
	case protocol.MessageNewEpoch:
		pm := &MessageNewEpoch{
			Epoch: uint64(v.Epoch),
		}
		msgWrapper.Msg = &MessageWrapper_MessageNewEpoch{pm}
	case protocol.MessageObserveReq:
		pm := &MessageObserveReq{
			Round: uint64(v.Round),
			Epoch: uint64(v.Epoch),
		}
		msgWrapper.Msg = &MessageWrapper_MessageObserveReq{pm}
	case protocol.MessageObserve:
		pm := &MessageObserve{
			Round:             uint64(v.Round),
			Epoch:             uint64(v.Epoch),
			SignedObservation: signedObservationToProtoMessage(v.SignedObservation),
		}
		msgWrapper.Msg = &MessageWrapper_MessageObserve{pm}
	case protocol.MessageReportReq:
		pm := &MessageReportReq{
			Round: uint64(v.Round),
			Epoch: uint64(v.Epoch),
		}
		for _, o := range v.AttributedSignedObservations {
			pm.AttributedSignedObservations = append(pm.AttributedSignedObservations,
				attributedSignedObservationToProtoMessage(o))
		}
		msgWrapper.Msg = &MessageWrapper_MessageReportReq{pm}
	case protocol.MessageReport:
		pm := &MessageReport{
			Epoch:  uint64(v.Epoch),
			Round:  uint64(v.Round),
			Report: attestedReportOneToProtoMessage(v.Report),
		}
		msgWrapper.Msg = &MessageWrapper_MessageReport{pm}
	case protocol.MessageFinal:
		msgWrapper.Msg = &MessageWrapper_MessageFinal{finalToProtoMessage(v)}
	case protocol.MessageFinalEcho:
		msgWrapper.Msg = &MessageWrapper_MessageFinalEcho{
			&MessageFinalEcho{Final: finalToProtoMessage(v.MessageFinal)},
		}
	default:
		return nil, errors.Errorf("Unable to serialize message of type %T", m)

	}
	return &msgWrapper, nil
}

func observationToProtoMessage(o observation.Observation) *Observation {
	return &Observation{Value: o.Marshal()}
}

func signedObservationToProtoMessage(o protocol.SignedObservation) *SignedObservation {
	sig := o.Signature
	if sig == nil {
		sig = []byte{}
	}
	return &SignedObservation{
		Observation: observationToProtoMessage(o.Observation),
		Signature:   sig,
	}
}

func attributedSignedObservationToProtoMessage(aso protocol.AttributedSignedObservation) *AttributedSignedObservation {
	return &AttributedSignedObservation{
		SignedObservation: signedObservationToProtoMessage(aso.SignedObservation),
		Observer:          uint32(aso.Observer),
	}
}

func (wrapper *MessageWrapper) fromProtoMessage() (protocol.Message, error) {
	switch msg := wrapper.Msg.(type) {
	case *MessageWrapper_MessageNewEpoch:
		return wrapper.GetMessageNewEpoch().fromProtoMessage()
	case *MessageWrapper_MessageObserveReq:
		return wrapper.GetMessageObserveReq().fromProtoMessage()
	case *MessageWrapper_MessageObserve:
		return wrapper.GetMessageObserve().fromProtoMessage()
	case *MessageWrapper_MessageReportReq:
		return wrapper.GetMessageReportReq().fromProtoMessage()
	case *MessageWrapper_MessageReport:
		return wrapper.GetMessageReport().fromProtoMessage()
	case *MessageWrapper_MessageFinal:
		return wrapper.GetMessageFinal().fromProtoMessage()
	case *MessageWrapper_MessageFinalEcho:
		return wrapper.GetMessageFinalEcho().fromProtoMessage()
	default:
		return nil, errors.Errorf("Unrecognised Msg type %T", msg)
	}
}

func (m *MessageNewEpoch) fromProtoMessage() (protocol.MessageNewEpoch, error) {
	if m == nil {
		return protocol.MessageNewEpoch{}, errors.New("Unable to extract a MessageNewEpoch value")
	}
	return protocol.MessageNewEpoch{
		Epoch: uint32(m.Epoch),
	}, nil
}

func (m *MessageObserveReq) fromProtoMessage() (protocol.MessageObserveReq, error) {
	if m == nil {
		return protocol.MessageObserveReq{}, errors.New("Unable to extract a MessageObserveReq value")
	}
	return protocol.MessageObserveReq{
		Round: uint8(m.Round),
		Epoch: uint32(m.Epoch),
	}, nil
}

func (m *MessageObserve) fromProtoMessage() (protocol.MessageObserve, error) {
	if m == nil {
		return protocol.MessageObserve{}, errors.New("Unable to extract a MessageObserve value")
	}
	so, err := m.SignedObservation.fromProtoMessage()
	if err != nil {
		return protocol.MessageObserve{}, nil
	}
	return protocol.MessageObserve{
		Epoch:             uint32(m.Epoch),
		Round:             uint8(m.Round),
		SignedObservation: so,
	}, nil
}

func (m *MessageReportReq) fromProtoMessage() (protocol.MessageReportReq, error) {
	if m == nil {
		return protocol.MessageReportReq{}, errors.New("Unable to extract a MessageReportReq value")
	}
	asos, err := AttributedSignedObservations(m.AttributedSignedObservations).fromProtoMessage()
	if err != nil {
		return protocol.MessageReportReq{}, err
	}
	return protocol.MessageReportReq{
		Round:                        uint8(m.Round),
		Epoch:                        uint32(m.Epoch),
		AttributedSignedObservations: asos,
	}, nil
}

func (o *Observation) fromProtoMessage() (observation.Observation, error) {
	if o == nil {
		return observation.Observation{}, errors.New("Unable to extract a Observation value")
	}
	obs, err := observation.UnmarshalObservation(o.Value)
	if err != nil {
		return observation.Observation{}, errors.Errorf(`could not deserialize bytes as `+
			`observation.Observation: "%v" from 0x%x`, err, o.Value)
	}
	return obs, nil
}

func (m *AttestedReportOne) fromProtoMessage() (protocol.AttestedReportOne, error) {
	if m == nil {
		return protocol.AttestedReportOne{}, errors.New("Unable to extract a AttestedReportOne value")
	}
	if m == nil {
		return protocol.AttestedReportOne{}, nil
	}
	aos := make([]protocol.AttributedObservation, len(m.AttributedObservations))
	for i, ao := range m.AttributedObservations {
		o, err := ao.Observation.fromProtoMessage()
		if err != nil {
			return protocol.AttestedReportOne{}, err
		}
		aos[i] = protocol.AttributedObservation{o, types.OracleID(ao.Observer)}
	}
	sig := m.Signature
	if sig == nil {
		sig = []byte{}
	}

	return protocol.AttestedReportOne{aos, sig}, nil
}

func (m *MessageReport) fromProtoMessage() (protocol.MessageReport, error) {
	if m == nil {
		return protocol.MessageReport{}, errors.New("Unable to extract a MessageReport value")
	}
	report, err := m.Report.fromProtoMessage()
	if err != nil {
		return protocol.MessageReport{}, err
	}

	return protocol.MessageReport{uint32(m.Epoch), uint8(m.Round), report}, nil
}

func (m *AttestedReportMany) fromProtoMessage() (protocol.AttestedReportMany, error) {
	if m == nil {
		return protocol.AttestedReportMany{}, errors.New("Unable to extract a AttestedReportMany value")
	}
	signatures := make([][]byte, len(m.Signatures))
	for i, sig := range m.Signatures {
		if sig == nil {
			sig = []byte{}
		}
		signatures[i] = sig
	}
	aos := protocol.AttributedObservations{}
	for _, v := range m.AttributedObservations {
		obs, err := v.Observation.fromProtoMessage()
		if err != nil {
			return protocol.AttestedReportMany{}, err
		}
		aos = append(aos, protocol.AttributedObservation{obs, types.OracleID(v.Observer)})
	}

	return protocol.AttestedReportMany{aos, signatures}, nil
}

func (m *MessageFinal) fromProtoMessage() (protocol.MessageFinal, error) {
	if m == nil {
		return protocol.MessageFinal{}, errors.New("Unable to extract a MessageFinal value")
	}
	report, err := m.Report.fromProtoMessage()
	if err != nil {
		return protocol.MessageFinal{}, nil
	}
	return protocol.MessageFinal{uint32(m.Epoch), uint8(m.Round), report}, nil
}

func (m *MessageFinalEcho) fromProtoMessage() (protocol.MessageFinalEcho, error) {
	if m == nil {
		return protocol.MessageFinalEcho{}, errors.New("Unable to extract a MessageFinalEcho value")
	}
	final, err := m.Final.fromProtoMessage()
	if err != nil {
		return protocol.MessageFinalEcho{}, err
	}
	return protocol.MessageFinalEcho{MessageFinal: final}, nil
}


type AttributedSignedObservations []*AttributedSignedObservation

func (ms AttributedSignedObservations) fromProtoMessage() ([]protocol.AttributedSignedObservation, error) {
	if ms == nil {
		
		
		return []protocol.AttributedSignedObservation{}, nil
	}
	observations := make([]protocol.AttributedSignedObservation, len(ms))
	for i, o := range ms {
		obs, err := o.fromProtoMessage()
		if err != nil {
			return nil, err
		}
		observations[i] = obs
	}
	return observations, nil
}

func (m *AttributedSignedObservation) fromProtoMessage() (protocol.AttributedSignedObservation, error) {
	if m == nil {
		return protocol.AttributedSignedObservation{}, errors.New("Unable to extract an AttributedSignedObservation value")
	}

	signedObservation, err := m.SignedObservation.fromProtoMessage()
	if err != nil {
		return protocol.AttributedSignedObservation{}, err
	}
	return protocol.AttributedSignedObservation{
		signedObservation,
		types.OracleID(m.Observer),
	}, nil
}


type SignedObservations []*SignedObservation

func (ms SignedObservations) fromProtoMessage() ([]protocol.SignedObservation, error) {
	if ms == nil {
		
		
		return []protocol.SignedObservation{}, nil
	}
	observations := make([]protocol.SignedObservation, len(ms))
	for i, o := range ms {
		obs, err := o.fromProtoMessage()
		if err != nil {
			return nil, err
		}
		observations[i] = obs
	}
	return observations, nil
}

func (m *SignedObservation) fromProtoMessage() (protocol.SignedObservation, error) {
	if m == nil {
		return protocol.SignedObservation{}, errors.New("Unable to extract an SignedObservation value")
	}
	sig := m.Signature
	if sig == nil {
		sig = []byte{}
	}
	obs, err := m.Observation.fromProtoMessage()
	if err != nil {
		return protocol.SignedObservation{}, err
	}
	return protocol.SignedObservation{obs, sig}, nil
}

func attestedReportOneToProtoMessage(v protocol.AttestedReportOne) *AttestedReportOne {
	sig := v.Signature
	if sig == nil {
		sig = []byte{}
	}
	pm := &AttestedReportOne{
		AttributedObservations: make([]*AttributedObservation, len(v.AttributedObservations)),
		Signature:              sig,
	}
	for i, val := range v.AttributedObservations {
		pm.AttributedObservations[i] = &AttributedObservation{
			Observation: &Observation{Value: val.Observation.Marshal()},
			Observer:    uint32(val.Observer),
		}
	}
	return pm
}

func attributedObservationsToProtoMessage(aos protocol.AttributedObservations) []*AttributedObservation {
	result := []*AttributedObservation{}
	for _, ao := range aos {
		result = append(result, &AttributedObservation{
			Observation: &Observation{Value: ao.Observation.Marshal()},
			Observer:    uint32(ao.Observer),
		})
	}

	return result
}

func finalToProtoMessage(v protocol.MessageFinal) *MessageFinal {
	pm := &MessageFinal{
		Epoch: uint64(v.Epoch),
		Round: uint64(v.Round),
		Report: &AttestedReportMany{
			AttributedObservations: attributedObservationsToProtoMessage(v.Report.AttributedObservations),
			Signatures:             make([][]byte, len(v.Report.Signatures)),
		},
	}
	for i, sig := range v.Report.Signatures {
		pm.Report.Signatures[i] = sig
	}
	return pm
}
