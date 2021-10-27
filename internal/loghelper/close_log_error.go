package loghelper

import (
	"io"

	"github.com/smartcontractkit/libocr/commontypes"
)

// Closes closer. If an error occurs, it is logged at INFO level together with
// msg
func CloseLogError(closer io.Closer, logger commontypes.Logger, msg string) {
	if err := closer.Close(); err != nil {
		logger.Info(msg, commontypes.LogFields{
			"error": err,
		})
	}
}
