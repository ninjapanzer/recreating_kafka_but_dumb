package tools

import (
	"fmt"
	"go_stream_events/internal"
	"log"
	"net"
	"time"
)

func StartMultiplexProducer() {
	for {
		//s := internal.NewSerde()
		msg := internal.NewProducerRegistrationMessage(
			internal.ProducerRegistration{TopicName: "topic1", Id: "identifier"})
		msg2 := internal.NewProducerRegistrationMessage(
			internal.ProducerRegistration{TopicName: "topic2", Id: "identifier"})
		msg3 := internal.NewProducerRegistrationMessage(
			internal.ProducerRegistration{TopicName: "topic2", Id: "identifier"})

		address := "localhost:8080"
		// Dial a TCP connection to the server
		producerConn, err := net.Dial("tcp", address)
		producerConn2, err := net.Dial("tcp", address)
		producerConn3, err := net.Dial("tcp", address)
		if err != nil {
			log.Println("waiting")
			time.Sleep(1 * time.Second)
			continue
		}

		producerConn.Write(msg.Bytes())
		producerConn2.Write(msg2.Bytes())
		producerConn3.Write(msg3.Bytes())

		for i := 0; i < 10; i++ {
			msg = internal.NewMessage(
				internal.Message{
					Timestamp: time.Now(),
					Payload:   "Hello" + fmt.Sprint(i),
				})

			producerConn.Write(msg.Bytes())
			producerConn2.Write(msg.Bytes())
			producerConn3.Write(msg.Bytes())
			time.Sleep(1 * time.Second)
		}

		consumerConn, _ := net.Dial("tcp", address)
		defer consumerConn.Close()

		producerConn.Close()
		producerConn2.Close()
		producerConn3.Close()

	}
}
