package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Network connections implement io.Reader and io.Writer
var _ interface {
	io.Reader
	io.Writer
} = (net.Conn)(nil)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	handle(conn)
}

func handle(conn net.Conn) {
	defer conn.Close()

	n, err := conn.Write([]byte("Hello World"))
	fmt.Printf("bytes-written: %d, error: %v\n", n, err)
}
