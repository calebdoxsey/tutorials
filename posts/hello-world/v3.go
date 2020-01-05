package main

import (
	"bufio"
	"os"
	"strconv"
)


func main() {
	n, _ := strconv.Atoi(os.Args[1])
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	str := []byte("hello world")
	for i := 0; i < n; i++ {
		w.Write(str)
	}
}
