package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	log.SetFlags(0)

	http.HandleFunc("/listen/", handleListen)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:5000"
	}
	log.Printf("starting http server addr=%s\n", addr)
	http.ListenAndServe(addr, nil)
}

// START PROXY OMIT

func proxy(dst, src net.Conn) error {
	var eg errgroup.Group
	eg.Go(func() error {
		_, err := io.Copy(dst, src)
		return err
	})
	eg.Go(func() error {
		_, err := io.Copy(src, dst)
		return err
	})
	return eg.Wait()
}

// END PROXY OMIT
