// Package tiff handles metadata blocks that use TIFF-style encoding.
package tiff

import (
	"encoding/binary"

	"github.com/rothskeller/photo-tools/metadata"
)

// TIFF is a handler for a TIFF-like file (or portion of file).
type TIFF struct {
	r      metadata.Reader
	enc    binary.ByteOrder
	ifd0   *IFD
	ranges rangelist
}

// IFD is a single Image File Directory contained within the TIFF-like block.
type IFD struct {
	t       *TIFF
	back    *IFD
	offset  uint32
	size    uint32
	tags    []*Tag
	next    uint32
	dirty   bool
	resize  bool
	nextIFD *IFD
}

// Tag is a single tag in an IFD.
type Tag struct {
	ifd    *IFD
	offset uint32
	tag    uint16
	ttype  uint16
	count  uint32
	doff   uint32
	data   []byte
	reader metadata.Reader
	toIFD  *IFD
}
