package main

import (
	"log"
	"os"

	"github.com/edsrzf/mmap-go"
)

//#include <fcntl.h>
import "C"

func server() {
	sem, err := sem_open("/shm-example", O_CREAT|O_EXCL, 0777, 0)
	if err != nil {
		panic(err)
	}
	defer sem_close(sem)
	defer sem_unlink("/shm-example")

	fd, err := shm_open("/shm-example", C.O_RDWR|C.O_CREAT|C.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	defer shm_unlink("/shm-example")

	err = fd.Truncate(256)
	if err != nil {
		panic(err)
	}

	m, err := mmap.Map(fd, mmap.RDWR|mmap.EXEC, 0)
	if err != nil {
		panic(err)
	}

	err = sem_wait(sem)
	if err != nil {
		panic(err)
	}

	log.Println(m)
}

func client() {
	sem, err := sem_open("/shm-example", 0, 0777, 1)
	if err != nil {
		panic(err)
	}
	defer sem_close(sem)

	fd, err := shm_open("/shm-example", C.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	m, err := mmap.Map(fd, mmap.RDWR|mmap.EXEC, 0)
	if err != nil {
		panic(err)
	}

	copy(m, "Hello World")
	sem_post(sem)
}

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatalln("expected mode")
	}
	switch os.Args[1] {
	case "server":
		server()
	case "client":
		client()
	}
}
