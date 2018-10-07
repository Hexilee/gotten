package unmarshalers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/Hexilee/gotten/headers"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type (
	FilenameStrategy int

	FileCtrBuilder struct {
		basePath      string           // Default: wd
		extensionName string           // .xxx; order: ContentDisposition > extensionName > ContentType
		discard       bool             // Default: false
		strategy      FilenameStrategy // Default: ContentDisposition
	}

	FileCtr struct {
		basePath      string           // Default: wd
		extensionName string           // .xxx; order: ContentDisposition > extensionName > ContentType
		discard       bool             // Default: false
		strategy      FilenameStrategy // Default: ContentDisposition
	}

	FileInfo struct {
		Hash     string
		Filename string
		FilePath string
		BasePath string
		Size     int64
		Ext      string
	}
)

const (
	// Filename
	ContentDisposition FilenameStrategy = iota // filename; as the same as the filename in Content-Disposition
	Hash                                       // hash.ext; md5
	HashHeader                                 // hash-filename
)

const (
	ZeroStr         = ""
	TempFilePattern = "gotten-*.tmp"
)

func (builder *FileCtrBuilder) SetBasePath(path string) *FileCtrBuilder {
	builder.basePath = path
	return builder
}

func (builder *FileCtrBuilder) SetExtName(name string) *FileCtrBuilder {
	builder.extensionName = name
	return builder
}

func (builder *FileCtrBuilder) Discard() *FileCtrBuilder {
	builder.discard = true
	return builder
}

func (builder *FileCtrBuilder) SetStrategy(strategy FilenameStrategy) *FileCtrBuilder {
	builder.strategy = strategy
	return builder
}

func (builder *FileCtrBuilder) Build() (ctr *FileCtr, err error) {
	if builder.basePath == ZeroStr {
		builder.basePath, err = os.Getwd()
	}

	if err == nil {
		ctr = &FileCtr{
			basePath:      builder.basePath,
			extensionName: builder.extensionName,
			discard:       builder.discard,
			strategy:      builder.strategy,
		}
	}
	return
}

func NewFileCtr() (ctr *FileCtr, err error) {
	var dir string
	dir, err = os.Getwd()
	ctr = &FileCtr{
		basePath: dir,
	}
	return
}

func (ctr FileCtr) Unmarshal(reader io.ReadCloser, header http.Header, v interface{}) (err error) {
	fileInfo, ok := v.(*FileInfo)
	if !ok {
		err = errors.New(MustPassPtrOfFileInfo)
	}

	if err == nil {
		contentDisposition := header.Get(headers.HeaderContentDisposition)
		contentType := header.Get(headers.HeaderContentType)
		fileInfo.Filename, fileInfo.Ext, err = ctr.filenameInfo(contentDisposition, contentType)
		if err == nil {
			if fileInfo.Filename == ZeroStr &&
				(ctr.strategy == ContentDisposition || ctr.strategy == HashHeader) {
				err = errors.New(ContentDispositionOrFilenameEmpty)
			}

			if err == nil {
				if ctr.discard {
					err = ctr.hashNotSave(reader, fileInfo)
				} else {
					err = ctr.hashAndSave(reader, fileInfo)
				}
			}
		}
	}
	reader.Close()
	return
}

// 1. if we can parse NOT EMPTY filename, getting ext and return;
// 2. then, if ctr.extensionName is not empty, ext = ctr.extensionName
// 3. then, ext = mime.ExtensionsByType(contentType)[0]
func (ctr FileCtr) filenameInfo(contentDisposition, contentType string) (filename, ext string, err error) {
	if contentDisposition != ZeroStr {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err == nil {
			filename = params["filename"]
			ext = path.Ext(filename)
		}
	}

	if err == nil && filename == ZeroStr {
		ext = ctr.extensionName
		if ext == ZeroStr {
			var exts []string
			exts, err = mime.ExtensionsByType(contentType)
			if err != nil {
				ext = exts[0]
			}
		}
	}
	return
}

func (ctr FileCtr) resolveFilePath(info *FileInfo) (err error) {
	switch ctr.strategy {
	case ContentDisposition:
	case Hash:
		info.Filename = info.Hash + info.Ext
	case HashHeader:
		info.Filename = info.Hash + info.Filename
	default:
		err = UnsupportedFilenameStrategyError(ctr.strategy)
	}

	if err == nil {
		info.FilePath = filepath.Join(ctr.basePath, info.Filename)
	}
	return
}

func (ctr FileCtr) hashNotSave(reader io.Reader, info *FileInfo) (err error) {
	info.Size, info.Hash, err = getSizeAndHash(reader)
	if err == nil {
		err = ctr.resolveFilePath(info)
	}
	return
}

func (ctr FileCtr) hashAndSave(reader io.Reader, info *FileInfo) (err error) {
	var tempFile *os.File
	var tempFilePath string
	tempFile, tempFilePath, err = newTempFile()
	teeReader := io.TeeReader(reader, tempFile)
	info.Size, info.Hash, err = getSizeAndHash(teeReader)
	tempFile.Close()
	if err == nil {
		err = ctr.resolveFilePath(info)
		if err == nil {
			err = os.Rename(tempFilePath, info.FilePath)
		} else {
			err = os.Remove(tempFilePath)
		}
	}
	return
}

func getSizeAndHash(reader io.Reader) (size int64, hash string, err error) {
	hashWriter := md5.New()
	size, err = io.Copy(hashWriter, reader)
	if err == nil {
		hash = fmt.Sprintf("%x", hashWriter.Sum(nil))
	}
	return
}

func newTempFile() (tempFile *os.File, filePath string, err error) {
	tempFile, err = ioutil.TempFile(ZeroStr, TempFilePattern)
	if err == nil {
		filePath = tempFile.Name()
	}
	return
}
