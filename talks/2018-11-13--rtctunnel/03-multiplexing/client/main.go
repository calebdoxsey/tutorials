package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/xtaci/smux"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalln(err)
	}
	handleConn(conn)
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// create a multiplexing session
	session, err := smux.Server(conn, nil)
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		// invert the relationship between client and server and receive a stream
		// which is also a net.Conn
		stream, err := session.AcceptStream()
		if err != nil {
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			handleStream(stream)
		}()
	}
	wg.Wait()
}

func handleStream(conn *smux.Stream) {
	defer conn.Close()

	var buf [4]byte
	io.ReadFull(conn, buf[:])
	fmt.Printf("[recv] id: %d, local-addr: %v, remote-addr: %v, data: %s\n",
		conn.ID(), conn.LocalAddr(), conn.RemoteAddr(), buf[:])
}
