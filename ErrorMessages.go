package gotten

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const (
	BaseUrlCannotBeEmpty          = "baseUrl cannot be empty"
	MustPassPtrToImpl             = "must pass the ptr of the service to be implemented"
	ServiceMustBeStruct           = "service must be struct"
	UnrecognizedHTTPMethod        = "http method is unrecognized"
	ParamTypeMustBePtrOfStruct    = "param type must be ptr of struct"
	ValueIsNotStringer            = "value is not a stringer"
	ValueIsNotString              = "value is not a string"
	ValueIsNotInt                 = "value is not a int"
	DuplicatedPathKey             = "duplicated path key"
	UnsupportedValueType          = "field type is unrecognized"
	UnrecognizedPathKey           = "path key is unrecognized"
	EmptyRequiredVariable         = "required variable is empty"
	UnsupportedFieldType          = "field type is unsupported"
	SomePathVarHasNoValue         = "some pathValue has no value"
	NoUnmarshalerFoundForResponse = "no unmarshaler found for response"
	ContentTypeConflict           = "content type conflict: "
)

func MustPassPtrToImplError(p reflect.Type) error {
	return errors.New(MustPassPtrToImpl + ": " + p.String())
}

func ServiceMustBeStructError(p reflect.Type) error {
	return errors.New(ServiceMustBeStruct + ": " + p.String())
}

func UnrecognizedHTTPMethodError(method string) error {
	return errors.New(UnrecognizedHTTPMethod + ": " + method)
}

func ParamTypeMustBePtrOfStructError(p reflect.Type) error {
	return errors.New(ParamTypeMustBePtrOfStruct + ": " + p.String())
}

func ValueIsNotStringerError(p reflect.Type) error {
	return errors.New(ValueIsNotStringer + ": " + p.String())
}

func ValueIsNotStringError(p reflect.Type) error {
	return errors.New(ValueIsNotString + ": " + p.String())
}

func ValueIsNotIntError(p reflect.Type) error {
	return errors.New(ValueIsNotInt + ": " + p.String())
}

func DuplicatedPathKeyError(key string) error {
	return errors.New(DuplicatedPathKey + ": " + key)
}

func UnsupportedValueTypeError(valueType string) error {
	return errors.New(UnsupportedValueType + ": " + valueType)
}

func UnrecognizedPathKeyError(key string) error {
	return errors.New(UnrecognizedPathKey + ": " + key)
}

func EmptyRequiredVariableError(name string) error {
	return errors.New(EmptyRequiredVariable + ": " + name)
}

func UnsupportedFieldTypeError(fieldType reflect.Type, valueType string) error {
	return errors.New(fmt.Sprintf(UnsupportedFieldType+": %s -> %s", fieldType, valueType))
}

func SomePathVarHasNoValueError(list PathKeyList) error {
	buf := strings.Builder{}
	buf.WriteString(SomePathVarHasNoValue)
	buf.WriteString(": ")
	for key := range list {
		buf.WriteString(" ")
		buf.WriteString(key)
	}
	return errors.New(buf.String())
}

func NoUnmarshalerFoundForResponseError(response *http.Response) error {
	return errors.New(fmt.Sprintf(NoUnmarshalerFoundForResponse+"%#v", response))
}

func ContentTypeConflictError(former string, latter string) error {
	return errors.New(fmt.Sprintf(ContentTypeConflict+" %s(former), %s(former)", former, latter))
}
