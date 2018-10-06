package gotten_test

import (
	"github.com/Hexilee/gotten"
	"github.com/stretchr/testify/assert"
	"net/http"
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
