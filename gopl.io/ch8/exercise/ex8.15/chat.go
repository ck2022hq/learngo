// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

const timeout = 10 * time.Second

//!+broadcaster
type client struct {
	receiver chan<- string // an outgoing message channel
	name     string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				select {
				case cli.receiver <- msg:
					// do nothing
				default:
					fmt.Println("skip message:", msg)
					// if receiver cannot receiver message, then skip it
				}
			}

		case cli := <-entering:
			var buffer bytes.Buffer
			buffer.WriteString("current user: [")
			for currentCli := range clients {
				buffer.WriteString(currentCli.name)
				buffer.WriteString(" ")
			}
			buffer.WriteString("]")
			cli.receiver <- buffer.String()
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.receiver)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	input := bufio.NewScanner(conn)
	ch <- "who are you? please type your name: "
	var who string
	if input.Scan() {
		who = input.Text()
	}

	ch <- "You are " + who
	messages <- who + " has arrived"
	cli := client{ch, who}
	entering <- cli

	timer := time.NewTimer(timeout)
	go func() {
		<-timer.C
		conn.Close()
	}()

	for input.Scan() {
		messages <- who + ": " + input.Text()
		timer.Reset(timeout)
	}

	// NOTE: ignoring potential errors from input.Err()

	leaving <- cli
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
		// add sleep to simulate receiver channel block
		time.Sleep(2 * time.Second)
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
