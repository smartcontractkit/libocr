package protocol

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"errors"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type SignedObservation struct {
	Observation types.Observation
	Signature   []byte
}

func MakeSignedObservation(
	repts types.ReportTimestamp,
	query types.Query,
	observation types.Observation,
	signer func(msg []byte) (sig []byte, err error),
) (
	SignedObservation,
	error,
) {
	payload := signedObservationWireMessage(repts, query, observation)
	sig, err := signer(payload)
	if err != nil {
		return SignedObservation{}, err
	}
	return SignedObservation{observation, sig}, nil
}

func (so SignedObservation) Equal(so2 SignedObservation) bool {
	return bytes.Equal(so.Observation, so2.Observation) &&
		bytes.Equal(so.Signature, so2.Signature)
}

func (so SignedObservation) Verify(repts types.ReportTimestamp, query types.Query, publicKey types.OffchainPublicKey) error {
	ok := ed25519.Verify(ed25519.PublicKey(publicKey), signedObservationWireMessage(repts, query, so.Observation), so.Signature)
	if !ok {
		return errors.New("SignedObservation has invalid signature")
	}

	return nil
}

func signedObservationWireMessage(repts types.ReportTimestamp, query types.Query, observation types.Observation) []byte {
	h := sha256.New()
	// ConfigDigest
	h.Write(repts.ConfigDigest[:])
	binary.Write(h, binary.BigEndian, repts.Epoch)
	h.Write([]byte{repts.Round})

	// Query
	binary.Write(h, binary.BigEndian, uint64(len(query)))
	h.Write(query)

	// Observation
	binary.Write(h, binary.BigEndian, uint64(len(observation)))
	h.Write(observation)

	return h.Sum(nil)
}

type AttributedSignedObservation struct {
	SignedObservation SignedObservation
	Observer          commontypes.OracleID
}

func (aso AttributedSignedObservation) Equal(aso2 AttributedSignedObservation) bool {
	return aso.SignedObservation.Equal(aso2.SignedObservation) &&
		aso.Observer == aso2.Observer
}
