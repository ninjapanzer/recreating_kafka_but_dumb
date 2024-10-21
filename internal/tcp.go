package internal

import (
	"context"
	"fmt"
	"github.com/cockroachdb/pebble"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type server struct {
	wg         sync.WaitGroup
	listener   net.Listener
	shutdown   chan struct{}
	connection chan net.Conn
	db         *pebble.DB
	logFile    *os.File
}

func newServer(address string) (*server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", address, err)
	}

	return &server{
		listener:   listener,
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
	}, nil
}

func (s *server) acceptConnections() {
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

func (s *server) handleConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			go s.handleConnection(conn)
		}
	}
}

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 1024)
	count, err := conn.Read(buff)
	request := []byte(strings.TrimSpace(string(buff[:count])))
	if err != nil {
		log.Fatalf("Failed")
	}
	out := append(request, ' ')
	out = append(out, []byte(time.Now().String())...)
	s.db.Set([]byte("key"), out, pebble.NoSync)
	logOut := append([]byte("key "), out...)
	logOut = append(logOut, '\n')
	s.logFile.Write(logOut)
	fmt.Fprintf(conn, string(request))
}

func (s *server) Start(db *pebble.DB, logFile *os.File) {
	s.wg.Add(2)
	s.db = db
	s.logFile = logFile
	go s.acceptConnections()
	go s.handleConnections()
}

func (s *server) Stop() {
	close(s.shutdown)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		fmt.Println("Timed out waiting for connections to finish.")
		return
	}
}

func Boot() {
	s, err := newServer(":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	opts := pebble.Options{}
	db, err := pebble.Open("demo", &opts)
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile("demo/event.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close() // Ensure the file is closed when done

	s.Start(db, logFile)

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down server...")
	s.Stop()
	fmt.Println("Server stopped.")

	get, i, err := s.db.Get([]byte("key"))
	if err != nil {
		return
	}
	defer i.Close()

	ctx := context.WithValue(context.Background(), "jet", "bpw")
	iops := pebble.IterOptions{}
	withContext, err := s.db.NewIterWithContext(ctx, &iops)
	if err != nil {
		return
	}
	withContext.First()
	for withContext.Valid() {
		key := withContext.Key()
		value := withContext.Value()
		fmt.Printf("Key: %s, Value: %s\n", key, value)

		withContext.Next()
	}
	fmt.Println(string(get))
}
