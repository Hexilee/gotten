package mock

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type (
	WorldHandler struct {
	}
)

func HelloHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "hello")
}

func (WorldHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "world")
}

func TestClientImpl_Do(t *testing.T) {
	clientBuilder := NewClientBuilder()
	clientBuilder.RegisterFunc("hello.me", HelloHandlerFunc)
	clientBuilder.Register("world.me", WorldHandler{})
	client := clientBuilder.Build()

	helloReq, err := http.NewRequest(http.MethodGet, "https://hello.me", nil)
	assert.Nil(t, err)
	worldReq, err := http.NewRequest(http.MethodGet, "https://world.me", nil)
	assert.Nil(t, err)
	wrongReq, err := http.NewRequest(http.MethodGet, "https://wrong.me", nil)
	assert.Nil(t, err)
	for _, testCase := range [] struct {
		request *http.Request
		err     error
		bodyStr string
	}{
		{helloReq, nil, "hello"},
		{worldReq, nil, "world"},
		{wrongReq, HostNotExistError("wrong.me"), ""},
	} {
		resp, err := client.Do(testCase.request)
		if err == nil {
			result, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, testCase.bodyStr, string(result))
		} else {
			fmt.Println(testCase.request.Host)
			assert.NotNil(t, testCase.err)
			assert.Equal(t, testCase.err.Error(), err.Error())
		}
	}
}
