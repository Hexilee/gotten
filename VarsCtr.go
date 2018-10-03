package gotten

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
)

const (
	PathKeyRegexp      = `\{[a-zA-Z_][0-9a-zA-Z_]*\}`
	PathKeyRegexpError = "regexp of path key is wrong: " + PathKeyRegexp
)

var (
	pathKeyRegexp *regexp.Regexp
)

func init() {
	var err error
	pathKeyRegexp, err = regexp.Compile(PathKeyRegexp)
	if err != nil {
		panic(PathKeyRegexpError)
	}
}

type (
	VarsConstructor interface {
		setValues(value reflect.Value) error
		getUrl() *url.URL
	}

	VarsParser struct {
		path         string
		pathSegments []string
		pathKeys     PathKeyList
		pathFields   map[string]*PathField  // pathKey as mapKey
		queryFields  map[string]*QueryField // fieldName as mapKey
	}

	VarsCtr struct {
		path         string
		pathSegments []string
		pathPairs    map[string]string
		pathFields   map[string]*PathField
		queryFields  map[string]*QueryField
		queryValues  url.Values
	}

	QueryField struct {
		key          string
		defaultValue string
		require      bool
		getValue     func(value reflect.Value) string
	}

	PathField struct {
		defaultValue string
		order        int
		getValue     func(value reflect.Value) string
	}

	PathKeyList map[string]bool
)

func newVarsParser(path string) *VarsParser {

	return &VarsParser{
		path:         path,
		pathSegments: make([]string, 0),
		pathKeys:     make(PathKeyList),
		pathFields:   make(map[string]*PathField),
		queryFields:  make(map[string]*QueryField),
	}
}

func (list PathKeyList) addKey(key string) (added bool) {
	if _, ok := list[key]; !ok {
		list[key] = true
		added = true
	}
	return
}

func (list PathKeyList) deleteKey(key string) (exist bool) {
	if _, ok := list[key]; ok {
		delete(list, key)
		exist = true
	}
	return
}

func (list PathKeyList) empty() bool {
	return len(list) == 0
}

func (parser *VarsParser) parse(paramType reflect.Type) (err error) {
	paramElem := paramType.Elem()
	if paramType.Kind() != reflect.Ptr || paramElem.Kind() != reflect.Struct {
		err = ParamTypeMustBePtrOfStructError(paramElem)
	}

	if err == nil {
		for i := 0; i < paramElem.NumField(); i++ {
			// TODO: parse vars
			//field := paramElem.Field(i)

		}
	}

	return
}

func (varsCtr *VarsCtr) getUrl() *url.URL {
	// TODO: *varsCtr.getUrl
	return nil
}

func (varsCtr *VarsCtr) setValues(value reflect.Value) error {
	// TODO: *UrlCtr.setValues
	return nil
}

func getValueFromVar(value reflect.Value) string {
	stringer, ok := value.Interface().(fmt.Stringer)
	if !ok {
		panic(ValueIsNotStringerError(value.Type()))
	}
	return stringer.String()
}

func getValueFromString(value reflect.Value) string {
	val, ok := value.Interface().(string)
	if !ok {
		panic(ValueIsNotStringError(value.Type()))
	}
	return val
}

func getValueFromInt(value reflect.Value) string {
	val, ok := value.Interface().(int)
	if !ok {
		panic(ValueIsNotStringError(value.Type()))
	}
	return strconv.Itoa(val)
}
