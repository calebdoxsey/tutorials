package main

import (
	"log"
	"net"
	"net/http"
	"strings"
)

func handleDial(w http.ResponseWriter, r *http.Request) {
	// a path like: /listen/1234
	port := r.URL.Path[strings.LastIndexByte(r.URL.Path, '/')+1:]

	// START DIAL OMIT

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade connection to websocket:", err) // OMIT
		return
	}
	defer ws.Close()

	srcc := &binaryWSConn{Conn: ws}
	dstc, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Println("failed to dial:", err)
		return
	}
	defer dstc.Close()

	// END DIAL OMIT

	err = proxy(dstc, srcc)
	if err != nil {
		log.Println("error handling connection:", err)
	}
}
