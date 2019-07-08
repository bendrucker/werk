package werk

import "time"

// Worker executes a single task
type Worker struct {
	timeout time.Duration
}

// Work is an interface that holds the data that should be processed by the worker
type Work interface{}

// WorkFunc is a function that receives and handles Work values
type WorkFunc func(Work)

// NewWorker initializes a new worker object
func NewWorker() *Worker {
	return &Worker{}
}

// Do executes the WorkFunc with the Work
func (w *Worker) Do(work Work, fn WorkFunc) {
	fn(work)
}
