package types

// ConfigDigestPrefix acts as a domain separator between different (typically
// chain-specific) methods of computing a ConfigDigest.
type ConfigDigestPrefix uint16

// This acts as the canonical "registry" of ConfigDigestPrefixes. Pick an unused
// prefix and add it to this list before you build an OffchainConfigDigester for
// whatever chain you're targeting.
const (
	_                     ConfigDigestPrefix = 0
	ConfigDigestPrefixEVM                    = 1
	_                     ConfigDigestPrefix = 0xFFFF
)

// An OffchainConfigDigester computes a ConfigDigest the same way as the
// contract, but *offchain*. This is used to ensure that the ConfigDigest
// returned from the contract was computed correctly and to prevent a malicious
// blockchain node from breaking domain separation between different protocol
// instances.
//
// All its functions should be thread-safe.
type OffchainConfigDigester interface {
	// Compute ConfigDigest for the given ContractConfig. The first two bytes of the
	// ConfigDigest must be the big-endian encoding of ConfigDigestPrefix!
	ConfigDigest(ContractConfig) (ConfigDigest, error)

	// This should return the same constant value on every invocation
	ConfigDigestPrefix() ConfigDigestPrefix
}
