package offchainreporting

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/managed"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/subprocesses"

	"golang.org/x/sync/semaphore"
)

// OracleArgs contains the configuration and services a caller must provide, in
// order to run the offchainreporting protocol.
type OracleArgs struct {
	// A factory for producing network endpoints. A network endpoints consists of
	// networking methods a consumer must implement to allow a node to
	// communicate with other participating nodes.
	BinaryNetworkEndpointFactory types.BinaryNetworkEndpointFactory

	// OCR2 doesn't use the V1 networking
	// // V1Bootstrappers is the list of bootstrap node addresses and IDs for the v1 stack
	// V1Bootstrappers []string

	// V2Bootstrappers is the list of bootstrap node addresses and IDs for the v2 stack
	V2Bootstrappers []commontypes.BootstrapperLocator

	// Tracks configuration changes
	ContractConfigTracker types.ContractConfigTracker

	// Interfaces with the OCR2Aggregator smart contract's transmission related logic
	ContractTransmitter types.ContractTransmitter

	// Database provides persistent storage
	Database types.Database

	// LocalConfig contains oracle-specific configuration details which are not
	// mandated by the on-chain configuration specification via OffchainAggregatoo.SetConfig
	LocalConfig types.LocalConfig

	// Logger logs stuff
	Logger commontypes.Logger

	// Used to send logs to a monitor
	MonitoringEndpoint commontypes.MonitoringEndpoint

	// Computes a config digest using purely offchain logic
	OffchainConfigDigester types.OffchainConfigDigester

	// OffchainKeyring contains the secret keys needed for the OCR protocol, and methods
	// which use those keys without exposing them to the rest of the application.
	OffchainKeyring types.OffchainKeyring

	OnchainKeyring types.OnchainKeyring

	ReportingPluginFactory types.ReportingPluginFactory
}

type Oracle struct {
	oracleArgs OracleArgs

	// Indicates whether the Oracle has been started, in a thread-safe way
	started *semaphore.Weighted

	// subprocesses tracks completion of all go routines on Oracle.Close()
	subprocesses subprocesses.Subprocesses

	// cancel sends a cancel message to all subprocesses, via a context.Context
	cancel context.CancelFunc
}

// NewOracle returns a newly initialized Oracle using the provided services
// and configuration.
func NewOracle(args OracleArgs) (*Oracle, error) {
	if err := SanityCheckLocalConfig(args.LocalConfig); err != nil {
		return nil, errors.Wrapf(err, "bad local config while creating new oracle")
	}
	return &Oracle{
		oracleArgs: args,
		started:    semaphore.NewWeighted(1),
	}, nil
}

// Start spins up a Oracle. Panics if called more than once.
func (o *Oracle) Start() error {
	if !o.started.TryAcquire(1) {
		defer o.Close()
		return fmt.Errorf("can only start an Oracle once")
	}

	logger := loghelper.MakeRootLoggerWithContext(o.oracleArgs.Logger)

	ctx, cancel := context.WithCancel(context.Background())
	o.cancel = cancel
	o.subprocesses.Go(func() {
		defer cancel()
		managed.RunManagedOracle(
			ctx,

			o.oracleArgs.V2Bootstrappers,
			o.oracleArgs.ContractConfigTracker,
			o.oracleArgs.ContractTransmitter,
			o.oracleArgs.Database,
			o.oracleArgs.LocalConfig,
			logger,
			o.oracleArgs.MonitoringEndpoint,
			o.oracleArgs.BinaryNetworkEndpointFactory,
			o.oracleArgs.OffchainConfigDigester,
			o.oracleArgs.OffchainKeyring,
			o.oracleArgs.OnchainKeyring,
			o.oracleArgs.ReportingPluginFactory,
		)
	})
	return nil
}

// Close shuts down an oracle. Can safely be called multiple times.
func (o *Oracle) Close() error {
	if o.cancel != nil {
		o.cancel()
	}
	// Wait for all subprocesses to shut down, before shutting down other resources.
	// (Wouldn't want anything to panic from attempting to use a closed resource.)
	o.subprocesses.Wait()
	return nil
}
