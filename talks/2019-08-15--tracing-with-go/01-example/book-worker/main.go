package main

import (
	"context"
	"database/sql"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pb"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pkg/deps"
	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deps.RegisterTracer("book-worker")
	c := deps.DialJobConsumer(ctx)

	w := &worker{
		db: deps.DialCockroach(ctx),
	}

	for {
		jobs, err := c.Read(ctx)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		for _, job := range jobs {
			job := job
			go func() {
				err = w.handle(job)
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

type worker struct {
	db *sql.DB
}

func (w *worker) handle(job *pb.Job) error {
	log.Info().Interface("job", job).Msg("processing job")
	var err error
	switch job.Type {
	case pb.Job_DOWNLOAD:
		err = w.download(job.Book)
	case pb.Job_CALCULATE_STATS:
		err = w.calculateStats(job.Book)
	case pb.Job_GET_REVIEWS:
		err = w.getReviews(job.Book)
	default:
		err = xerrors.New("unknown job type")
	}

	status := "OK"
	if err != nil {
		status = err.Error()
	}
	if _, dberr := w.db.ExecContext(context.TODO(), `
UPSERT INTO book_job_status (book_id, job_type, status) VALUES ($1, $2, $3)
`, job.Book.Id, job.Type.String(), status); dberr != nil {
		log.Warn().Err(dberr).Msg("failed to update job status")
	}

	return err
}

func (w *worker) calculateStats(book *pb.Book) error {
	return xerrors.New("not implemented")
}

func (w *worker) getReviews(book *pb.Book) error {
	return xerrors.New("not implemented")
}
