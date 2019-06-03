package server

import (
	"net"
	"sync"
	"time"
)

// LimitListener returns a Listener that accepts at most n simultaneous
// connections from the provided Listener.
func LimitListener(l net.Listener, n int) net.Listener {
	return &limitListener{l, make(chan struct{}, n)}
}

type limitListener struct {
	net.Listener
	sem chan struct{}
}

func (l *limitListener) acquire() { l.sem <- struct{}{} }
func (l *limitListener) release() { <-l.sem }

func (l *limitListener) Accept() (net.Conn, error) {
	l.acquire()
	c, err := l.Listener.Accept()
	if err != nil {
		l.release()
		return nil, err
	}

	if tc, ok := c.(*net.TCPConn); ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(3 * time.Minute)
	}

	return &limitListenerConn{Conn: c, release: l.release}, nil
}

type limitListenerConn struct {
	net.Conn
	releaseOnce sync.Once
	release     func()
}

func (l *limitListenerConn) Close() error {
	err := l.Conn.Close()
	l.releaseOnce.Do(l.release)
	return err
}
