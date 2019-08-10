package main

import (
	"context"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pkg/jobs"
	"github.com/opencensus-integrations/redigo/redis"
	"log"
)

func main() {
	log.SetFlags(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := redis.DialWithContext(ctx, "tcp", "localhost:6379")
	if err != nil {
		log.Fatalln(err)
	}

	sub, err := jobs.NewSubscriber(ctx, conn.(redis.ConnWithContext))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		job, err := sub.Receive(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("JOB", job)
	}
}
