package main

import (
	"fmt"
	"io"
	"strings"
)

type limitReader struct {
	reader io.Reader
	n      int
}

// Read implements io.Reader
func (lr *limitReader) Read(p []byte) (n int, err error) {
	if lr.n <= 0 {
		return 0, io.EOF
	}

	if len(p) > lr.n {
		p = p[0:lr.n]
	}

	return lr.reader.Read(p)
}

func LimitReader(reader io.Reader, n int) io.Reader {
	return &limitReader{reader, n}
}

func main() {
	// there is a io.LimitedReader
	reader := LimitReader(strings.NewReader("hello world"), 5)

	var b []byte = make([]byte, 20)
	n, err := reader.Read(b)

	fmt.Printf("%s\n", string(b))
	fmt.Printf("n=%d\nerr=%v\n", n, err)
}
