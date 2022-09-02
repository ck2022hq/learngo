// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 148.

// Fetch saves the contents of a URL into a local file.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

var cancel chan struct{} = make(chan struct{})
var responses chan string = make(chan string)
var wg sync.WaitGroup

//!+
// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) (filename string, n int64, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}
	req.Cancel = cancel

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	n, err = io.Copy(f, resp.Body)
	// Close file, but prefer error from Copy, if any.
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return local, n, err
}

//!-

func canceled() bool {
	select {
	case <-cancel:
		return true
	default:
		return false
	}
}

func mirroredQuery(urls []string) {
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			local, n, err := fetch(url)
			if err != nil {
				fmt.Printf("fetch %s: %v\n", url, err)
			} else {
				responses <- fmt.Sprintf("%s => %s (%d bytes).\n", url, local, n)
			}
		}(url)
	}
}

func main() {
	mirroredQuery(os.Args[1:])

	go func() {
		for {
			if canceled() {
				wg.Wait()
				fmt.Println("close responses")
				close(responses)
				break
			} else {
				time.Sleep(50 * time.Millisecond)
				fmt.Println("sleep 50ms")
			}
		}
	}()

loop:
	for {
		select {
		case resp := <-responses:
			fmt.Println(resp)
			fmt.Println("close cancel")
			close(cancel)
		case <-cancel:
			fmt.Println("into cancel")
			for range responses {
				// do nothing
			}
			fmt.Println("finish clean")
			break loop
		}
	}

	fmt.Println("finish loop")
}
