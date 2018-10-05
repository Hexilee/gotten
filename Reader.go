package gotten

import "io"

type (
	Reader interface {
		io.Reader
		Empty() bool
	}

	ReadCloser interface {
		io.ReadCloser
		Empty() bool
	}

	ReaderImpl struct {
		reader io.Reader
		empty  bool
	}
)

func (reader ReaderImpl) Read(p []byte) (n int, err error) {
	return reader.reader.Read(p)
}

func (reader ReaderImpl) Close() (err error) {
	if x, ok := reader.reader.(io.Closer); ok {
		x.Close()
	}
	return
}

func (reader ReaderImpl) Empty() bool {
	return reader.empty
}

func newReadCloser(reader io.Reader, empty bool) ReadCloser {
	return &ReaderImpl{reader, empty}
}
