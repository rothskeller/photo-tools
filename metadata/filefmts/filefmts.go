// Package filefmts defines the interface for file format handlers, and a
// factory function for them.
package filefmts

import (
	"fmt"
	"io"
	"os"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/filefmts/jpeg"
	"github.com/rothskeller/photo-tools/metadata/filefmts/tiff"
	"github.com/rothskeller/photo-tools/metadata/filefmts/xmp"
)

// FileFormat is the interface satisfied by all file format handlers.
type FileFormat interface {
	// Provider returns the metadata provider for the file.
	Provider() metadata.Provider
	// Dirty returns whether the metadata from the file have been changed
	// since they were read (and therefore need to be saved).
	Dirty() bool
	// Save writes the entire file to the supplied writer, including all
	// revised metadata.
	Save(out io.Writer) error
}

// HandlerFor returns a file format handler appropriate for the type of the
// specified file, or nil if there is no handler for the file type.  It returns
// an error if the file cannot be read, or if the handler for its type finds a
// problem with it.
func HandlerFor(fh *os.File) (f FileFormat, err error) {
	if f, err := jpeg.Read(reader{fh}); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	} else if f != nil {
		return f, nil
	}
	if f, err := xmp.Read(reader{fh}); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	} else if f != nil {
		return f, nil
	}
	if f, err := tiff.Read(reader{fh}); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	} else if f != nil {
		return f, nil
	}
	return nil, nil
}

type reader struct {
	*os.File
}

func (r reader) Size() (size int64) {
	var (
		offset int64
		err    error
	)
	if offset, err = r.Seek(0, io.SeekCurrent); err != nil {
		panic(err)
	}
	if size, err = r.Seek(0, io.SeekEnd); err != nil {
		panic(err)
	}
	if _, err = r.Seek(offset, io.SeekStart); err != nil {
		panic(err)
	}
	return size
}
