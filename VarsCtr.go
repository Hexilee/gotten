package gotten

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	PathKeyRegexp = `\{[a-zA-Z_][0-9a-zA-Z_]*\}`
	ZeroStr       = ""
)

var (
	pathKeyRegexp, _ = regexp.Compile(PathKeyRegexp)
)

type (
	VarsConstructor interface {
		setValues(value reflect.Value) error
		getUrl() *url.URL
	}

	VarsParser struct {
		regex       *regexp.Regexp
		path        string
		pathKeys    PathKeyList
		pathFields  map[string]*PathField  // pathKey as mapKey
		queryFields map[string]*QueryField // fieldName as mapKey
	}

	VarsCtr struct {
		regex       *regexp.Regexp
		path        string
		pathFields  map[string]*PathField
		queryFields map[string]*QueryField
		queryValues url.Values
	}

	QueryField struct {
		key          string
		defaultValue string
		require      bool
		getValueFunc func(value reflect.Value) string
	}

	PathField struct {
		defaultValue string
		key          string
		value        string
		getValueFunc func(value reflect.Value) string
	}

	PathKeyList map[string]bool
)

func newVarsParser(path string) (*VarsParser, error) {
	pathKeys, err := getPathKeys(pathKeyRegexp, path)
	return &VarsParser{
		regex:       pathKeyRegexp,
		path:        path,
		pathKeys:    pathKeys,
		pathFields:  make(map[string]*PathField),
		queryFields: make(map[string]*QueryField),
	}, err
}

func getPathKeys(regex *regexp.Regexp, path string) (pathKeys PathKeyList, err error) {
	pathKeys = make(PathKeyList)
	patterns := regex.FindAllString(path, -1)
	for _, pattern := range patterns {
		key := getKeyFromPattern(pattern)
		if !pathKeys.addKey(key) {
			err = DuplicatedPathKeyError(key)
			break
		}
	}
	return
}

func (field PathField) getValue() (val string, err error) {
	val = field.value
	if !field.hasValue() {
		val = field.defaultValue
		if !field.hasDefaultValue() {
			err = EmptyPathVariableError(field.key)
		}
	}
	return
}

func (field QueryField) hasDefaultValue() bool {
	return field.defaultValue == ZeroStr
}

func (field PathField) hasDefaultValue() bool {
	return field.defaultValue == ZeroStr
}

func (field PathField) hasValue() bool {
	return field.value == ZeroStr
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
			field := paramElem.Field(i)
			if fieldExportable(field.Name) {
				fieldType := field.Type
				fieldTag := field.Tag
				valueType := fieldTag.Get(KeyType)
				key := processKey(fieldTag.Get(KeyKey), valueType, field.Name)
				defaultValue := fieldTag.Get(KeyDefault)
				require, parseErr := processRequired(fieldTag.Get(KeyRequire))
				if err = parseErr; err == nil {

					switch valueType {
					case TypePath:
						if exist := parser.pathKeys.deleteKey(key); !exist {
							err = UnrecognizedPathKeyError(key)
							break
						}

						pathField := &PathField{
							defaultValue: defaultValue,
							key:          key,
						}
						pathField.getValueFunc, err = FirstValueGetterFunc(fieldType, TypePath)
						if err == nil {
							parser.pathFields[key] = pathField
						}

					case TypeQuery:
						queryField := &QueryField{
							defaultValue: defaultValue,
							key:          key,
							require:      require,
						}
						queryField.getValueFunc, err = FirstValueGetterFunc(fieldType, TypePath)
					// TODO: TypeHeader
					case TypeHeader:
					case TypeJSON:
					case TypeMultipart:
					case TypeForm:
					case TypeXML:
					default:
						err = UnrecognizedFieldTypeError(fieldTag.Get(KeyType))
					}
				}
			}
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

// str: {key}
func (varsCtr VarsCtr) findAndReplace(pattern string) string {
	key := getKeyFromPattern(pattern)
	return varsCtr.pathFields[key].value
}

func (varsCtr VarsCtr) genPath() string {
	return varsCtr.regex.ReplaceAllStringFunc(varsCtr.path, varsCtr.findAndReplace)
}

func getValueFromStringer(value reflect.Value) string {
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
		panic(ValueIsNotIntError(value.Type()))
	}
	return strconv.Itoa(val)
}

func getKeyFromPattern(pattern string) string {
	return strings.TrimRight(strings.TrimLeft(pattern, "{"), "}")
}

func processKey(rawKey, valueType, fieldName string) (key string) {
	key = rawKey
	if key == ZeroStr {
		switch valueType {
		case TypePath:
			fallthrough
		case TypeQuery:
			fallthrough
		case TypeMultipart:
			fallthrough
		case TypeForm:
			key = strcase.ToSnake(fieldName)
		case TypeHeader:
			key = strcase.ToScreamingKebab(fieldName)
		case TypeJSON:
		case TypeXML:
		default:
		}
	}
	return
}

func processRequired(raw string) (required bool, err error) {
	if raw != ZeroStr {
		required, err = strconv.ParseBool(raw)
	}
	return
}

func fieldExportable(fieldName string) bool {
	return unicode.IsUpper(bytes.Runes([]byte{fieldName[0]})[0])
}
