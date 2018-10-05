package gotten

import "io"

type (
	Reader interface {
		io.Reader
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

func (reader ReaderImpl) Empty() bool {
	return reader.empty
}

func newReader(reader io.Reader, empty bool) Reader {
	return &ReaderImpl{reader, empty}
}
