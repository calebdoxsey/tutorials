package main

import (
	"syscall"

	"github.com/gopherjs/gopherjs/ext"
)

var vfs = newVirtualFileSystem()

// START OMIT

func init() {
	ext.RegisterSyscallHandler(syscall.SYS_OPEN, vfs.Open)
	ext.RegisterSyscallHandler(syscall.SYS_CLOSE, vfs.Close)
	ext.RegisterSyscallHandler(syscall.SYS_WRITE, vfs.Write)
	ext.RegisterSyscallHandler(syscall.SYS_READ, vfs.Read)

	// ignore fcnt
	ext.RegisterSyscallHandler(syscall.SYS_FCNTL, func(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
		return 0, 0, 0
	})
}

// END OMIT
