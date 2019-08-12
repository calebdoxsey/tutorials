package main

import (
	"context"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pb"
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pkg/deps"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/status"
	"html/template"
	"log"
	"net/http"
)

//go:generate go run github.com/mjibson/esc -o tpl_gen.go -pkg main tpl

func main() {
	log.SetFlags(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deps.RegisterTracer("book-www")

	cc := deps.DialGRPC(ctx, "localhost:5100")
	srv := &server{
		client: pb.NewBookBackendClient(cc),
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", srv.index)

	h := nethttp.Middleware(opentracing.GlobalTracer(), r)

	addr := "127.0.0.1:5000"
	log.Println("starting book-www on", addr)
	err := http.ListenAndServe(addr, h)
	if err != nil {
		log.Fatalln(err)
	}
}

var (
	index = template.Must(template.New("index.gohtml").Parse(FSMustString(false, "/tpl/index.gohtml")))
)

type server struct {
	client pb.BookBackendClient
}

func (srv *server) index(w http.ResponseWriter, r *http.Request) {
	res, err := srv.client.ListBooks(r.Context(), &pb.ListBooksRequest{})
	if err != nil {
		http.Error(w, err.Error(), runtime.HTTPStatusFromCode(status.Code(err)))
		return
	}

	err = index.Execute(w, res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
