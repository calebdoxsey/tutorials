package main

import (
	"syscall"

	"github.com/gopherjs/gopherjs/js"
)

// START OMIT

// Write a file: http://man7.org/linux/man-pages/man2/write.2.html
//
//       ssize_t write(int fd, const void *buf, size_t count);
//
func (vfs *virtualFileSystem) Write(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	fd := a1
	buf := uint8ArrayToBytes(a2)
	cnt := a3
	// OMIT
	js.Global.Get("console").Call("log", "::WRITE", fd, buf) // OMIT
	// OMIT
	// write to stdout/stdin    OMIT
	switch fd { // OMIT
	case uintptr(syscall.Stdout), uintptr(syscall.Stderr): // OMIT
		js.Global.Get("document").Call("write", "<pre>"+string(buf)+"</pre>") // OMIT
		return cnt, 0, 0                                                      // OMIT
	} // OMIT

	// find our file descriptor
	ref, ok := vfs.fds[fd]
	if !ok {
		return uintptr(minusOne), 0, syscall.EBADF
	}

	// append to the file data and move the cursor
	ref.file.data = append(ref.file.data[ref.pos:], buf...)
	ref.pos += len(buf)

	return cnt, 0, 0
}

// END OMIT
