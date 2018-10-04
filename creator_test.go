package gotten_test

import (
	"github.com/Hexilee/gotten"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	SampleService struct {
		Get func(params *GetParams) (*GetResult, error) `path:"/post/{year}/{month}/{day}"`
	}

	GetParams struct {
		Year  int `type:"path"`
		Month int `type:"path"`
		Day   int `type:"path"`
		Page  int `type:"query" default:"1"`
		Limit int `type:"query" default:"15"`
	}

	GetResult struct {
		Status gotten.Status
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
	service.Get(&GetParams{2018, 10, 1, 1, 10})
}
