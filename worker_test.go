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

	_, err := worker.Do(context.TODO(), Work{Value: "foo"}, func(ctx context.Context, v interface{}) (interface{}, error) {
		assert.Equal(t, "foo", v.(string))
		return nil, nil
	})

	assert.NoError(t, err)
}

func TestWorkerError(t *testing.T) {
	worker := NewWorker()

	_, err := worker.Do(context.TODO(), Work{Value: "foo"}, func(ctx context.Context, v interface{}) (interface{}, error) {
		return nil, errors.New("oops")
	})

	assert.EqualError(t, err, "oops")
}

func TestWorkerTimeout(t *testing.T) {
	worker := NewWorker()
	timeout := time.Duration(100)

	_, err := worker.Do(context.TODO(), Work{"foo", timeout}, func(ctx context.Context, v interface{}) (interface{}, error) {
		time.Sleep(time.Duration(200))
		assert.Error(t, ctx.Err())
		return nil, nil
	})

	assert.Equal(t, context.DeadlineExceeded, err)
}
