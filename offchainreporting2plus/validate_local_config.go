package offchainreporting2plus

import (
	"errors"
	"fmt"
	"time"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

func boundTimeDuration(
	duration time.Duration,
	durationName string,
	min time.Duration,
	max time.Duration,
) error {
	if !(min <= duration && duration <= max) {
		return fmt.Errorf("%s must be between %v and %v, but is currently %v",
			durationName, min, max, duration)
	}
	return nil
}

func SanityCheckLocalConfig(c types.LocalConfig) (err error) {
	if c.DevelopmentMode == types.EnableDangerousDevelopmentMode {
		return nil
	}

	err = errors.Join(err,
		boundTimeDuration(
			c.BlockchainTimeout,
			"blockchain timeout",
			1*time.Second, 20*time.Second,
		))
	err = errors.Join(err,
		boundTimeDuration(
			c.ContractConfigTrackerPollInterval,
			"contract config tracker poll interval",
			1*time.Second, 120*time.Second,
		))
	err = errors.Join(err,
		boundTimeDuration(
			c.ContractConfigLoadTimeout,
			"contract config load timeout",
			1*time.Second, 1*time.Hour,
		))
	err = errors.Join(err,
		boundTimeDuration(
			c.ContractTransmitterTransmitTimeout,
			"contract transmitter transmit timeout",
			1*time.Second, 1*time.Minute,
		))
	err = errors.Join(err,
		boundTimeDuration(
			c.DatabaseTimeout,
			"database timeout",
			100*time.Millisecond, 10*time.Second,
		))
	err = errors.Join(err,
		boundTimeDuration(
			c.DefaultMaxDurationInitialization,
			"DefaultMaxDurationInitialization",
			1*time.Second, 1*time.Hour,
		))

	const minContractConfigConfirmations = 1
	const maxContractConfigConfirmations = 100
	if !(minContractConfigConfirmations <= c.ContractConfigConfirmations && c.ContractConfigConfirmations <= maxContractConfigConfirmations) {
		err = errors.Join(err, fmt.Errorf(
			"contract config block-depth confirmation threshold must be between %v and %v, but is currently %v",
			minContractConfigConfirmations,
			maxContractConfigConfirmations,
			c.ContractConfigConfirmations))

	}

	return err
}
