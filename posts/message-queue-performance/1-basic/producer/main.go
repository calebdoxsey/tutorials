package main

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	err = c.Send("XADD", "basic", "*", "message", "example")
	if err != nil {
		log.Fatalln(err)
	}
	err = c.Flush()
	if err != nil {
		log.Fatalln(err)
	}
}
