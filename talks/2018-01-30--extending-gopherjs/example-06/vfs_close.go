package main

import (
	"syscall"

	"github.com/gopherjs/gopherjs/js"
)

// Close a file: http://man7.org/linux/man-pages/man2/close.2.html
//
//       int close(int fd);
//
func (vfs *virtualFileSystem) Close(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	fd := a1

	js.Global.Get("console").Call("log", "::CLOSE", fd)

	// See if the file descriptor exists. If it doesn't, return an error
	_, ok := vfs.fds[fd]
	if !ok {
		return uintptr(minusOne), 0, syscall.EBADF
	}

	// Close the file descriptor by removing it
	delete(vfs.fds, fd)
	return 0, 0, 0
}
