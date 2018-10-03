package gotten

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	BaseUrlCannotBeEmpty       = "baseUrl cannot be empty"
	MustPassPtrToImpl          = "must pass the ptr of the service to be implemented"
	ServiceMustBeStruct        = "service must be struct"
	UnrecognizedHTTPMethod     = "http method is unrecognized"
	ParamTypeMustBePtrOfStruct = "param type must be ptr of struct"
	ValueIsNotStringer         = "value is not a stringer"
	ValueIsNotString           = "value is not a string"
	ValueIsNotInt              = "value is not a int"
	DuplicatedPathKey          = "duplicated path key"
	UnrecognizedFieldType      = "field type is unrecognized"
	UnrecognizedPathKey        = "path key is unrecognized"
	EmptyPathVariable          = "path variable is empty"
	UnsupportedFieldType       = "field type is unsupported"
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

func UnrecognizedFieldTypeError(fieldType string) error {
	return errors.New(UnrecognizedFieldType + ": " + fieldType)
}

func UnrecognizedPathKeyError(key string) error {
	return errors.New(UnrecognizedPathKey + ": " + key)
}

func EmptyPathVariableError(key string) error {
	return errors.New(EmptyPathVariable + ": " + key)
}

func UnsupportedFieldTypeError(fieldType reflect.Type, valueType string) error {
	return errors.New(fmt.Sprintf(UnrecognizedFieldType+": %s -> %s", fieldType, valueType))
}
