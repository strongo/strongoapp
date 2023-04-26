package delaying

import "context"

type Function interface {
	EnqueueWork(c context.Context, args ...interface{}) error
	EnqueueWorkMulti(c context.Context, args ...[]interface{}) error
}

type function struct {
	enqueueWork      func(c context.Context, args ...interface{}) error
	enqueueWorkMulti func(c context.Context, args ...[]interface{}) error
}

func NewFunction(
	enqueueWork func(c context.Context, args ...interface{}) error,
	enqueueWorkMulti func(c context.Context, args ...[]interface{}) error,
) Function {
	return function{
		enqueueWork:      enqueueWork,
		enqueueWorkMulti: enqueueWorkMulti,
	}
}

func (f function) EnqueueWork(c context.Context, args ...interface{}) error {
	return f.enqueueWork(c, args...)
}

func (f function) EnqueueWorkMulti(c context.Context, args ...[]interface{}) error {
	return f.enqueueWorkMulti(c, args...)
}

func MustRegisterFunc(key string, i any) Function {
	if funcImpl == nil {
		panic("No implementation has been registered. Application should call delaying.Init(...) before using delaying.MustRegisterFunc()")
	}
	return funcImpl(key, i)
}

var funcImpl func(key string, i any) Function

func Init(f func(key string, i any) Function) {
	if f == nil {
		panic("f is nil")
	}
	funcImpl = f
}
