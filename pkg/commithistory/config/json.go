package config

import (
	"encoding/json"
	"io"
)

// Marshaller :
type Marshaller interface {
	Marshal(w io.Writer, ob interface{}) error
}

// Unmarshaller :
type Unmarshaller interface {
	Unmarshal(r io.Reader, ob interface{}) error
}

// MarshalUnmarshaller :
type MarshalUnmarshaller interface {
	Marshaller
	Unmarshaller
}

// JSONModule :
type JSONModule struct {
}

// Marshal :
func (m *JSONModule) Marshal(w io.Writer, ob interface{}) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(ob)
}

// Unmarshal :
func (m *JSONModule) Unmarshal(r io.Reader, ob interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(ob)
}
