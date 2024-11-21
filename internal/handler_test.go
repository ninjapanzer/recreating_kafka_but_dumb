package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

var serde CborSerde
var messageBuf bytes.Buffer
var consumerBuf bytes.Buffer
var producerBuf bytes.Buffer
var pollBuf bytes.Buffer

func TestMain(m *testing.M) {
	// Setup code (executed before tests)
	serde := NewSerde()

	messageBuf := bytes.Buffer{}
	message := Message{
		Timestamp: time.Now(),
		Payload:   "Hello",
	}
	serde.EncodeCbor(&messageBuf, message)

	consumerBuf := bytes.Buffer{}
	consumer := Consumer{
		Position:  0,
		TopicName: "topic",
	}
	serde.EncodeCbor(&consumerBuf, consumer)

	producerBuf := bytes.Buffer{}
	producer := Producer{
		TopicName: "topic",
	}
	serde.EncodeCbor(&producerBuf, producer)

	pollBuf := bytes.Buffer{}
	poll := Poll{
		Offset: 0,
		Limit:  100,
	}
	serde.EncodeCbor(&pollBuf, poll)

	// Run tests
	code := m.Run()

	// Teardown code (executed after tests)
	fmt.Println("Tearing down tests...")

	// Exit with the appropriate status
	os.Exit(code)
}

func TestParseMessage(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		msgType         uint32
		message         []byte
		ctx             context.Context
		expectedMode    interface{}
		expectedMessage string
	}{
		{1, producerBuf.Bytes(), ctx, "Producer", "internal.ProducerRegistration"},
		{2, consumerBuf.Bytes(), ctx, "Consumer", "internal.ConsumerRegistration"},
		{3, messageBuf.Bytes(), ctx, nil, "internal.Message"},
		{4, pollBuf.Bytes(), ctx, nil, "internal.Poll"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.expectedMessage), func(t *testing.T) {
			v, ctx, _ := parseMessage(test.msgType, test.message, test.ctx, serde)
			if ctx.Value("mode") != test.expectedMode && ctx.Value("mode").(string) != test.expectedMode {
				t.Error("Failed to setup producer context")
			}

			if reflect.TypeOf(v).String() != test.expectedMessage {
				t.Error("Failed parse message")
			}
		})
	}
}
