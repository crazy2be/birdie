package birdie

import (
	"net"
	"crypto/tls"
)

type Listener struct {
	li net.Listener
}

func Listen(laddr string) (l *Listener, err error) {
	l = &Listener{}
	l.li, err = net.Listen("tcp", laddr)
	return
}

func ListenTLS(laddr string, config *tls.Config) (l *Listener, err error) {
	l = &Listener{}
	l.li, err = tls.Listen("tcp", laddr, config)
	return
}

type Conn struct {
	con net.Conn
}

func (l *Listener) Accept() (c *Conn, err error) {
	con, err := l.li.Accept()
	if err != nil {
		return nil, err
	}
	return &Conn{con: con}, nil
}

func (l *Listener) Close() (err error) {
	return l.li.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.li.Addr()
}