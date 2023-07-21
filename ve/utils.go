package ve

import (
	"bytes"
	"encoding/json"
	"io"
)

func UnmarshalUseNumber(data []byte, v interface{}) error {
	decoder := NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(v)
}

func NewDecoder(reader io.Reader) *json.Decoder {
	return json.NewDecoder(reader)
}
