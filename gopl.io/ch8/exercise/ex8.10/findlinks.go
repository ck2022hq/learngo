// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var logger *log.Logger

func init() {
	file := "./" + time.Now().Format("2022") + "_log" + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	logger = log.New(logFile, "[orcale_query]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
}

var maxDepth = flag.Int("depth", 2, "depth")

// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func Extract(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Cancel = done

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

//!-Extract

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

type item struct {
	url   string
	depth int
}

func crawl(elem item) []item {
	if elem.depth > *maxDepth {
		return []item{}
	}

	logger.Println("start: ", elem)

	select {
	case <-done:
		return nil
	case tokens <- struct{}{}:
		tokens <- struct{}{} // acquire a token
	}
	defer func() { <-tokens }() // release the token

	list, err := Extract(elem.url)
	logger.Println("end: ", elem)

	if err != nil {
		logger.Print(err)
		return nil
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

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

//!-sema

//!+
func main() {
	flag.Parse()

	// Cancel crawl when input is detected.
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		logger.Println("close done")
		close(done)
	}()

	worklist := make(chan []item)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	go func() { worklist <- transform(flag.Args(), 0) }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)

	var wg sync.WaitGroup

	go func() {
		for {
			// when finish, wait all task done and close worklist
			if cancelled() {
				wg.Wait()
				close(worklist)
				break
			} else {
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

loop:
	for ; n > 0 && !cancelled(); n-- {
		select {
		case list := <-worklist:
			for _, link := range list {
				if cancelled() {
					break loop
				}

				if !seen[link.url] {
					seen[link.url] = true
					n++
					wg.Add(1)
					go func(link item) {
						defer func() {
							wg.Done()
						}()
						worklist <- crawl(link)
					}(link)
				}
			}
		case <-done:
			// Drain worklist to allow existing goroutines to finish.
			// otherwise will cause goroutine leak
			for range worklist {
				// do nothing
			}
			break loop
		}
	}

	logger.Println("exit for loop")
}

//!-
