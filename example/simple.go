package example

import (
	"github.com/Hexilee/gotten"
	"net/http"
	"time"
)

type (
	SimpleParams struct {
		Id   int
		Page int
	}

	Item struct {
		TypeId      int
		IId         int
		Name        string
		Description string
	}

	ExpectResult   []*Item
	ObjectNotFound struct {
		Key         string
		Reason      string
		Description string
	}

	SimpleService struct {
		GetItems func(SimpleParams) (gotten.Response, error) `method:"GET";path:"itemType/{id}"`
	}
)

var (
	creator, err = gotten.NewBuilder().
			SetBaseUrl("https://api.sample.com").
			AddCookie(&http.Cookie{Name: "clientcookieid", Value: "121", Expires: time.Now().Add(111 * time.Second)}).
			Build()

	simpleServiceImpl = new(SimpleService)
)

func init() {
	err := creator.Impl(simpleServiceImpl)
	if err != nil {
		panic(err)
	}
}
