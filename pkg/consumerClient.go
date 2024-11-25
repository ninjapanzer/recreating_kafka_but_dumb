package pkg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"go_stream_events/internal"
	"io"
	"log"
	"net"
)

type ConsumerClient interface {
	Connect() error
	Poll() ([]internal.Message, error)
}

type DefaultConsumerClient struct {
	conn    net.Conn
	serde   internal.CborSerde
	Address string
	Limit   uint64
}

func NewConsumerClient(address string) ConsumerClient {
	return &DefaultConsumerClient{
		Address: address,
		Limit:   uint64(100),
	}
}

func (d *DefaultConsumerClient) Connect() error {
	var err error
	d.serde = internal.NewSerde()
	msg := internal.NewConsumerRegistrationMessage(
		internal.ConsumerRegistration{TopicName: "topic1", ConsumerName: "conn1", Offset: 22777})
	d.conn, err = net.Dial("tcp", d.Address)

	_, err = d.conn.Write(msg.Bytes())

	return err
}
func (d *DefaultConsumerClient) Poll() ([]internal.Message, error) {
	if d.conn == nil {
		return nil, errors.New("not Connected")
	}

	pollReq := internal.NewPoll(
		internal.Poll{Offset: 0, Limit: d.Limit})
	_, err := d.conn.Write(pollReq.Bytes())

	var messages = make([]internal.Message, 0)
	for i := uint64(0); i < d.Limit; i++ {
		var msgLength uint32
		err = binary.Read(d.conn, binary.BigEndian, &msgLength)
		if err == io.EOF {
			log.Printf("Connection Closed from %v", d.conn.RemoteAddr())
			break
		} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			continue
		}
		if err != nil {
			return nil, err
		}

		data := make([]byte, msgLength)
		_, _ = io.ReadFull(d.conn, data)
		v := internal.Message{}
		err = d.serde.DecodeCbor(bytes.NewReader(data[4:]), &v)
		if err != nil {
			return nil, err
		}
		messages = append(messages, v)
	}

	return messages, err
}
