package jobs

import (
	"context"
	"github.com/badgerodon/go-redis-streams/consumer"
	"github.com/badgerodon/go-redis-streams/producer"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pb"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
	"time"
)

const (
	groupName    = "workers"
	streamName   = "jobs"
	consumerName = "worker"
	maxLen       = 100
)

// A Consumer is a job consumer.
type Consumer struct {
	consumer    *consumer.Consumer
	outstanding map[*pb.Job]string
}

// NewConsumer creates a new Consumer.
func NewConsumer(client *redis.Client) *Consumer {
	c := consumer.New(client, groupName, consumerName,
		consumer.WithBlock(time.Second*10),
		consumer.WithStream(streamName))
	return &Consumer{
		consumer:    c,
		outstanding: make(map[*pb.Job]string),
	}
}

// Read reads a job from the job stream.
func (c *Consumer) Read(ctx context.Context) ([]*pb.Job, error) {
	for {
		msgs, err := c.consumer.Read(ctx)
		if err != nil {
			return nil, err
		}
		var jobs []*pb.Job
		for _, msg := range msgs {
			payload, ok := msg.Values["payload"]
			if !ok {
				log.Warn().Msg("invalid message in job stream: no payload detected")
				err = c.consumer.Ack(ctx, msg)
				if err != nil {
					return nil, xerrors.Errorf("error acking non-job message: %w", err)
				}
				continue
			}

			var bs []byte
			if v, ok := payload.([]byte); ok {
				bs = v
			} else if v, ok := payload.(string); ok {
				bs = []byte(v)
			} else {
				log.Warn().Msgf("invalid message in job stream: invalid payload type detected: %T", payload)
				err = c.consumer.Ack(ctx, msg)
				if err != nil {
					return nil, xerrors.Errorf("error acking non-job message: %w", err)
				}
				continue
			}

			job := new(pb.Job)
			err = proto.Unmarshal(bs, job)
			if err != nil {
				log.Warn().Err(err).Msg("invalid message in job stream: not a protobuf job")
				err = c.consumer.Ack(ctx, msg)
				if err != nil {
					return nil, xerrors.Errorf("error acking invalid job message: %w", err)
				}
				continue
			}

			// all good, so add the job to the list and the outstanding map
			c.outstanding[job] = msg.ID
			jobs = append(jobs, job)
		}
		if len(jobs) > 0 {
			return jobs, nil
		}
	}
}

// Ack marks a job as completed.
func (c *Consumer) Ack(ctx context.Context, jobs ...*pb.Job) error {
	msgs := make([]consumer.Message, 0, len(jobs))
	for _, job := range jobs {
		if id, ok := c.outstanding[job]; ok {
			msgs = append(msgs, consumer.Message{
				Stream: streamName,
				ID:     id,
			})
		}
	}
	err := c.consumer.Ack(ctx, msgs...)
	if err != nil {
		return xerrors.Errorf("error acking jobs: %w", err)
	}
	for _, job := range jobs {
		delete(c.outstanding, job)
	}
	return nil
}

// A Producer writes jobs to the job stream.
type Producer struct {
	producer *producer.Producer
}

// NewProducer creates a new Producer.
func NewProducer(client *redis.Client) *Producer {
	p := producer.New(client, streamName,
		producer.WithMaxLenApprox(maxLen))
	return &Producer{
		producer: p,
	}
}

// Write writes the job to the jobs stream.
func (p *Producer) Write(ctx context.Context, job *pb.Job) error {
	bs, err := proto.Marshal(job)
	if err != nil {
		return xerrors.Errorf("error marshaling job: %w", err)
	}
	_, err = p.producer.Write(ctx, producer.WithField("payload", bs))
	if err != nil {
		return xerrors.Errorf("error writing job to redis jobs stream: %w", err)
	}
	return nil
}
