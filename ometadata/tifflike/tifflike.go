// Package tifflike handles metadata blocks that use TIFF-style encoding.
package tifflike

import (
	"encoding/binary"
	"io"
)

// TIFFLike is a handler for a TIFF-like file (or portion of file).
type TIFFLike struct {
	r      tiffReader
	enc    binary.ByteOrder
	ifd0   *IFD
	ranges rangelist
}

// IFD is a single Image File Directory contained within the TIFF-like block.
type IFD struct {
	t       *TIFFLike
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
	toIFD  *IFD
}

// tiffReader is the interface that the reader passed to NewTIFFLike must
// satisfy.
type tiffReader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}
