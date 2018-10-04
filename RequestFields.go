package gotten

import (
	"fmt"
	"io"
	"reflect"
)

const (
	// support types: fmt.Stringer, int, string
	TypePath = "path"

	// support types: fmt.Stringer, int, string
	TypeQuery = "query"

	// support types: fmt.Stringer, int, string, io.Reader, PartFile
	TypeMultipart = "part"

	// support types: fmt.Stringer, io.Reader, string, struct
	TypeJSON = "json"

	// support types: fmt.Stringer, io.Reader, string, struct
	TypeXML = "xml"

	// support types: fmt.Stringer, int, string
	TypeForm = "form"
)

type (
	PartFile string

	//	PathVar interface {
	//		fmt.Stringer
	//	}
	//	QueryVar interface {
	//		fmt.Stringer
	//	}
	//	PartVar interface {
	//		fmt.Stringer
	//	}
	//
	//	PathStr string
	//
	//	PathInt int
	//
	//	QueryStr string
	//
	//	QueryInt int
	//
	//	PartStr string
	//
	//	PartInt int
	//
	//	PartReader io.Reader
	//

	Types struct {
		stringer fmt.Stringer
		reader   io.Reader
	}
)

//
var (
	types        = Types{}
	partFile     = PartFile("")
	StringerType = reflect.ValueOf(types).FieldByName("stringer").Type()
	ReaderType   = reflect.ValueOf(types).FieldByName("reader").Type()
	PartFileType = reflect.TypeOf(partFile)
	IntType      = reflect.TypeOf(int(1))
	StringType   = reflect.TypeOf("")
)

// for TypePath, TypeQuery, TypeHeader and TypeForm
func FirstValueGetterFunc(fieldType reflect.Type, valueType string) (getValueFunc func(value reflect.Value) string, err error) {
	switch fieldType {
	case IntType:
		getValueFunc = getValueFromInt
	case StringType:
		getValueFunc = getValueFromString
	case StringerType:
		getValueFunc = getValueFromStringer
	default:
		err = UnsupportedFieldTypeError(fieldType, valueType)
	}
	return
}
