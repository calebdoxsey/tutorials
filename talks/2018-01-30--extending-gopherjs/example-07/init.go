package main

import (
	"syscall"

	"github.com/gopherjs/gopherjs/ext"
)

// START OMIT

var vfs interface {
	Close(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
	Open(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
	Read(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
	Write(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
} = newVirtualFileSystem()

func init() {
	ext.RegisterSyscallHandler(syscall.SYS_CLOSE, vfs.Close)
	ext.RegisterSyscallHandler(syscall.SYS_OPEN, vfs.Open)
	ext.RegisterSyscallHandler(syscall.SYS_READ, vfs.Read)
	ext.RegisterSyscallHandler(syscall.SYS_WRITE, vfs.Write)

	// ignore fcntl
	ext.RegisterSyscallHandler(syscall.SYS_FCNTL, func(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
		return 0, 0, 0
	})
}

// END OMIT
