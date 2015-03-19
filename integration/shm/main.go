package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sysv_mq"
	"github.com/edsrzf/mmap-go"
)

func send(mq *sysv_mq.MessageQueue, offset, size int) error {
	var data [16]byte
	binary.BigEndian.PutUint64(data[:], uint64(offset))
	binary.BigEndian.PutUint64(data[8:], uint64(size))
	return mq.SendBytes(data[:], 1, 0)
}

func recv(mq *sysv_mq.MessageQueue) (offset, size int, err error) {
	data, _, e := mq.ReceiveBytes(1, 0)
	if err != nil {
		err = e
		return
	}
	if len(data) < 16 {
		err = fmt.Errorf("expected offset and size")
		return
	}
	offset = int(binary.BigEndian.Uint64(data[:8]))
	size = int(binary.BigEndian.Uint64(data[8:]))
	return
}

func server(mq *sysv_mq.MessageQueue) {
	fd, err := shm_open("/shm-example", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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

	for {
		offset, sz, err := recv(mq)
		if err != nil {
			panic(err)
		}
		log.Println("recv", offset, sz)

		log.Println(string(m[offset : offset+sz]))
	}
}

func client(mq *sysv_mq.MessageQueue) {
	fd, err := shm_open("/shm-example", os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	m, err := mmap.Map(fd, mmap.RDWR|mmap.EXEC, 0)
	if err != nil {
		panic(err)
	}

	offset := 0
	for i := 0; i < 3; i++ {
		data := "Hello World"
		copy(m[offset:], data)

		log.Println("sending", offset, len(data))
		send(mq, offset, len(data))
		offset += len(data)
	}
}

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatalln("expected mode")
	}

	mq, err := sysv_mq.NewMessageQueue(&sysv_mq.QueueConfig{
		Key:     1001,
		MaxSize: 8 * 1024,
		Mode:    sysv_mq.IPC_CREAT | 0600,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer mq.Close()

	switch os.Args[1] {
	case "server":
		server(mq)
	case "client":
		client(mq)
	}
}
