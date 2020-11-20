package protocol

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var reportTypes = getReportTypes()


type AttributedObservation struct {
	Observation observation.Observation
	Observer    types.OracleID
}


func (o AttributedObservation) Equal(o2 AttributedObservation) bool {
	return (o.Observer == o2.Observer) && o.Observation.Equal(o2.Observation)
}

type AttributedObservations []AttributedObservation

func (aos AttributedObservations) Median() (observation.Observation, error) {
	if len(aos) == 0 {
		return observation.Observation{}, errors.Errorf(
			"can't take median of empty list")
	}
	return aos[len(aos)/2].Observation, nil
}

func (aos AttributedObservations) Equal(aos2 AttributedObservations) bool {
	if len(aos) != len(aos2) {
		return false
	}
	for i := range aos {
		if !aos[i].Equal(aos2[i]) {
			return false
		}
	}
	return true
}



func (aos AttributedObservations) observers() (rv common.Hash, err error) {
	if len(aos) > 32 {
		return rv, errors.Errorf("too many values! can only handle 32, got %d",
			len(aos))
	}
	for i, ao := range aos {
		id := int(ao.Observer)
		if id < 0 || i > 31 {
			return [32]byte{}, errors.Errorf(
				"Oracle index %d for %#+v is out of range", id, ao)
		}
		rv[i] = byte(uint8(id))
	}
	return rv, nil
}

func (aos AttributedObservations) onChainObservations() (rv []*big.Int) {
	for _, ao := range aos {
		rv = append(rv, ao.Observation.GoEthereumValue())
	}
	return rv
}




func (aos AttributedObservations) OnChainReport(repctx ReportContext) ([]byte, error) {
	observers, err := aos.observers()
	if err != nil {
		return nil, errors.Wrapf(err, "while collating observers for onChainReport")
	}
	return reportTypes.Pack(repctx.DomainSeparationTag(), observers, aos.onChainObservations())
}



type AttestedReportOne struct {
	AttributedObservations AttributedObservations
	Signature              []byte
}

func MakeAttestedReportOne(
	aos AttributedObservations,
	repctx ReportContext,
	signer func([]byte) ([]byte, error),
) (AttestedReportOne, error) {
	onchainReport, err := aos.OnChainReport(repctx)
	if err != nil {
		return AttestedReportOne{}, errors.Wrapf(err, "while serializing on-chain report")
	}
	sig, err := signer(onchainReport)
	if err != nil {
		return AttestedReportOne{}, errors.Wrapf(err, "while signing on-chain report")
	}

	return AttestedReportOne{aos, sig}, nil
}

func (rep AttestedReportOne) Equal(rep2 AttestedReportOne) bool {
	return rep.AttributedObservations.Equal(rep2.AttributedObservations) &&
		bytes.Equal(rep.Signature, rep2.Signature)
}



func (c *AttestedReportOne) Verify(repctx ReportContext, a types.OnChainSigningAddress) (err error) {
	report, err := c.AttributedObservations.OnChainReport(repctx)
	if err != nil {
		return err
	}
	var dummyID types.OracleID
	address := map[types.OnChainSigningAddress]types.OracleID{a: dummyID}
	_, err = signature.VerifyOnChain(report, c.Signature, address)
	return err
}




type AttestedReportMany struct {
	AttributedObservations AttributedObservations
	Signatures             [][]byte
}

func (rep AttestedReportMany) Equal(c2 AttestedReportMany) bool {
	if !rep.AttributedObservations.Equal(c2.AttributedObservations) ||
		len(rep.Signatures) != len(c2.Signatures) {
		return false
	}
	for i := range rep.Signatures {
		if !bytes.Equal(rep.Signatures[i], c2.Signatures[i]) {
			return false
		}
	}
	return true
}

func (rep *AttestedReportMany) collateSignatures() (rs, ss [][32]byte, vs [32]byte) {
	for i, sig := range rep.Signatures {
		rs = append(rs, common.BytesToHash(sig[:32]))
		ss = append(ss, common.BytesToHash(sig[32:64]))
		vs[i] = sig[64]
	}
	return rs, ss, vs
}

func (rep *AttestedReportMany) TransmissionArgs(repctx ReportContext) (report []byte, rs,
	ss [][32]byte, vs [32]byte, err error) {
	report, err = rep.AttributedObservations.OnChainReport(repctx)
	if err != nil {
		return nil, nil, nil, [32]byte{}, errors.Wrapf(err,
			"while constructing report for on-chain transmission")
	}
	rs, ss, vs = rep.collateSignatures()
	return report, rs, ss, vs, nil
}




func (rep *AttestedReportMany) VerifySignatures(
	repctx ReportContext,
	as signature.EthAddresses,
) error {
	report, err := rep.AttributedObservations.OnChainReport(repctx)
	if err != nil {
		return errors.Wrapf(err,
			"while serializing report to check signatures on it")
	}
	seen := make(map[types.OracleID]bool)
	for _, sig := range rep.Signatures {
		if oid, err := signature.VerifyOnChain(report, sig, as); err != nil {
			return errors.Wrapf(err,
				"while checking a signature on a report, 0x%x", sig)
		} else {
			if seen[oid] {
				return errors.Errorf("oracle #%d signed more than once", oid)
			}
			seen[oid] = true
		}
	}
	return nil
}

func getReportTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "rawReportContext", Type: mustNewType("bytes32")},
		{Name: "rawObservers", Type: mustNewType("bytes32")},
		{Name: "observations", Type: mustNewType("int192[]")},
	})
}
