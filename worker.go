package werk

import (
	"context"
	"time"
)

// Worker executes a single task and can be reused
type Worker struct{}

// Work is an struct that holds the data that should be processed by the worker and an optional timeout
// The zero value of Work is valid. If the Timeout is 0, it will be ignored.
type Work struct {
	Value   interface{}
	Timeout time.Duration
}

// WorkFunc is the signature of functions that can handle Work.
// A WorkFunc accepts a context and an interface that will hold the
// Value specified in Work. A WorkFunc may return a single result and an error.
type WorkFunc func(context.Context, interface{}) (interface{}, error)

// NewWorker initializes a new worker object
func NewWorker() *Worker {
	return &Worker{}
}

type workResult struct {
	value interface{}
	err   error
}

// Do invokes the WorkFunc with the Context and Work. The function will be invoked in a new goroutine.
// Do will block until either:
// a) the fn returns a result
// b) work.Timeout has passed
// c) the supplied context is canceled
func (w *Worker) Do(ctx context.Context, work Work, fn WorkFunc) (interface{}, error) {
	if work.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, work.Timeout)
		defer cancel()
	}

	done := make(chan workResult)
	go func() {
		v, err := fn(ctx, work.Value)
		done <- workResult{v, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-done:
		if r.err != nil {
			return nil, r.err
		}

		return r.value, nil
	}
}
