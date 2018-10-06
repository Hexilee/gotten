package gotten

import (
	"github.com/stretchr/testify/assert"
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
