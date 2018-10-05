package gotten

import (
	"bytes"
	"fmt"
	"github.com/Hexilee/gotten/headers"
	"github.com/iancoleman/strcase"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	VarsController interface {
		setValues(value reflect.Value) error
		getUrl() (*url.URL, error)
		getBody() (io.Reader, error)
		getHeader() http.Header
		getContentType() string
	}

	VarsParser struct {
		regex        *regexp.Regexp
		path         string
		pathKeys     PathKeyList
		contentType  string
		fieldTable   []*Field
		ioFieldTable []*IOField
	}

	VarsCtr struct {
		regex            *regexp.Regexp
		path             string
		contentType      string
		fieldTable       []*Field
		ioFieldTable     []*IOField
		pathValues       map[string]string
		queryValues      url.Values
		formValues       url.Values
		multipartValues  map[string]string
		multipartFiles   map[string]string
		multipartReaders map[string]MultipartReader
		header           http.Header
		body             io.Reader
	}

	// TypePath, TypeQuery, TypeForm, TypeHeader, TypeCookie, TypeMultipart(except io.Reader)
	Field struct {
		key          string
		name         string
		defaultValue string
		valueType    string
		require      bool
		fieldType    reflect.Type
		// can only called by getValue
		getValueFunc func(value reflect.Value) (string, error)
	}

	// TypeJSON, TypeXML, TypeMultipart(io.Reader)
	IOField struct {
		key          string
		name         string
		defaultValue string
		valueType    string
		require      bool
		// can only called by getValue
		getReaderFunc func(value reflect.Value) (Reader, error)
	}

	MultipartReader struct {
		reader io.Reader
		header http.Header
	}

	PathKeyList map[string]bool
)

func newVarsParser(path string) (*VarsParser, error) {
	// fieldTable and ioFieldTables will be made after num of fields is known
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
	val, err = field.getValueFunc(value)
	if err == nil && val == ZeroStr {
		val = field.defaultValue
		if !field.hasDefaultValue() && field.require {
			err = EmptyRequiredVariableError(field.name)
		}
	}
	return
}

func (field IOField) getValue(value reflect.Value) (val Reader, err error) {
	val, err = field.getReaderFunc(value)
	if err == nil && val.Empty() {
		val = newReadCloser(bytes.NewBufferString(field.defaultValue), !field.hasDefaultValue())
		if val.Empty() && field.require {
			err = EmptyRequiredVariableError(field.name)
		}
	}
	return
}

func (field Field) hasDefaultValue() bool {
	return field.defaultValue != ZeroStr
}

func (field IOField) hasDefaultValue() bool {
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

// can only be called by parse
func (parser *VarsParser) addField(index int, valueType string, field reflect.StructField) (err error) {
	fieldType := field.Type
	fieldTag := field.Tag
	key := processKey(fieldTag.Get(KeyKey), valueType, field.Name)
	defaultValue := fieldTag.Get(KeyDefault)
	require, parseErr := processRequired(fieldTag.Get(KeyRequire))
	if err = parseErr; err == nil {
		parser.fieldTable[index] = &Field{
			key:          key,
			name:         field.Name,
			defaultValue: defaultValue,
			valueType:    valueType,
			fieldType:    fieldType,
			require:      require,
		}
		switch valueType {
		case TypePath:
			if exist := parser.pathKeys.deleteKey(key); !exist {
				err = UnrecognizedPathKeyError(key)
			}
			parser.fieldTable[index].require = true // path is always required
			fallthrough
		case TypeQuery:
			fallthrough
		case TypeHeader:
			parser.fieldTable[index].getValueFunc, err = getValueGetterFunc(fieldType, TypePath)
		case TypeForm:
			parser.fieldTable[index].getValueFunc, err = getValueGetterFunc(fieldType, TypePath)
		case TypeMultipart:
			parser.fieldTable[index].getValueFunc, err = getMultipartValueGetterFunc(fieldType, TypePath)
			// TODO: TypeCookie
		default:
		}
	}
	return
}

// can only be called by parse
func (parser *VarsParser) addIOField(index int, valueType string, field reflect.StructField) (err error) {
	fieldType := field.Type
	fieldTag := field.Tag
	key := processKey(fieldTag.Get(KeyKey), valueType, field.Name)
	defaultValue := fieldTag.Get(KeyDefault)
	require, parseErr := processRequired(fieldTag.Get(KeyRequire))
	if err = parseErr; err == nil {
		parser.ioFieldTable[index] = &IOField{
			key:          key,
			name:         field.Name,
			defaultValue: defaultValue,
			valueType:    valueType,
			require:      require,
		}

		switch valueType {
		case TypeJSON:
			parser.ioFieldTable[index].getReaderFunc, err = getJSONReaderGetterFunc(fieldType, valueType)
		case TypeXML:
			parser.ioFieldTable[index].getReaderFunc, err = getXMLReaderGetterFunc(fieldType, valueType)
		case TypeMultipart:
			parser.ioFieldTable[index].getReaderFunc = getReaderFromReader
		default:
		}
	}
	return
}

func (parser *VarsParser) parse(paramType reflect.Type) (err error) {
	paramElem := paramType.Elem()
	if paramType.Kind() != reflect.Ptr || paramElem.Kind() != reflect.Struct {
		err = ParamTypeMustBePtrOfStructError(paramElem)
	}
	if err == nil {
		parser.fieldTable = make([]*Field, paramElem.NumField())
		parser.ioFieldTable = make([]*IOField, paramElem.NumField())
		for i := 0; i < paramElem.NumField(); i++ {
			field := paramElem.Field(i)
			if fieldExportable(field.Name) {
				valueType := field.Tag.Get(KeyType)
				switch valueType {
				case TypePath:
					fallthrough
				case TypeQuery:
					fallthrough
				case TypeCookie:
					fallthrough
				case TypeHeader:
					err = parser.addField(i, valueType, field)
				case TypeForm:
					err = parser.checkContentType(headers.MIMEApplicationForm)
					if err == nil {
						err = parser.addField(i, valueType, field)
					}
				case TypeJSON:
					err = parser.checkContentType(headers.MIMEApplicationJSONCharsetUTF8)
					if err == nil {
						err = parser.addIOField(i, valueType, field)
					}
				case TypeXML:
					err = parser.checkContentType(headers.MIMEApplicationXMLCharsetUTF8)
					if err == nil {
						err = parser.addIOField(i, valueType, field)
					}
				case TypeMultipart:
					err = parser.checkContentType(headers.MIMEMultipartForm)
					if err == nil {
						if field.Type == ReaderType {
							err = parser.addIOField(i, valueType, field)
						} else {
							err = parser.addField(i, valueType, field)
						}
					}
				default:
					err = UnsupportedValueTypeError(valueType)
				}
				if err != nil {
					break
				}
			}
		}
		if err == nil && !parser.pathKeys.empty() {
			err = SomePathVarHasNoValueError(parser.pathKeys)
		}
	}

	return
}

// can only be called by checkContentType
func (parser *VarsParser) setContentType(contentType string) {
	parser.contentType = contentType
}

// TODO: test it
// can only be called by parse()
func (parser *VarsParser) checkContentType(contentType string) (err error) {
ContentTypeSwitch:
	switch contentType {
	case headers.MIMEMultipartForm:
		switch parser.contentType {
		case headers.MIMEApplicationForm:
			err = ContentTypeConflictError(parser.contentType, contentType)
			break ContentTypeSwitch
		case ZeroStr:
			fallthrough
		case headers.MIMEApplicationJSONCharsetUTF8:
			fallthrough
		case headers.MIMEApplicationXMLCharsetUTF8:
			parser.setContentType(contentType)
		case headers.MIMEMultipartForm:
		default:
			panic("Unsupported content type of parser: " + contentType)
		}
	case headers.MIMEApplicationForm:
		switch parser.contentType {
		case headers.MIMEMultipartForm:
			err = ContentTypeConflictError(parser.contentType, contentType)
			break ContentTypeSwitch
		case ZeroStr:
			fallthrough
		case headers.MIMEApplicationJSONCharsetUTF8:
			fallthrough
		case headers.MIMEApplicationXMLCharsetUTF8:
			parser.setContentType(contentType)
		case headers.MIMEApplicationForm:
		default:
			panic("Unsupported content type of parser: " + contentType)
		}
	case headers.MIMEApplicationJSONCharsetUTF8:
		fallthrough
	case headers.MIMEApplicationXMLCharsetUTF8:
		switch parser.contentType {
		case headers.MIMEApplicationJSONCharsetUTF8:
			fallthrough
		case headers.MIMEApplicationXMLCharsetUTF8:
			err = ContentTypeConflictError(parser.contentType, contentType)
			break ContentTypeSwitch
		case ZeroStr:
			parser.setContentType(contentType)
		case headers.MIMEApplicationForm:
		case headers.MIMEMultipartForm:
		default:
			panic("Unsupported content type of parser: " + contentType)
		}
	default:
		panic("Unsupported content type: " + contentType)
	}
	return
}

func (parser *VarsParser) Build() VarsController {
	varsCtr := &VarsCtr{
		regex:        parser.regex,
		path:         parser.path,
		contentType:  parser.contentType,
		fieldTable:   parser.fieldTable,
		ioFieldTable: parser.ioFieldTable,
		pathValues:   make(map[string]string),
		queryValues:  make(url.Values),
		header:       make(http.Header),
	}

	if varsCtr.contentType != ZeroStr {
		varsCtr.body = bytes.NewBuffer(make([]byte, 0))
	}

	if varsCtr.contentType == headers.MIMEMultipartForm {
		varsCtr.multipartValues = make(map[string]string)
		varsCtr.multipartReaders = make(map[string]MultipartReader)
	}

	if varsCtr.contentType == headers.MIMEApplicationForm {
		varsCtr.formValues = make(url.Values)
	}

	return varsCtr
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

func (varsCtr VarsCtr) getContentType() string {
	return varsCtr.contentType
}

func (varsCtr VarsCtr) getHeader() http.Header {
	return varsCtr.header
}

func (varsCtr VarsCtr) getBody() (body io.Reader, err error) {
	switch varsCtr.contentType {
	case headers.MIMEMultipartForm:
		body, err = varsCtr.getMultipartBody()
	case headers.MIMEApplicationForm:
		body = bytes.NewBufferString(varsCtr.formValues.Encode())
	case headers.MIMEApplicationXMLCharsetUTF8:
		fallthrough
	case headers.MIMEApplicationJSONCharsetUTF8:
		body = varsCtr.body
	}
	return
}

func (varsCtr VarsCtr) getMultipartBody() (body io.ReadWriter, err error) {
	body = bytes.NewBufferString("")
	var partWriter io.Writer
	writer := multipart.NewWriter(body)
	for key, val := range varsCtr.multipartValues {
		writer.WriteField(key, val)
	}

	for key, reader := range varsCtr.multipartReaders {
		if x, ok := reader.reader.(io.Closer); ok {
			defer x.Close()
		}

		if partWriter, err = writer.CreateFormField(key); err == nil {
			_, err = io.Copy(partWriter, reader.reader)
		}

		if err != nil {
			break
		}
	}

	for key, path := range varsCtr.multipartFiles {
		file, err := os.Open(path)
		defer file.Close()
		if err == nil {
			if partWriter, err = writer.CreateFormFile(key, filepath.Base(file.Name())); err == nil {
				_, err = io.Copy(partWriter, file)
			}
		}

		if err != nil {
			break
		}
	}

	writer.Close()
	return
}

func (varsCtr *VarsCtr) setValuesByFields(value reflect.Value) (err error) {
	for i, field := range varsCtr.fieldTable {
		if field != nil {
			fieldValue := value.Field(i)
			var val string
			switch field.valueType {
			case TypePath:
				varsCtr.pathValues[field.key], err = field.getValue(fieldValue)
			case TypeQuery:
				val, err = field.getValue(fieldValue)
				if err == nil {
					varsCtr.queryValues.Add(field.key, val)
				}
			case TypeHeader:
				val, err = field.getValue(fieldValue)
				if err == nil {
					varsCtr.header.Add(field.key, val)
				}
			case TypeForm:
				val, err = field.getValue(fieldValue)
				if err == nil {
					varsCtr.formValues.Add(field.key, val)
				}
			case TypeMultipart:
				val, err = field.getValue(fieldValue)
				if field.fieldType == FilePathType {
					varsCtr.multipartFiles[field.key] = val
				} else {
					varsCtr.multipartValues[field.key] = val
				}
			default:
				panic(UnsupportedValueTypeError(field.valueType))
			}
			if err != nil {
				break
			}
		}
	}
	return
}

func (varsCtr *VarsCtr) setValuesByIOFields(value reflect.Value) (err error) {
	for i, field := range varsCtr.ioFieldTable {
		if field != nil {
			fieldValue := value.Field(i)
			var reader Reader
			switch field.valueType {
			case TypeJSON:
				switch varsCtr.contentType {
				case headers.MIMEApplicationJSONCharsetUTF8:
					reader, err = field.getValue(fieldValue)
					varsCtr.body = reader
				case headers.MIMEApplicationForm:
					var data []byte
					reader, err = field.getValue(fieldValue)
					fmt.Printf("%#v\n", field)
					if err == nil && !reader.Empty() {
						data, err = ioutil.ReadAll(reader)
						varsCtr.formValues.Add(field.key, string(data))
					}
				case headers.MIMEMultipartForm:
					header := make(http.Header)
					header.Add(headers.HeaderContentType, headers.MIMEApplicationJavaScriptCharsetUTF8)
					reader, err = field.getValue(fieldValue)
					varsCtr.multipartReaders[field.key] = MultipartReader{reader, header}
				default:
					panic("Unsupported content type: " + varsCtr.contentType)
				}
			case TypeXML:
				switch varsCtr.contentType {
				case headers.MIMEApplicationXMLCharsetUTF8:
					reader, err = field.getValue(fieldValue)
					varsCtr.body = reader
				case headers.MIMEApplicationForm:
					var data []byte
					if err == nil && !reader.Empty() {
						data, err = ioutil.ReadAll(reader)
						varsCtr.formValues.Add(field.key, string(data))
					}
				case headers.MIMEMultipartForm:
					header := make(http.Header)
					header.Add(headers.HeaderContentType, headers.MIMEApplicationXMLCharsetUTF8)
					reader, err = field.getValue(fieldValue)
					varsCtr.multipartReaders[field.key] = MultipartReader{reader, header}
				default:
					panic("Unsupported content type: " + varsCtr.contentType)
				}
			case TypeMultipart:
				switch varsCtr.contentType {
				case headers.MIMEMultipartForm:
					header := make(http.Header)
					header.Add(headers.HeaderContentType, headers.MIMEOctetStream)
					reader, err = field.getValue(fieldValue)
					varsCtr.multipartReaders[field.key] = MultipartReader{reader, header}
				default:
					panic("Unsupported content type: " + varsCtr.contentType)
				}
			default:
				panic(UnsupportedValueTypeError(field.valueType))
			}
			if err != nil {
				break
			}
		}
	}
	return
}

func (varsCtr *VarsCtr) setValues(ptr reflect.Value) (err error) {
	value := ptr.Elem()
	err = varsCtr.setValuesByFields(value)
	if err == nil {
		err = varsCtr.setValuesByIOFields(value)
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
		case TypeJSON:
			fallthrough
		case TypeXML:
			fallthrough
		case TypeForm:
			key = strcase.ToSnake(fieldName)
		case TypeHeader:
			key = strcase.ToScreamingKebab(fieldName)
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
