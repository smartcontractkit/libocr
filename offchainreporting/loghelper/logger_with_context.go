package loghelper

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

type LoggerWithContext interface {
	types.Logger
	MakeChild(extraContext types.LogFields) LoggerWithContext
	ErrorIfNotCanceled(msg string, ctx context.Context, fields types.LogFields)
}

type loggerWithContextImpl struct {
	logger  types.Logger
	context types.LogFields
}

// MakeRootLoggerWithContext creates a base logger by wrapping a types.Logger.
// NOTE! Most loggers should extend an existing LoggerWithContext using MakeChild!
func MakeRootLoggerWithContext(logger types.Logger) LoggerWithContext {
	return loggerWithContextImpl{logger, types.LogFields{}}
}

func (l loggerWithContextImpl) Trace(msg string, fields types.LogFields) {
	l.logger.Trace(msg, merge(l.context, fields))
}

func (l loggerWithContextImpl) Debug(msg string, fields types.LogFields) {
	l.logger.Debug(msg, merge(l.context, fields))
}

func (l loggerWithContextImpl) Info(msg string, fields types.LogFields) {
	l.logger.Info(msg, merge(l.context, fields))
}

func (l loggerWithContextImpl) Warn(msg string, fields types.LogFields) {
	l.logger.Warn(msg, merge(l.context, fields))
}

func (l loggerWithContextImpl) Error(msg string, fields types.LogFields) {
	l.logger.Error(msg, merge(l.context, fields))
}

func (l loggerWithContextImpl) ErrorIfNotCanceled(msg string, ctx context.Context, fields types.LogFields) {
	if ctx.Err() != context.Canceled {
		l.logger.Error(msg, merge(l.context, fields))
	} else {
		l.logger.Debug("logging as debug due to context cancelation: "+msg, merge(l.context, fields))
	}
}

// MakeChild is the preferred way to create a new specialised logger.
// It will reuse the base types.Logger and create a new extended context.
func (l loggerWithContextImpl) MakeChild(extra types.LogFields) LoggerWithContext {
	return loggerWithContextImpl{
		l.logger,
		merge(l.context, extra),
	}
}

// Helpers

// merge will create a new LogFields and add all the properties from extras on it.
// Key conflicts are resolved by prefixing the pkey for the new value with underscores until there's no overwrite.
func merge(extras ...types.LogFields) types.LogFields {
	base := types.LogFields{}
	for _, extra := range extras {
		for k, v := range extra {
			add(base, k, v)
		}
	}
	return base
}

// add (key, val) to base. If base already has key, then the old key will be
// left in place and the new key will be prefixed with underscore.
func add(base types.LogFields, key string, val interface{}) {
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
