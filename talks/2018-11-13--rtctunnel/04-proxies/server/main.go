package main

import (
	"net"
)

func main() {
	li, err := net.Listen("tcp", "127.0.0.1:8001")
	if err != nil {
		panic(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			panic(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	var buf [4]byte
	conn.Read(buf[:])
}
