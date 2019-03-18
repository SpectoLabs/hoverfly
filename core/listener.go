package hoverfly

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// StoppableListener - wrapper for tcp listener that can stop
type StoppableListener struct {
	*net.TCPListener
	stop chan int
}

// NewStoppableListener returns new StoppableListener listener
func NewStoppableListener(l net.Listener) (*StoppableListener, error) {
	tcpListener, ok := l.(*net.TCPListener)

	if !ok {
		return nil, fmt.Errorf("failed to wrap listener")
	}

	sl := &StoppableListener{}
	sl.TCPListener = tcpListener
	sl.stop = make(chan int)

	return sl, nil
}

// Accept - TCPListener waits for the next call, implements default interface method
func (sl *StoppableListener) Accept() (net.Conn, error) {
	for {
		sl.SetDeadline(time.Now().Add(time.Second))

		newConn, err := sl.TCPListener.Accept()

		select {
		case <-sl.stop:
			log.Debug("Stopping listener")
			return nil, fmt.Errorf("Stopping listener")
		default:
			// continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

// Stop - stops listener
func (sl *StoppableListener) Stop() {
	close(sl.stop)
}
