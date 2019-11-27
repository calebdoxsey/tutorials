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

	err = c.Send("XREAD", "COUNT", 1, "BLOCK", 1000, "STREAMS", "basic", 0)
	if err != nil {
		log.Fatalln(err)
	}
	err = c.Flush()
	if err != nil {
		log.Fatalln(err)
	}
}
