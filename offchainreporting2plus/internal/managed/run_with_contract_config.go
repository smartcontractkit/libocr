package managed

import (
	"context"
	"math/rand"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
)

// A retriableFn is a function that is managed by runWithContractConfig.
// On error, it can indicate whether it would like to be retried automatically
// via the the retry bool. If err is nil, retry is ignored.
type retriableFn func(context.Context, loghelper.LoggerWithContext, types.ContractConfig) (err error, retry bool)

type retryParams struct {
	InitialSleep      time.Duration
	SleepMultiplier   float64       // should be greater than 1 to achieve exponential backoff
	MaxRelativeJitter float64       // keep this between 0 and 1
	MaxSleep          time.Duration // sleep will never be longer than this
}

func defaultRetryParams() retryParams {
	return retryParams{
		InitialSleep:      time.Second,
		SleepMultiplier:   2,
		MaxRelativeJitter: 0.1,
		MaxSleep:          2 * time.Minute,
	}
}

// runWithContractConfig runs fn with a contractConfig and manages its lifecycle
// as contractConfigs change according to contractConfigTracker. It also saves
// and restores contract configs using database.
func runWithContractConfig(
	ctx context.Context,

	contractConfigTracker types.ContractConfigTracker,
	database types.ConfigDatabase,
	fn retriableFn,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	offchainConfigDigester types.OffchainConfigDigester,
	retryParams retryParams,
) {
	rwcc := runWithContractConfigState{
		ctx,

		types.ConfigDigest{},
		contractConfigTracker,
		database,
		fn,
		localConfig,
		logger,
		retryParams,

		prefixCheckConfigDigester{offchainConfigDigester},
		func() {},
		subprocesses.Subprocesses{},
		subprocesses.Subprocesses{},
	}
	rwcc.run()
}

type runWithContractConfigState struct {
	ctx context.Context

	configDigest          types.ConfigDigest
	contractConfigTracker types.ContractConfigTracker
	database              types.ConfigDatabase
	fn                    retriableFn
	localConfig           types.LocalConfig
	logger                loghelper.LoggerWithContext
	retryParams           retryParams

	configDigester prefixCheckConfigDigester
	fnCancel       context.CancelFunc
	fnSubs         subprocesses.Subprocesses
	otherSubs      subprocesses.Subprocesses
}

func (rwcc *runWithContractConfigState) run() {
	// Restore config from database, so that we can run even if the ethereum node
	// isn't working.
	rwcc.restoreFromDatabase()

	// Only start tracking config after we attempted to load config from db
	chNewConfig := make(chan types.ContractConfig, 5)
	rwcc.otherSubs.Go(func() {
		TrackConfig(rwcc.ctx, rwcc.configDigester, rwcc.contractConfigTracker, rwcc.configDigest, rwcc.localConfig, rwcc.logger, chNewConfig)
	})

	for {
		select {
		case change := <-chNewConfig:
			rwcc.logger.Info("runWithContractConfig: switching between configs", commontypes.LogFields{
				"oldConfigDigest": rwcc.configDigest.Hex(),
				"newConfigDigest": change.ConfigDigest.Hex(),
			})
			rwcc.configChanged(change)
		case <-rwcc.ctx.Done():
			rwcc.logger.Info("runWithContractConfig: winding down", nil)
			rwcc.fnSubs.Wait()
			rwcc.otherSubs.Wait()
			rwcc.logger.Info("runWithContractConfig: exiting", nil)
			return // Exit managed event loop altogether
		}
	}
}

func (rwcc *runWithContractConfigState) restoreFromDatabase() {
	var contractConfig *types.ContractConfig
	ok := rwcc.otherSubs.BlockForAtMost(
		rwcc.ctx,
		rwcc.localConfig.DatabaseTimeout,
		func(ctx context.Context) {
			contractConfig = loadConfigFromDatabase(ctx, rwcc.database, rwcc.logger)
		},
	)
	if !ok {
		rwcc.logger.Error("runWithContractConfig: database timed out while attempting to restore configuration", commontypes.LogFields{
			"timeout": rwcc.localConfig.DatabaseTimeout,
		})
		return
	}

	if contractConfig == nil {
		rwcc.logger.Info("runWithContractConfig: found no configuration to restore", commontypes.LogFields{})
		return
	}

	rwcc.configChanged(*contractConfig)
}

// We assume that contractConfig has already been validated by the OffchainConfigDigester.
func (rwcc *runWithContractConfigState) configChanged(contractConfig types.ContractConfig) {
	// Cease any operation from earlier configs
	rwcc.logger.Info("runWithContractConfig: winding down old configuration", commontypes.LogFields{
		"oldConfigDigest": rwcc.configDigest,
		"newConfigDigest": contractConfig.ConfigDigest,
	})
	rwcc.fnCancel()
	rwcc.fnSubs.Wait()
	rwcc.logger.Info("runWithContractConfig: closed old configuration", commontypes.LogFields{
		"oldConfigDigest": rwcc.configDigest,
		"newConfigDigest": contractConfig.ConfigDigest,
	})

	rwcc.configDigest = contractConfig.ConfigDigest

	logger := rwcc.logger.MakeChild(commontypes.LogFields{"configDigest": contractConfig.ConfigDigest})

	fnCtx, fnCancel := context.WithCancel(rwcc.ctx)
	rwcc.fnCancel = fnCancel
	rwcc.fnSubs.Go(func() {
		defer fnCancel()
		retryOnError(fnCtx, logger, rwcc.retryParams, contractConfig, rwcc.fn)
	})

	writeCtx, writeCancel := context.WithTimeout(rwcc.ctx, rwcc.localConfig.DatabaseTimeout)
	defer writeCancel()
	if err := rwcc.database.WriteConfig(writeCtx, contractConfig); err != nil {
		rwcc.logger.ErrorIfNotCanceled("runWithContractConfig: error writing new config to database", writeCtx, commontypes.LogFields{
			"configDigest": contractConfig.ConfigDigest,
			"config":       contractConfig,
			"error":        err,
		})
	}

}

func retryOnError(ctx context.Context, logger loghelper.LoggerWithContext, retryParams retryParams, contractConfig types.ContractConfig, fn retriableFn) {
	sleep := retryParams.InitialSleep

	for retry := 0; ; retry++ {
		retryLogger := logger.MakeChild(commontypes.LogFields{"retry": retry})
		retryLogger.Info("runWithContractConfig: running function", nil)
		// We intentionally don't pass retryLogger to fn, because we don't want to include the "retry"
		// in every log line logged from fn.
		err, retry := fn(ctx, logger, contractConfig)
		if err == nil {
			retryLogger.Info("runWithContractConfig: function exited without error. not retrying", nil)
			return
		}
		if !retry {
			retryLogger.ErrorIfNotCanceled("runWithContractConfig: function exited with non-retriable error. not retrying", ctx, commontypes.LogFields{
				"error": err,
			})
			return
		}
		retryLogger.ErrorIfNotCanceled("runWithContractConfig: function returned error. sleeping before retrying", ctx, commontypes.LogFields{
			"error": err,
			"sleep": sleep.String(),
		})

		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			retryLogger.Info("runWithContractConfig: context expired while sleeping before retrying. not retrying", nil)
			return
		}

		sleep = time.Duration(float64(sleep) * retryParams.SleepMultiplier)
		if sleep > retryParams.MaxSleep {
			sleep = retryParams.MaxSleep
		}
		// Subtract jitter so we always respect MaxSleep
		jitterFactor := 1 - rand.Float64()*retryParams.MaxRelativeJitter
		sleep = time.Duration(float64(sleep) * jitterFactor)
	}
}
