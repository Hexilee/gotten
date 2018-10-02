package gotten

import (
	"net/url"
	"reflect"
)

type (
	VarsConstructor interface {
		setValues(value reflect.Value) error
		getUrl() *url.URL
	}

	VarsParser struct {
		path    string
		paths   map[string]*PathField
		queries map[string]*QueryField
	}

	VarsCtr struct {
		path        string
		paths       map[string]*PathField
		queries     map[string]*QueryField
		pathPairs   map[string]string
		queryValues url.Values
	}

	QueryField struct {
		key          string
		defaultValue string
		require      bool
		getValue     func(value reflect.Value) string
	}

	PathField struct {
		key          string
		defaultValue string
		getValue     func(value reflect.Value) string
	}
)

func newVarsParser(path string) *VarsParser {
	return &VarsParser{
		path:    path,
		paths:   make(map[string]*PathField),
		queries: make(map[string]*QueryField),
	}
}

func (varsCtr *VarsCtr) getUrl() *url.URL {
	// TODO: *varsCtr.getUrl
	return nil
}

func (varsCtr *VarsCtr) setValues(value reflect.Value) error {
	// TODO: *UrlCtr.setValues
	return nil
}
