package gotten_test

import (
	"github.com/Hexilee/gotten"
	"github.com/Hexilee/gotten/headers"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
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
		}) (*http.Request, error) `method:"GO"`
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
		}) (*http.Request, error) `path:"{id}/{id}"`
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

func TestNoUnmarshalerFoundForResponseError(t *testing.T) {
	type TextService struct {
		Get func(*struct{}) (gotten.Response, error) `path:"/text"`
	}

	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		SetClient(mockClient).
		Build()

	assert.Nil(t, err)
	service := new(TextService)
	assert.Nil(t, creator.Impl(service))
	assert.NotNil(t, service.Get)

	_, err = service.Get(nil)
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), gotten.NoUnmarshalerFoundForResponse))
}

func TestRequestError(t *testing.T) {
	type TextService struct {
		Get func(*struct{}) (gotten.Response, error) `path:"/text"`
	}

	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		Build()

	assert.Nil(t, err)
	service := new(TextService)
	assert.Nil(t, creator.Impl(service))
	assert.NotNil(t, service.Get)
	_, err = service.Get(nil)
	assert.Error(t, err)
}

func TestResolveMultipartError(t *testing.T) {
	type MultipartFile struct {
		Reader io.Reader `type:"part"`
	}

	type UploadService struct {
		Upload func(params *MultipartFile) (*http.Request, error)
	}

	service := new(UploadService)

	file, err := os.Open("testAssets/Concurrency-in-Go.pdf")
	assert.Nil(t, err)
	file.Close()

	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		Build()
	assert.Nil(t, err)

	creator.Impl(service)
	assert.NotNil(t, service.Upload)
	_, err = service.Upload(&MultipartFile{file})
	assert.Error(t, err)
	assert.Equal(t, "read testAssets/Concurrency-in-Go.pdf: file already closed", err.Error())

	type NotExistFile struct {
		File gotten.FilePath `type:"part"`
	}

	type NotExistFileService struct {
		Upload func(params *NotExistFile) (*http.Request, error)
	}

	wrongService := new(NotExistFileService)
	creator, err = gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		Build()
	assert.Nil(t, err)

	creator.Impl(wrongService)
	assert.NotNil(t, wrongService.Upload)
	_, err = wrongService.Upload(&NotExistFile{"Concurrency-in-Go.pdf"})

	assert.Error(t, err)
	assert.Equal(t, "open Concurrency-in-Go.pdf: no such file or directory", err.Error())
}

func TestUnsupportedFuncTypeError(t *testing.T) {
	var wrongService struct {
		WrongMethod func(*struct {
			Id int `type:"query"`
		}) (*http.Response, error)
	}
	creator, err := gotten.NewBuilder().SetBaseUrl("https://mock.io").Build()
	assert.Nil(t, err)
	err = creator.Impl(&wrongService)
	assert.NotNil(t, err)
	assert.Equal(t, gotten.UnsupportedFuncTypeError(reflect.TypeOf(wrongService.WrongMethod)), err)
}
