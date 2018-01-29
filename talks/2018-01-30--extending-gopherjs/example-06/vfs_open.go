package main

import (
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
	file, ok := vfs.files[pathname]
	if !ok {
		if flags&os.O_CREATE != 0 {
			file = new(virtualFile)
			vfs.files[pathname] = file
		} else {
			return 0, 0, syscall.ENOENT
		}
	}

	// Truncate it if we passed O_TRUNC
	if flags&os.O_TRUNC != 0 {
		file.data = nil
	}

	// Generate a file descriptor, and store a reference in the map
	fd := vfs.nextFD
	vfs.nextFD++
	vfs.fds[fd] = &virtualFileReference{
		file: file,
	}

	// Return the file descriptor
	return fd, 0, 0
}
