package delaying

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	assert.Nil(t, funcImpl)
	var called bool
	f := func(key string, i any) Function {
		called = true
		return function{}
	}
	Init(f)
	assert.NotNil(t, funcImpl)
	funcImpl("key", func() {})
	assert.True(t, called)
	funcImpl = nil
}

func TestMustRegisterFunc(t *testing.T) {
	funcImpl = func(key string, i any) Function {
		return function{}
	}
	doSomething := func() {
	}
	MustRegisterFunc("key", doSomething)
}

func TestWith(t *testing.T) {
	type args struct {
		queue string
		path  string
		delay time.Duration
	}
	tests := []struct {
		name         string
		args         args
		want         params
		expectsPanic bool
	}{
		{"full", args{"queue1", "path1", 3}, params{"queue1", "path1", 3}, false},
		{"empty_queue", args{"", "path1", 3}, params{}, true},
		{"empty_path", args{"queue1", "", 3}, params{}, true},
		{"negative_delay", args{"queue1", "path1", -1}, params{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectsPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("With() did not panic")
					}
				}()
			}
			params := With(tt.args.queue, tt.args.path, tt.args.delay)
			if !tt.expectsPanic {
				assert.Equalf(t, tt.want, params, "With(%v, %v, %v)", tt.args.queue, tt.args.path, tt.args.delay)
			}
		})
	}
}

func TestNewFunction(t *testing.T) {
	t.Run("EnqueueWork", func(t *testing.T) {
		var singleArgs []interface{}
		enqueueWork := func(c context.Context, params Params, args ...any) error {
			singleArgs = args
			return nil
		}
		enqueueWorkMulti := func(c context.Context, params Params, args ...[]any) error {
			panic("unexpected call")
		}
		f := NewFunction("EnqueueWorkTest", func() {}, enqueueWork, enqueueWorkMulti)
		assert.NotNil(t, f)
		assert.Nil(t, singleArgs)
		err := f.EnqueueWork(context.Background(), With("queue1", "path1", 0), 1, 2, 3)
		assert.Nil(t, err)
		assert.Equal(t, []interface{}{1, 2, 3}, singleArgs)
	})
	t.Run("EnqueueWorkMulti", func(t *testing.T) {
		var multiArgs [][]interface{}
		enqueueWork := func(c context.Context, params Params, args ...any) error {
			panic("unexpected call")
		}
		enqueueWorkMulti := func(c context.Context, params Params, args ...[]any) error {
			multiArgs = args
			return nil
		}
		f := NewFunction("EnqueueWorkMultiTest", func() {}, enqueueWork, enqueueWorkMulti)
		err := f.EnqueueWorkMulti(context.Background(), With("queue1", "path1", 0), []any{1, 2}, []any{3, 4})
		assert.Nil(t, err)
		assert.Equal(t, [][]any{{1, 2}, {3, 4}}, multiArgs)
	})
}
