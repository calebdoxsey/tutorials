package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	start := time.Now()
	var xs []int
	for i := 0; i < 10000; i++ {
		xs = append(xs, rand.Intn(1000000))
	}
	xs = ConcurrentMergeSort(xs)
	end := time.Now()

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "sorted %d elements in %s\n", len(xs), end.Sub(start))
	for _, x := range xs {
		fmt.Fprintf(&buf, "%8d\n", x)
	}
	js.Global.Get("document").Call("write", buf.String())
}
