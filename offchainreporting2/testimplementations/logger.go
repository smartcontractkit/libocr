package testimplementations

import (
	"bytes"
	"io"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/smartcontractkit/libocr/commontypes"
)

type DevnullLogger struct{}

func (l DevnullLogger) Trace(msg string, fields commontypes.LogFields) {
}

func (l DevnullLogger) Debug(msg string, fields commontypes.LogFields) {
}

func (l DevnullLogger) Info(msg string, fields commontypes.LogFields) {
}

func (l DevnullLogger) Warn(msg string, fields commontypes.LogFields) {
}

func (l DevnullLogger) Error(msg string, fields commontypes.LogFields) {
}

type Logger struct {
	logger *logrus.Logger
}

func MakeLogger() Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	return Logger{
		logger,
	}
}

func MakePrefixLogger(prefix string) Logger {
	logger := MakeLogger()
	logger.logger.Out = prefixWriter{prefix, logger.logger.Out} // XXX: Does this create a loop?
	return logger
}

type prefixWriter struct {
	prefix string
	writer io.Writer
}

func (p prefixWriter) Write(po []byte) (n int, err error) {
	return p.writer.Write(
		bytes.Join([][]byte{[]byte(p.prefix), []byte(" "), po}, nil),
	)
}

func (l Logger) Trace(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (l Logger) Debug(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l Logger) Info(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l Logger) Warn(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Warn(msg)
}

func (l Logger) Error(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Error(msg)
}

type sLog struct {
	Msg    string
	Fields commontypes.LogFields
}

type MemLogger struct {
	lock sync.RWMutex
	logs []sLog
}

func (m *MemLogger) Trace(msg string, fields commontypes.LogFields) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.logs = append(m.logs, sLog{msg, fields})
}

func (m *MemLogger) Debug(msg string, fields commontypes.LogFields) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.logs = append(m.logs, sLog{msg, fields})
}

func (m *MemLogger) Info(msg string, fields commontypes.LogFields) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.logs = append(m.logs, sLog{msg, fields})
}

func (m *MemLogger) Warn(msg string, fields commontypes.LogFields) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.logs = append(m.logs, sLog{msg, fields})
}

func (m *MemLogger) Error(msg string, fields commontypes.LogFields) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.logs = append(m.logs, sLog{msg, fields})
}

// Pop returns the last log line
func (m *MemLogger) Pop() sLog {
	m.lock.RLock()
	defer m.lock.RUnlock()
	pos := len(m.logs) - 1
	rv := m.logs[pos]
	m.logs = m.logs[:pos]
	return rv
}

// Show returns the current list of logs
func (m *MemLogger) Show() (rv []sLog) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	rv = append(rv, m.logs...)
	return
}

// Peek returns the last logs
func (m *MemLogger) Peek() sLog {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.logs[len(m.logs)-1]
}

type LogChannel <-chan sLog

type ChanLogger struct {
	logger *logrus.Logger
	logCh  chan sLog
}

func NewChanLogger() (*ChanLogger, LogChannel) {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	chLog := make(chan sLog, 1000*1000)
	return &ChanLogger{logger, chLog}, chLog
}

func (c *ChanLogger) Trace(msg string, fields commontypes.LogFields) {
	c.logCh <- sLog{msg, fields}
}

func (c *ChanLogger) Debug(msg string, fields commontypes.LogFields) {
	c.logCh <- sLog{msg, fields}
}

func (c *ChanLogger) Info(msg string, fields commontypes.LogFields) {
	c.logCh <- sLog{msg, fields}
}

func (c *ChanLogger) Warn(msg string, fields commontypes.LogFields) {
	c.logCh <- sLog{msg, fields}
}

func (c *ChanLogger) Error(msg string, fields commontypes.LogFields) {
	c.logCh <- sLog{msg, fields}
}
