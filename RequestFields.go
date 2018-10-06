package gotten

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// value types
const (
	// support types: fmt.Stringer, int, string
	TypeHeader = "header"

	// support types: fmt.Stringer, int, string
	TypePath = "path"

	// support types: fmt.Stringer, int, string
	TypeQuery = "query"

	// support types: fmt.Stringer, int, string
	TypeForm = "form"

	// support types: fmt.Stringer, int, string
	TypeCookie = "cookie"

	// support types: fmt.Stringer, int, string, Reader, FilePath
	TypeMultipart = "part"

	// support types: fmt.Stringer, Reader, string, struct, slice, map
	TypeJSON = "json"

	// support types: fmt.Stringer, Reader, string, struct, slice, map
	TypeXML = "xml"
)

type (
	FilePath string
)

//
var (
	filePath     = FilePath("")
	FilePathType = reflect.TypeOf(filePath)
	IntType      = reflect.TypeOf(int(1))
	StringType   = reflect.TypeOf("")
)

func getMultipartValueGetterFunc(fieldType reflect.Type, valueType string) (getValueFunc func(value reflect.Value) (string, error), err error) {
	switch fieldType {
	case IntType:
		getValueFunc = getValueFromInt
	case StringType:
		getValueFunc = getValueFromString
	case StringerType:
		getValueFunc = getValueFromStringer
	case FilePathType:
		getValueFunc = getValueFromFilePath
	default:
		err = UnsupportedFieldTypeError(fieldType, valueType)
	}
	return
}

// for TypePath, TypeQuery, TypeHeader and TypeForm
func getValueGetterFunc(fieldType reflect.Type, valueType string) (getValueFunc func(value reflect.Value) (string, error), err error) {
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

// for TypeJSON, TypeXML, field type cannot be struct, map and slice
func getReaderGetterFunc(fieldType reflect.Type, valueType string) (getValueFunc func(value reflect.Value) (Reader, error), err error) {
	switch fieldType {
	case StringType:
		getValueFunc = getReaderFromString
	case StringerType:
		getValueFunc = getReaderFromStringer
	case ReaderType:
		getValueFunc = getReaderFromReader
	default:
		err = UnsupportedFieldTypeError(fieldType, valueType)
	}
	return
}

// can only be called by parse
func getJSONReaderGetterFunc(fieldType reflect.Type, valueType string) (getValueFunc func(value reflect.Value) (Reader, error), err error) {
	fieldKind := fieldType.Kind()
	switch fieldKind {
	case reflect.Ptr:
		fallthrough
	case reflect.Struct:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		getValueFunc = getMarshalReaderGetterFunc(json.Marshal)
	default:
		getValueFunc, err = getReaderGetterFunc(fieldType, valueType)
	}
	return
}

// can only be called by parse
func getXMLReaderGetterFunc(fieldType reflect.Type, valueType string) (getValueFunc func(value reflect.Value) (Reader, error), err error) {
	fieldKind := fieldType.Kind()
	switch fieldKind {
	case reflect.Ptr:
		fallthrough
	case reflect.Struct:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		getValueFunc = getMarshalReaderGetterFunc(xml.Marshal)
	default:
		getValueFunc, err = getReaderGetterFunc(fieldType, valueType)
	}
	return
}

func getMarshalReaderGetterFunc(marshalFunc func(obj interface{}) ([]byte, error)) func(value reflect.Value) (Reader, error) {
	return func(value reflect.Value) (Reader, error) {
		data, err := marshalFunc(value.Interface())
		return newReadCloser(bytes.NewBuffer(data), value.IsNil()), err
	}
}

func getValueFromStringer(value reflect.Value) (str string, err error) {
	stringer, ok := value.Interface().(fmt.Stringer)
	if !ok {
		stringer = ZeroStringer
	}
	str = stringer.String()
	return
}

func getValueFromFilePath(value reflect.Value) (str string, err error) {
	filePath, ok := value.Interface().(FilePath)
	if !ok {
		panic(value.Type().String() + " is not FilePath")
	}
	str = string(filePath)
	return
}

func getValueFromString(value reflect.Value) (str string, err error) {
	val, ok := value.Interface().(string)
	if !ok {
		panic(ValueIsNotStringError(value.Type()))
	}
	str = val
	return
}

func getValueFromInt(value reflect.Value) (str string, err error) {
	val, ok := value.Interface().(int)
	if !ok {
		panic(ValueIsNotIntError(value.Type()))
	}

	if val != ZeroInt {
		str = strconv.Itoa(val)
	}
	return
}

func getReaderFromStringer(value reflect.Value) (reader Reader, err error) {
	stringer, ok := value.Interface().(fmt.Stringer)
	if !ok {
		stringer = ZeroStringer
	}

	str := stringer.String()
	reader = newReadCloser(bytes.NewBufferString(str), str == ZeroStr)
	return
}

func getReaderFromString(value reflect.Value) (reader Reader, err error) {
	val, ok := value.Interface().(string)
	if !ok {
		panic(ValueIsNotStringError(value.Type()))
	}
	reader = newReadCloser(bytes.NewBufferString(val), val == ZeroStr)
	return
}

func getReaderFromReader(value reflect.Value) (reader Reader, err error) {
	ioReader, ok := value.Interface().(io.Reader)
	if !ok {
		ioReader = ZeroReader
	}

	reader = newReadCloser(ioReader, value.IsNil())
	return
}
