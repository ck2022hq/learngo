package main

import (
	"fmt"
	"io"
	"os"
)

type countingWriter struct {
	writer io.Writer
	cnt    int64
}

func CountingWriter(writer io.Writer) (io.Writer, *int64) {
	cw := &countingWriter{writer: writer, cnt: 0}
	return cw, &(cw.cnt)
}

func (cw *countingWriter) Write(p []byte) (n int, err error) {
	n, err = cw.writer.Write(p)
	cw.cnt += int64(n)
	return
}

func main() {
	cw, n := CountingWriter(os.Stdout)
	cw.Write([]byte("hello world!\n"))
	cw.Write([]byte("newline\n"))
	fmt.Printf("cnt=%d\n", *n)
}
