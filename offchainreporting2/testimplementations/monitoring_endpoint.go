package testimplementations

import (
	"log"

	"github.com/smartcontractkit/libocr/commontypes"
)

type MonitoringEndpoint struct{}

var _ commontypes.MonitoringEndpoint = (*MonitoringEndpoint)(nil)

func (MonitoringEndpoint) SendLog(logbytes []byte) {
	log.Printf("MonitoringEndpoint: sending to monitoring service: %#+v\n", logbytes)
}

type DevnullMonitoringEndpoint struct{}

var _ commontypes.MonitoringEndpoint = (*DevnullMonitoringEndpoint)(nil)

func (DevnullMonitoringEndpoint) SendLog([]byte) {}
