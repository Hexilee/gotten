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
		GetPosts func(params *GetPostsParams) (gotten.Response, error) `path:"/post/{year}/{month}/{day}"`
		AddPost  func(params *AddPostParams) (gotten.Response, error)  `method:"POST" path:"/post/{year}/{month}/{day}"`
	}

	GetPostsParams struct {
		Year  int `type:"path"`
		Month int `type:"path"`
		Day   int `type:"path"`
		Page  int `type:"query" default:"1"`
		Limit int `type:"query" default:"15"`
	}

	AddPostParams struct {
		Year  int       `type:"path"`
		Month int       `type:"path"`
		Day   int       `type:"path"`
		Post  *TestPost `type:"json"`
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
	assert.NotNil(t, service.GetPosts)
	resp, err := service.GetPosts(&GetPostsParams{2018, 10, 1, 1, 10})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	results := make([]TestPost, 0)
	assert.Nil(t, resp.Unmarshal(&results))

	data, err := json.Marshal(&results)
	assert.Nil(t, err)
	assert.Equal(t, "[{\"author\":\"Hexilee\",\"title\":\"Start!\",\"content\":\"Hello world!\"}]", string(data))

	resp, err = service.AddPost(&AddPostParams{
		Year:  2018,
		Month: 10,
		Day:   1,
		Post:  &TestPost{"Hexilee", "AddPost Test", "Success!"},
	})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
	var addedResult AddedData
	assert.Nil(t, resp.Unmarshal(&addedResult))
	assert.True(t, addedResult.Success)
	assert.Equal(t, 2018, addedResult.Year)
	assert.Equal(t, 10, addedResult.Month)
	assert.Equal(t, 1, addedResult.Day)
	assert.Equal(t, 2, addedResult.Order)
}
