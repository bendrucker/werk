// Package werk provides a worker pool that can concurrently process tasks up to
// a specified pool size. A worker pool is useful for limiting concurrency.
// Operations in werk are blocking, but allow early termination via configurable
// timeouts and context propagation.
package werk

import (
	"context"
	"errors"
	"time"
)

// Pool represents a fixed-size worker pool. A Pool must be created with NewPool.
type Pool struct {
	size    int
	options *PoolOptions
	ready   workers
}

// PoolOptions includes optional configuration for pool behavior
type PoolOptions struct {
	// AcquireTimeout specifies the maximum wait when calling Acquire
	AcquireTimeout time.Duration
}

type workers chan *Worker

// NewPool initializes a new Pool object. The options argument can be nil.
func NewPool(size int, options *PoolOptions) *Pool {
	if options == nil {
		options = &PoolOptions{}
	}

	p := &Pool{
		options: options,
		size:    size,
		ready:   make(workers, size),
	}

	for i := 0; i < p.size; i++ {
		p.ready <- NewWorker()
	}

	return p
}

// Size returns the originally specified Pool size.
func (p *Pool) Size() int {
	return p.size
}

// Available returns the number of workers that are ready to receive work.
func (p *Pool) Available() int {
	return len(p.ready)
}

// ErrAcquireTimeout is returned by Acquire when a worker is not available
// within the AcquireTimeout specified in PoolOptions.
var ErrAcquireTimeout = errors.New("Acquire timeout")

// Acquire returns a ready worker from the pool, blocking until one is
// available, the context is canceled, or the AcquireTimeout passes. Acquire
// will return ErrAcquireTimeout if the acquire timeout passes.
func (p *Pool) Acquire(ctx context.Context) (*Worker, error) {
	timer := time.NewTimer(p.options.AcquireTimeout)
	if p.options.AcquireTimeout == 0 {
		timer.Stop()
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return nil, ErrAcquireTimeout
	case w := <-p.ready:
		return w, nil

	}
}

// Free returns a worker to the pool.
func (p *Pool) Free(worker *Worker) {
	p.ready <- worker
}

// Do acquires a worker, executes the specified function/work in a new
// goroutine, and frees the worker. Do will block until the fn is done, times
// out, or the context is canceled.
func (p *Pool) Do(ctx context.Context, work Work, fn WorkFunc) (interface{}, error) {
	worker, err := p.Acquire(ctx)

	if err != nil {
		return nil, err
	}

	defer p.Free(worker)
	return worker.Do(ctx, work, fn)
}
