package gotten_test

import (
	"github.com/Hexilee/gotten"
	"github.com/Hexilee/gotten/headers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestBaseUrlCannotBeEmpty(t *testing.T) {
	_, err := gotten.NewBuilder().Build()
	assert.NotNil(t, err)
	assert.Equal(t, gotten.BaseUrlCannotBeEmpty, err.Error())
}

func TestMustPassPtrToImplError(t *testing.T) {
	var wrongService int
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(wrongService)
	assert.NotNil(t, err)
	assert.Equal(t, gotten.MustPassPtrToImplError(reflect.TypeOf(wrongService)).Error(), err.Error())
}

func TestServiceMustBeStructError(t *testing.T) {
	var wrongService int
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongService)
	assert.NotNil(t, err)
	assert.Equal(t, gotten.ServiceMustBeStructError(reflect.TypeOf(wrongService)).Error(), err.Error())
}

func TestUnrecognizedHTTPMethodError(t *testing.T) {
	var wrongService struct {
		WrongMethod func(*struct {
			Id int `type:"query"`
		}) (*http.Response, error) `method:"GO"`
	}
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongService)
	assert.NotNil(t, err)
	assert.Equal(t, gotten.UnrecognizedHTTPMethodError("GO").Error(), err.Error())
}

func TestDuplicatedPathKeyError(t *testing.T) {
	var wrongService struct {
		WrongMethod func(*struct {
			Id int `type:"path"`
		}) (*http.Response, error) `path:"{id}/{id}"`
	}
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongService)
	assert.NotNil(t, err)
	assert.Equal(t, gotten.DuplicatedPathKeyError("id").Error(), err.Error())
}

func TestEmptyRequiredVariableError(t *testing.T) {
	type rightParam struct {
		Id int `type:"path"`
	}
	var rightService struct {
		Get func(param *rightParam) (gotten.Response, error) `path:"{id}"`
	}
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	assert.Nil(t, creator.Impl(&rightService))
	assert.NotNil(t, rightService.Get)
	_, err = rightService.Get(&rightParam{})
	assert.NotNil(t, err)
	assert.Equal(t, gotten.EmptyRequiredVariableError("Id").Error(), err.Error())
}

func TestEmptyRequiredVariableError2(t *testing.T) {
	type rightParam struct {
		Value string `type:"json" require:"true"`
	}
	var rightService struct {
		Get func(param *rightParam) (gotten.Response, error) `path:"id"`
	}
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	assert.Nil(t, creator.Impl(&rightService))
	assert.NotNil(t, rightService.Get)
	_, err = rightService.Get(&rightParam{})
	assert.NotNil(t, err)
	assert.Equal(t, gotten.EmptyRequiredVariableError("Value").Error(), err.Error())
}

func TestUnrecognizedPathKeyError(t *testing.T) {
	type NotExistPathParam struct {
		Name string `type:"path"`
	}
	var wrongService struct {
		Get func(param *NotExistPathParam) (gotten.Response, error) `path:"{id}"`
	}

	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongService)
	assert.Error(t, err)
	assert.Equal(t, gotten.UnrecognizedPathKeyError("name"), err)
}

func TestParamTypeMustBePtrOfStructError(t *testing.T) {
	type rightParam struct {
		Id int `type:"path"`
	}
	var wrongService struct {
		Get func(param rightParam) (gotten.Response, error) `path:"{id}"`
	}
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongService)
	assert.Error(t, err)
	assert.Equal(t, gotten.ParamTypeMustBePtrOfStructError(reflect.TypeOf(rightParam{})), err)
}

func TestUnsupportedValueTypeError(t *testing.T) {
	type wrongParam struct {
		Id string `type:"id"`
	}
	var wrongService struct {
		Get func(param *wrongParam) (gotten.Response, error) `path:"{id}"`
	}
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongService)
	assert.Error(t, err)
	assert.Equal(t, gotten.UnsupportedValueTypeError("id"), err)
}

func TestSomePathVarHasNoValueError(t *testing.T) {
	type wrongParam struct {
		Id string `type:"path"`
	}
	var wrongService struct {
		Get func(param *wrongParam) (gotten.Response, error) `path:"{id}/{name}"`
	}
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongService)
	assert.Error(t, err)
	assert.Equal(t, gotten.SomePathVarHasNoValueError(gotten.PathKeyList{"name": true}), err)
}

func TestContentTypeConflictError(t *testing.T) {
	type wrongParamOne struct {
		Id   string `type:"form"`
		Name string `type:"part"`
	}

	var wrongServiceOne struct {
		Get func(param *wrongParamOne) (gotten.Response, error) `path:"{id}/{name}"`
	}
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongServiceOne)
	assert.Error(t, err)
	assert.Equal(t, gotten.ContentTypeConflictError(headers.MIMEApplicationForm, headers.MIMEMultipartForm), err)

	type wrongParamTwo struct {
		Name string `type:"part"`
		Id   string `type:"form"`
	}

	var wrongServiceTwo struct {
		Get func(param *wrongParamTwo) (gotten.Response, error) `path:"{id}/{name}"`
	}
	creator, err = gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongServiceTwo)
	assert.Error(t, err)
	assert.Equal(t, gotten.ContentTypeConflictError(headers.MIMEMultipartForm, headers.MIMEApplicationForm), err)

	type wrongParamThree struct {
		Name string `type:"xml"`
		Id   string `type:"json"`
	}

	var wrongServiceThree struct {
		Get func(param *wrongParamThree) (gotten.Response, error) `path:"{id}/{name}"`
	}
	creator, err = gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongServiceThree)
	assert.Error(t, err)
	assert.Equal(t, gotten.ContentTypeConflictError(headers.MIMEApplicationXMLCharsetUTF8, headers.MIMEApplicationJSONCharsetUTF8), err)
}
