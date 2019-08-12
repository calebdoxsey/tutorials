package main

import (
	"context"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pb"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pkg/deps"
	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
	"net/url"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deps.RegisterTracer("book-worker")
	c := deps.DialJobConsumer(ctx)

	for {
		jobs, err := c.Read(ctx)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		for _, job := range jobs {
			job := job
			go func() {
				err = handle(job)
				if err != nil {
					log.Error().Err(err).Msg("failed to process job")
				}
				err = c.Ack(ctx, job)
				if err != nil {
					log.Fatal().Err(err).Send()
				}
			}()
		}
	}
}

func handle(job *pb.Job) error {
	log.Info().Interface("job", job).Msg("processing job")
	var err error
	switch job.Type {
	case pb.Job_DOWNLOAD:
		err = download(job.Book)
	case pb.Job_CALCULATE_STATS:
		err = calculateStats(job.Book)
	case pb.Job_GET_REVIEWS:
		err = getReviews(job.Book)
	default:
		err = xerrors.New("unknown job type")
	}
	return err
}

func download(book *pb.Book) error {
	u, err := url.Parse(book.Url)
	if err != nil {
		return xerrors.Errorf("invalid book url: %w", err)
	}

	log.Info().Interface("url", u).Msg("received book download job")

	return xerrors.New("not implemented")
}

func calculateStats(book *pb.Book) error {
	return xerrors.New("not implemented")
}

func getReviews(book *pb.Book) error {
	return xerrors.New("not implemented")
}
