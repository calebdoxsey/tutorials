package main

import (
	"syscall"

	"github.com/gopherjs/gopherjs/js"
)

// Read a file
//
//        ssize_t read(int fd, void *buf, size_t count);
//
func (vfs *virtualFileSystem) Read(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	fd := a1
	buf := a2
	cnt := int(a3)

	js.Global.Get("console").Call("log", "::READ", fd, buf, cnt) // OMIT

	// find our file descriptor
	ref, ok := vfs.fds[fd]
	if !ok {
		return uintptr(minusOne), 0, syscall.EBADF
	}

	// copy the data in the file into the buffer
	for i := 0; i < cnt && ref.pos < len(ref.file.data); i++ {
		js.InternalObject(buf).SetIndex(i, ref.file.data[ref.pos])
		r1++
		ref.pos++
	}

	return r1, 0, 0
}
