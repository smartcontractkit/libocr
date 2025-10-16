package common

import (
	"context"
	"fmt"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
)

const ReportingPluginTimeoutWarningGracePeriod = 100 * time.Millisecond

func CallPlugin[T any](
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	logFields commontypes.LogFields,
	name string,
	maxDuration time.Duration,
	f func(context.Context) (T, error),
) (T, bool) {
	pluginCtx, cancel := context.WithTimeout(ctx, maxDuration)
	defer cancel()

	ins := loghelper.NewIfNotStopped(
		maxDuration+ReportingPluginTimeoutWarningGracePeriod,
		func() {
			logger.MakeChild(logFields).Warn(fmt.Sprintf("call to ReportingPlugin.%s is taking too long", name), commontypes.LogFields{
				"maxDuration": maxDuration.String(),
				"gracePeriod": ReportingPluginTimeoutWarningGracePeriod.String(),
			})
		},
	)

	result, err := f(pluginCtx)

	ins.Stop()

	if err != nil {
		logger.MakeChild(logFields).ErrorIfNotCanceled(fmt.Sprintf("call to ReportingPlugin.%s errored", name), ctx, commontypes.LogFields{
			"error": err,
		})
		// failed to get data, nothing to be done
		var zero T
		return zero, false
	}

	return result, true
}

// Unlike CallPlugin, CallPluginFromBackground only uses the "recommendedMaxDuration" to warn
// if the call takes longer than recommended, but does not use it for context expiration
// purposes. Context expiration is solely controlled by the passed ctx.
func CallPluginFromBackground[T any](
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	logFields commontypes.LogFields,
	name string,
	recommendedMaxDuration time.Duration,
	f func(context.Context) (T, error),
) (T, bool) {
	ins := loghelper.NewIfNotStopped(
		recommendedMaxDuration+ReportingPluginTimeoutWarningGracePeriod,
		func() {
			logger.MakeChild(logFields).Warn(fmt.Sprintf("call to ReportingPlugin.%s is taking longer than recommended", name), commontypes.LogFields{
				"recommendedMaxDuration": recommendedMaxDuration.String(),
				"gracePeriod":            ReportingPluginTimeoutWarningGracePeriod.String(),
			})
		},
	)

	result, err := f(ctx)

	ins.Stop()

	if err != nil {
		logger.MakeChild(logFields).ErrorIfNotCanceled(fmt.Sprintf("call to ReportingPlugin.%s errored", name), ctx, commontypes.LogFields{
			"error": err,
		})
		// failed to get data, nothing to be done
		var zero T
		return zero, false
	}

	return result, true
}
