package main

import "github.com/gopherjs/gopherjs/js"

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
