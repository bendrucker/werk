package werk

import "context"

// Pool represents a fixed-size worker pool. A Pool must be created with NewPool.
type Pool struct {
	size  int
	ready workers
}

type workers chan *Worker

// NewPool initializes a new Pool object
func NewPool(size int) *Pool {
	p := &Pool{
		size:  size,
		ready: make(workers, size),
	}

	for i := 0; i < p.size; i++ {
		p.ready <- NewWorker()
	}

	return p
}

// Size returns the originally specified Pool size
func (p *Pool) Size() int {
	return p.size
}

// Available returns the number of workers that are ready to receive work
func (p *Pool) Available() int {
	return len(p.ready)
}

// Acquire returns a ready worker from the pool, blocking until one is available
func (p *Pool) Acquire() *Worker {
	return <-p.ready
}

// Free returns a worker to the pool
func (p *Pool) Free(worker *Worker) {
	p.ready <- worker
}

// Do acquires a worker, executes the specified function/work in a new goroutine, and frees the worker.
// Do will block until the fn is done, times out, or the context is canceled.
func (p *Pool) Do(ctx context.Context, work Work, fn WorkFunc) (interface{}, error) {
	worker := p.Acquire()
	defer p.Free(worker)

	return worker.Do(ctx, work, fn)
}
