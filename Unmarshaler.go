package gotten

import (
	"io"
	"io/ioutil"
	"net/http"
)

type (
	ReadUnmarshaler interface {
		Unmarshal(reader io.ReadCloser, header http.Header, v interface{}) error
	}

	ReadUnmarshalFunc func(reader io.ReadCloser, header http.Header, v interface{}) error

	ReaderAdapter struct {
		unmarshaler Unmarshaler
	}

	Unmarshaler interface {
		Unmarshal(data []byte, v interface{}) error
	}

	UnmarshalFunc func(data []byte, v interface{}) error
)

func UnmarshalAdapter(fn UnmarshalFunc) Unmarshaler {
	return fn
}

//func ReaderFuncAdapter(fn ReadUnmarshalFunc) ReadUnmarshaler {
//	return fn
//}

func NewReaderAdapter(unmarshaler Unmarshaler) ReadUnmarshaler {
	return &ReaderAdapter{unmarshaler}
}

func (fn UnmarshalFunc) Unmarshal(data []byte, v interface{}) error {
	return fn(data, v)
}

func (fn ReadUnmarshalFunc) Unmarshal(reader io.ReadCloser, header http.Header, v interface{}) error {
	return fn(reader, header, v)
}

func (adapter *ReaderAdapter) Unmarshal(reader io.ReadCloser, header http.Header, v interface{}) (err error) {
	var body []byte
	body, err = ioutil.ReadAll(reader)
	reader.Close()
	if err == nil {
		err = adapter.unmarshaler.Unmarshal(body, v)
	}
	return
}
