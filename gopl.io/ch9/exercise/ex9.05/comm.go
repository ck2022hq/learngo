package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	num = 0
	mtx sync.Mutex
)

func get() int {
	mtx.Lock()
	defer mtx.Unlock()
	return num
}

func incr(count int) {
	mtx.Lock()
	num += count
	mtx.Unlock()
}

func Comm(ch1, ch2 chan int, done chan struct{}) {
	go intercomm(ch1, ch2, done)
	go intercomm(ch2, ch1, done)

	select {
	case <-done:
		// do nothing
	}
}

func isDone(done chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func intercomm(in chan<- int, out <-chan int, done chan struct{}) {
	for !isDone(done) {
		select {
		case in <- 1:
		case <-out:
			incr(1)
		}
	}
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	done := make(chan struct{})

	go Comm(ch1, ch2, done)
	time.Sleep(5 * time.Second)
	close(done)

	fmt.Println("num = ", get())
}
