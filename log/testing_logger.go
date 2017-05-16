package log

import (
	"testing"
	"golang.org/x/net/context"
)

type TestingLogger struct {
	t *testing.T
}

func (_ TestingLogger) Name() string {
	return "TestingLogger"
}

func NewTestingLogger(t *testing.T) *TestingLogger {
	return &TestingLogger{t: t}
}

func (logger *TestingLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.t.Logf("Debug: " + format, args...)
}

// Infof is like Debugf, but at Info level.
func (logger *TestingLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	logger.t.Logf("Info: " + format, args...)
}

// Warningf is like Debugf, but at Warning level.
func (logger *TestingLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	logger.t.Logf("Warning: " + format, args...)
}

// Errorf is like Debugf, but at Error level.
func (logger *TestingLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.t.Logf("Error: " + format, args...)
}

// Criticalf is like Debugf, but at Critical level.
func (logger *TestingLogger) Criticalf(ctx context.Context, format string, args ...interface{}) {
	logger.t.Logf("Critical: " + format, args...)
}

