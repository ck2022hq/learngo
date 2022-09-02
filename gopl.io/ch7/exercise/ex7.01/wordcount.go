package main

import (
	"bufio"
	"fmt"
	"strings"
)

type Stat struct {
	nWord int
	nLine int
}

// another implement: https://github.com/ray-g/gopl/blob/master/ch07/ex7.01/counter.go
func (s *Stat) Write(p []byte) (int, error) {
	s.nLine += strings.Count(string(p), "\n") + 1
	for {
		advance, token, _ := bufio.ScanWords(p, true)
		if token == nil {
			break
		}
		s.nWord++
		p = p[advance:]
	}
	return 0, nil
}

func main() {
	var s Stat
	s.Write([]byte("hello world.\nthis is a go file"))
	s.Write([]byte("hello world.\nthis is a go file\n"))

	fmt.Println(s)
}
