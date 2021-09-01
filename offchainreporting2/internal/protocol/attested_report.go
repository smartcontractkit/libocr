package protocol

import (
	"bytes"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/pkg/errors"
)

// AttributedObservation succinctly atrributes a value reported to an oracle
type AttributedObservation struct {
	Observation types.Observation
	Observer    commontypes.OracleID
}

func (o AttributedObservation) Equal(o2 AttributedObservation) bool {
	return (o.Observer == o2.Observer) && bytes.Equal(o.Observation, o2.Observation)
}

type AttributedObservations []AttributedObservation

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

// AttestedReportOne is the collated report oracles sign off on, after they've
// verified the individual signatures in a report-req sent by the current leader
type AttestedReportOne struct {
	Skip      bool
	Report    types.Report
	Signature []byte
}

func MakeAttestedReportOneSkip() AttestedReportOne {
	return AttestedReportOne{true, nil, nil}
}

func MakeAttestedReportOneNoskip(
	repctx types.ReportContext,
	report types.Report,
	signer func(types.ReportContext, types.Report) ([]byte, error),
) (AttestedReportOne, error) {
	sig, err := signer(repctx, report)
	if err != nil {
		return AttestedReportOne{}, errors.Wrapf(err, "while signing on-chain report")
	}

	return AttestedReportOne{false, report, sig}, nil
}

func (rep AttestedReportOne) Equal(rep2 AttestedReportOne) bool {
	return rep.Skip == rep2.Skip && bytes.Equal(rep.Report, rep2.Report) && bytes.Equal(rep.Signature, rep2.Signature)
}

func (rep AttestedReportOne) EqualExceptSignature(rep2 AttestedReportOne) bool {
	return rep.Skip == rep2.Skip && bytes.Equal(rep.Report, rep2.Report)
}

// Verify is used by the leader to check the signature a process attaches to its
// report message (the c.Sig value.)
func (aro *AttestedReportOne) Verify(contractSigner types.OnchainKeyring, publicKey types.OnchainPublicKey, repctx types.ReportContext) (err error) {
	if aro.Skip {
		if len(aro.Report) != 0 || len(aro.Signature) != 0 {
			return fmt.Errorf("AttestedReportOne with Skip=true has non-empty Report or Signature")
		}
	} else {
		ok := contractSigner.Verify(publicKey, repctx, aro.Report, aro.Signature)
		if !ok {
			return fmt.Errorf("failed to verify signature on AttestedReportOne")
		}
	}
	return nil
}

// AttestedReportMany consists of attributed observations with aggregated
// signatures from the oracles which have sent this report to the current epoch
// leader.
//

type AttestedReportMany struct {
	ReportData           []byte
	AttributedSignatures []types.AttributedOnChainSignature
}

func (rep AttestedReportMany) Equal(c2 AttestedReportMany) bool {
	if !bytes.Equal(rep.ReportData, c2.ReportData) {
		return false
	}

	if len(rep.AttributedSignatures) != len(c2.AttributedSignatures) {
		return false
	}

	for i := range rep.AttributedSignatures {
		if !rep.AttributedSignatures[i].Equal(c2.AttributedSignatures[i]) {
			return false
		}
	}

	return true
}

// VerifySignatures checks that all the signatures (c.Signatures) come from the
// addresses in the map "as", and returns a list of which oracles they came
// from.
func (rep *AttestedReportMany) VerifySignatures(
	numSignatures int,
	contractSigner types.OnchainKeyring,
	oracleIdentities []config.OracleIdentity,
	repctx types.ReportContext,
) error {
	if numSignatures != len(rep.AttributedSignatures) {
		return fmt.Errorf("wrong number of signatures, expected %v and got %v", numSignatures, len(rep.AttributedSignatures))
	}
	seen := make(map[commontypes.OracleID]bool)
	for i, sig := range rep.AttributedSignatures {
		if seen[sig.Signer] {
			return fmt.Errorf("duplicate Signature by %v", sig.Signer)
		}
		seen[sig.Signer] = true
		if len(oracleIdentities) <= int(sig.Signer) {
			return fmt.Errorf("signer out of bounds: %v", sig.Signer)
		}
		if !contractSigner.Verify(oracleIdentities[sig.Signer].OnChainPublicKey, repctx, rep.ReportData, sig.Signature) {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify", i, sig.Signer, oracleIdentities[sig.Signer].OnChainPublicKey)
		}
	}
	return nil
}
