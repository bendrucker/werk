package werk

import "context"

// Pool represents a fixed-size worker pool
type Pool struct {
	count int
	ready workers
}

type workers chan *Worker

// NewPool initializes a new Pool object
func NewPool(count int) *Pool {
	return &Pool{
		count: count,
		ready: make(workers, count),
	}
}

// Start readies the number of workers specified in "count"
func (p *Pool) Start() *Pool {
	for i := 0; i < p.count; i++ {
		p.worker()
	}

	return p
}

func (p *Pool) worker() {
	p.ready <- NewWorker()
}

// Available returns the number of workers that are ready to receive work
func (p *Pool) Available() int {
	return len(p.ready)
}

// Acquire returns a ready worker from the pool, blocking until one is available
func (p *Pool) Acquire() *Worker {
	return <-p.ready
}

// Free returns a worker to pool
func (p *Pool) Free(worker *Worker) {
	p.ready <- worker
}

// Do acquires a worker, executes the specified function/work, and frees the worker
func (p *Pool) Do(work Work, fn WorkFunc) {
	worker := p.Acquire()
	defer p.Free(worker)

	worker.Do(context.TODO(), work, fn)
}
