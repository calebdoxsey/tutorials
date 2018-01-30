package main

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/yamux"
)

func handleListen(w http.ResponseWriter, r *http.Request) {
	// a path like: /listen/1234
	port := r.URL.Path[strings.LastIndexByte(r.URL.Path, '/')+1:]

	// START LISTEN OMIT

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade connection to websocket:", err) // OMIT
		return
	}
	defer ws.Close()

	wsc := &binaryWSConn{Conn: ws}

	dst, err := yamux.Client(wsc, yamux.DefaultConfig())
	if err != nil {
		log.Println("failed to create session", err) // OMIT
		return
	}
	defer dst.Close()

	src, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Println("failed to create new TCP listener:", err) // OMIT
		return
	}
	defer src.Close()

	// END LISTEN OMIT

	// if the "server" disconnects, close the listener too
	go func() {
		for range time.Tick(time.Second) {
			if dst.IsClosed() {
				src.Close()
				return
			}
		}
	}()

	log.Println("started listener", src.Addr())
	defer log.Println("closed listener", src.Addr())

	for {
		// START ACCEPT OMIT
		srcc, err := src.Accept()
		if err != nil {
			log.Println("error accepting connection:", err)
			break
		}

		dstc, err := dst.Open()
		if err != nil {
			srcc.Close()
			log.Println("error opening connection:", err)
			break
		}

		go func() {
			defer srcc.Close()
			defer dstc.Close()
			err := proxy(dstc, srcc)
			if err != nil {
				log.Println("error handling connection:", err)
			}
		}()
		// END ACCEPT OMIT
	}
}
