package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Create("/tmp/hello.txt")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(f, "Example 06\n")
	f.Close()

	f, err = os.Open("/tmp/hello.txt")
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, f)
	f.Close()
}
