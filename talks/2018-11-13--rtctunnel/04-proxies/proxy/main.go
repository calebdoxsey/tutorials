package main

import (
	"io"
	"net"
)

func main() {
	li, err := net.Listen("tcp", "127.0.0.1:8000")
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

func handle(client net.Conn) {
	defer client.Close()

	server, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		panic(err)
	}
	defer server.Close()

	go io.Copy(client, server)
	io.Copy(server, client)
}
