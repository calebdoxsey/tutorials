package main

/*
#cgo LDFLAGS: -lrt
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <stdlib.h>
*/
import "C"
import (
	"os"
	"unsafe"
)

func shm_open(name string, oflag, mode int64) (*os.File, error) {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))
	c_oflag := C.int(oflag)
	c_mode := C.mode_t(mode)
	c_res, err := C.shm_open(c_name, c_oflag, c_mode)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(c_res), name), nil
}

func shm_unlink(name string) error {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))
	_, err := C.shm_unlink(c_name)
	return err
}
