package jobs

import (
	"context"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pb"
	"github.com/golang/protobuf/proto"
	"golang.org/x/xerrors"
)

// RedisConn is a redis connection.
type RedisConn interface {
	SendContext(ctx context.Context, commandName string, args ...interface{}) error
	FlushContext(context.Context) error
	ReceiveContext(context.Context) (reply interface{}, err error)
}

// A JobType is the type of job.
type JobType string

// Job Types
const (
	JobTypeDownload       JobType = "download"
	JobTypeCalculateStats JobType = "calculate-stats"
	JobTypeGetReviews     JobType = "get-reviews"
)

// A Subscriber subscribes to jobs in the job queue.
type Subscriber struct {
	conn RedisConn
}

// NewSubscriber creates a new Subscriber.
func NewSubscriber(ctx context.Context, conn RedisConn) (*Subscriber, error) {
	sub := &Subscriber{conn: conn}

	err := sub.conn.SendContext(ctx, "SUBSCRIBE", "jobs")
	if err != nil {
		return nil, xerrors.Errorf("failed to send SUBSCRIBE command for job queue: %w", err)
	}

	err = sub.conn.FlushContext(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to flush SUBSCRIBE command for job queue: %w", err)
	}

	return sub, nil
}

// Receive receives a job from the job queue.
func (sub *Subscriber) Receive(ctx context.Context) (*pb.Job, error) {
	for {
		reply, err := sub.conn.ReceiveContext(ctx)
		if err != nil {
			return nil, xerrors.Errorf("failed to receive from job queue: %w", err)
		}

		switch msg := reply.(type) {
		case []interface{}:
			if len(msg) != 3 {
				return nil, xerrors.Errorf("expected array of 3 elements as reply from subscribe, got: %v", msg)
			}

			replyType, ok := msg[0].([]byte)
			if !ok {
				return nil, xerrors.Errorf("expected first element of array to be the reply type, but got: %T", msg[0])
			}

			if string(replyType) == "message" {
				payload, ok := msg[2].([]byte)
				if !ok {
					return nil, xerrors.Errorf("expected payload to be a []byte, but got: %T", msg[2])
				}
				var job pb.Job
				err = proto.Unmarshal(payload, &job)
				if err != nil {
					return nil, xerrors.Errorf("failed to unmarshal job in job queue: %w", err)
				}
				return &job, nil
			}
		default:
			return nil, xerrors.Errorf("unexpected message type (%T) in job queue: %v", msg, msg)
		}
	}
}

// Close unsubscribes.
func (sub *Subscriber) Close() error {
	err := sub.conn.SendContext(context.Background(), "UNSUBSCRIBE", "jobs")
	if err != nil {
		return xerrors.Errorf("failed to send UNSUBSCRIBE request to job queue: %w", err)
	}
	err = sub.conn.FlushContext(context.Background())
	if err != nil {
		return xerrors.Errorf("failed to flush UNSUBSCRIBE request to job queue: %w", err)
	}
	return nil
}

// A Publisher publishes jobs to the job queue.
type Publisher struct {
	conn RedisConn
}

// NewPublisher creates a new Publisher.
func NewPublisher(conn RedisConn) *Publisher {
	return &Publisher{
		conn: conn,
	}
}

// Publish submits a job to the job queue.
func (pub *Publisher) Publish(ctx context.Context, jobType JobType, book *pb.Book) error {
	job := &pb.Job{
		Type: string(jobType),
		Book: book,
	}

	data, err := proto.Marshal(job)
	if err != nil {
		return xerrors.Errorf("failed to marshal job: %w", err)
	}

	err = pub.conn.SendContext(ctx, "PUBLISH", "jobs", data)
	if err != nil {
		return xerrors.Errorf("failed to publish job: %w", err)
	}

	err = pub.conn.FlushContext(ctx)
	if err != nil {
		return xerrors.Errorf("failed to flush job: %w", err)
	}

	return nil
}
