package tracing

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type OffchainConfigDigester struct {
	oracleID OracleID
	tracer   *Tracer
	backend  types.OffchainConfigDigester
}

var _ types.OffchainConfigDigester = (*OffchainConfigDigester)(nil)

func MakeOffchainConfigDigester(tracer *Tracer, oracleID OracleID, backend types.OffchainConfigDigester) *OffchainConfigDigester {
	return &OffchainConfigDigester{
		oracleID,
		tracer,
		backend,
	}
}

func (o *OffchainConfigDigester) ConfigDigest(cfg types.ContractConfig) (types.ConfigDigest, error) {
	digest, err := o.backend.ConfigDigest(cfg)
	o.tracer.Append(NewConfigDigest(o.oracleID, cfg, digest, err))
	return digest, err
}

func (o *OffchainConfigDigester) ConfigDigestPrefix() types.ConfigDigestPrefix {
	prefix := o.backend.ConfigDigestPrefix()
	o.tracer.Append(NewConfigDigestPrefix(o.oracleID, prefix))
	return prefix
}
