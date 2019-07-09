package werk

import (
	"context"
	"time"
)

// Worker executes a single task
type Worker struct {
	timeout time.Duration
}

// Work is an struct that holds the data that should be processed by the worker and an optional timeout
// The zero-value of work is valid—if the Timeout is 0, it will be ignored.
type Work struct {
	Value   interface{}
	Timeout time.Duration
}

// WorkFunc is a function that receives and handles Work values
type WorkFunc func(context.Context, interface{}) error

// NewWorker initializes a new worker object
func NewWorker() *Worker {
	return &Worker{}
}

// Do executes the WorkFunc with the Work
func (w *Worker) Do(ctx context.Context, work Work, fn WorkFunc) error {
	if work.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, work.Timeout)
		defer cancel()
	}

	done := make(chan error, 1)
	go func() {
		err := fn(ctx, work.Value)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
