// Package jpeg handles JPEG files.
package jpeg

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/exif"
	"github.com/rothskeller/photo-tools/metadata/iptc"
	"github.com/rothskeller/photo-tools/metadata/xmp"
)

// NewHandler returns a handler for a JPEG photo Item.
func NewHandler(path string) (h *JPEG) {
	return &JPEG{path: path}
}

// JPEG is a JPEG file handler.
type JPEG struct {
	path     string
	exif     *exif.EXIF
	iptc     *iptc.IPTC
	xmp      *xmp.XMP
	problems []string
}

// EXIF returns the EXIF block of the JPEG, if it has one.
func (h *JPEG) EXIF() *exif.EXIF {
	if h != nil {
		return h.exif
	}
	return nil
}

// IPTC returns the IPTC block of the JPEG, if it has one.
func (h *JPEG) IPTC() *iptc.IPTC {
	if h != nil {
		return h.iptc
	}
	return nil
}

// XMP returns the XMP block of the JPEG, if it has one.  If create is true, and
// it doesn't have one, one is created and returned.
func (h *JPEG) XMP(create bool) *xmp.XMP {
	if h != nil {
		if h.xmp == nil {
			h.xmp = xmp.New()
		}
		return h.xmp
	}
	return nil
}

func (h *JPEG) log(f string, a ...interface{}) {
	h.problems = append(h.problems, fmt.Sprintf(f, a...))
}

// Problems returns the accumulated set of problems encountered.
func (h *JPEG) Problems() (problems []string) {
	problems = append(problems, h.problems...)
	problems = append(problems, h.exif.Problems...)
	problems = append(problems, h.iptc.Problems...)
	problems = append(problems, h.xmp.Problems...)
	return problems
}
