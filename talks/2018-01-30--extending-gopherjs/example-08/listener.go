package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/gopherjs/gopherjs/ext"
	"github.com/gopherjs/gopherjs/js"
)

var listeners map[int]net.Listener

func init() {
	ext.RegisterListenFunc(func(network, laddr string) (net.Listener, error) {
		switch network {
		case "tcp", "tcp4":
		default:
			return nil, errors.New("unsupported network")
		}

		host, port, err := net.SplitHostPort(laddr)
		if err != nil {
			return nil, err
		}

		switch host {
		case "127.0.0.1", "localhost", "0.0.0.0":
		default:
			return nil, errors.New("only localhost addresses are supported")
		}

		iport, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}

		return newMessagePortListener()

		warn("not supported %s %d", host, port)
		panic(errors.New("network access is not supported by GopherJS"))
	})
}

type messagePortListener struct {
	messagePort *js.Object
	networkPort int
	incoming    chan net.Conn
	onClose     func()
}

func newMessagePortListener(networkPort int, messagePort *js.Object, onClose func()) *messagePortListener {
	l := &messagePortListener{
		messagePort: messagePort,
		networkPort: networkPort,
		incoming:    make(chan net.Conn, 64),
		onClose:     onClose,
	}
	l.messagePort.Set("onmessage", func(evt *js.Object) {
		method := evt.Get("data").Index(0).String()
		switch method {
		case "close":
			l.Close()
		case "connection":
			connPort := int(evt.Get("data").Index(1).Int64())
			connMessagePort := evt.Get("data").Index(2)
			conn := newAckedMessagePortConn(l.networkPort, connPort, connMessagePort, nil)
			select {
			case l.incoming <- conn:
			default:
				conn.Close()
			}
		default:
			panic(fmt.Sprintf("method %s not implemented", method))
		}
	})
	return l
}

func (l *messagePortListener) Accept() (net.Conn, error) {
	conn := <-l.incoming
	return conn, nil
}

func (l *messagePortListener) Close() error {
	if l.messagePort != nil {
		l.messagePort.Call("postMessage", []interface{}{
			"close",
		})
		l.messagePort.Call("close")
	}
	l.messagePort = nil
	if l.onClose != nil {
		l.onClose()
		l.onClose = nil
	}
	return nil
}

func (l *messagePortListener) Addr() net.Addr {
	return addr{l.networkPort}
}
