package tools

import (
	"bytes"
	"encoding/binary"
	"go_stream_events/internal"
	"io"
	"log"
	"net"
	"time"
)

func StartMultiplexConsumer() {
	msg := internal.NewConsumerRegistrationMessage(
		internal.ConsumerRegistration{TopicName: "topic1", ConsumerName: "conn1", Offset: 22777})
	msg2 := internal.NewConsumerRegistrationMessage(
		internal.ConsumerRegistration{TopicName: "topic2", ConsumerName: "conn2", Offset: 0})
	var consumerConn, consumerConn2 net.Conn
	for {
		var err error
		address := "localhost:8080"
		// Dial a TCP connection to the server
		consumerConn, err = net.Dial("tcp", address)
		consumerConn2, err = net.Dial("tcp", address)
		if err != nil {
			log.Println("waiting")
			time.Sleep(1 * time.Second)
			continue
		}
		consumerConn.Write(msg.Bytes())
		consumerConn2.Write(msg2.Bytes())
		break
	}
	for {
		pollReq := internal.NewPoll(
			internal.Poll{Offset: 22777, Limit: 100})

		write, err := consumerConn.Write(pollReq.Bytes())
		if err != nil {
			return
		}
		log.Println("write len", write)

		consumerConn.SetReadDeadline(time.Now().Add(5 * time.Second))
		for i := 0; i < 100; i++ {
			var msgLength uint32
			err = binary.Read(consumerConn, binary.BigEndian, &msgLength)
			if err == io.EOF {
				log.Printf("Connection Closed from %v", consumerConn.RemoteAddr())
				break
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if err != nil {
				log.Println(err)
			}
			serde := internal.NewSerde()
			data := make([]byte, msgLength)
			_, _ = io.ReadFull(consumerConn, data)
			v := internal.Message{}
			err = serde.DecodeCbor(bytes.NewReader(data[4:]), &v)
			if err != nil {
				log.Println(err)
			}
			log.Println(v)
		}

		time.Sleep(1 * time.Second)
	}
	err := consumerConn.Close()
	err = consumerConn2.Close()
	if err != nil {
		return
	}
}
