package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Create a listener on port 8000
	li, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	for {
		// Accept the next connection
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	// Reads correspond with Writes on the client side
	var buf [11]byte
	n, err := conn.Read(buf[:])
	fmt.Printf("bytes-read: %d, data: %s, error: %v\n", n, buf[:], err)
}
