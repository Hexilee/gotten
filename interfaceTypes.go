package gotten

import (
	"fmt"
	"io"
	"reflect"
)

type (
	Types struct {
		stringer fmt.Stringer
		reader   io.Reader
		error    error
		response Response
	}
)

var (
	types        = Types{}
	typesValue   = reflect.ValueOf(types)
	StringerType = typesValue.FieldByName("stringer").Type()
	ReaderType   = typesValue.FieldByName("reader").Type()
	ErrorType    = typesValue.FieldByName("error").Type()
	ResponseType = typesValue.FieldByName("response").Type()
)
