package loghelper

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
)

type LoggerWithContext interface {
	commontypes.Logger
	MakeChild(extraContext commontypes.LogFields) LoggerWithContext
	ErrorIfNotCanceled(msg string, ctx context.Context, fields commontypes.LogFields)
}

type loggerWithContextImpl struct {
	logger  commontypes.Logger
	context commontypes.LogFields
}

// MakeRootLoggerWithContext creates a base logger by wrapping a commontypes.Logger.
// NOTE! Most loggers should extend an existing LoggerWithContext using MakeChild!
func MakeRootLoggerWithContext(logger commontypes.Logger) LoggerWithContext {
	return loggerWithContextImpl{logger, commontypes.LogFields{}}
}

func (l loggerWithContextImpl) Trace(msg string, fields commontypes.LogFields) {
	l.logger.Trace(msg, Merge(l.context, fields))
}

func (l loggerWithContextImpl) Debug(msg string, fields commontypes.LogFields) {
	l.logger.Debug(msg, Merge(l.context, fields))
}

func (l loggerWithContextImpl) Info(msg string, fields commontypes.LogFields) {
	l.logger.Info(msg, Merge(l.context, fields))
}

func (l loggerWithContextImpl) Warn(msg string, fields commontypes.LogFields) {
	l.logger.Warn(msg, Merge(l.context, fields))
}

func (l loggerWithContextImpl) Error(msg string, fields commontypes.LogFields) {
	l.logger.Error(msg, Merge(l.context, fields))
}

func (l loggerWithContextImpl) ErrorIfNotCanceled(msg string, ctx context.Context, fields commontypes.LogFields) {
	if ctx.Err() != context.Canceled {
		l.logger.Error(msg, Merge(l.context, fields))
	} else {
		l.logger.Debug("logging as debug due to context cancelation: "+msg, Merge(l.context, fields))
	}
}

// MakeChild is the preferred way to create a new specialised logger.
// It will reuse the base commontypes.Logger and create a new extended context.
func (l loggerWithContextImpl) MakeChild(extra commontypes.LogFields) LoggerWithContext {
	return loggerWithContextImpl{
		l.logger,
		Merge(l.context, extra),
	}
}

// Helpers

// Merge will create a new LogFields and add all the properties from extras on it.
// Key conflicts are resolved by prefixing the pkey for the new value with underscores until there's no overwrite.
func Merge(extras ...commontypes.LogFields) commontypes.LogFields {
	base := commontypes.LogFields{}
	for _, extra := range extras {
		for k, v := range extra {
			add(base, k, v)
		}
	}
	return base
}

// add (key, val) to base. If base already has key, then the old key will be
// left in place and the new key will be prefixed with underscore.
func add(base commontypes.LogFields, key string, val interface{}) {
	for {
		_, found := base[key]
		if found {
			key = "_" + key
			continue
		}
		base[key] = val
		return
	}
}
