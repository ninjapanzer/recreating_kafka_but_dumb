package internal

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Server struct {
	wg          sync.WaitGroup
	listener    net.Listener
	shutdown    chan struct{}
	connection  chan net.Conn
	logFile     *os.File
	dbInstances DBCache
	eventStore  EventStore
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	err = tc.SetKeepAlive(true)
	err = tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, err
}

func newPersistentServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", address, err)
	}

	return &Server{
		listener:   tcpKeepAliveListener{listener.(*net.TCPListener)},
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
	}, nil
}

func newServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", address, err)
	}

	return &Server{
		listener:   listener,
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
	}, nil
}

func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			s.connection <- conn
		}
	}
}

func (s *Server) handleConnections(handler func(conn net.Conn)) {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			go handler(conn)
		}
	}
}
