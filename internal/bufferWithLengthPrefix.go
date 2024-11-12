package internal

import (
	"bytes"
	"encoding/binary"
	"io"
)

type BufferWithLengthPrefix struct {
	io.Reader
	buf bytes.Buffer
}

func (b *BufferWithLengthPrefix) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

func (b *BufferWithLengthPrefix) Reset() {
	b.buf.Reset()
}

func (b *BufferWithLengthPrefix) Bytes() []byte {
	// Get the raw bytes of the buffer
	rawBytes := b.buf.Bytes()

	// Create a new slice with enough space for the length prefix and the raw bytes
	totalLength := len(rawBytes) + 4
	result := make([]byte, totalLength)

	// Write the length prefix (4 bytes, big-endian)
	binary.BigEndian.PutUint32(result[:4], uint32(len(rawBytes)))

	// Copy the actual data from the buffer
	copy(result[4:], rawBytes)

	return result
}
