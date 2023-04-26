package delaying

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVoid(t *testing.T) {
	f := VoidWithLog("test", func() {})
	assert.NotNil(t, f)
	t.Run("EnqueueWork", func(t *testing.T) {
		assert.Nil(t, f.EnqueueWork(context.Background(), params{queue: "queue1"}, 1, 2, 3))
	})
	t.Run("EnqueueWorkMulti", func(t *testing.T) {
		assert.Nil(t, f.EnqueueWorkMulti(context.Background(), params{queue: "queue1"}, []any{1, 2}, []any{3, 4}))
	})
}
