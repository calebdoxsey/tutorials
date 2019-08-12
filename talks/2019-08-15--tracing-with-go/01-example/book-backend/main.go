package main

import (
	"context"
	"database/sql"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pb"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pkg/deps"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pkg/jobs"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type server struct {
	p  *jobs.Producer
	db *sql.DB
}

func (s server) AddBook(ctx context.Context, req *pb.AddBookRequest) (*pb.AddBookResponse, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
INSERT INTO book (url) VALUES ($1) RETURNING id;
`, req.Url).Scan(&id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "error inserting book into database: %v", err)
	}

	err = s.p.Write(ctx, &pb.Job{
		Type: pb.Job_DOWNLOAD,
		Book: &pb.Book{Id: id, Url: req.Url},
	})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "error publishing download job: %v", err)
	}

	return &pb.AddBookResponse{Id: id}, nil
}

func (s server) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	res := &pb.GetBookResponse{Book: new(pb.Book)}
	err := s.db.QueryRowContext(ctx, `
SELECT id, url FROM book WHERE id = $1
`, req.Id).Scan(&res.Book.Id, &res.Book.Url)
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "no book found with that id")
	} else if err != nil {
		return nil, status.Errorf(codes.Unknown, "error getting book: %v", err)
	}

	var stats pb.BookStats
	err = s.db.QueryRowContext(ctx, `
SELECT number_of_words, longest_word FROM book_stat WHERE book_id = $1
`, res.Book.Id).Scan(&stats.NumberOfWords, &stats.LongestWord)
	if err == sql.ErrNoRows {

	} else if err != nil {
		return nil, status.Errorf(codes.Unknown, "error getting book stats: %v", err)
	} else {
		res.Stats = &stats
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT username, review FROM book_review WHERE book_id = $1
`, res.Book.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "error getting book reviews: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var review pb.BookReview
		err := rows.Scan(&review.Username, &review.Review)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "error listing books: %v", err)
		}
		res.Reviews = append(res.Reviews, &review)
	}
	err = rows.Err()
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "error listing books: %v", err)
	}

	return res, nil
}

func (s server) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	res := &pb.ListBooksResponse{}
	rows, err := s.db.QueryContext(ctx, `
SELECT id, url FROM book
`)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "error listing books: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var book pb.Book
		err := rows.Scan(&book.Id, &book.Url)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "error listing books: %v", err)
		}
		res.Books = append(res.Books, &book)
	}

	if rows.Err() != nil {
		return nil, status.Errorf(codes.Unknown, "error listing books: %v", rows.Err())
	}

	return res, nil
}

func main() {
	log.SetFlags(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deps.RegisterTracer("book-backend")
	db := deps.DialCockroach(ctx)
	p := deps.DialJobProducer(ctx)

	addr := "127.0.0.1:5100"
	log.Println("starting book-backend on", addr)
	li, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	srv := deps.NewGRPCServer()
	pb.RegisterBookBackendServer(srv, &server{
		p:  p,
		db: db,
	})
	err = srv.Serve(li)
	if err != nil {
		log.Fatalln(err)
	}
}
