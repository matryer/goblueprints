package nilobjects

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Codec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

type jsonCodec struct{}

func (j *jsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonCodec) Decode(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

var JSON Codec = (*jsonCodec)(nil)

type stringCodec struct{}

func (j *stringCodec) Encode(v interface{}) ([]byte, error) {
	s, ok := v.(fmt.Stringer)
	if !ok {
		return nil, errors.New("v must be fmt.Stringer")
	}

}

func (j *stringCodec) Decode(b []byte, v interface{}) error {

}

var String Codec = (*stringCodec)(nil)
