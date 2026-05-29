package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
)

var _ commontypes.Logger = DevnullLogger{}

type DevnullLogger struct{}

func (l DevnullLogger) Trace(msg string, fields commontypes.LogFields)    {}
func (l DevnullLogger) Debug(msg string, fields commontypes.LogFields)    {}
func (l DevnullLogger) Info(msg string, fields commontypes.LogFields)     {}
func (l DevnullLogger) Warn(msg string, fields commontypes.LogFields)     {}
func (l DevnullLogger) Error(msg string, fields commontypes.LogFields)    {}
func (l DevnullLogger) Critical(msg string, fields commontypes.LogFields) {}

var _ prometheus.Registerer = DevnullRegisterer{}

type DevnullRegisterer struct{}

func (r DevnullRegisterer) Register(prometheus.Collector) error  { return nil }
func (r DevnullRegisterer) MustRegister(...prometheus.Collector) {}
func (r DevnullRegisterer) Unregister(prometheus.Collector) bool { return false }
