package unmarshalers

import (
	"errors"
	"github.com/Hexilee/gotten/headers"
	"io"
	"mime"
	"net/http"
	"path"
)

type (
	Filename int
	FileUnmarshaler struct {
		basePath      string // Default: wd
		extensionName string // ZeroStr means keep the same extension name with Content-Disposition or Content-Type
		discard       bool
		filename      Filename // Default: ContentDisposition
	}
)

const (
	// Filename
	ContentDisposition Filename = iota // as the same as the filename in Content-Disposition
	Hash                               // md5
	HashHeader                         // hash-filename
)

const (
	ZeroStr = ""
)

func (unmarshaler FileUnmarshaler) Unmarshal(reader io.Reader, header http.Header, v interface{}) (err error) {
	var filename string
	var ext string
	var hash string
	contentDisposition := header.Get(headers.HeaderContentDisposition)
	if contentDisposition != ZeroStr {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err == nil {
			filename = params["filename"]
			ext = path.Ext(filename)
		}
	}

	if err == nil {
		if (contentDisposition == ZeroStr || filename == ZeroStr) &&
			(unmarshaler.filename == ContentDisposition || unmarshaler.filename == HashHeader) {
			err = errors.New(ContentDispositionOrFilenameEmpty)
		}

		if err == nil {
			switch unmarshaler.filename {
			case ContentDisposition:
				_ = ext
				_ = hash
			}
		}
	}
	return
}
