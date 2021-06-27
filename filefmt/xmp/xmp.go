// Package xmp contains the file type handler for XMP sidecar files.
package xmp

import (
	"os"
	"path/filepath"

	"github.com/rothskeller/photo-tools/metadata/exif"
	"github.com/rothskeller/photo-tools/metadata/iptc"
	"github.com/rothskeller/photo-tools/metadata/xmp"
)

// XMP is the file format handler for an XMP sidecar file.
type XMP struct {
	path     string
	xmp      *xmp.XMP
	problems []string
}

// NewHandler creates a new XMP file handler for the file at the specified path.
func NewHandler(path string) *XMP {
	return &XMP{path: path}
}

// ReadMetadata reads the metadata for the file.  Callers should check
// Problems() after this call and after any calls to query methods in
// EXIF, IPTC, or XMP.
func (h *XMP) ReadMetadata() {
	var (
		by  []byte
		err error
	)
	if by, err = os.ReadFile(h.path); err != nil {
		h.problems = []string{err.Error()}
		return
	}
	h.xmp = xmp.Parse(by)
}

// EXIF returns nil, since XMP files never have EXIF blocks.
func (h *XMP) EXIF() *exif.EXIF { return nil }

// IPTC returns nil, since XMP files never have IPTC blocks.
func (h *XMP) IPTC() *iptc.IPTC { return nil }

// XMP returns the XMP metadata.
func (h *XMP) XMP(create bool) *xmp.XMP { return h.xmp }

// Problems returns the accumulated set of problems encountered by the
// handler.
func (h *XMP) Problems() (problems []string) {
	problems = append(problems, h.problems...)
	problems = append(problems, h.xmp.Problems...)
	return problems
}

// Dirty returns whether there are any unsaved changes to the metadata.
func (h *XMP) Dirty() bool {
	if len(h.Problems()) != 0 {
		return false
	}
	return h.xmp.Dirty()
}

// SaveMetadata writes the metadata for the photo.  It returns any error
// that occurs.
func (h *XMP) SaveMetadata() (err error) {
	var (
		tempfn string
		fh     *os.File
		by     []byte
	)
	if len(h.Problems()) != 0 {
		panic("XMP SaveMetadata after parse failures")
	}
	if !h.xmp.Dirty() {
		return nil
	}
	tempfn = filepath.Dir(h.path) + "/." + filepath.Base(h.path) + ".TEMP"
	if fh, err = os.Create(tempfn); err != nil {
		return err
	}
	defer fh.Close()
	defer os.Remove(tempfn)
	if by, err = h.xmp.Render(); err != nil {
		return err
	}
	if _, err = fh.Write(by); err != nil {
		return err
	}
	if err = fh.Close(); err != nil {
		return err
	}
	if err = os.Rename(tempfn, h.path); err != nil {
		return err
	}
	return nil
}
