package main

import (
	"os"
	"syscall"

	"github.com/gopherjs/gopherjs/ext"
	"github.com/gopherjs/gopherjs/js"
)

var minusOne = -1

var vfs = newVirtualFileSystem()

func init() {
	ext.RegisterSyscallHandler(syscall.SYS_OPEN, vfs.Open)
	ext.RegisterSyscallHandler(syscall.SYS_CLOSE, vfs.Close)
	ext.RegisterSyscallHandler(syscall.SYS_WRITE, vfs.Write)
	ext.RegisterSyscallHandler(syscall.SYS_READ, vfs.Read)

	// ignore fcnt
	ext.RegisterSyscallHandler(syscall.SYS_FCNTL, func(a1 uintptr, a2 uintptr, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
		return 0, 0, 0
	})
}

// Our virtual file system contains files and references to files
// A file is just a slice of bytes
// A reference also tracks the position within the file

type virtualFile struct {
	data []byte
}

type virtualFileReference struct {
	file *virtualFile
	pos  int
}

type virtualFileSystem struct {
	// use maps to track files and references to files
	files  map[string]*virtualFile
	fds    map[uintptr]*virtualFileReference
	nextFD uintptr
}

func newVirtualFileSystem() *virtualFileSystem {
	return &virtualFileSystem{
		files:  make(map[string]*virtualFile),
		fds:    make(map[uintptr]*virtualFileReference),
		nextFD: 1000,
	}
}

// Open a file: http://man7.org/linux/man-pages/man2/open.2.html
//
//        int open(const char *pathname, int flags, mode_t mode);
//
func (vfs *virtualFileSystem) Open(a1 uintptr, a2 uintptr, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	pathname := uint8ArrayToString(a1)
	flags := int(a2)
	mode := os.FileMode(a3)

	js.Global.Get("console").Call("log", "::OPEN", pathname, flags, os.FileMode(mode).String())

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

// Close a file: http://man7.org/linux/man-pages/man2/close.2.html
//
//       int close(int fd);
//
func (vfs *virtualFileSystem) Close(a1 uintptr, a2 uintptr, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
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

// Read a file
//
//        ssize_t read(int fd, void *buf, size_t count);
//
func (vfs *virtualFileSystem) Read(a1 uintptr, a2 uintptr, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	fd := a1
	buf := a2
	cnt := int(a3)

	js.Global.Get("console").Call("log", "::READ", fd, buf, cnt)

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

// Write a file: http://man7.org/linux/man-pages/man2/write.2.html
//
//       ssize_t write(int fd, const void *buf, size_t count);
//
func (vfs *virtualFileSystem) Write(a1 uintptr, a2 uintptr, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	fd := a1
	buf := uint8ArrayToBytes(a2)
	cnt := a3

	js.Global.Get("console").Call("log", "::WRITE", fd, buf)

	// write to stdout/stdin
	switch fd {
	case uintptr(syscall.Stdout), uintptr(syscall.Stderr):
		js.Global.Get("document").Call("write", "<pre>"+string(buf)+"</pre>")
		return cnt, 0, 0
	}

	// find our file descriptor
	ref, ok := vfs.fds[fd]
	if !ok {
		return uintptr(minusOne), 0, syscall.EBADF
	}

	// append to the file data and move the cursor
	ref.file.data = append(ref.file.data[ref.pos:], buf...)
	ref.pos += len(buf)

	return uintptr(minusOne), 0, syscall.EACCES
}
