package main

import (
	"io"
	"net/http"

	"github.com/hashicorp/yamux"

	"github.com/gopherjs/gopherjs/js"
)

// START OMIT

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World!")
	})
}

func main() {
	ws := js.Global.Get("WebSocket").New("ws://localhost:5000/listen/5001")
	traceWS(ws) // OMIT

	conn := newWSConn(ws)
	defer conn.Close()

	li, err := yamux.Server(conn, yamux.DefaultConfig())
	if err != nil {
		panic(err)
	}
	defer li.Close()

	err = http.Serve(li, nil)
	if err != nil {
		panic(err)
	}
}

// END OMIT

func traceWS(ws *js.Object) {
	ws.Call("addEventListener", "open", func(evt *js.Object) {
		js.Global.Get("console").Call("log", "open", evt)
	})
	ws.Call("addEventListener", "message", func(evt *js.Object) {
		enc := js.Global.Get("TextDecoder").New()
		msg := enc.Call("decode", evt.Get("data"))

		js.Global.Get("console").Call("log", "message", msg)
	})
	ws.Call("addEventListener", "error", func(evt *js.Object) {
		js.Global.Get("console").Call("log", "error", evt)
	})
	ws.Call("addEventListener", "close", func(evt *js.Object) {
		js.Global.Get("console").Call("log", "close", evt)
	})
}
