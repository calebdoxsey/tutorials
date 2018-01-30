package main

import (
	"syscall"

	"github.com/gopherjs/gopherjs/js"

	"github.com/gopherjs/gopherjs/ext"
)

func init() {
	ext.RegisterSyscallHandler(syscall.SYS_ACCEPT, Accept)
	ext.RegisterSyscallHandler(syscall.SYS_BIND, Bind)
	ext.RegisterSyscallHandler(syscall.SYS_CLOSE, Close)
	ext.RegisterSyscallHandler(syscall.SYS_FCNTL, FCNTL)
	ext.RegisterSyscallHandler(syscall.SYS_GETSOCKNAME, GetSockName)
	ext.RegisterSyscallHandler(syscall.SYS_LISTEN, Listen)
	ext.RegisterSyscallHandler6(syscall.SYS_SETSOCKOPT, SetSockOpt)
	ext.RegisterSyscallHandler(syscall.SYS_SOCKET, Socket)
	// SYS_CONNECT
}

var minusOne = -1

type socket struct {
	fd        uintptr
	boundPort int
}

var (
	listeners         = map[int]*socket{}
	sockets           = map[int]*socket{}
	nextFD    uintptr = 1000
)

// int accept4(int sockfd, struct sockaddr *addr, socklen_t *addrlen, int flags);
func Accept(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	info("Accept %d %d %d", a1, a2, a3)
	return 0, 0, 0
}

// int bind(int sockfd, const struct sockaddr *addr, socklen_t addrlen);
func Bind(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	fd := a1
	buf := js.InternalObject(a2)
	js.Global.Get("console").Call("log", buf)
	info("Bind fd=%d %d", fd, a3)
	return 0, 0, 0
}

func Close(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	info("Close %d %d %d", a1, a2, a3)
	return uintptr(minusOne), 0, syscall.EACCES
}

// FCNTL performs one of the operations described below on the open file descriptor fd. The operation is determined by cmd.
//
// int fcntl(int fd, int cmd, ... /* arg */ );
//
func FCNTL(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	//fd := a1
	cmd := int(a2)

	switch cmd {
	case syscall.F_GETFL:
		return 0, 0, 0
	case syscall.F_SETFD, syscall.F_SETFL:
		return 0, 0, 0
	}

	warn("fcntl cmd `%d` not supported", cmd)

	return uintptr(minusOne), 0, syscall.EACCES
}

func GetSockName(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	info("GetSockName %d %d %d", a1, a2, a3)
	return 0, 0, 0
}

func Listen(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	fd := a1
	info("Listen fd=%d %d %d", fd, a2, a3)
	return 0, 0, 0
}

func SetSockOpt(a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	info("SetSockOpt %d %d %d %d %d %d", a1, a2, a3, a4, a5, a6)
	return 0, 0, 0
}

// Socket creates an endpoint for communication and returns a descriptor.
//
//     int socket(int domain, int type, int protocol);
//
func Socket(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	domain := int(a1)
	socketType := int(a2)
	//protocol := int(a3)

	info("Socket  domain=%d type=%d", domain, socketType)

	if domain != syscall.AF_INET && domain != syscall.AF_INET6 {
		warn("socket domain `%03x` not supported", domain)
		return uintptr(minusOne), 0, syscall.EACCES
	}

	if socketType != syscall.SOCK_STREAM {
		warn("socket type `%03x` not supported", socketType)
		return uintptr(minusOne), 0, syscall.EACCES
	}

	fd := nextFD
	nextFD++

	return fd, 0, 0
}
