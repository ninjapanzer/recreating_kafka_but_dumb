package internal

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/cockroachdb/pebble"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
		if ctx.Value("mode") != nil {
			v = Message{}
		} else if msgType == 1 {
			v = ProducerRegistration{}
			ctx = context.WithValue(ctx, "mode", "Producer")
		} else if msgType == 2 {
			v = ConsumerRegistration{}
			ctx = context.WithValue(ctx, "mode", "Consumer")
		}
		err := serde.DecodeCbor(bytes.NewReader(data[4:]), &v)
		switch m := v.(type) {
		case ProducerRegistration:
			ctx = context.WithValue(ctx, "topic", m.TopicName)
		case ConsumerRegistration:
			ctx = context.WithValue(ctx, "topic", m.TopicName)
		case Message:
			var topic string = ctx.Value("topic").(string)
			var pdb = s.dbInstances.Get(topic)
			if pdb == nil {
				pdb, err = pebble.Open(topic, &pebble.Options{})
				if err != nil {
					log.Println("pebble", err)
					continue
				}
				s.dbInstances.Set(topic, pdb)
			}
			err := pdb.Set([]byte("key"), []byte(m.Payload), nil)
			var writer = s.eventStore.Get(topic)
			if writer == nil {
				h, err := os.OpenFile(topic+".log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					log.Println("open", err)
				}
				writer = s.eventStore.Set(topic, h)
			}
			writer.Write([]byte(m.Payload))
			if err != nil {
				log.Println("error writing", err)
			}
		}
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error decoding message:", err)
			}
			break
		}
		log.Println(v, ctx.Value("topic"))
	}
	ctx.Done()
}

func (s *Server) Start2(logFile *os.File) {
	s.wg.Add(2)
	s.logFile = logFile
	s.dbInstances.store = make(map[string]*pebble.DB)
	s.eventStore.store = make(map[string]*EventWriter)
	go s.acceptConnections()
	go s.handleConnections(s.handleConnection2)
}

func (s *Server) Stop2() {
	close(s.shutdown)
	s.listener.Close()
	log.Println("Closing DB Connections")
	for topic, instance := range s.dbInstances.store {
		instance.Flush()
		instance.Close()
		log.Println("Closing Connection for ", topic)
	}

	for topic, writer := range s.eventStore.store {
		writer.handle.Sync()
		writer.handle.Close()
		log.Println("Closing handle for ", topic)
	}

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
	s.Stop2()
	log.Println("Server stopped.")
}
