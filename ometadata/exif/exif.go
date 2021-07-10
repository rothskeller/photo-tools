// Package exif handles EXIF metadata blocks.
package exif

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/tifflike"
)

// EXIF is a an EXIF parser and generator.
type EXIF struct {
	artist            []string
	dateTime          metadata.DateTime
	dateTimeDigitized metadata.DateTime
	dateTimeOriginal  metadata.DateTime
	imageDescription  string
	gpsCoords         metadata.GPSCoords
	userComment       string
	Problems          []string

	tl      *tifflike.TIFFLike
	ifd0    *tifflike.IFD
	exifIFD *tifflike.IFD
	gpsIFD  *tifflike.IFD
}
