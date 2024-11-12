package internal

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (s Server) handleConnection2(conn net.Conn) {
	log.Printf("New connection from %v", conn.RemoteAddr())
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	serde := NewSerde()
	for {
		ctx.Deadline()
		var msgLength uint32
		if err := binary.Read(conn, binary.BigEndian, &msgLength); err != nil {
			if err == io.EOF {
				log.Printf("Connection Closed from %v", conn.RemoteAddr())
				break
			} else {
				log.Println("bwoke", err)
				ctx.Done()
			}
		}
		var v interface{}
		data := make([]byte, msgLength)
		_, _ = io.ReadFull(conn, data)
		msgType := binary.BigEndian.Uint32(data[:4])
		if msgType == 1 {
			v = ProducerRegistration{}
		} else if msgType == 2 {
			v = ConsumerRegistration{}
		}
		err := serde.DecodeCbor(bytes.NewReader(data[4:]), &v)
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error decoding message:", err)
			}
			break
		}
		log.Println(v)
	}
	ctx.Done()
}

func (s *Server) Start2(logFile *os.File) {
	s.wg.Add(2)
	s.logFile = logFile
	go s.acceptConnections()
	go s.handleConnections(s.handleConnection2)
}

func (s *Server) Stop2() {
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
		log.Println("Timed out waiting for connections to finish.")
		return
	}
}

func BootV2() {
	s, err := newServer(":8080")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	logFile, err := os.OpenFile("demo/event.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close() // Ensure the file is closed when done

	s.Start2(logFile)

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	s.Stop()
	log.Println("Server stopped.")
}
