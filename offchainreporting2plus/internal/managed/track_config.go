package managed

import (
	"context"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
)

type trackConfigState struct {
	ctx context.Context
	// in
	configDigester prefixCheckConfigDigester
	configTracker  types.ContractConfigTracker
	localConfig    types.LocalConfig
	logger         loghelper.LoggerWithContext
	// out
	chChanges chan<- types.ContractConfig
	// local
	subprocesses subprocesses.Subprocesses
	configDigest types.ConfigDigest
}

func (state *trackConfigState) run() {
	// Check immediately after startup
	tCheckLatestConfigDetails := time.After(0)

	chNotify := state.configTracker.Notify()

	for {
		select {
		case _, ok := <-chNotify:
			if ok {
				// Check immediately for new config
				tCheckLatestConfigDetails = time.After(0 * time.Second)
				state.logger.Info("TrackConfig: ContractConfigTracker.Notify() fired", nil)
			} else {
				chNotify = nil
				state.logger.Error("TrackConfig: ContractConfigTracker.Notify() was closed, which should never happen. Will ignore ContractConfigTracker.Notify() from now", nil)
			}
		case <-tCheckLatestConfigDetails:
			change, awaitingConfirmation := state.checkLatestConfigDetails()
			state.logger.Debug("TrackConfig: checking latestConfigDetails", nil)

			// poll more rapidly if we're awaiting confirmation
			if awaitingConfirmation {
				wait := 15 * time.Second
				if state.localConfig.ContractConfigTrackerPollInterval < wait {
					wait = state.localConfig.ContractConfigTrackerPollInterval
				}
				tCheckLatestConfigDetails = time.After(wait)
				state.logger.Info("TrackConfig: awaiting confirmation of new config", commontypes.LogFields{
					"wait": wait,
				})
			} else {
				tCheckLatestConfigDetails = time.After(state.localConfig.ContractConfigTrackerPollInterval)
			}

			if change != nil {
				state.configDigest = change.ConfigDigest
				state.logger.Info("TrackConfig: returning config", commontypes.LogFields{
					"configDigest": change.ConfigDigest.Hex(),
				})
				select {
				case state.chChanges <- *change:
				case <-state.ctx.Done():
				}
			}
		case <-state.ctx.Done():
		}

		// ensure prompt exit
		select {
		case <-state.ctx.Done():
			state.logger.Debug("TrackConfig: winding down", nil)
			state.subprocesses.Wait()
			state.logger.Debug("TrackConfig: exiting", nil)
			return
		default:
		}
	}
}

func (state *trackConfigState) checkLatestConfigDetails() (
	latestConfigDetails *types.ContractConfig,
	awaitingConfirmation bool,
) {
	ctx, cancel := context.WithTimeout(state.ctx, state.localConfig.ContractConfigLoadTimeout)
	defer cancel()

	blockheight, err := state.configTracker.LatestBlockHeight(ctx)
	if err != nil {
		state.logger.ErrorIfNotCanceled("TrackConfig: error during LatestBlockHeight()", ctx, commontypes.LogFields{
			"error": err,
		})
		return nil, false
	}

	changedInBlock, latestConfigDigest, err := state.configTracker.LatestConfigDetails(ctx)
	if err != nil {
		state.logger.ErrorIfNotCanceled("TrackConfig: error during LatestConfigDetails()", ctx, commontypes.LogFields{
			"error": err,
		})
		return nil, false
	}
	if latestConfigDigest == (types.ConfigDigest{}) {
		state.logger.Warn("TrackConfig: LatestConfigDetails() returned a zero configDigest. Looks like the contract has not been configured", commontypes.LogFields{
			"configDigest": latestConfigDigest,
		})
		return nil, false
	}
	if state.configDigest == latestConfigDigest {
		return nil, false
	}
	if !state.localConfig.SkipContractConfigConfirmations && blockheight < changedInBlock+uint64(state.localConfig.ContractConfigConfirmations)-1 {
		return nil, true
	}

	contractConfig, err := state.configTracker.LatestConfig(ctx, changedInBlock)
	if err != nil {
		state.logger.ErrorIfNotCanceled("TrackConfig: error during LatestConfig()", ctx, commontypes.LogFields{
			"error": err,
		})
		return nil, true
	}

	if latestConfigDigest != contractConfig.ConfigDigest {
		state.logger.Error("TrackConfig: received config change with ConfigDigest mismatch", commontypes.LogFields{
			"error":              err,
			"contractConfig":     contractConfig,
			"latestConfigDigest": latestConfigDigest,
		})
		return nil, false
	}

	// Ignore configs where the configDigest doesn't match, they might have
	// been corrupted somehow.
	if err := state.configDigester.CheckContractConfig(ctx, contractConfig); err != nil {
		state.logger.ErrorIfNotCanceled("TrackConfig: configDigester returned error, likely due to a corrupted contractConfig", state.ctx, commontypes.LogFields{
			"error":          err,
			"contractConfig": contractConfig,
		})
		return nil, false
	}

	return &contractConfig, false
}

// Uses configTracker to track the latest configuration. When a new configuration is detected, it
// is validated with configDigester and sent to chChanges.
func TrackConfig(
	ctx context.Context,

	configDigester prefixCheckConfigDigester,
	configTracker types.ContractConfigTracker,
	initialConfigDigest types.ConfigDigest,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,

	chChanges chan<- types.ContractConfig,
) {
	state := trackConfigState{
		ctx,
		// in
		configDigester,
		configTracker,
		localConfig,
		logger,
		//out
		chChanges,
		// local
		subprocesses.Subprocesses{},
		initialConfigDigest,
	}
	state.run()
}
