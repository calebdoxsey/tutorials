package beanstalk

import (
	"context"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pkg/waitfor"
)

func Dial() (*beanstalk.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := waitfor.TCP(ctx, "localhost:11300")
	if err != nil {
		return nil, err
	}
	return beanstalk.Dial("tcp", "localhost:11300")
}
