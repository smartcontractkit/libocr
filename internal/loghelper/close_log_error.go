package loghelper

import (
	"io"

	"github.com/smartcontractkit/libocr/commontypes"
)

// Closes closer. If an error occurs, it is logged at WARN level together with
// msg
func CloseLogError(closer io.Closer, logger commontypes.Logger, msg string) {
	if err := closer.Close(); err != nil {
		logger.Warn(msg, commontypes.LogFields{
			"error": err,
		})
	}
}

// Closes closer if success is false. If an error occurs, it is logged at WARN level together with
// msg
//
// Useful for deferred closing in a constructor like this:
//
//	func newFoo(logger commontypes.Logger) (*Foo, err) {
//		success := false
//
//		bar, err := newBar()
//		if err != nil {
//			return nil, err
//		}
//		defer loghelper.CloseLogErrorUnlessSuccess(&success, bar, logger, "failed to close bar in failed newFoo")
//
//		baz, err := newBaz()
//		if err != nil {
//			return nil, err
//		}
//		defer loghelper.CloseLogErrorUnlessSuccess(&success, baz, logger, "failed to close baz in failed newFoo")
//
//		success = true
//		return &Foo{
//			bar,
//			baz
//		}
//	}
func CloseLogErrorUnlessSuccess(success *bool, closer io.Closer, logger commontypes.Logger, msg string) {
	if success != nil && *success {
		// no failure, return early
		return
	}
	if success == nil {
		logger.Debug("CloseLogErrorUnlessSuccess: got nil success value. this should not happen, assuming it means we did not succeed", nil)
	}
	if err := closer.Close(); err != nil {
		logger.Warn(msg, commontypes.LogFields{
			"error": err,
		})
	}
}
