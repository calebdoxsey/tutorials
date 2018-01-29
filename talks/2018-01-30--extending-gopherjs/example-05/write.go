package main

import (
	"syscall"

	"github.com/gopherjs/gopherjs/ext"
	"github.com/gopherjs/gopherjs/js"
)

// START OMIT

var minusOne = -1

func init() {
	ext.RegisterSyscallHandler(syscall.SYS_WRITE, func(fd, buf, count uintptr) (r1, r2 uintptr, err syscall.Errno) {
		switch fd {
		case uintptr(syscall.Stdout), uintptr(syscall.Stderr):
			js.Global.Get("document").Call("write", "<pre>"+uint8ArrayToString(buf)+"</pre>")
			return count, 0, 0
		}
		return uintptr(minusOne), 0, syscall.EACCES
	})
}

// END OMIT
