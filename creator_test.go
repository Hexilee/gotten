package gotten_test

import (
	"encoding/json"
	"github.com/Hexilee/gotten"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type (
	SampleService struct {
		Get func(params *GetParams) (gotten.Response, error) `path:"/post/{year}/{month}/{day}"`
	}

	GetParams struct {
		Year  int `type:"path"`
		Month int `type:"path"`
		Day   int `type:"path"`
		Page  int `type:"query" default:"1"`
		Limit int `type:"query" default:"15"`
	}
)

func TestCreator_Impl(t *testing.T) {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		SetClient(mockClient).
		Build()

	assert.Nil(t, err)
	service := new(SampleService)
	assert.Nil(t, creator.Impl(service))
	assert.NotNil(t, service.Get)
	resp, err := service.Get(&GetParams{2018, 10, 1, 1, 10})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	results := make([]TestPost, 0)
	assert.Nil(t, resp.Unmarshal(&results))

	data, err := json.Marshal(&results)
	assert.Nil(t, err)
	assert.Equal(t, "[{\"author\":\"Hexilee\",\"title\":\"Start!\",\"content\":\"Hello world!\"}]", string(data))
}
