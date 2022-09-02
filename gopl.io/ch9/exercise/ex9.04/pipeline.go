package pipeline

import (
	"sync"
)

func pipeline(num int) {
	var chans []chan int = make([]chan int, num)
	for i := 0; i < num; i++ {
		chans[i] = make(chan int, 1)
	}

	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			intercomm(chans[i], chans[(i+1)%num])
		}(i)
	}
	wg.Wait()
}

func intercomm(in chan<- int, out <-chan int) {
	in <- 1
	<-out
}
