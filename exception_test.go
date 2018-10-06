package gotten

import (
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
