package internal

import (
	"encoding/binary"
	"fmt"
	"log"
)

type ProtocolMessage interface {
	SetMessageType(uint32)
	Write([]byte) (int, error)
}

type TCPMessage struct {
	BufferWithLengthPrefix
	msgType uint32
}

func (msg *TCPMessage) SetMessageType(msgType uint32) {
	msg.msgType = msgType
}

func (msg *TCPMessage) Write(p []byte) (n int, err error) {
	if msg.msgType == 0 {
		log.Fatal("MessageTypeUnset 0")
	}
	totalLength := len(p) + 4
	result := make([]byte, totalLength)
	binary.BigEndian.PutUint32(result[:4], msg.msgType)
	copy(result[4:], p)
	fmt.Print(p)
	return msg.BufferWithLengthPrefix.Write(result)
}

//
//func (b *TCPMessage) Reset() {
//	b.buf.Reset()
//}

//func (b *TCPMessage) Bytes() []byte {
//	return b.buf.Bytes()
//}

//func (b *TCPMessage) Write(p []byte) (n int, err error) {
//	fmt.Printf("Writing %d bytes: %s\n", len(p), string(p))
//	return
//}

//func (b *TCPMessage) Read(p []byte) (n int, err error) {
//	n, err = b.Read(p)
//	if err != nil && err != io.EOF {
//		return n, err
//	}
//	fmt.Printf("Read %d bytes: %s\n", n, string(p[:n]))
//	return n, err
//}
