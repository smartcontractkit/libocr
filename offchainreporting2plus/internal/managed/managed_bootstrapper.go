package managed

import (
	"context"
	"fmt"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/netconfig"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

// RunManagedBootstrapper runs a "managed" bootstrapper. It handles
// configuration updates on the contract.
func RunManagedBootstrapper(
	ctx context.Context,

	bootstrapperFactory types.BootstrapperFactory,
	v2bootstrappers []commontypes.BootstrapperLocator,
	contractConfigTracker types.ContractConfigTracker,
	database types.ConfigDatabase,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	offchainConfigDigester types.OffchainConfigDigester,
) {
	runWithContractConfig(
		ctx,

		contractConfigTracker,
		database,
		func(ctx context.Context, logger loghelper.LoggerWithContext, contractConfig types.ContractConfig) (err error, retry bool) {
			config, err := netconfig.NetConfigFromContractConfig(contractConfig)
			if err != nil {
				return fmt.Errorf("ManagedBootstrapper: error while decoding ContractConfig: %w", err), false
			}

			bootstrapper, err := bootstrapperFactory.NewBootstrapper(config.ConfigDigest, config.PeerIDs, v2bootstrappers, config.F)
			if err != nil {
				logger.Error("ManagedBootstrapper: error during NewBootstrapper", commontypes.LogFields{
					"error":           err,
					"peerIDs":         config.PeerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return fmt.Errorf("ManagedBootstrapper: error during NewBootstrapper: %w", err), true
			}

			if err := bootstrapper.Start(); err != nil {
				return fmt.Errorf("ManagedBootstrapper: error during bootstrapper.Start(): %w", err), true
			}
			defer loghelper.CloseLogError(
				bootstrapper,
				logger,
				"ManagedBootstrapper: error during bootstrapper.Close()",
			)

			<-ctx.Done()

			return nil, false
		},
		localConfig,
		logger,
		offchainConfigDigester,
		defaultRetryParams(),
	)
}
