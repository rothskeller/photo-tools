// Package filefmts defines the interface for file format handlers, and a
// factory function for them.
package filefmts

import (
	"os"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/filefmts/jpeg"
	"github.com/rothskeller/photo-tools/metadata/filefmts/xmp"
)

// FileFormat is the interface satisfied by all file format handlers.
type FileFormat interface {
	// Provider returns the metadata provider for the file.
	Provider() metadata.Provider
}

// HandlerFor returns a file format handler appropriate for the type of the
// specified file, or nil if there is no handler for the file type.  It returns
// an error if the file cannot be read, or if the handler for its type finds a
// problem with it.
func HandlerFor(fh *os.File) (f FileFormat, err error) {
	if f, err := jpeg.Read(fh); err != nil {
		return nil, err
	} else if f != nil {
		return f, nil
	}
	if f, err := xmp.Read(fh); err != nil {
		return nil, err
	} else if f != nil {
		return f, nil
	}
	return nil, nil
}
