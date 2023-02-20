package gopool 

import "golang.org/x/sync/semaphore"
import "context"
import "sync"

type WorkerPool struct {
	sem *semaphore.Weighted
	wg *sync.WaitGroup
	size int64 
} 



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


func Spawn[T any](pool *WorkerPool, f func() T)  chan T {



	outCh := make(chan T) 
	
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
		w.Done() 
		close(out)
	}()

	return out 
}
