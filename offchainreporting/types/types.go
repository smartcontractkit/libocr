// Package types contains the types and interfaces a consumer of the OCR library needs to be aware of
package types

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/curve25519"
)

// OracleID is an index over the oracles, used as a succinct attribution to an
// oracle in communication with the on-chain contract. It is not a cryptographic
// commitment to the oracle's private key, like a public key is.
type OracleID int

type ConfigDigest [16]byte

func (c ConfigDigest) Hex() string {
	return fmt.Sprintf("%x", c[:])
}

func BytesToConfigDigest(b []byte) (ConfigDigest, error) {
	configDigest := ConfigDigest{}

	if len(b) != len(configDigest) {
		return ConfigDigest{}, fmt.Errorf("Cannot convert bytes to ConfigDigest. bytes have wrong length %v", len(b))
	}

	if n := copy(configDigest[:], b); n != len(configDigest) {
		panic("copy returned wrong length")
	}

	return configDigest, nil
}

// BinaryNetworkEndpoint contains the network methods a consumer must implement
// SendTo and Broadcast must not block. They should buffer messages and
// (optionally) drop the oldest buffered messages if the buffer reaches capacity.
//
// The protocol trusts the sender in BinaryMessageWithSender. Implementors of
// this interface are responsible for securely authenticating that messages come
// from their indicated senders.
//
// All its functions should be thread-safe.
type BinaryNetworkEndpoint interface {
	// SendTo(msg, to) sends msg to "to"
	SendTo(payload []byte, to OracleID)
	// Broadcast(msg) sends msg to all oracles
	Broadcast(payload []byte)
	// Receive returns channel which carries all messages sent to this oracle.
	Receive() <-chan BinaryMessageWithSender
	// Start starts the endpoint
	Start() error
	// Close stops the endpoint
	Close() error
}

// Bootstrapper helps nodes find each other on the network level by providing
// peer-discovery services.
//
// All its functions should be thread-safe.
type Bootstrapper interface {
	Start() error
	Close() error
}

// BinaryNetworkEndpointFactory creates permissioned BinaryNetworkEndpoints.
//
// All its functions should be thread-safe.
type BinaryNetworkEndpointFactory interface {
	MakeEndpoint(cd ConfigDigest, peerIDs []string, bootstrappers []string,
		failureThreshold int, tokenBucketRefillRate float64, tokenBucketSize int,
	) (BinaryNetworkEndpoint, error)
	PeerID() string
}

// BootstrapperFactory creates permissioned Bootstrappers.
//
// All its functions should be thread-safe.
type BootstrapperFactory interface {
	MakeBootstrapper(cd ConfigDigest, peerIDs []string, bootstrappers []string, failureThreshold int) (Bootstrapper, error)
}

// BinaryMessageWithSender contains the information from a Receive() channel
// message: The binary representation of the message, and the ID of its sender.
type BinaryMessageWithSender struct {
	Msg    []byte
	Sender OracleID
}

// Observation is the type returned by the DataSource.Observe method. Represents
// an int192 at time of writing
type Observation *big.Int

// DataSource implementations must be thread-safe. Observe may be called by many different threads concurrently.
type DataSource interface {
	// Observe queries the data source. Returns a value or an error.
	// Must not block indefinitely.
	Observe(context.Context) (Observation, error)
}

// MonitoringEndpoint is where the OCR protocol sends monitoring output
//
// All its functions should be thread-safe.
type MonitoringEndpoint interface {
	SendLog(log []byte)
}

// ContractTransmitter sends new reports to the OffchainAggregator smart contract.
//
// All its functions should be thread-safe.
type ContractTransmitter interface {

	// Transmit sends the report to the on-chain OffchainAggregator smart contract's Transmit method
	Transmit(
		ctx context.Context,
		report []byte, // wire-formatted report to transmit on-chain
		rs, ss [][32]byte, vs [32]byte, // Signatures; i'th elt's are i'th (v,r,s)
	) error

	LatestTransmissionDetails(
		ctx context.Context,
	) (
		configDigest ConfigDigest,
		epoch uint32,
		round uint8,
		latestAnswer Observation,
		latestTimestamp time.Time,
		err error,
	)

	FromAddress() common.Address
}

// ContractConfigTracker tracks OffchainAggregator.ConfigSet events emitted from blockchain.
//
// All its functions should be thread-safe.
type ContractConfigTracker interface {
	SubscribeToNewConfigs(ctx context.Context) (ContractConfigSubscription, error)
	LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ConfigDigest, err error)
	ConfigFromLogs(ctx context.Context, changedInBlock uint64) (ContractConfig, error)

	// LatestBlockHeight returns the height of the most recent block in the chain.
	LatestBlockHeight(ctx context.Context) (blockheight uint64, err error)
}

type ContractConfigSubscription interface {
	// May be closed by sender at any time
	Configs() <-chan ContractConfig
	Close()
}

type ContractConfig struct {
	ConfigDigest         ConfigDigest
	Signers              []common.Address // TODO: use OnChainSigningAddress?
	Transmitters         []common.Address
	Threshold            uint8
	EncodedConfigVersion uint64
	Encoded              []byte
}

// OffChainPublicKey is the public key used to cryptographically identify an
// oracle in inter-oracle communications.
type OffchainPublicKey ed25519.PublicKey

// OnChainSigningAddress is the public key used to cryptographically identify an
// oracle to the on-chain smart contract.
type OnChainSigningAddress common.Address

// SharedSecretEncryptionPublicKey is the public key used to receive an encrypted
// version of the secret shared amongst all oracles on a common contract.
type SharedSecretEncryptionPublicKey [curve25519.PointSize]byte // X25519

// PrivateKeys contains the secret keys needed for the OCR protocol, and methods
// which use those keys without exposing them to the rest of the application.
// There are three key pairs to track, here:
//
// - The on-chain signing key pair (secp256k1), used to sign contract reports
//
// - The off-chain key signing key pair (Ed25519), used to sign observations
//
// - The config encryption key (X25519), used to decrypt the symmetric
// key which encrypts the offchain configuration data passed through the OffchainAggregator
// smart contract.
//
// All its functions should be thread-safe.
type PrivateKeys interface {

	// SignOnChain returns an ethereum-style ECDSA secp256k1 signature on msg. See
	// signature.OnChainPrivateKey.Sign for the logic it needs to implement
	SignOnChain(msg []byte) (signature []byte, err error)

	// SignOffChain returns an EdDSA-Ed25519 signature on msg. See
	// signature.OffChainPrivateKey.Sign for the logic it needs to implement
	SignOffChain(msg []byte) (signature []byte, err error)

	// ConfigDiffieHellman multiplies base, as a representative of a Curve 25519
	// point, by a secret scalar, which is also the scalar to multiply
	// curve25519.BasePoint to, in order to get PublicKeyConfig
	ConfigDiffieHellman(base *[curve25519.ScalarSize]byte) (sharedPoint *[curve25519.PointSize]byte, err error)

	// PublicKeyAddressOnChain returns the address corresponding to the
	// public component of the keypair used in SignOnChain
	PublicKeyAddressOnChain() OnChainSigningAddress

	// PublicKeyOffChain returns the pbulic component of the keypair used in SignOffChain
	PublicKeyOffChain() OffchainPublicKey

	// PublicKeyConfig returns the public component of the keypair used in ConfigKeyShare
	PublicKeyConfig() [curve25519.PointSize]byte
}
