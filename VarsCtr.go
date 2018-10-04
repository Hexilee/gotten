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
	ZeroInt       = 0
)

var (
	pathKeyRegexp, _ = regexp.Compile(PathKeyRegexp)
)

type (
	VarsConstructor interface {
		setValues(value reflect.Value) error
		getUrl() (*url.URL, error)
	}

	VarsParser struct {
		regex       *regexp.Regexp
		path        string
		pathKeys    PathKeyList
		fieldTables []*Field
	}

	VarsCtr struct {
		regex       *regexp.Regexp
		path        string
		fieldTables []*Field
		pathValues  map[string]string
		queryValues url.Values
	}

	Field struct {
		key          string
		name         string
		defaultValue string
		valueType    string
		require      bool
		getValueFunc func(value reflect.Value) string
	}

	//QueryField struct {
	//	key          string
	//	defaultValue string
	//	require      bool
	//	getValueFunc func(value reflect.Value) string
	//}
	//
	//PathField struct {
	//	defaultValue string
	//	key          string
	//	value        string
	//	getValueFunc func(value reflect.Value) string
	//}

	PathKeyList map[string]bool
)

func newVarsParser(path string) (*VarsParser, error) {
	pathKeys, err := getPathKeys(pathKeyRegexp, path)
	return &VarsParser{
		regex:    pathKeyRegexp,
		path:     path,
		pathKeys: pathKeys,
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

func (field Field) getValue(value reflect.Value) (val string, err error) {
	val = field.getValueFunc(value)
	if val == ZeroStr {
		val = field.defaultValue
		if !field.hasDefaultValue() && field.require {
			err = EmptyRequiredVariableError(field.key)
		}
	}
	return
}

func (field Field) hasDefaultValue() bool {
	return field.defaultValue != ZeroStr
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
		parser.fieldTables = make([]*Field, paramElem.NumField())
	FieldCycle:
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
					parser.fieldTables[i] = &Field{
						key:          key,
						name:         field.Name,
						defaultValue: defaultValue,
						valueType:    valueType,
						require:      require,
					}

					switch valueType {
					case TypePath:
						if exist := parser.pathKeys.deleteKey(key); !exist {
							err = UnrecognizedPathKeyError(key)
							break FieldCycle
						}
						parser.fieldTables[i].require = true // path is always required
						fallthrough
					case TypeQuery:
						fallthrough
					case TypeHeader:
						fallthrough
					case TypeForm:
						parser.fieldTables[i].getValueFunc, err = FirstValueGetterFunc(fieldType, TypePath)
						if err != nil {
							break FieldCycle
						}
						// TODO: TypeHeader
					case TypeJSON:
					case TypeMultipart:
					case TypeXML:
					default:
						err = UnrecognizedFieldTypeError(fieldTag.Get(KeyType))
						break FieldCycle
					}
				}
			}
		}
		if err == nil && !parser.pathKeys.empty() {
			err = SomePathVarHasNoValueError(parser.pathKeys)
		}
	}

	return
}

func (parser *VarsParser) Builder() VarsConstructor {
	return &VarsCtr{
		regex:       parser.regex,
		path:        parser.path,
		fieldTables: parser.fieldTables,
		pathValues:  make(map[string]string),
		queryValues: make(url.Values),
	}
}

func (varsCtr VarsCtr) getUrl() (result *url.URL, err error) {
	path := varsCtr.regex.ReplaceAllStringFunc(varsCtr.path, varsCtr.findAndReplace)
	result, err = url.Parse(path)
	if err == nil {
		query := varsCtr.queryValues.Encode()
		result.RawQuery = query
	}
	return
}

func (varsCtr *VarsCtr) setValues(ptr reflect.Value) (err error) {
	value := ptr.Elem()
RangeCycle:
	for i, field := range varsCtr.fieldTables {
		if field != nil {
			fieldValue := value.Field(i)
			switch field.valueType {
			case TypePath:
				varsCtr.pathValues[field.key], err = field.getValue(fieldValue)
				if err != nil {
					break RangeCycle
				}
			case TypeQuery:
				val := ""
				val, err = field.getValue(fieldValue)
				if err != nil {
					break RangeCycle
				}
				varsCtr.queryValues.Add(field.key, val)
			}
		}
	}
	return
}

// str: {key}
func (varsCtr VarsCtr) findAndReplace(pattern string) string {
	key := getKeyFromPattern(pattern)
	return varsCtr.pathValues[key]
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

func getValueFromInt(value reflect.Value) (str string) {
	val, ok := value.Interface().(int)
	if !ok {
		panic(ValueIsNotIntError(value.Type()))
	}

	if val != ZeroInt {
		str = strconv.Itoa(val)
	}
	return
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
