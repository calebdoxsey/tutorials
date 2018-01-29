package main

import (
	"fmt"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

// START OMIT

func main() {
	for range time.Tick(time.Second) {
		fmt.Println("Hello World")
	}
}

// END OMIT

func uint8ArrayToString(buf uintptr) string {
	array := js.InternalObject(buf)
	slice := make([]byte, array.Length())
	js.InternalObject(slice).Set("$array", array)
	return string(slice)
}
