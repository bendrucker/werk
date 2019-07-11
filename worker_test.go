package werk

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerSuccess(t *testing.T) {
	worker := NewWorker()

	err := worker.Do(context.TODO(), Work{Value: "foo"}, func(ctx context.Context, v interface{}) error {
		assert.Equal(t, "foo", v.(string))
		return nil
	})

	assert.NoError(t, err)
}

func TestWorkerError(t *testing.T) {
	worker := NewWorker()

	err := worker.Do(context.TODO(), Work{Value: "foo"}, func(ctx context.Context, v interface{}) error {
		return errors.New("oops")
	})

	assert.EqualError(t, err, "oops")
}

func TestWorkerTimeout(t *testing.T) {
	worker := NewWorker()
	timeout := time.Duration(100)

	err := worker.Do(context.TODO(), Work{"foo", timeout}, func(ctx context.Context, v interface{}) error {
		time.Sleep(time.Duration(200))
		assert.Error(t, ctx.Err())
		return nil
	})

	assert.Equal(t, context.DeadlineExceeded, err)
}
