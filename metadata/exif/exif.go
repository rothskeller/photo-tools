// Package exif handles EXIF metadata blocks.
package exif

import (
	"encoding/binary"

	"github.com/rothskeller/photo-tools/metadata"
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

	offset  uint32
	buf     []byte
	enc     binary.ByteOrder
	ifd0    *ifdt
	exifIFD *ifdt
	gpsIFD  *ifdt
	ranges  [][]uint32
}

type ifdt struct {
	offset uint32
	size   uint32
	tags   []*tagt
	next   uint32
	dirty  bool
}

type tagt struct {
	offset uint32
	tag    uint16
	ttype  uint16
	count  uint32
	doff   uint32
	data   []byte
}
