package gotten

import (
	"github.com/Hexilee/gotten/headers"
	"io"
	"net/http"
	"net/url"
)

type (
	Response interface {
		StatusCode() int
		Header() http.Header
		Body() io.ReadCloser
		ContentType() string

		Cookies() []*http.Cookie

		// Location returns the URL of the response's "Location" header,
		// if present. Relative redirects are resolved relative to
		// the Response's Request. ErrNoLocation is returned if no
		// Location header is present.
		Location() (*url.URL, error)

		// ProtoAtLeast reports whether the HTTP protocol used
		// in the response is at least major.minor.
		ProtoAtLeast(major, minor int) bool

		Unmarshal(ptr interface{}) error
	}

	ResponseImpl struct {
		*http.Response
		unmarshaler ReadUnmarshaler
	}
)

func newResponse(resp *http.Response, unmarshaler ReadUnmarshaler) Response {
	return &ResponseImpl{resp, unmarshaler}
}

func (resp ResponseImpl) StatusCode() int {
	return resp.Response.StatusCode
}

func (resp ResponseImpl) Header() http.Header {
	return resp.Response.Header
}

func (resp ResponseImpl) Body() io.ReadCloser {
	return resp.Response.Body
}

func (resp ResponseImpl) ContentType() string {
	return resp.Response.Header.Get(headers.HeaderContentType)
}

func (resp ResponseImpl) Unmarshal(ptr interface{}) error {
	defer resp.Body().Close()
	return resp.unmarshaler.Unmarshal(resp.Body(), resp.Header(), ptr)
}
