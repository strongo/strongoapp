package delaying

import (
	"context"
	"strings"
	"time"
)

func With(queue, path string, delay time.Duration) Params {
	if strings.TrimSpace(queue) == "" {
		panic("queue is empty")
	}
	if strings.TrimSpace(path) == "" {
		panic("path is empty")
	}
	if delay < 0 {
		panic("delay is negative")
	}
	return params{queue: queue, path: path, delay: delay}
}

type Params interface {
	Queue() string
	Path() string
	Delay() time.Duration
}

var _ Params = params{}

type params struct {
	queue string
	path  string
	delay time.Duration
}

func (p params) Queue() string {
	return p.queue
}

func (p params) Path() string {
	return p.path
}

func (p params) Delay() time.Duration {
	return p.delay
}

type Function interface {
	ID() string
	Implementation() any
	EnqueueWork(c context.Context, params Params, args ...interface{}) error
	EnqueueWorkMulti(c context.Context, params Params, args ...[]interface{}) error
}

type function struct {
	id               string
	implementation   any
	enqueueWork      func(c context.Context, params Params, args ...interface{}) error
	enqueueWorkMulti func(c context.Context, params Params, args ...[]interface{}) error
}

func NewFunction(
	id string,
	implementation any,
	enqueueWork func(c context.Context, params Params, args ...interface{}) error,
	enqueueWorkMulti func(c context.Context, params Params, args ...[]interface{}) error,
) Function {
	if implementation == nil {
		panic("implementation is nil")
	}
	if enqueueWork == nil {
		panic("enqueueWork is nil")
	}
	if enqueueWorkMulti == nil {
		panic("enqueueWorkMulti is nil")
	}
	return function{
		id:               id,
		implementation:   implementation,
		enqueueWork:      enqueueWork,
		enqueueWorkMulti: enqueueWorkMulti,
	}
}

func (f function) ID() string {
	return f.id
}

func (f function) Implementation() any {
	return f.implementation
}

func (f function) EnqueueWork(c context.Context, params Params, args ...interface{}) error {
	return f.enqueueWork(c, params, args...)
}

func (f function) EnqueueWorkMulti(c context.Context, params Params, args ...[]interface{}) error {
	return f.enqueueWorkMulti(c, params, args...)
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
