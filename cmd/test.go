package main

import (
	"fmt"
	"go_stream_events/internal"
	"net"
	"time"
)

func main() {
	s := internal.NewSerde()
	buff := internal.TCPMessage{}
	buff.SetMessageType(1)
	s.EncodeCbor(&buff, internal.ProducerRegistration{TopicName: "topic1", Id: "identifier"})

	address := "localhost:8080"

	// Dial a TCP connection to the server
	producerConn, _ := net.Dial("tcp", address)
	defer producerConn.Close()

	producerConn.Write(buff.Bytes())

	buff.Reset()

	s.EncodeCbor(&buff, internal.ProducerRegistration{TopicName: "topic2", Id: "identifier2"})

	producerConn.Write(buff.Bytes())

	buff = internal.TCPMessage{}
	buff.SetMessageType(2)
	s.EncodeCbor(&buff, internal.ConsumerRegistration{TopicName: "topic2", Offset: 0})
	fmt.Println(buff.Bytes())

	time.Sleep(1 * time.Second)
	consumerConn, _ := net.Dial("tcp", address)
	defer consumerConn.Close()

	consumerConn.Write(buff.Bytes())
}
