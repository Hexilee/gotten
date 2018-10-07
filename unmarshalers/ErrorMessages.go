package unmarshalers

import (
	"errors"
	"fmt"
)

const (
	ContentDispositionOrFilenameEmpty = "Content-Disposition or filename is empty"
	UnsupportedFilenameStrategy       = "filename strategy is not supported"
	MustPassPtrOfFileInfo             = "parameter must be Ptr of FileInfo"
)

func UnsupportedFilenameStrategyError(strategy FilenameStrategy) error {
	return errors.New(fmt.Sprintf(UnsupportedFilenameStrategy+": %d", strategy))
}
