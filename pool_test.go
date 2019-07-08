package werk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	pool := NewPool(10).Start()

	assert.Equal(t, pool.Available(), 10)

	pool.Do(Work{"hello", 0}, WorkFunc(func(_ context.Context, v interface{}) {
		assert.Equal(t, "hello", v.(string))
		assert.Equal(t, 9, pool.Available())
	}))

	assert.Equal(t, pool.Available(), 10)
}

func TestPoolTimeout(t *testing.T) {
	pool := NewPool(10).Start()

	timeout := time.Duration(10) * time.Millisecond
	start := time.Now()

	pool.Do(Work{"hello", timeout}, WorkFunc(func(ctx context.Context, v interface{}) {
		<-ctx.Done()
		err := ctx.Err()
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Equal(t, timeout, time.Since(start))
	}))
}
