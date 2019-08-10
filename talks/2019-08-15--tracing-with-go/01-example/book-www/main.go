package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	beanstalkConn *beanstalk.Conn
)

func main() {
	log.SetFlags(0)

	conn, err := beanstalk.Dial("tcp", "localhost:11300")
	if err != nil {
		log.Fatalln(err)
	}
	beanstalkConn = conn

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", index)

	addr := "127.0.0.1:5000"
	log.Println("starting", os.Args[0], "on", addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalln(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	jobID, err := (&beanstalk.Tube{
		Conn: beanstalkConn,
		Name: "book-stats",
	}).Put([]byte("HELLO"), 0, 0, time.Minute*10)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(jobID)
}
