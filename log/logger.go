package log

import "golang.org/x/net/context"

type Logger interface {
	Name() string
	Debugf(c context.Context, format string, args ...interface{})
	Infof(c context.Context, format string, args ...interface{})
	Warningf(c context.Context, format string, args ...interface{})
	Errorf(c context.Context, format string, args ...interface{})
	Criticalf(c context.Context, format string, args ...interface{})
}

var _loggers []Logger

func AddLogger(logger Logger) {
	for _, l := range _loggers {
		if l.Name() == logger.Name() {
			panic("Duplicate logger name: " + logger.Name())
		}
	}
	_loggers = append(_loggers)
}

func Debugf(c context.Context, format string, args ...interface{}) {
	for _, l := range _loggers {
		l.Debugf(c, format, args...)
	}
}

func Infof(c context.Context, format string, args ...interface{}) {
	for _, l := range _loggers {
		l.Infof(c, format, args...)
	}
}

func Warningf(c context.Context, format string, args ...interface{}) {
	for _, l := range _loggers {
		l.Warningf(c, format, args...)
	}
}

func Errorf(c context.Context, format string, args ...interface{}) {
	for _, l := range _loggers {
		l.Errorf(c, format, args...)
	}
}

func Criticalf(c context.Context, format string, args ...interface{}) {
	for _, l := range _loggers {
		l.Criticalf(c, format, args...)
	}
}
