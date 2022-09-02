// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 276.

// Package memo provides a concurrency-safe memoization a function of
// a function.  Requests for different keys proceed in parallel.
// Concurrent requests for the same key block until the first completes.
// This implementation uses a Mutex.
package memo

import (
	"fmt"
	"log"
	"sync"
)

// Func is the type of the function to memoize.
type Func func(key string, cancel <-chan struct{}) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

//!+
type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

func (memo *Memo) Print() {
	log.Println("*********memo**********")
	for k := range memo.cache {
		log.Println(k)
	}
	log.Println("*********memo**********")
}

func (memo *Memo) Get(key string, cancel <-chan struct{}) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		log.Println("building key:", key)

		value, err := memo.f(key, cancel)
		if !canceled(cancel) {
			e.res.value, e.res.err = value, err
			close(e.ready) // broadcast ready condition
		} else {
			close(e.ready)
			memo.mu.Lock()
			delete(memo.cache, key)
			memo.mu.Unlock()
			return nil, fmt.Errorf("request canceled")
		}
	} else {
		// This is a repeat request for this key.
		memo.mu.Unlock()

		log.Println("getting key:", key)

		select {
		case <-e.ready:
			// wait for ready condition
		case <-cancel:
			return nil, fmt.Errorf("request canceled")
			// cancel get request
		}
	}
	return e.res.value, e.res.err
}

func canceled(cancel <-chan struct{}) bool {
	select {
	case <-cancel:
		return true
	default:
		return false
	}
}

//!-
