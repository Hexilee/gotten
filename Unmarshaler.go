package gotten

import (
	"io"
	"io/ioutil"
	"net/http"
)

type (
	ReaderUnmarshaler interface {
		Unmarshal(reader io.Reader, header http.Header, v interface{}) error
	}

	ReaderUnmarshalerFunc func(reader io.Reader, header http.Header, v interface{}) error

	ReaderAdapter struct {
		unmarshaler Unmarshaler
	}

	Unmarshaler interface {
		Unmarshal(data []byte, v interface{}) error
	}

	UnmarshalFunc func(data []byte, v interface{}) error
)

func (fn UnmarshalFunc) Unmarshal(data []byte, v interface{}) error {
	return fn(data, v)
}

func (fn ReaderUnmarshalerFunc) Unmarshal(reader io.Reader, header http.Header, v interface{}) error {
	return fn(reader, header, v)
}

func NewReaderAdapter(unmarshaler Unmarshaler) ReaderUnmarshaler {
	return &ReaderAdapter{unmarshaler}
}

func (adapter *ReaderAdapter) Unmarshal(reader io.Reader, header http.Header, v interface{}) (err error) {
	var body []byte
	body, err = ioutil.ReadAll(reader)
	if err == nil {
		err = adapter.unmarshaler.Unmarshal(body, v)
	}
	return
}
