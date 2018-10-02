package gotten

import (
	"github.com/Hexilee/gotten/headers"
	"net/http"
	"net/textproto"
)

type (
	CheckerFactory struct {
		statusSet StatusSet
		headerSet HeaderSet
	}

	CheckerFunc func(*http.Response) bool

	Checker interface {
		Check(*http.Response) bool
	}

	HeaderSet map[string]map[string]bool
	StatusSet map[int]bool
)

func (fn CheckerFunc) Check(resp *http.Response) bool {
	return fn(resp)
}

func (set HeaderSet) add(key string, values ...string) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	for _, value := range values {
		if set[key] == nil {
			set[key] = map[string]bool{value: true}
		} else {
			set[key][value] = true
		}
	}
}

func (set HeaderSet) contain(key, value string) (contain bool) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	if set[key] != nil {
		_, contain = set[key][value]
	}
	return
}

func (set StatusSet) add(statuses ...int) {
	for _, status := range statuses {
		set[status] = true
	}
}

func (set StatusSet) contain(status int) (contain bool) {
	_, contain = set[status]
	return
}

func (factory *CheckerFactory) When(key string, values ...string) *CheckerFactory {
	if factory.headerSet == nil {
		factory.headerSet = make(HeaderSet)
	}
	factory.headerSet.add(key, values...)
	return factory
}

func (factory *CheckerFactory) WhenStatuses(statuses ...int) *CheckerFactory {
	if factory.statusSet == nil {
		factory.statusSet = make(StatusSet)
	}
	factory.statusSet.add(statuses...)
	return factory
}

func (factory *CheckerFactory) WhenContentType(values ...string) *CheckerFactory {
	return factory.When(headers.HeaderContentType, values...)
}

func (factory *CheckerFactory) statusChecker(response *http.Response) bool {
	return factory.statusSet != nil && factory.statusSet.contain(response.StatusCode)
}

func (factory *CheckerFactory) headerChecker(response *http.Response) (ok bool) {
	ok = true
	if factory.headerSet != nil {
		for key := range factory.headerSet {
			value := response.Header.Get(key)
			if !factory.headerSet.contain(key, value) {
				ok = false
				break
			}
		}
	}
	return
}

func (factory *CheckerFactory) Create() (checker CheckerFunc) {
	flag := 1
	if factory.statusSet != nil {
		flag <<= 1
	}

	if factory.headerSet != nil {
		flag <<= 2
	}

	switch flag {
	case 1:
		checker = Any
	case 1 << 1:
		checker = factory.statusChecker
	case 1 << 2:
		checker = factory.headerChecker
	case 1 << 3:
		checker = func(response *http.Response) bool {
			return factory.statusChecker(response) && factory.headerChecker(response)
		}
	default:
	}
	return checker
}

func Any(_ *http.Response) bool {
	return true
}
