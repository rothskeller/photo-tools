// Package filefmt handles selecting the correct file format handler for a file.
package filefmt

import (
	"path/filepath"
	"strings"

	"github.com/rothskeller/photo-tools/filefmt/ifc"
	"github.com/rothskeller/photo-tools/filefmt/jpeg"
	"github.com/rothskeller/photo-tools/filefmt/xmp"
)

// FileHandler is the interface satisfied by a handler returned from HandlerFor.
type FileHandler = ifc.FileHandler

// HandlerFor returns the FileHandler for the file at the specified path, or nil
// if the file format isn't supported.
func HandlerFor(path string) (handler ifc.FileHandler) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		return jpeg.NewHandler(path)
	case ".xmp":
		return xmp.NewHandler(path)
	}
	return nil
}
