package main

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	// START OMIT
	client := &http.Client{Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil { // OMIT
				return nil, err // OMIT
			} // OMIT
			ws := js.Global.Get("WebSocket").New("ws://" + host + ":5000/dial/" + port)
			traceWS(ws) // OMIT
			conn := newWSConn(ws)
			return conn, nil
		},
	}}

	resp, err := client.Get("http://127.0.0.1:5001/")
	if err != nil { // OMIT
		panic(err) // OMIT
	} // OMIT
	defer resp.Body.Close() // OMIT

	//...

	bs, _ := ioutil.ReadAll(resp.Body)
	if err != nil { // OMIT
		panic(err) // OMIT
	} // OMIT
	js.Global.Get("document").Call("write", string(bs))
	// END OMIT
}

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
