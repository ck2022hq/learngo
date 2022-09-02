package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func connect(tz, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Fprintf(os.Stdout, "%s: %s\n", tz, scanner.Text())
	}
	fmt.Println(tz, " done")
	if scanner.Err() != nil {
		fmt.Printf("can't read from %s: %s", tz, scanner.Err())
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s location=addr...\n", os.Args[0])
		return
	}

	for _, arg := range os.Args[1:] {
		result := strings.Split(arg, "=")
		if len(result) != 2 {
			fmt.Printf("invalid args:%s\n", arg)
			continue
		}
		go connect(result[0], result[1])
	}

	for {
		time.Sleep(time.Second)
	}
}
