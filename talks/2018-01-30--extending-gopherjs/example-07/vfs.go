package main

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

var minusOne = -1

// Our virtual file system contains files and references to files
// A file is just a slice of bytes
// A reference also tracks the position within the file

type (
	virtualFileReference struct {
		changed bool
		path    string
		data    []byte
		pos     int
	}
	virtualFileSystem struct {
		db     *js.Object
		fds    map[uintptr]*virtualFileReference
		nextFD uintptr
	}
)

func newVirtualFileSystem() *virtualFileSystem {
	// START OMIT
	type Result struct {
		vfs *virtualFileSystem
		err error
	}
	c := make(chan Result, 1)
	req := js.Global.Get("indexedDB").Call("open", "vfs")
	req.Set("onerror", func(evt *js.Object) { // OMIT
		c <- Result{err: errors.New(evt.String())} // OMIT
	}) // OMIT
	req.Set("onupgradeneeded", func(evt *js.Object) {
		db := evt.Get("target").Get("result")
		db.Call("createObjectStore", "files")
	})
	req.Set("onsuccess", func(evt *js.Object) {
		c <- Result{vfs: &virtualFileSystem{
			db:     req.Get("result"),
			fds:    make(map[uintptr]*virtualFileReference),
			nextFD: 1000,
		}}
	})
	res := <-c
	if res.err != nil {
		panic(res.err)
	}
	// END OMIT
	return res.vfs
}
