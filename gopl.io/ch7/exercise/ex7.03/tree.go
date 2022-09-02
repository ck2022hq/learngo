package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type tree struct {
	value       int
	left, right *tree
}

func (t *tree) String() string {
	var buffer bytes.Buffer
	t.travel(&buffer, 0)
	return buffer.String()
}

func (t *tree) travel(buffer *bytes.Buffer, nSpace int) {
	if t == nil {
		return
	}

	buffer.WriteString(strings.Repeat(" ", nSpace))
	buffer.WriteString(strconv.FormatInt(int64(t.value), 10))
	buffer.WriteString("\n")
	t.left.travel(buffer, nSpace+2)
	t.right.travel(buffer, nSpace+2)
}

// Sort sorts values in place.
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	fmt.Println("********************************")
	fmt.Println(root)
	fmt.Println("********************************")
	appendValues(values[:0], root)
}

// appendValues appends the elements of t to values in order
// and returns the resulting slice.
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		// Equivalent to return &tree{value: value}.
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

func main() {
	data := make([]int, 50)
	for i := range data {
		data[i] = rand.Int() % 50
	}
	Sort(data)
	if !sort.IntsAreSorted(data) {
		fmt.Printf("not sorted: %v\n", data)
	}
}
