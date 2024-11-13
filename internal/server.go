package internal

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	"log"
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
	dbInstances DBInstances
	eventStore  EventStore
}

type DBInstances struct {
	wg      sync.WaitGroup
	dbMutex sync.RWMutex
	store   map[string]*pebble.DB
}

type EventWriter struct {
	wg     sync.WaitGroup
	mu     sync.RWMutex
	handle *os.File
}

func (ew *EventWriter) Write(b []byte) {
	ew.wg.Add(1)
	go func([]byte) {
		defer ew.wg.Done()
		ew.mu.Lock()
		ew.handle.Write(b)
		ew.handle.Write([]byte("\n"))
		ew.handle.Sync()
		log.Println("writing ", b)
		ew.mu.Unlock()
	}(b)
}

type EventStore struct {
	wg     sync.WaitGroup
	mu     sync.RWMutex
	writer map[string]*EventWriter
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

func newTCPServer(address string) (*Server, error) {
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
