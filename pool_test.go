package werk

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleNewPool() {
	pool := NewPool(10).Start()

	_ = pool.Do(context.TODO(), Work{Value: "beep boop"}, func(ctx context.Context, v interface{}) error {
		fmt.Println("value:", v)
		fmt.Println("workers:", pool.Available())
		return nil
	})

	err := pool.Do(context.TODO(), Work{Value: "beep boop"}, func(ctx context.Context, v interface{}) error {
		return errors.New("oops")
	})

	fmt.Println("err:", err)

	// Output:
	// value: beep boop
	// workers: 9
	// err: oops
}

func ExamplePool_Do_timeout() {
	pool := NewPool(10).Start()
	work := Work{
		Value:   "foo",
		Timeout: time.Duration(100),
	}

	err := pool.Do(context.TODO(), work, func(ctx context.Context, v interface{}) error {
		time.Sleep(time.Duration(200))
		return nil
	})

	fmt.Println("err:", err)

	// Output: err: context deadline exceeded
}

func TestPool(t *testing.T) {
	pool := NewPool(10).Start()

	assert.Equal(t, pool.Available(), 10)

	pool.Do(context.TODO(), Work{"hello", 0}, func(_ context.Context, v interface{}) error {
		assert.Equal(t, "hello", v.(string))
		assert.Equal(t, 9, pool.Available())
		return nil
	})

	assert.Equal(t, pool.Available(), 10)
}

func TestPoolTimeout(t *testing.T) {
	pool := NewPool(10).Start()

	timeout := time.Duration(10) * time.Millisecond
	errors := make(chan error, 1)

	pool.Do(context.TODO(), Work{"hello", timeout}, func(ctx context.Context, v interface{}) error {
		<-ctx.Done()
		errors <- ctx.Err()
		return nil
	})

	err := <-errors
	assert.Equal(t, context.DeadlineExceeded, err)
}
