package gotten_test

import (
	"bytes"
	"fmt"
	"github.com/Hexilee/gotten"
	"github.com/Hexilee/gotten/headers"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFormRequest(t *testing.T) {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		SetClient(mockClient).
		Build()
	assert.Nil(t, err)
	var service AllTypesService
	assert.Nil(t, creator.Impl(&service))
	req, err := service.FormParamsRequest(&FormParams{
		JsonBeforeForm: TestSerializationObject,
		Int:            TestInt,
		String:         TestString,
		Stringer:       TestStringer,
		JsonAfterForm:  TestSerializationObject,
		XmlAfterForm:   TestSerializationObject,
	})
	assert.Nil(t, err)
	req.ParseForm()
	assert.Equal(t, TestJSON, req.PostFormValue("json_before_form"))
	assert.Equal(t, TestJSON, req.PostFormValue("json_after_form"))
	assert.Equal(t, TestXML, req.PostFormValue("xml_after_form"))
	assert.Equal(t, TestString, req.PostFormValue("int"))
	assert.Equal(t, TestString, req.PostFormValue("string"))
	assert.Equal(t, TestString, req.PostFormValue("stringer"))
	assert.Equal(t, headers.MIMEApplicationForm, req.Header.Get(headers.HeaderContentType))
	assert.Equal(t, "/form", req.URL.Path)

	req, err = service.FormParamsWithDefaultRequest(&FormParamsWithDefault{})
	assert.Nil(t, err)
	req.ParseForm()
	assert.Equal(t, TestJSON, req.PostFormValue("json_before_form"))
	assert.Equal(t, TestJSON, req.PostFormValue("json_after_form"))
	assert.Equal(t, TestXML, req.PostFormValue("xml_after_form"))
	assert.Equal(t, TestString, req.PostFormValue("int"))
	assert.Equal(t, TestString, req.PostFormValue("string"))
	assert.Equal(t, TestString, req.PostFormValue("stringer"))
	assert.Equal(t, headers.MIMEApplicationForm, req.Header.Get(headers.HeaderContentType))
	assert.Equal(t, "/form", req.URL.Path)
}

func TestMultipartRequest(t *testing.T) {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		SetClient(mockClient).
		Build()
	assert.Nil(t, err)
	var service AllTypesService
	assert.Nil(t, creator.Impl(&service))
	req, err := service.MultipartRequest(&MultipartParams{
		JsonBeforeForm: TestSerializationObject,
		Int:            TestInt,
		String:         TestString,
		Stringer:       TestStringer,
		Reader:         getTestReader(),
		JsonAfterForm:  TestSerializationObject,
		XmlAfterForm:   TestSerializationObject,
	})
	assert.Nil(t, err)
	req.ParseMultipartForm(2 << 32)
	assert.Equal(t, TestJSON, req.PostFormValue("json_before_form"))
	assert.Equal(t, TestJSON, req.PostFormValue("json_after_form"))
	assert.Equal(t, TestXML, req.PostFormValue("xml_after_form"))
	assert.Equal(t, TestString, req.PostFormValue("int"))
	assert.Equal(t, TestString, req.PostFormValue("string"))
	assert.Equal(t, TestString, req.PostFormValue("stringer"))
	assert.Equal(t, TestString, req.PostFormValue("reader"))
	assert.True(t, strings.HasPrefix(req.Header.Get(headers.HeaderContentType), headers.MIMEMultipartForm))
	assert.Equal(t, "/multipart", req.URL.Path)
}

type (
	AllTypesService struct {
		FormParamsRequest            func(*FormParams) (*http.Request, error)                        `method:"POST" path:"/form"`
		FormParamsWithDefaultRequest func(withDefault *FormParamsWithDefault) (*http.Request, error) `method:"POST" path:"/form"`
		MultipartRequest             func(*MultipartParams) (*http.Request, error)                   `method:"POST" path:"/multipart"`
	}

	FormParams struct {
		JsonBeforeForm *SerializationStruct `type:"json"`
		Int            int                  `type:"form"`
		String         string               `type:"form"`
		Stringer       fmt.Stringer         `type:"form"`
		JsonAfterForm  *SerializationStruct `type:"json"`
		XmlAfterForm   *SerializationStruct `type:"xml" `
	}

	FormParamsWithDefault struct {
		JsonBeforeForm *SerializationStruct `type:"json" default:"{\"data\":\"1\"}"`
		Int            int                  `type:"form" default:"1"`
		String         string               `type:"form" default:"1"`
		Stringer       fmt.Stringer         `type:"form" default:"1"`
		JsonAfterForm  *SerializationStruct `type:"json" default:"{\"data\":\"1\"}"`
		XmlAfterForm   *SerializationStruct `type:"xml" default:"<SerializationStruct><Data>1</Data></SerializationStruct>"`
	}

	MultipartParams struct {
		JsonBeforeForm *SerializationStruct `type:"json"`
		Int            int                  `type:"part"`
		String         string               `type:"part"`
		Stringer       fmt.Stringer         `type:"part"`
		Reader         io.Reader            `type:"part"`
		JsonAfterForm  *SerializationStruct `type:"json"`
		XmlAfterForm   *SerializationStruct `type:"xml" `
	}

	SerializationStruct struct {
		Data string `json:"data"`
	}
)

var (
	TestInt                 = 1
	TestString              = "1"
	TestStringer            = bytes.NewBufferString(TestString)
	TestSerializationObject = &SerializationStruct{TestString}
	TestJSON                = "{\"data\":\"1\"}"
	TestXML                 = "<SerializationStruct><Data>1</Data></SerializationStruct>"
)

func getTestReader() io.Reader {
	return bytes.NewBufferString(TestString)
}
