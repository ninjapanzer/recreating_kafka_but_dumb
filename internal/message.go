package internal

import (
	"time"
)

type Message struct {
	Timestamp time.Time `codec:"1"`
	Payload   string    `codec:"payload"`
}

type ConsumerRegistration struct {
	TopicName string `codec:"topic,string"`
	Offset    uint64 `codec:"offset,uint64"`
}

type ConsumerRequest struct {
	Offset uint64 `codec:"offset,uint64"`
}

type ProducerRegistration struct {
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

func NewProducerMessage(msg Message) *CborMessage {
	m := &CborMessage{
		protocolMessage: &TCPMessage{},
	}
	m.protocolMessage.SetMessageType(3)
	NewSerde().EncodeCbor(m.protocolMessage, msg)
	m.bufferWithLengthPrefix.Write(m.protocolMessage.Bytes())
	return m
}
