package main

import (
	"os"
	"strconv"
)

func main() {
	n, _ := strconv.Atoi(os.Args[1])
	str := []byte("hello world")
	for i := 0; i < n; i++ {
		os.Stdout.Write(str)
	}
}
