//go:generate protoc -I. --go_out=.  ./offchainreporting2_messages.proto
//go:generate protoc -I. --go_out=.  ./offchainreporting2_telemetry.proto

package serialization
