package mock

import (
	"errors"
	"github.com/Hexilee/gotten"
	"net/http"
	"net/http/httptest"
)

const (
	HostNotExist = "host not exist"
)

type (
	ClientBuilder struct {
		// key: host or host:port
		services map[string]http.Handler
	}

	ClientImpl struct {
		services map[string]http.Handler
	}
)

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		services: make(map[string]http.Handler),
	}
}

// base: host or host:port
func (builder *ClientBuilder) Register(base string, handler http.Handler) {
	builder.services[base] = handler
}

func (builder *ClientBuilder) RegisterFunc(base string, handler http.HandlerFunc) {
	builder.services[base] = handler
}

func (builder *ClientBuilder) Build() gotten.Client {
	return &ClientImpl{builder.services}
}

func (client ClientImpl) Do(req *http.Request) (resp *http.Response, err error) {
	handler, ok := client.services[req.URL.Host]
	if !ok {
		err = HostNotExistError(req.URL.Host)
	}

	if err == nil {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		resp = recorder.Result()
	}
	return
}

func HostNotExistError(host string) error {
	return errors.New(HostNotExist + ": " + host)
}
