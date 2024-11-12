package internal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cockroachdb/pebble"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func (s Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	for {
		buff := make([]byte, 1024)
		count, err := conn.Read(buff)
		serde := NewSerde()
		request := []byte(strings.TrimSpace(string(buff[:count])))
		out := append(request, ' ')
		out = append(out, []byte(time.Now().String())...)
		switch value := ctx.Value("registered").(type) {
		case nil:
			var producer ProducerRegistration
			var consumer ConsumerRegistration
			err = serde.DecodeCbor(bytes.NewReader(request), &producer)
			if err != nil {
			} else {
				ctx = context.WithValue(ctx, "registered", producer)
				goto registered
			}
			err = serde.DecodeCbor(bytes.NewReader(request), &consumer)
			if err != nil {
			} else {
				ctx = context.WithValue(ctx, "registered", consumer)
			}
		registered:
			;
		case ProducerRegistration:
			ctx = context.WithValue(ctx, "type", "Producer")
			var db *pebble.DB
			opts := pebble.Options{}
			db, err = pebble.Open(value.TopicName, &opts)
			if err != nil {
				log.Println("Producer")
				log.Println(err)
				ctx.Done()
			}
			db.Set([]byte("key"), out, pebble.NoSync)
			ctx.Done()
			db.Close()
		case ConsumerRegistration:
			ctx = context.WithValue(ctx, "type", "Consumer")
			var db *pebble.DB
			opts := pebble.Options{}
			opts.ReadOnly = true
			db, err = pebble.Open(value.TopicName, &opts)
			if err != nil {
				log.Println("Producer")
				log.Println(err)
				ctx.Done()
			}
			get, _, _ := db.Get([]byte("key"))
			ctx.Done()
			db.Close()
			fmt.Println(string(get))
		}

		logOut := append([]byte("key "), out...)
		logOut = append(logOut, '\n')
		s.logFile.Write(logOut)
		fmt.Fprintf(conn, string(request))
	}
}

func (s *Server) Start(logFile *os.File) {
	s.wg.Add(2)
	s.logFile = logFile
	go s.acceptConnections()
	go s.handleConnections(s.handleConnection)
}

func (s *Server) Stop() {
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

func BootV1() {
	s, err := newServer(":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logFile, err := os.OpenFile("demo/event.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close() // Ensure the file is closed when done

	s.Start(logFile)

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down server...")
	s.Stop()
	fmt.Println("Server stopped.")
}
