package gotten_test

import (
	"bytes"
	"fmt"
	"github.com/Hexilee/gotten"
	"github.com/Hexilee/gotten/headers"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
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

func TestXMLRequest(t *testing.T) {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		SetClient(mockClient).
		Build()
	assert.Nil(t, err)
	var service AllTypesService
	assert.Nil(t, creator.Impl(&service))
	req, err := service.XMLAllRequest(&XMLAllParams{
		Int:      TestInt,
		Xml:      TestSerializationObject,
		String:   TestXML,
		Stringer: bytes.NewBufferString(TestXML),
		Reader:   bytes.NewBufferString(TestXML),
	})
	assert.Nil(t, err)
	req.ParseMultipartForm(2 << 32)
	assert.Equal(t, TestString, req.PostFormValue("int"))
	assert.Equal(t, TestXML, req.PostFormValue("xml"))
	assert.Equal(t, TestXML, req.PostFormValue("string"))
	assert.Equal(t, TestXML, req.PostFormValue("stringer"))
	assert.Equal(t, TestXML, req.PostFormValue("reader"))
	assert.True(t, strings.HasPrefix(req.Header.Get(headers.HeaderContentType), headers.MIMEMultipartForm))
	assert.Equal(t, "/xml", req.URL.Path)

	req, err = service.XMLAllWithDefaultRequest(&XMLAllWithDefaultParams{})
	assert.Nil(t, err)
	req.ParseMultipartForm(2 << 32)
	assert.Equal(t, TestString, req.PostFormValue("int"))
	assert.Equal(t, TestXML, req.PostFormValue("xml"))
	assert.Equal(t, TestXML, req.PostFormValue("string"))
	assert.Equal(t, TestXML, req.PostFormValue("stringer"))
	assert.Equal(t, TestXML, req.PostFormValue("reader"))
	assert.True(t, strings.HasPrefix(req.Header.Get(headers.HeaderContentType), headers.MIMEMultipartForm))
	assert.Equal(t, "/xml", req.URL.Path)

	req, err = service.XMLSingleRequest(&XMLSingleParams{TestSerializationObject})
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(req.Body)
	assert.Nil(t, err)
	assert.Equal(t, TestXML, string(body))
	assert.Equal(t, headers.MIMEApplicationXMLCharsetUTF8, req.Header.Get(headers.HeaderContentType))
	assert.Equal(t, "/xml", req.URL.Path)
}

func TestJSONRequest(t *testing.T) {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		SetClient(mockClient).
		Build()
	assert.Nil(t, err)
	var service AllTypesService
	assert.Nil(t, creator.Impl(&service))
	req, err := service.JSONSingleParamsRequest(&JSONSingleParams{TestSerializationObject})
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(req.Body)
	assert.Nil(t, err)
	assert.Equal(t, TestJSON, string(body))
	assert.Equal(t, headers.MIMEApplicationJSONCharsetUTF8, req.Header.Get(headers.HeaderContentType))
	assert.Equal(t, "/json", req.URL.Path)
}

func TestHeadersAllRequest(t *testing.T) {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		SetClient(mockClient).
		Build()
	assert.Nil(t, err)
	var service AllTypesService
	assert.Nil(t, creator.Impl(&service))
	req, err := service.HeadersAllParamsRequest(&HeadersAllParams{TestString, TestString, TestString})
	assert.Nil(t, err)
	assert.Equal(t, TestString, req.Header.Get("HOST"))
	assert.Equal(t, TestString, req.Header.Get("LOCATION"))
	assert.Equal(t, TestString, req.Header.Get(headers.HeaderContentType))
}

func TestCookieAllParams(t *testing.T) {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		SetClient(mockClient).
		Build()
	assert.Nil(t, err)
	var service AllTypesService
	assert.Nil(t, creator.Impl(&service))
	req, err := service.CookieAllParamsRequest(&CookieAllParams{
		Ga:         TestString,
		GaId:       TestInt,
		QscSession: TestStringer,
	})
	assert.Nil(t, err)
	cookie, err := req.Cookie("ga")
	assert.Nil(t, err)
	assert.Equal(t, TestString, cookie.Value)
	cookie, err = req.Cookie("ga_id")
	assert.Nil(t, err)
	assert.Equal(t, TestString, cookie.Value)
	cookie, err = req.Cookie("qsc_session")
	assert.Nil(t, err)
	assert.Equal(t, TestString, cookie.Value)
}

type (
	AllTypesService struct {
		FormParamsRequest            func(*FormParams) (*http.Request, error)                        `method:"POST" path:"/form"`
		FormParamsWithDefaultRequest func(withDefault *FormParamsWithDefault) (*http.Request, error) `method:"POST" path:"/form"`
		MultipartRequest             func(*MultipartParams) (*http.Request, error)                   `method:"POST" path:"/multipart"`
		XMLAllRequest                func(*XMLAllParams) (*http.Request, error)                      `method:"POST" path:"/xml"`
		XMLAllWithDefaultRequest     func(*XMLAllWithDefaultParams) (*http.Request, error)           `method:"POST" path:"/xml"`
		XMLSingleRequest             func(*XMLSingleParams) (*http.Request, error)                   `method:"POST" path:"/xml"`
		JSONSingleParamsRequest      func(*JSONSingleParams) (*http.Request, error)                  `method:"POST" path:"/json"`
		HeadersAllParamsRequest      func(*HeadersAllParams) (*http.Request, error)                  `path:"headers"`
		CookieAllParamsRequest       func(*CookieAllParams) (*http.Request, error)                   `path:"cookie"`
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

	XMLAllParams struct {
		Int      int                  `type:"part"`
		Xml      *SerializationStruct `type:"xml"`
		String   string               `type:"xml"`
		Stringer fmt.Stringer         `type:"xml"`
		Reader   io.Reader            `type:"xml"`
	}

	XMLAllWithDefaultParams struct {
		Int      int                  `type:"part" default:"1"`
		Xml      *SerializationStruct `type:"xml" default:"<SerializationStruct><Data>1</Data></SerializationStruct>"`
		String   string               `type:"xml" default:"<SerializationStruct><Data>1</Data></SerializationStruct>"`
		Stringer fmt.Stringer         `type:"xml" default:"<SerializationStruct><Data>1</Data></SerializationStruct>"`
		Reader   io.Reader            `type:"xml" default:"<SerializationStruct><Data>1</Data></SerializationStruct>"`
	}

	XMLSingleParams struct {
		Xml *SerializationStruct `type:"xml"`
	}

	JSONSingleParams struct {
		Json *SerializationStruct `type:"json"`
	}

	HeadersAllParams struct {
		Host        string `type:"header"`
		Location    string `type:"header"`
		ContentType string `type:"header"`
	}

	CookieAllParams struct {
		Ga         string       `type:"cookie"`
		GaId       int          `type:"cookie"`
		QscSession fmt.Stringer `type:"cookie"`
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
