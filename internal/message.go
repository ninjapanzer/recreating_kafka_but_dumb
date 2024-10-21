package internal

import (
	"time"
)

type Message struct {
	Timestamp time.Time `codec:"1"`
	Topic     string    `codec:"topic,omitempty"`
	Payload   string    `codec:"payload"`
}
