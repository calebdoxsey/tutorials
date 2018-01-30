package main

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World")
	})
	_, err := net.Listen("tcp4", "127.0.0.1:80")
	if err != nil {
		panic(err)
	}
	// var eg errgroup.Group
	// eg.Go(func() error {
	// 	return (&http.Server{Handler: http.DefaultServeMux}).Serve(ln)
	// })
	// eg.Go(func() error {
	// 	res, err := http.Get("http://127.0.0.1")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer res.Body.Close()

	// 	bs, err := ioutil.ReadAll(res.Body)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	js.Global.Get("document").Call("write", string(bs))
	// 	return nil
	// })
	// err = eg.Wait()
	// if err != nil {
	// 	panic(err)
	// }
}

func info(msg string, args ...interface{}) {
	js.Global.Get("console").Call("info", fmt.Sprintf(msg, args...))
}

func warn(msg string, args ...interface{}) {
	js.Global.Get("console").Call("error", fmt.Sprintf("warning: "+msg, args...))
}

func uint8ArrayToBytes(buf uintptr) []byte {
	array := js.InternalObject(buf)
	slice := make([]byte, array.Length())
	js.InternalObject(slice).Set("$array", array)
	return slice
}

func uint8ArrayToString(buf uintptr) string {
	array := js.InternalObject(buf)
	slice := make([]byte, array.Length())
	js.InternalObject(slice).Set("$array", array)
	return string(slice)
}

func toBytes(obj *js.Object) []byte {
	return js.Global.Get("Uint8Array").New(obj).Interface().([]byte)
}
