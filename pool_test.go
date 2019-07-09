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

	pool.Do(Work{"hello", 0}, WorkFunc(func(_ context.Context, v interface{}) error {
		assert.Equal(t, "hello", v.(string))
		assert.Equal(t, 9, pool.Available())
		return nil
	}))

	assert.Equal(t, pool.Available(), 10)
}

func TestPoolTimeout(t *testing.T) {
	pool := NewPool(10).Start()

	timeout := time.Duration(10) * time.Millisecond
	errors := make(chan error, 1)

	pool.Do(Work{"hello", timeout}, WorkFunc(func(ctx context.Context, v interface{}) error {
		<-ctx.Done()
		errors <- ctx.Err()
		return nil
	}))

	err := <-errors
	assert.Equal(t, context.DeadlineExceeded, err)
}
