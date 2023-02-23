package workpool
import "golang.org/x/sync/semaphore"
import "context"
import "sync"

type WorkerPool struct {
	sem *semaphore.Weighted
	wg *sync.WaitGroup
	size int64 
} 

// run function as goroutine onto pool 
func (pool *WorkerPool) Run(f func()) {
	 
	ctx := context.Background()

	pool.wg.Add(1)
	pool.sem.Acquire(ctx,1)
	go func() {
		defer pool.sem.Release(1)
		defer pool.wg.Done() 
		f() 
	}()
}

// Spawns function onto worker pool and returns a future / channel 
func Spawn[T any](pool *WorkerPool, f func() T)  chan T {



	outCh := make(chan T, 1) 	
	f1 := func() {
		outCh <- f()
	}

	pool.Run(f1)
	return outCh 
} 


func Create(n_workers int) WorkerPool {
	w := int64(n_workers)
	sem := semaphore.NewWeighted(w)
	var wg sync.WaitGroup
	return WorkerPool{sem, &wg, w}
}  


// blocks until all pending tasks on pool are completed 
func (pool *WorkerPool) Wait() {
	pool.wg.Wait()
}


/* map list of tasks on to worker pool and outputs results as stream */
func Map[T, U any](wp *WorkerPool, ch chan T, f func(x T) U) chan U {
	out := make(chan U, wp.size)
	


	var w sync.WaitGroup 


	for x := range ch {


		w.Add(1)
		wp.Run(func() {
			out <- f(x)
			w.Done() 
		})
	}

	go func() {
		w.Wait() 
		close(out)
	}()

	return out 
}
