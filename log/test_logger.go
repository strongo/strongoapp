package log

import (
	"golang.org/x/net/context"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

type LogMessage struct {
	Level LogLevel
	Format string
	Args []interface{}
}

type TestLogger struct {
	name string
	Messages []LogMessage
}

func (logger *TestLogger) Name() string {
	return logger.name
}
func (logger *TestLogger) add(level LogLevel, format string, args ...interface{}) {
	logger.Messages = append(logger.Messages, LogMessage{Level: level, Format: format, Args: args})
}

func (logger *TestLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.add(LevelDebug, format, args...)
}

// Infof is like Debugf, but at Info level.
func (logger *TestLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	logger.add(LevelInfo, format, args...)
}

// Warningf is like Debugf, but at Warning level.
func (logger *TestLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	logger.add(LevelWarning, format, args...)
}

// Errorf is like Debugf, but at Error level.
func (logger *TestLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.add(LevelError, format, args...)
}

// Criticalf is like Debugf, but at Critical level.
func (logger *TestLogger) Criticalf(ctx context.Context, format string, args ...interface{}) {
	logger.add(LevelCritical, format, args...)
}
