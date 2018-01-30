package main

import (
	"context"
	"io"
	"sync"
	"time"
)

// A ChannelReader reads bytes from a channel and buffers them
type ChannelReader struct {
	c        <-chan []byte
	buf      []byte
	deadline time.Time
}

// NewChannelReader creates a new ChannelReader
func NewChannelReader(c <-chan []byte) *ChannelReader {
	return &ChannelReader{
		c: c,
	}
}

// Read reads from the channel. It should not be called by multiple goroutines
func (r *ChannelReader) Read(b []byte) (sz int, err error) {
	if len(b) == 0 {
		return 0, io.ErrShortBuffer
	}

	for {
		if len(r.buf) > 0 {
			if len(r.buf) <= len(b) {
				sz = len(r.buf)
				copy(b, r.buf)
				r.buf = nil
			} else {
				copy(b, r.buf)
				r.buf = r.buf[len(b):]
				sz = len(b)
			}
			return sz, nil
		}

		var ok bool
		if r.deadline.IsZero() {
			r.buf, ok = <-r.c
		} else {
			timer := time.NewTimer(r.deadline.Sub(time.Now()))
			defer timer.Stop()

			select {
			case r.buf, ok = <-r.c:
			case <-timer.C:
				return 0, context.DeadlineExceeded
			}
		}
		if len(r.buf) == 0 && !ok {
			return 0, io.EOF
		}
	}
}

// SetDeadline sets the deadline to read to the channel
func (r *ChannelReader) SetDeadline(deadline time.Time) {
	r.deadline = deadline
}

// A ChannelWriter writes slices of bytes to a channel
type ChannelWriter struct {
	c        chan<- []byte
	deadline time.Time
}

// NewChannelWriter creates a new ChannelWriter
func NewChannelWriter(c chan<- []byte) *ChannelWriter {
	return &ChannelWriter{
		c: c,
	}
}

func (w *ChannelWriter) Close() error {
	c := w.c
	w.c = nil
	if c != nil {
		close(c)
	}
	return nil
}

// Write writes p to the channel
func (w *ChannelWriter) Write(b []byte) (sz int, err error) {
	select {
	case w.c <- b:
		return len(b), nil
	default:
	}

	if w.deadline.IsZero() {
		w.c <- b
		return len(b), nil
	}

	timer := time.NewTimer(w.deadline.Sub(time.Now()))
	defer timer.Stop()

	select {
	case w.c <- b:
		return len(b), nil
	case <-timer.C:
		return 0, context.DeadlineExceeded
	}
}

// SetDeadline sets the deadline to write to the channel
func (w *ChannelWriter) SetDeadline(t time.Time) {
	w.deadline = t
}

type queue struct {
	signal chan struct{}

	mu sync.Mutex
	q  [][]byte
}

func newQueue() *queue {
	return &queue{
		signal: make(chan struct{}, 1),
	}
}

func (q *queue) enqueue(msg []byte) {
	q.mu.Lock()
	q.q = append(q.q, msg)
	q.mu.Unlock()

	select {
	case q.signal <- struct{}{}:
	default:
	}
}

func (q *queue) dequeue(ctx context.Context) ([]byte, error) {
	for {
		var msg []byte
		q.mu.Lock()
		if len(q.q) > 0 {
			msg = q.q[0]
			q.q = q.q[1:]
		}
		q.mu.Unlock()

		if msg != nil {
			return msg, nil
		}

		select {
		case <-q.signal:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

type inOrderQueue struct {
	signal chan struct{}

	mu sync.Mutex
	q  map[uint16][]byte
}

func newInOrderQueue() *inOrderQueue {
	return &inOrderQueue{
		signal: make(chan struct{}, 1),
		q:      make(map[uint16][]byte),
	}
}

func (q *inOrderQueue) enqueue(id uint16, msg []byte) {
	q.mu.Lock()
	q.q[id] = msg
	q.mu.Unlock()

	select {
	case q.signal <- struct{}{}:
	default:
	}
}

func (q *inOrderQueue) dequeue(ctx context.Context, id uint16) ([]byte, error) {
	for {
		q.mu.Lock()
		msg, ok := q.q[id]
		if ok {
			delete(q.q, id)
		}
		q.mu.Unlock()

		if ok {
			return msg, nil
		}

		select {
		case <-q.signal:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

type inOrderQueueReader struct {
	q       *inOrderQueue
	buf     []byte
	counter uint16
}

func newInOrderQueueReader(q *inOrderQueue) *inOrderQueueReader {
	return &inOrderQueueReader{
		q: q,
	}
}

func (r *inOrderQueueReader) Read(ctx context.Context, b []byte) (sz int, err error) {
	if len(b) == 0 {
		return 0, nil
	}

	for {
		if len(r.buf) > 0 {
			if len(r.buf) <= len(b) {
				sz = len(r.buf)
				copy(b, r.buf)
				r.buf = nil
			} else {
				sz = len(b)
				copy(b, r.buf)
				r.buf = r.buf[len(b):]
			}
			return sz, nil
		}

		var err error
		r.buf, err = r.q.dequeue(ctx, r.counter)
		if err != nil {
			return 0, err
		}
		r.counter++
	}
}
