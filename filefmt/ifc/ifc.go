package ifc

import (
	"github.com/rothskeller/photo-tools/metadata/exif"
	"github.com/rothskeller/photo-tools/metadata/iptc"
	"github.com/rothskeller/photo-tools/metadata/xmp"
)

// FileHandler is the interface for handlers of various different photo file
// types.
type FileHandler interface {
	// ReadMetadata reads the metadata for the photo.  Callers should check
	// Problems() after this call and after any calls to query methods in
	// EXIF, IPTC, or XMP.
	ReadMetadata()
	// EXIF returns the EXIF metadata, if the photo has any.
	EXIF() *exif.EXIF
	// IPTC returns the IPTC metadata, if the photo has any.
	IPTC() *iptc.IPTC
	// XMP returns the XMP metadata, if the photo has any.  If create is
	// is true, and the photo supports XMP metadata, an XMP block will be
	// created if none already exists.
	XMP(create bool) *xmp.XMP
	// Problems returns the accumulated set of problems encountered by the
	// handler.
	Problems() []string
	// SaveMetadata writes the metadata for the photo.  It returns any error
	// that occurs.
	SaveMetadata() error
}
