package main

import (
	"fmt"
	"go_stream_events/internal"
	"net"
	"time"
)

func main() {
	//s := internal.NewSerde()
	msg := internal.NewProducerRegistrationMessage(
		internal.ProducerRegistration{TopicName: "topic1", Id: "identifier"})
	msg2 := internal.NewProducerRegistrationMessage(
		internal.ProducerRegistration{TopicName: "topic2", Id: "identifier"})
	msg3 := internal.NewProducerRegistrationMessage(
		internal.ProducerRegistration{TopicName: "topic2", Id: "identifier"})

	address := "localhost:8080"
	// Dial a TCP connection to the server
	producerConn, _ := net.Dial("tcp", address)
	producerConn2, _ := net.Dial("tcp", address)
	producerConn3, _ := net.Dial("tcp", address)
	defer producerConn.Close()
	defer producerConn2.Close()
	defer producerConn3.Close()

	producerConn.Write(msg.Bytes())
	producerConn2.Write(msg2.Bytes())
	producerConn3.Write(msg3.Bytes())

	for i := 0; i < 10000; i++ {
		msg = internal.NewProducerMessage(
			internal.Message{
				Timestamp: time.Now(),
				Payload:   "Hello" + fmt.Sprint(i),
			})

		producerConn.Write(msg.Bytes())
		producerConn2.Write(msg.Bytes())
		producerConn3.Write(msg.Bytes())
	}

	consumerConn, _ := net.Dial("tcp", address)
	defer consumerConn.Close()

	msg = internal.NewConsumerRegistrationMessage(
		internal.ConsumerRegistration{TopicName: "topic1", Offset: 0})
	consumerConn.Write(msg.Bytes())
}
