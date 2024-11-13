package internal

import (
	"github.com/ugorji/go/codec"
	"io"
)

type CborSerde struct {
	ch codec.CborHandle
}

func NewSerde() CborSerde {
	cs := CborSerde{}
	ch := codec.CborHandle{}
	ch.ErrorIfNoField = true
	ch.TimeRFC3339 = false
	ch.SkipUnexpectedTags = true
	cs.ch = ch
	return cs
}

func (cs CborSerde) DecodeCbor(buffer io.Reader, v interface{}) error {
	dec := codec.NewDecoder(buffer, &cs.ch)
	err := dec.Decode(v)
	return err
}

func (cs CborSerde) EncodeCbor(buffer io.Writer, v interface{}) error {
	enc := codec.NewEncoder(buffer, &cs.ch)
	err := enc.Encode(v)
	return err
}
