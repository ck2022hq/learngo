// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
package main

import (
	"fmt"
	"log"
	"os"

	"learngo/gopl.io/ch8/exercise/ex8.06/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

type item struct {
	url   string
	depth int
}

func crawl(elem item) []item {
	if elem.depth > 1 {
		return []item{}
	}

	fmt.Println(elem)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(elem.url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}

	return transform(list, elem.depth+1)
}

func transform(list []string, depth int) []item {
	elems := make([]item, 0, len(list))
	for _, l := range list {
		elems = append(elems, item{l, depth})
	}
	return elems
}

//!-sema

//!+
func main() {
	worklist := make(chan []item)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	go func() { worklist <- transform(os.Args[1:], 0) }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link.url] {
				seen[link.url] = true
				n++
				go func(link item) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

//!-
