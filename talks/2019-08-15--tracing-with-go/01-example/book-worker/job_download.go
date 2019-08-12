package main

import (
	"github.com/calebdoxsey/tutorials/talks/2019-08-15--tracing-with-go/01-example/pb"
	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
	"net/url"
)

func (w *worker) download(book *pb.Book) error {
	log.Info().Interface("book", book).Msg("received book download job")

	u, err := url.Parse(book.Url)
	if err != nil {
		return xerrors.Errorf("invalid book url: %w", err)
	}

	switch u.Host {
	case "www.gutenberg.org":

	default:
		return xerrors.Errorf("unsupported host: %v", u.Host)
	}

	return xerrors.New("not implemented")
}
