// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	var wg sync.WaitGroup
	var ch chan string = make(chan string)
	var done chan bool = make(chan bool)

	quit := false
	go func(ch chan string) {
		for input.Scan() {
			content := input.Text()
			log.Println("receiver: ", content)
			ch <- content
		}
		if input.Err() != nil {
			fmt.Println("input error:", input.Err())
		}
		log.Println("finish scan...")
		done <- true
	}(ch)

	for !quit {
		select {
		case <-done:
			log.Println("done")
			quit = true
		case <-time.After(10 * time.Second):
			log.Println("time out, finish")
			quit = true
		case content := <-ch:
			wg.Add(1)
			go func(content string) {
				defer wg.Done()
				echo(c, input.Text(), 1*time.Second)
			}(content)
		}
	}

	wg.Wait()
	if conn, ok := c.(*net.TCPConn); ok {
		conn.CloseWrite()
	} else {
		c.Close()
	}
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
