package gotten

import (
	"github.com/Hexilee/gotten/headers"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetMultipartValueGetterFunc(t *testing.T) {
	_, err := getMultipartValueGetterFunc(ReaderType, TypeMultipart)
	assert.Error(t, err)
	assert.Equal(t, UnsupportedFieldTypeError(ReaderType, TypeMultipart), err)
}

func TestGetValueGetterFunc(t *testing.T) {
	_, err := getValueGetterFunc(ReaderType, TypePath)
	assert.Error(t, err)
	assert.Equal(t, UnsupportedFieldTypeError(ReaderType, TypePath), err)
}

func TestGetReaderGetterFunc(t *testing.T) {
	_, err := getReaderGetterFunc(FilePathType, TypeJSON)
	assert.Error(t, err)
	assert.Equal(t, UnsupportedFieldTypeError(FilePathType, TypeJSON), err)
}

func TestGetValueFromFilePath(t *testing.T) {
	reader := 1
	defer func() {
		err := recover()
		assert.Equal(t, reflect.TypeOf(reader).String()+" is not FilePath", err)
	}()
	getValueFromFilePath(reflect.ValueOf(reader))
}

func TestGetValueFromString(t *testing.T) {
	reader := 1
	defer func() {
		err := recover()
		assert.Equal(t, ValueIsNotStringError(reflect.TypeOf(reader)), err)
	}()
	getValueFromString(reflect.ValueOf(reader))
}

func TestGetValueFromInt(t *testing.T) {
	reader := "1"
	defer func() {
		err := recover()
		assert.Equal(t, ValueIsNotIntError(reflect.TypeOf(reader)), err)
	}()
	getValueFromInt(reflect.ValueOf(reader))
}

func TestGetReaderFromString(t *testing.T) {
	reader := 1
	defer func() {
		err := recover()
		assert.Equal(t, ValueIsNotStringError(reflect.TypeOf(reader)), err)
	}()
	getReaderFromString(reflect.ValueOf(reader))
}

func TestCheckContentType1(t *testing.T) {
	contentType := headers.MIMEMultipartForm
	parser := new(VarsParser)
	parser.contentType = headers.MIMETextPlain
	defer func() {
		err := recover()
		assert.Equal(t, "Unsupported content type of parser: "+parser.contentType, err)
	}()
	parser.checkContentType(contentType)
}

func TestCheckContentType2(t *testing.T) {
	parser := new(VarsParser)
	parser.contentType = headers.MIMETextPlain
	contentType := headers.MIMEApplicationForm
	defer func() {
		err := recover()
		assert.Equal(t, "Unsupported content type of parser: "+parser.contentType, err)
	}()

	parser.checkContentType(contentType)
}

func TestCheckContentType3(t *testing.T) {
	contentType := headers.MIMEApplicationJSONCharsetUTF8
	parser := new(VarsParser)
	parser.contentType = headers.MIMETextPlain
	defer func() {
		err := recover()
		assert.Equal(t, "Unsupported content type of parser: "+parser.contentType, err)
	}()
	parser.checkContentType(contentType)
}

func TestCheckContentType4(t *testing.T) {
	contentType := headers.MIMETextPlain
	defer func() {
		err := recover()
		assert.Equal(t, "Unsupported content type of parser: "+contentType, err)
	}()
	parser := new(VarsParser)
	parser.contentType = headers.MIMEApplicationJSONCharsetUTF8
	parser.checkContentType(contentType)
}

func TestSetValuesByFields(t *testing.T) {

}
