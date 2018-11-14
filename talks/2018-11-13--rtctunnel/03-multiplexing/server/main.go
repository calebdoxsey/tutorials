package main

import (
	"fmt"
	"log"
	"net"

	"github.com/xtaci/smux"
)

func main() {
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

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// create a multiplexing session
	session, err := smux.Client(conn, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < 4; i++ {
		stream, err := session.OpenStream()
		if err != nil {
			break
		}
		go handleStream(stream)
	}
}

func handleStream(stream *smux.Stream) {
	defer stream.Close()

	buf := []byte("ping")
	fmt.Printf("[send] id: %d, local-addr: %v, remote-addr: %v, data: %s\n",
		stream.ID(), stream.LocalAddr(), stream.RemoteAddr(), buf)
	stream.Write(buf)
}
