package main

import (
	"errors"
	"os"
	"syscall"

	"github.com/gopherjs/gopherjs/js"
)

// Open a file: http://man7.org/linux/man-pages/man2/open.2.html
//
//        int open(const char *pathname, int flags, mode_t mode);
//
func (vfs *virtualFileSystem) Open(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	pathname := uint8ArrayToString(a1)
	flags := int(a2)
	mode := os.FileMode(a3)

	js.Global.Get("console").Call("log", "::OPEN", pathname, flags, os.FileMode(mode).String()) // OMIT

	// See if the file exists, if it doesn't, and we passed O_CREATE, create it

	// START OMIT

	type Result struct {
		res *js.Object
		err error
	}
	c := make(chan Result, 1)
	tx := vfs.db.Call("transaction", js.S{"files"}, "readonly")
	req := tx.Call("objectStore", "files").Call("get", pathname)
	req.Set("onsuccess", func(evt *js.Object) {
		c <- Result{
			res: evt.Get("target").Get("result"),
		}
	})
	req.Set("onerror", func(evt *js.Object) {
		c <- Result{
			err: errors.New(evt.Get("target").Get("error").String()),
		}
	})
	res := <-c
	if res.err != nil {
		return 0, 0, syscall.EACCES
	}

	// END OMIT

	ref := &virtualFileReference{
		path: pathname,
	}
	if bs, ok := res.res.Interface().([]byte); ok && bs != nil {
		ref.data = bs
	} else {
		if flags&os.O_CREATE == 0 {
			return 0, 0, syscall.ENOENT
		}
	}

	// Truncate it if we passed O_TRUNC
	if flags&os.O_TRUNC != 0 {
		ref.data = nil
	}

	// Generate a file descriptor, and store a reference in the map
	fd := vfs.nextFD
	vfs.nextFD++
	vfs.fds[fd] = ref

	// Return the file descriptor
	return fd, 0, 0
}
