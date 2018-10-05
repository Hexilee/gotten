package gotten

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
)

const (
	// support types: fmt.Stringer, int, string
	TypeHeader = "header"

	// support types: fmt.Stringer, int, string
	TypePath = "path"

	// support types: fmt.Stringer, int, string
	TypeQuery = "query"

	// support types: fmt.Stringer, int, string
	TypeForm = "form"

	// support types: fmt.Stringer, int, string, Reader, PartFile
	TypeMultipart = "part"

	// support types: fmt.Stringer, Reader, string, struct, slice, map
	TypeJSON = "json"

	// support types: fmt.Stringer, Reader, string, struct, slice, map
	TypeXML = "xml"

	// support types: string, *http.Cookie
	TypeCookie = "cookie"
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
	//	PartReader Reader
	//
)

//
var (
	partFile     = PartFile("")
	PartFileType = reflect.TypeOf(partFile)
	IntType      = reflect.TypeOf(int(1))
	StringType   = reflect.TypeOf("")
)

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
func getJSONReaderGetterFunc(fieldKind reflect.Kind, fieldType reflect.Type, valueType string) (getValueFunc func(value reflect.Value) (Reader, error), err error) {
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
func getXMLReaderGetterFunc(fieldKind reflect.Kind, fieldType reflect.Type, valueType string) (getValueFunc func(value reflect.Value) (Reader, error), err error) {
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
		data, err := marshalFunc(value)
		return newReader(bytes.NewBuffer(data), false), err
	}
}

func getValueFromStringer(value reflect.Value) (str string, err error) {
	stringer, ok := value.Interface().(fmt.Stringer)
	if !ok {
		panic(ValueIsNotStringerError(value.Type()))
	}
	str = stringer.String()
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
		panic(ValueIsNotStringerError(value.Type()))
	}

	str := stringer.String()
	reader = newReader(bytes.NewBufferString(str), str == ZeroStr)
	return
}

func getReaderFromString(value reflect.Value) (reader Reader, err error) {
	val, ok := value.Interface().(string)
	if !ok {
		panic(ValueIsNotStringError(value.Type()))
	}
	reader = newReader(bytes.NewBufferString(val), val == ZeroStr)
	return
}

func getReaderFromInt(value reflect.Value) (reader Reader, err error) {
	val, ok := value.Interface().(int)
	if !ok {
		panic(ValueIsNotIntError(value.Type()))
	}

	reader = newReader(bytes.NewBufferString(strconv.Itoa(val)), val == ZeroInt)
	return
}

func getReaderFromReader(value reflect.Value) (reader Reader, err error) {
	reader, ok := value.Interface().(Reader)
	if !ok {
		panic(value.Type().String() + " is not a Reader")
	}
	return
}
