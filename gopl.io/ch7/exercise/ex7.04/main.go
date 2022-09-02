package main

import (
	"fmt"
	"io"
)

type stringReader struct {
	s string
}

func (sr *stringReader) Read(p []byte) (n int, err error) {
	n = copy(p, sr.s)
	sr.s = sr.s[n:]
	if len(sr.s) == 0 {
		err = io.EOF
	}
	return
}

func NewReader(s string) io.Reader {
	return &stringReader{s}
}

func main() {
	reader := NewReader("hello world")
	var p []byte = make([]byte, 20)
	n, _ := reader.Read(p)
	fmt.Printf("%s\n", string(p))
	fmt.Printf("n=%d\n", n)
}
