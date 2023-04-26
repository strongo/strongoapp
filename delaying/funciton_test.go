package delaying

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestNewFunction(t *testing.T) {
	var single []interface{}
	var multi [][]interface{}
	enqueueWork := func(c context.Context, args ...any) error {
		single = args
		return nil
	}
	enqueueWorkMulti := func(c context.Context, args ...[]any) error {
		multi = args
		return nil
	}
	f := NewFunction(enqueueWork, enqueueWorkMulti)
	assert.NotNil(t, f)
	assert.Nil(t, single)
	assert.Nil(t, multi)
	t.Run("EnqueueWork", func(t *testing.T) {
		err := f.EnqueueWork(context.Background(), 1, 2, 3)
		assert.Nil(t, err)
		assert.Equal(t, []interface{}{1, 2, 3}, single)
	})
	t.Run("EnqueueWorkMulti", func(t *testing.T) {
		err := f.EnqueueWorkMulti(context.Background(), []any{1, 2}, []any{3, 4})
		assert.Nil(t, err)
		assert.Equal(t, [][]any{{1, 2}, {3, 4}}, multi)
	})
}
