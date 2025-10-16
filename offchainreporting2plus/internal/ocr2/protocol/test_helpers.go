package protocol

import (
	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

// Used only for testing
type XXXUnknownMessageType struct{}

// Conform to protocol.Message interface
func (XXXUnknownMessageType) CheckSize(types.ReportingPluginLimits) bool { return true }
func (XXXUnknownMessageType) process(*oracleState, commontypes.OracleID) {}
