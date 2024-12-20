package internal

import (
	"time"
)

type ApiContract struct{}

type Message struct {
	ApiContract
	Timestamp time.Time `codec:"1"`
	Payload   string    `codec:"payload"`
}

type Poll struct {
	ApiContract
	Offset uint64 `codec:"offset"`
	Limit  uint64 `codec:"limit"`
}

type ConsumerRegistration struct {
	ApiContract
	TopicName    string `codec:"topic,string"`
	ConsumerName string `codec:"consumer,string"`
	Offset       uint64 `codec:"offset,uint64"`
}

type ConsumerRequest struct {
	ApiContract
	Offset uint64 `codec:"offset,uint64"`
}

type ProducerRegistration struct {
	ApiContract
	TopicName string `codec:"topic,string"`
	Id        string `codec:"id,string"`
}

type CborMessage struct {
	bufferWithLengthPrefix BufferWithLengthPrefix
	protocolMessage        ProtocolMessage
}

func (msg *CborMessage) Bytes() []byte {
	return msg.bufferWithLengthPrefix.Bytes()
}

func NewProducerRegistrationMessage(msg ProducerRegistration) *CborMessage {
	m := &CborMessage{
		protocolMessage: &TCPMessage{},
	}
	m.protocolMessage.SetMessageType(1)
	NewSerde().EncodeCbor(m.protocolMessage, msg)
	m.bufferWithLengthPrefix.Write(m.protocolMessage.Bytes())
	return m
}

func NewConsumerRegistrationMessage(msg ConsumerRegistration) *CborMessage {
	m := &CborMessage{
		protocolMessage: &TCPMessage{},
	}
	m.protocolMessage.SetMessageType(2)
	NewSerde().EncodeCbor(m.protocolMessage, msg)
	m.bufferWithLengthPrefix.Write(m.protocolMessage.Bytes())
	return m
}

func NewMessage(msg Message) *CborMessage {
	m := &CborMessage{
		protocolMessage: &TCPMessage{},
	}
	m.protocolMessage.SetMessageType(3)
	NewSerde().EncodeCbor(m.protocolMessage, msg)
	m.bufferWithLengthPrefix.Write(m.protocolMessage.Bytes())
	return m
}

func NewPoll(msg Poll) *CborMessage {
	m := &CborMessage{
		protocolMessage: &TCPMessage{},
	}
	m.protocolMessage.SetMessageType(4)
	NewSerde().EncodeCbor(m.protocolMessage, msg)
	m.bufferWithLengthPrefix.Write(m.protocolMessage.Bytes())
	return m
}
