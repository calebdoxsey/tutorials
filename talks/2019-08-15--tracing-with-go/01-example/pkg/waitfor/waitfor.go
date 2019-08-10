package waitfor

import (
	"context"
	"net"
	"time"

	"golang.org/x/xerrors"
)

// TCP will wait for a connection to the given tcp server.
func TCP(ctx context.Context, addr string) error {
	for {
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
		if xerrors.Is(err, context.Canceled) || xerrors.Is(err, context.DeadlineExceeded) {
			return err
		} else if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}
}
