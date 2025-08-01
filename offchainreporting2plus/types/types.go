// Package types contains the types and interfaces a consumer of the OCR library needs to be aware of
package types

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"golang.org/x/crypto/curve25519"
)

type BinaryNetworkEndpointLimits struct {
	MaxMessageLength          int
	MessagesRatePerOracle     float64
	MessagesCapacityPerOracle int
	BytesRatePerOracle        float64
	BytesCapacityPerOracle    int
}

// 2x one per priority
type BinaryNetworkEndpoint2Config struct {
	BinaryNetworkEndpointLimits

	// Buffer sizes specified below override the values set in PeerConfig.
	OverrideIncomingMessageBufferSize int
	OverrideOutgoingMessageBufferSize int
}

// BinaryNetworkEndpointFactory creates permissioned BinaryNetworkEndpoint instances.
//
// All its functions should be thread-safe.
type BinaryNetworkEndpointFactory interface {
	// f is a remnant of P2Pv1 and is ignored.
	NewEndpoint(
		cd ConfigDigest,
		peerIDs []string,
		v2bootstrappers []commontypes.BootstrapperLocator,
		f int,
		limits BinaryNetworkEndpointLimits,
	) (commontypes.BinaryNetworkEndpoint, error)
	PeerID() string
}

// BinaryNetworkEndpoint2Factory creates permissioned BinaryNetworkEndpoint2 instances.
//
// All its functions should be thread-safe.
type BinaryNetworkEndpoint2Factory interface {
	NewEndpoint(
		cd ConfigDigest,
		peerIDs []string,
		v2bootstrappers []commontypes.BootstrapperLocator,
		defaultPriorityConfig BinaryNetworkEndpoint2Config,
		lowPriorityConfig BinaryNetworkEndpoint2Config,
	) (BinaryNetworkEndpoint2, error)
	PeerID() string
}

// BootstrapperFactory creates permissioned Bootstrappers.
//
// All its functions should be thread-safe.
type BootstrapperFactory interface {
	// f is a remnant of P2Pv1 and is ignored.
	NewBootstrapper(cd ConfigDigest, peerIDs []string,
		v2bootstrappers []commontypes.BootstrapperLocator,
		f int,
	) (commontypes.Bootstrapper, error)
}

type Query []byte

type AttributedQuery struct {
	Query    Query
	Proposer commontypes.OracleID
}

type Observation []byte

type AttributedObservation struct {
	Observation Observation
	Observer    commontypes.OracleID
}

func (ao AttributedObservation) Equal(other AttributedObservation) bool {
	return bytes.Equal(ao.Observation, other.Observation) && ao.Observer == other.Observer
}

// ReportTimestamp is the logical timestamp of a report.
type ReportTimestamp struct {
	ConfigDigest ConfigDigest
	Epoch        uint32
	Round        uint8
}

// ReportContext is the contextual data sent to contract along with the report
// itself.
type ReportContext struct {
	ReportTimestamp
	// A hash over some data that is exchanged during execution of the offchain
	// protocol. The data itself is not needed onchain, but we still want to
	// include it in the signature that goes onchain.
	ExtraHash [32]byte
}

type Report []byte

type AttributedOnchainSignature struct {
	Signature []byte
	Signer    commontypes.OracleID
}

func (as AttributedOnchainSignature) Equal(other AttributedOnchainSignature) bool {
	return bytes.Equal(as.Signature, other.Signature) && as.Signer == other.Signer
}

type ReportingPluginFactory interface {
	// Creates a new reporting plugin instance. The instance may have
	// associated goroutines or hold system resources, which should be
	// released when its Close() function is called.
	NewReportingPlugin(context.Context, ReportingPluginConfig) (ReportingPlugin, ReportingPluginInfo, error)
}

type ReportingPluginConfig struct {
	ConfigDigest ConfigDigest

	// OracleID (index) of the oracle executing this ReportingPlugin instance.
	OracleID commontypes.OracleID

	// N is the total number of nodes.
	N int

	// F is an upper bound on the number of faulty nodes.
	F int

	// Encoded configuration for the contract
	OnchainConfig []byte

	// Encoded configuration for the ReportingPlugin disseminated through the
	// contract. This value is only passed through the contract, but otherwise
	// ignored by it.
	OffchainConfig []byte

	// Estimate of the duration between rounds. You should not rely on this
	// value being accurate. Rounds might occur more or less frequently than
	// estimated.
	//
	// This value is intended for estimating the load incurred by a
	// ReportingPlugin before running it and for configuring caches.
	EstimatedRoundInterval time.Duration

	// Maximum duration the ReportingPlugin's functions are allowed to take
	MaxDurationQuery                        time.Duration
	MaxDurationObservation                  time.Duration
	MaxDurationReport                       time.Duration
	MaxDurationShouldAcceptFinalizedReport  time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration
}

// A ReportingPlugin allows plugging custom logic into the OCR protocol. The OCR
// protocol handles cryptography, networking, ensuring that a sufficient number
// of nodes is in agreement about any report, transmitting the report to the
// contract, etc... The ReportingPlugin handles application-specific logic. To
// do so, the ReportingPlugin defines a number of callbacks that are called by
// the OCR protocol logic at certain points in the protocol's execution flow.
// The report generated by the ReportingPlugin must be in a format understood by
// contract that the reports are transmitted to.
//
// Roughly speaking, the protocol works as follows: A designated leader (fixed
// per epoch) broadcasts a request for observations containing an
// application-specific query to all followers. (Note the leader is also a
// follower of itself, i.e. one node acts as leader and all nodes act as
// followers.) Followers make signed observations and send them back to the
// leader. The leader collects the observations and broadcasts a report request
// containing the original query and the collected signed observations to all
// followers. At this stage, followers decide whether a report should be
// created. If followers decide not to create a report, the round ends here.
// Otherwise, followers construct a signed report (ultimately destined for the
// target contract) and send it to the leader. The leader collects the signed
// report(s) and sends out a final message containing a report together with a
// sufficient number of signatures. The followers echo the final message amongst
// each other prevent a malicious leader from selectively excluding particular
// nodes. Each follower independently decides whether it wishes to accept the
// final report for transmission. If it accepts it for transmission, the
// follower will wait for some time (according to a shared schedule) before
// attempting to transmit the report. After this time, the follower will check
// one last time whether to broadcast the transmit transaction before sending
// it. (Due to its brevity, this description skips over lots of details and edge
// cases. This is just to give a rough idea of the protocol flow.)
//
// We assume that each correct node participating in the protocol instance will
// be running the same ReportingPlugin implementation. However, not all nodes
// may be correct; up to f nodes be faulty in arbitrary ways (aka byzantine
// faults). For example, faulty nodes could be down, have intermittent
// connectivity issues, send garbage messages, or be controlled by an adversary.
//
// For a protocol round where everything is working correctly, the leader will
// start by invoking Query. Followers will call Observation and Report. If a
// sufficient number of followers agree on a report, ShouldAcceptFinalizedReport
// will be called as well. If ShouldAcceptFinalizedReport returns true,
// ShouldTransmitAcceptedReport will be called. However, a ReportingPlugin must
// also correctly handle the case where faults occur.
//
// In particular, a ReportingPlugin must deal with cases where:
//
// - only a subset of the functions on the ReportingPlugin are invoked for a
// given round
//
// - an arbitrary number of epochs and rounds has been skipped between
// invocations of the ReportingPlugin
//
// - the observation returned by Observation is not included in the list of
// AttributedObservations passed to Report
//
// - a query or observation is malformed. (For defense in depth, it is also
// strongly recommended that malformed reports are handled gracefully.)
//
// - instances of the ReportingPlugin run by different oracles have different
// call traces. E.g., the ReportingPlugin's Observation function may have been
// invoked on node A, but not on node B. Or Observation may have been invoked on
// A and B, but with different queries.
//
// All functions on a ReportingPlugin should be thread-safe.
//
// All functions that take a context as their first argument may still do cheap
// computations after the context expires, but should stop any blocking
// interactions with outside services (APIs, database, ...) and return as
// quickly as possible. (Rough rule of thumb: any such computation should not
// take longer than a few ms.) A blocking function may block execution of the
// entire protocol instance!
//
// For a given OCR2 protocol instance, there can be many (consecutive) instances
// of a ReportingPlugin, e.g. due to software restarts. If you need
// ReportingPlugin state to survive across restarts, you should persist it. A
// ReportingPlugin instance will only ever serve a single protocol instance.
// When we talk about "instance" below, we typically mean ReportingPlugin
// instances, not protocol instances.
type ReportingPlugin interface {
	// Query creates a Query that is sent from the leader to all follower nodes
	// as part of the request for an observation. Be careful! A malicious leader
	// could equivocate (i.e. send different queries to different followers.)
	// Many applications will likely be better off always using an empty query
	// if the oracles don't need to coordinate on what to observe (e.g. in case
	// of a price feed) or the underlying data source offers an (eventually)
	// consistent view to different oracles (e.g. in case of observing a
	// blockchain).
	//
	// You may assume that the sequence of epochs and the sequence of rounds
	// within an epoch are strictly monotonically increasing during the lifetime
	// of an instance of this interface.
	Query(context.Context, ReportTimestamp) (Query, error)

	// Observation gets an observation from the underlying data source. Returns
	// a value or an error.
	//
	// You may assume that the sequence of epochs and the sequence of rounds
	// within an epoch are strictly monotonically increasing during the lifetime
	// of an instance of this interface.
	Observation(context.Context, ReportTimestamp, Query) (Observation, error)

	// Decides whether a report (destined for the contract) should be generated
	// in this round. If yes, also constructs the report.
	//
	// You may assume that the sequence of epochs and the sequence of rounds
	// within an epoch are strictly monotonically increasing during the lifetime
	// of an instance of this interface. This function will always be called
	// with at least 2f+1 AttributedObservations from distinct oracles.
	Report(context.Context, ReportTimestamp, Query, []AttributedObservation) (bool, Report, error)

	// Decides whether a report should be accepted for transmission. Any report
	// passed to this function will have been signed by a quorum of oracles.
	//
	// Don't make assumptions about the epoch/round order in which this function
	// is called.
	ShouldAcceptFinalizedReport(context.Context, ReportTimestamp, Report) (bool, error)

	// Decides whether the given report should actually be broadcast to the
	// contract. This is invoked just before the broadcast occurs. Any report
	// passed to this function will have been signed by a quorum of oracles and
	// been accepted by ShouldAcceptFinalizedReport.
	//
	// Don't make assumptions about the epoch/round order in which this function
	// is called.
	//
	// As mentioned above, you should gracefully handle only a subset of a
	// ReportingPlugin's functions being invoked for a given report. For
	// example, due to reloading persisted pending transmissions from the
	// database upon oracle restart, this function  may be called with reports
	// that no other function of this instance of this interface has ever
	// been invoked on.
	ShouldTransmitAcceptedReport(context.Context, ReportTimestamp, Report) (bool, error)

	// If Close is called a second time, it may return an error but must not
	// panic. This will always be called when a ReportingPlugin is no longer
	// needed, e.g. on shutdown of the protocol instance or shutdown of the
	// oracle node. This will only be called after any calls to other functions
	// of the ReportingPlugin will have completed.
	Close() error
}

const (
	twoHundredFiftySixMiB   = 256 * 1024 * 1024     // 256 MiB
	MaxMaxQueryLength       = twoHundredFiftySixMiB // 256 MiB
	MaxMaxObservationLength = twoHundredFiftySixMiB // 256 MiB
	MaxMaxReportLength      = twoHundredFiftySixMiB // 256 MiB
)

// Limits for data returned by the ReportingPlugin.
// Used for computing rate limits and defending against outsized messages.
// Messages are checked against these values during (de)serialization. Be
// careful when changing these values, they could lead to different versions
// of a ReportingPlugin being unable to communicate with each other.
type ReportingPluginLimits struct {
	MaxQueryLength       int
	MaxObservationLength int
	MaxReportLength      int
}

type ReportingPluginInfo struct {
	// Used for debugging purposes.
	Name string

	// If true, quorum requirements are adjusted so that only a single report
	// will reach a quorum of signatures for any (epoch, round) tuple.
	UniqueReports bool

	Limits ReportingPluginLimits
}

// Account is a human-readable account identifier, e.g. an Ethereum address
type Account string

// ContractTransmitter sends new reports to the OCR2Aggregator smart contract.
//
// All its functions should be thread-safe.
type ContractTransmitter interface {

	// Transmit sends the report to the on-chain OCR2Aggregator smart
	// contract's Transmit method.
	//
	// In most cases, implementations of this function should store the
	// transmission in a queue/database/..., but perform the actual
	// transmission (and potentially confirmation) of the transaction
	// asynchronously.
	Transmit(
		context.Context,
		ReportContext,
		Report,
		[]AttributedOnchainSignature,
	) error

	// LatestConfigDigestAndEpoch returns the logically latest configDigest and
	// epoch for which a report was successfully transmitted.
	LatestConfigDigestAndEpoch(
		context.Context,
	) (
		configDigest ConfigDigest,
		epoch uint32,
		err error,
	)

	// Account from which the transmitter invokes the contract
	FromAccount(context.Context) (Account, error)
}

// ContractConfigTracker tracks configuration changes of the OCR contract
// (on-chain).
//
// All its functions should be thread-safe.
type ContractConfigTracker interface {
	// Notify may optionally emit notification events when the contract's
	// configuration changes. This is purely used as an optimization reducing
	// the delay between a configuration change and its enactment. Implementors
	// who don't care about this may simply return a nil channel.
	//
	// The returned channel should never be closed.
	Notify() <-chan struct{}

	// LatestConfigDetails returns information about the latest configuration,
	// but not the configuration itself.
	LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ConfigDigest, err error)

	// LatestConfig returns the latest configuration.
	LatestConfig(ctx context.Context, changedInBlock uint64) (ContractConfig, error)

	// LatestBlockHeight returns the height of the most recent block in the chain.
	LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error)
}

type ContractConfig struct {
	ConfigDigest          ConfigDigest
	ConfigCount           uint64
	Signers               []OnchainPublicKey
	Transmitters          []Account
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

// OffchainPublicKey is the public key used to cryptographically identify an
// oracle in inter-oracle communications.
type OffchainPublicKey [ed25519.PublicKeySize]byte

// OnchainPublicKey is the public key used to cryptographically identify an
// oracle to the on-chain smart contract.
type OnchainPublicKey []byte

// ConfigEncryptionPublicKey is the public key used to receive an encrypted
// version of the secret shared amongst all oracles on a common contract.
type ConfigEncryptionPublicKey [curve25519.PointSize]byte // X25519

// OffchainKeyring contains the secret keys needed for the OCR protocol, and methods
// which use those keys without exposing them to the rest of the application.
// There are two key pairs to track, here:
//
// First, the off-chain key signing key pair (Ed25519), used to sign observations.
//
// Second, the config encryption key (X25519), used to decrypt the symmetric
// key which encrypts the offchain configuration data passed through the OCR2Aggregator
// smart contract.
//
// All its functions should be thread-safe.
type OffchainKeyring interface {
	// OffchainSign returns an EdDSA-Ed25519 signature on msg produced using the
	// standard library's ed25519.Sign function.
	OffchainSign(msg []byte) (signature []byte, err error)

	// ConfigDiffieHellman multiplies point with the secret key (i.e. scalar)
	// that ConfigEncryptionPublicKey corresponds to.
	ConfigDiffieHellman(point [curve25519.PointSize]byte) (sharedPoint [curve25519.PointSize]byte, err error)

	// OffchainPublicKey returns the public component of the keypair used in SignOffchain.
	OffchainPublicKey() OffchainPublicKey

	// ConfigEncryptionPublicKey returns the public component of the keypair used in ConfigDiffieHellman.
	ConfigEncryptionPublicKey() ConfigEncryptionPublicKey

	// OffchainKeyring has no Verify method because we always use ed25519 and
	// can thus use the standard library's ed25519.Verify. OffchainKeyring does
	// include a Sign function to prevent passing the private key itself to this
	// library.
}

// OnchainKeyring provides cryptographic signatures that need to be verifiable
// on the targeted blockchain. The underlying cryptographic primitives may be
// different on each chain; for example, on Ethereum one would use ECDSA over
// secp256k1 and Keccak256, whereas on Solana one would use Ed25519 and SHA256.
//
// All its functions should be thread-safe.
type OnchainKeyring interface {
	// PublicKey returns the public key of the keypair used by Sign.
	PublicKey() OnchainPublicKey

	// Sign returns a signature over ReportContext and Report.
	//
	// Reports may contain secret information.
	// Implementations of this function should be careful to not leak
	// the report's contents, e.g. by logging them or including them in
	// returned errors.
	Sign(ReportContext, Report) (signature []byte, err error)

	// Verify verifies a signature over ReportContext and Report allegedly
	// created from OnchainPublicKey.
	//
	// Implementations of this function must gracefully handle malformed or
	// adversarially crafted inputs.
	Verify(_ OnchainPublicKey, _ ReportContext, _ Report, signature []byte) bool

	// Maximum length of a signature
	MaxSignatureLength() int
}
