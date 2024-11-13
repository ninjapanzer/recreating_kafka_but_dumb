package internal

import (
	"bytes"
	"encoding/binary"
	"log"
)

type ProtocolMessage interface {
	SetMessageType(uint32)
	Write([]byte) (int, error)
	Bytes() []byte
}

type TCPMessage struct {
	buf     bytes.Buffer
	msgType uint32
}

func (msg *TCPMessage) SetMessageType(msgType uint32) {
	msg.msgType = msgType
}

func (msg *TCPMessage) Write(p []byte) (n int, err error) {
	if msg.msgType == 0 {
		log.Fatal("MessageTypeUnset 0")
	}

	return msg.buf.Write(p)
}

func (msg *TCPMessage) Bytes() []byte {
	if msg.msgType == 0 {
		log.Fatal("MessageTypeUnset 0")
	}

	totalLength := msg.buf.Len() + 4
	result := make([]byte, totalLength)
	binary.BigEndian.PutUint32(result[:4], msg.msgType)
	copy(result[4:], msg.buf.Bytes())
	return result
}
