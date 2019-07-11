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
	pool := NewPool(10)

	result, _ := pool.Do(context.TODO(), Work{Value: "beep boop"}, func(ctx context.Context, v interface{}) (interface{}, error) {
		fmt.Println("value:", v)
		fmt.Println("workers:", pool.Available())
		return "borp", nil
	})

	fmt.Println("result:", result)

	_, err := pool.Do(context.TODO(), Work{Value: "beep boop"}, func(ctx context.Context, v interface{}) (interface{}, error) {
		return nil, errors.New("oops")
	})

	fmt.Println("err:", err)

	// Output:
	// value: beep boop
	// workers: 9
	// result: borp
	// err: oops
}

func ExamplePool_Do_timeout() {
	pool := NewPool(10)
	work := Work{
		Value:   "foo",
		Timeout: time.Duration(100),
	}

	_, err := pool.Do(context.TODO(), work, func(ctx context.Context, v interface{}) (interface{}, error) {
		time.Sleep(time.Duration(200))
		// returned values received after a timeout are ignored
		return nil, errors.New("inner err")
	})

	fmt.Println("err:", err)

	// Output: err: context deadline exceeded
}

func TestPool(t *testing.T) {
	pool := NewPool(10)

	assert.Equal(t, pool.Size(), 10)
	assert.Equal(t, pool.Available(), 10)

	result, _ := pool.Do(context.TODO(), Work{"hello", 0}, func(_ context.Context, v interface{}) (interface{}, error) {
		assert.Equal(t, "hello", v.(string))
		assert.Equal(t, 9, pool.Available())
		return "woo", nil
	})

	assert.Equal(t, "woo", result)
	assert.Equal(t, pool.Available(), 10)
}

func TestPoolTimeout(t *testing.T) {
	pool := NewPool(10)

	timeout := time.Duration(10) * time.Millisecond
	errors := make(chan error, 1)

	_, _ = pool.Do(context.TODO(), Work{"hello", timeout}, func(ctx context.Context, v interface{}) (interface{}, error) {
		<-ctx.Done()
		errors <- ctx.Err()
		return nil, nil
	})

	err := <-errors
	assert.Equal(t, context.DeadlineExceeded, err)
}
