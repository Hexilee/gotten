package gotten

import (
	"github.com/Hexilee/gotten/headers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	TestResponse = &http.Response{
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			headers.HeaderContentType: {"text/html"},
		},
	}
)

func TestCheckerFactory_Create(t *testing.T) {
	assert.True(t, new(CheckerFactory).Create().Check(TestResponse))
	assert.True(t, new(CheckerFactory).WhenStatuses(http.StatusOK).Create().Check(TestResponse))
	assert.True(t, new(CheckerFactory).WhenStatuses(http.StatusOK, http.StatusAccepted).Create().Check(TestResponse))
	assert.True(t, new(CheckerFactory).WhenContentType("text/html").Create().Check(TestResponse))
	assert.True(t, new(CheckerFactory).WhenContentType("text/html", "text/xml").Create().Check(TestResponse))
	assert.True(t, new(CheckerFactory).WhenStatuses(http.StatusOK).WhenContentType("text/html").Create().Check(TestResponse))
	assert.False(t, new(CheckerFactory).WhenStatuses(http.StatusAccepted).Create().Check(TestResponse))
	assert.False(t, new(CheckerFactory).WhenContentType("text/xml").Create().Check(TestResponse))
	assert.False(t, new(CheckerFactory).WhenStatuses(http.StatusAccepted).WhenContentType("text/html").Create().Check(TestResponse))
	assert.False(t, new(CheckerFactory).WhenStatuses(http.StatusOK).WhenContentType("text/xml").Create().Check(TestResponse))
	assert.False(t, new(CheckerFactory).WhenStatuses(http.StatusAccepted).WhenContentType("text/xml").Create().Check(TestResponse))
}
