package gopool

import (
	"fmt"
	"math/rand"
	"testing"
)

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func newWorkerPool() WorkerPool {
	return Create(50)
}

func TestMap(t *testing.T) {
	in := make(chan string, 24)

	pool := newWorkerPool()
	for i := 0; i < 20; i++ {
		in <- RandomString(40)
	}

	close(in)

	f := func(x string) string {
		return x + "_" + x
	}

	out := Map(&pool, in, f)

	for x := range out {
		fmt.Println(x)
	}

}

func runN(f func(), n int) {
	for i := 0; i < n; i++ {
		f()
	}
}

func TestWait(t *testing.T) {
	pool := newWorkerPool()
	f := func() { fmt.Println("hello") }
	runN(func() { pool.Run(f) }, 5)

	randStr := func() string { return RandomString(40) }

	o := Spawn(&pool, randStr)

	x := <-o
	pool.Wait()
	fmt.Println(x)
}
