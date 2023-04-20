package utils

import (
	"github.com/goccy/go-json"
	"io"
)

func ToJsonString(input interface{}) string {
	val, _ := json.Marshal(input)
	return string(val)
}

func Marshal(input interface{}) ([]byte, error) {
	return json.Marshal(input)
}

func Decoder(input io.ReadCloser) *json.Decoder {
	return json.NewDecoder(input)
}

