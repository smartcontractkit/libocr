//go:generate ../../../gogeneratehelpers/protoc -I. --go_out=.  ./offchainreporting2_messages.proto
//go:generate ../../../gogeneratehelpers/protoc -I. --go_out=.  ./offchainreporting2_telemetry.proto

package serialization
