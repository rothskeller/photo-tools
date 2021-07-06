// Package exif handles EXIF metadata blocks.
package exif

import (
	"bytes"
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/tifflike"
)

const (
	tagExifIFDOffset uint16 = 0x8769
	tagGPSIFDOffset  uint16 = 0x8825
)

// Parse parses an EXIF block and returns the parse results.  offset is the
// starting offset of the EXIF block within the image file; it is used in
// problem messages.
func Parse(buf []byte, offset uint32) (exif *EXIF) {
	var err error

	exif = new(EXIF)
	if exif.tl, err = tifflike.NewTIFFLike(bytes.NewReader(buf)); err != nil {
		exif.Problems = append(exif.Problems, err.Error())
		return exif
	}
	exif.ifd0 = exif.tl.IFD0()
	if tag := exif.ifd0.Tag(tagExifIFDOffset); tag != nil {
		if exif.exifIFD, err = tag.AsIFD(); err != nil {
			exif.Problems = append(exif.Problems, err.Error())
			exif.ifd0 = nil
			return exif
		}
	}
	if tag := exif.ifd0.Tag(tagGPSIFDOffset); tag != nil {
		if exif.gpsIFD, err = tag.AsIFD(); err != nil {
			exif.Problems = append(exif.Problems, err.Error())
			exif.ifd0, exif.exifIFD = nil, nil
			return exif
		}
	}
	exif.getArtist()
	exif.getDateTime()
	exif.getImageDescription()
	if exif.exifIFD != nil {
		exif.getUserComment()
		exif.getDateTimeDigitized()
		exif.getDateTimeOriginal()
	}
	if exif.gpsIFD != nil {
		exif.getGPSCoords()
	}
	return exif
}

func (p *EXIF) log(f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	p.Problems = append(p.Problems, fmt.Sprintf("EXIF: %s", s))
}
