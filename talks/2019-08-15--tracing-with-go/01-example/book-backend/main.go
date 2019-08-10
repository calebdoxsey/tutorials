package main

import (
	"context"
	"database/sql"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pb"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pkg/jobs"
	"github.com/opencensus-integrations/redigo/redis"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pub *jobs.Publisher
	db  *sql.DB
}

func (s server) AddBook(ctx context.Context, req *pb.AddBookRequest) (*pb.AddBookResponse, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
INSERT INTO book (url) VALUES (?) RETURNING id;
`).Scan(&id)
	if err != nil {
		return nil, xerrors.Errorf("error inserting book into database: %w", err)
	}

	return &pb.AddBookResponse{Id: id}, nil
}

func (s server) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	panic("implement me")
}

func (s server) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	panic("implement me")
}

func main() {
	log.SetFlags(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := redis.DialWithContext(ctx, "tcp", "localhost:6379")
	if err != nil {
		log.Fatalln(err)
	}

	pub := jobs.NewPublisher(conn.(redis.ConnWithContext))

	addr := "127.0.0.1:5100"
	li, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	srv := grpc.NewServer()
	pb.RegisterBookBackendServer(srv, &server{pub: pub})
	err = srv.Serve(li)
	if err != nil {
		log.Fatalln(err)
	}
}
