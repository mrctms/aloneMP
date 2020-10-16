package util

import (
	"net"
	"time"
)

type TcpConn struct {
	net.Conn
}

func (t *TcpConn) Read(b []byte) (n int, err error) {
	t.Conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer t.resetDeadLine()
	return t.Conn.Read(b)
}

func (t *TcpConn) Write(b []byte) (n int, err error) {
	t.Conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer t.resetDeadLine()
	return t.Conn.Write(b)
}

func (t *TcpConn) resetDeadLine() {
	var zero time.Time
	t.Conn.SetDeadline(zero)
}
