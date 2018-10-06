[![Coverage Status](https://coveralls.io/repos/github/Hexilee/gotten/badge.svg)](https://coveralls.io/github/Hexilee/gotten)
[![Go Report Card](https://goreportcard.com/badge/github.com/Hexilee/gotten)](https://goreportcard.com/report/github.com/Hexilee/gotten)
[![Build Status](https://travis-ci.org/Hexilee/gotten.svg?branch=master)](https://travis-ci.org/Hexilee/gotten)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Hexilee/gotten/blob/master/LICENSE)
[![Documentation](https://godoc.org/github.com/Hexilee/gotten?status.svg)](https://godoc.org/github.com/Hexilee/gotten)

#### Usage

```go
package example

import (
	"fmt"
	"github.com/Hexilee/gotten"
	"net/http"
	"time"
)

type (
	SimpleParams struct {
		Id   int `type:"path"`
		Page int `type:"query"`
	}

	Item struct {
		TypeId      int
		IId         int
		Name        string
		Description string
	}

	SimpleService struct {
		GetItems func(*SimpleParams) (gotten.Response, error) `method:"GET";path:"itemType/{id}"`
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

func InYourFunc() {
	resp, err := simpleServiceImpl.GetItems(&SimpleParams{1, 1})
	if err == nil && resp.StatusCode() == http.StatusOK {
		result := make([]*Item, 0) 
		err = resp.Unmarshal(&result)
		fmt.Printf("%#v\n", result)
	}
}
```