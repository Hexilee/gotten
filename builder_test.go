package gotten_test

import (
	"github.com/Hexilee/gotten"
	"github.com/Hexilee/gotten/headers"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

type (
	EmptyParams struct {
	}

	EmptyService struct {
		EmptyGet func(*EmptyParams) (*http.Request, error)
	}
)

func TestBuilder(t *testing.T) {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		AddCookie(&http.Cookie{Name: "ga", Value: TestString}).
		AddCookies([]*http.Cookie{
			{Name: "ga_id", Value: TestString},
			{Name: "qsc_session", Value: TestString},
		}).AddHeader("HOST", "mock.io").
		SetHeader("HOST", "hexilee.me").
		Build()

	assert.Nil(t, err)
	var service EmptyService
	assert.Nil(t, creator.Impl(&service))
	req, err := service.EmptyGet(&EmptyParams{})
	assert.Nil(t, err)

	cookie, err := req.Cookie("ga_id")
	assert.Nil(t, err)
	assert.Equal(t, TestString, cookie.Value)

	cookie, err = req.Cookie("ga")
	assert.Nil(t, err)
	assert.Equal(t, TestString, cookie.Value)

	cookie, err = req.Cookie("qsc_session")
	assert.Nil(t, err)
	assert.Equal(t, TestString, cookie.Value)

	assert.Equal(t, "hexilee.me", req.Header.Get("HOST"))
}

func TestBuilder_AddUnmarshalFunc(t *testing.T) {
	type TextService struct {
		Get func(*struct{}) (gotten.Response, error) `path:"/text"`
	}

	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		AddUnmarshalFunc(func(data []byte, v interface{}) (err error) {
			var success bool
			success, err = strconv.ParseBool(string(data))
			if err == nil {
				value := reflect.ValueOf(v)
				if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Bool {
					value.Elem().SetBool(success)
				}
			}
			return
		}, new(gotten.CheckerFactory).WhenContentType(headers.MIMETextPlain).Create()).
		SetClient(mockClient).
		Build()

	assert.Nil(t, err)
	service := new(TextService)
	assert.Nil(t, creator.Impl(service))
	assert.NotNil(t, service.Get)

	resp, err := service.Get(nil)
	assert.Nil(t, err)
	var success bool
	assert.Nil(t, resp.Unmarshal(&success))
	assert.True(t, success)
}

func TestBuilder_AddReaderUnmarshalerFunc(t *testing.T) {
	type TextService struct {
		Get func(*struct{}) (gotten.Response, error) `path:"/text"`
	}

	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		AddReadUnmarshalFunc(func(body io.ReadCloser, _ http.Header, v interface{}) (err error) {
			var data []byte
			data, err = ioutil.ReadAll(body)
			body.Close()
			if err == nil {
				var success bool
				success, err = strconv.ParseBool(string(data))
				if err == nil {
					value := reflect.ValueOf(v)
					if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Bool {
						value.Elem().SetBool(success)
					}
				}
			}
			return
		}, new(gotten.CheckerFactory).WhenContentType(headers.MIMETextPlain).Create()).
		SetClient(mockClient).
		Build()

	assert.Nil(t, err)
	service := new(TextService)
	assert.Nil(t, creator.Impl(service))
	assert.NotNil(t, service.Get)

	resp, err := service.Get(nil)
	assert.Nil(t, err)
	var success bool
	assert.Nil(t, resp.Unmarshal(&success))
	assert.True(t, success)
}
