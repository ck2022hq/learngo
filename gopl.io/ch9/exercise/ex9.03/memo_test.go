// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package memo_test

import (
	"log"
	"testing"
	"time"

	memo "learngo/gopl.io/ch9/exercise/ex9.03"
)

func Test(t *testing.T) {
	var cancel = make(chan struct{})
	m := memo.New(memo.HttpGetBody)

	go func() {
		log.Println("start sleep")
		time.Sleep(10 * time.Second)
		log.Println("finish sleep 10s, close channel")
		close(cancel)
	}()

	memo.Sequential(t, m, cancel)
	m.Print()
}

func TestConcurrent(t *testing.T) {
	var cancel = make(chan struct{})

	go func() {
		log.Println("start sleep")
		time.Sleep(10 * time.Second)
		log.Println("finish sleep 10s, close channel")
		close(cancel)
	}()

	m := memo.New(memo.HttpGetBody)
	memo.Concurrent(t, m, cancel)
	m.Print()
}
