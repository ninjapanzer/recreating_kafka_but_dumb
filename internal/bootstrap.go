package internal

import (
	"bytes"
	"github.com/ugorji/go/codec"
	"log"
	"time"
)

type Proto struct {
	Id string `codec:"id,string"`
}

func Bootstrap() {
	var ch codec.CborHandle

	ch.TimeRFC3339 = false
	ch.SkipUnexpectedTags = true

	buff := bytes.Buffer{}

	encoder := codec.NewEncoder(&buff, &ch)

	var msg Message
	msg = Message{
		Timestamp: time.Now(),
		Payload:   "hi",
	}

	err := encoder.Encode(msg)
	if err != nil {
		log.Print("Decoding error:", err)
	}
	log.Printf("%v\n", buff.Bytes())
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
