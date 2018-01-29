package main

import (
	"errors"
	"syscall"

	"github.com/gopherjs/gopherjs/js"
)

// Close a file: http://man7.org/linux/man-pages/man2/close.2.html
//
//       int close(int fd);
//
func (vfs *virtualFileSystem) Close(a1, a2, a3 uintptr) (r1, r2 uintptr, errno syscall.Errno) {
	fd := a1

	js.Global.Get("console").Call("log", "::CLOSE", fd)

	// See if the file descriptor exists. If it doesn't, return an error
	ref, ok := vfs.fds[fd]
	if !ok {
		return uintptr(minusOne), 0, syscall.EBADF
	}

	// START OMIT

	// flush the data to the db
	if ref.changed {
		type Result struct {
			err error
		}
		c := make(chan Result, 1)
		tx := vfs.db.Call("transaction", js.S{"files"}, "readwrite")
		req := tx.Call("objectStore", "files").Call("put", ref.data, ref.path)
		req.Set("onsuccess", func(evt *js.Object) {
			c <- Result{}
		})
		req.Set("onerror", func(evt *js.Object) {
			c <- Result{
				err: errors.New(evt.Get("target").Get("error").String()),
			}
		})
		res := (<-c)
		if res.err != nil {
			return 0, 0, syscall.EACCES
		}
	}

	// END OMIT

	// Close the file descriptor by removing it
	delete(vfs.fds, fd)

	return 0, 0, 0
}
