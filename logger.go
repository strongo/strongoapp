package strongo

import "golang.org/x/net/context"

type Logger interface {
	Debugf(c context.Context, format string, args ...interface{})
	Infof(c context.Context, format string, args ...interface{})
	Warningf(c context.Context, format string, args ...interface{})
	Errorf(c context.Context, format string, args ...interface{})
	Criticalf(c context.Context, format string, args ...interface{})
}
