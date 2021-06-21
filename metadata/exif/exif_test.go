package exif

import (
	"bytes"
	"testing"
)

// The minimal case is a TIFF header with an empty IFD0.
var minimal = []byte{
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // TIFF header
	/* 0004 */ 0x00, 0x00, 0x00, 0x08, // Pointer to IFD0
	/* 0008 */ 0x00, 0x00, // Count of zero tags
	/* 000A */ 0x00, 0x00, 0x00, 0x00, // no next IFD
}

func TestMinimal(t *testing.T) {
	exif := Parse(minimal, 0)
	if exif == nil {
		t.Error("Parse failed")
		return
	}
	out := exif.Render(1000)
	if &out[0] != &minimal[0] || len(out) != len(minimal) {
		t.Error("Render rewrote minimal")
	}
}

func TestMinimalRewrite(t *testing.T) {
	exif := Parse(minimal, 0)
	exif.ifd0.dirty = true
	out := exif.Render(1000)
	if &out[0] == &minimal[0] {
		t.Error("Render didn't rewrite minimal")
	}
	if !bytes.Equal(out, minimal) {
		t.Error("Render rewrite of minimal changed it")
	}
}

func TestMinimalSetExif(t *testing.T) {
	exif := Parse(minimal, 0)
	exif.SetImageUniqueID("unique")
	out := exif.Render(1000)
	if &out[0] != &minimal[0] || len(out) != len(minimal) {
		t.Error("Render rewrote minimal")
	}
}

func TestMinimalSetIFD0(t *testing.T) {
	exif := Parse(minimal, 0)
	exif.SetImageDescription("desc")
	out := exif.Render(1000)
	if !bytes.Equal(out, minimalSetIFDExpected) {
		t.Error("wrong output")
	}
}

var minimalSetIFDExpected = []byte{
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // TIFF header
	/* 0004 */ 0x00, 0x00, 0x00, 0x08, // Pointer to IFD0
	/* 0008 */ 0x00, 0x01, // Count of one tag
	/* 000A */ 0x01, 0x0E, 0x00, 0x02, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x1A, // Tag
	/* 0016 */ 0x00, 0x00, 0x00, 0x00, // no next IFD
	/* 001A */ 'd', 'e', 's', 'c', 0x00, // text for tag
}

// The minexif case is a TIFF header with an empty ExifIFD and an IFD0 that is
// empty except for the pointer to ExifIFD.  Just for coverage, we make it
// little endian.
var minexif = []byte{
	/* 0000 */ 0x49, 0x49, 0x2A, 0x00, // TIFF header
	/* 0004 */ 0x08, 0x00, 0x00, 0x00, // Pointer to IFD0
	/* 0008 */ 0x01, 0x00, // Count of one tag
	/* 000A */ 0x69, 0x87, 0x04, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, // ExifOffset
	/* 0016 */ 0x00, 0x00, 0x00, 0x00, // no next IFD
	/* 001A */ 0x00, 0x00, // Count of zero tags
	/* 001C */ 0x00, 0x00, 0x00, 0x00, // no next IFD
}

func TestMinExif(t *testing.T) {
	exif := Parse(minexif, 0)
	if exif == nil {
		t.Error("Parse failed")
		return
	}
	out := exif.Render(1000)
	if &out[0] != &minexif[0] || len(out) != len(minexif) {
		t.Error("Render rewrote minexif")
	}
}

func TestMinExifRewrite(t *testing.T) {
	exif := Parse(minexif, 0)
	exif.ifd0.dirty = true
	exif.exifIFD.dirty = true
	out := exif.Render(1000)
	if &out[0] == &minexif[0] {
		t.Error("Render didn't rewrite minexif")
	}
	if !bytes.Equal(out, minexif) {
		t.Error("Render rewrite of minexif changed it")
	}
}

func TestMinExifSet(t *testing.T) {
	exif := Parse(minexif, 0)
	exif.SetImageDescription("desc")
	exif.SetImageUniqueID("unique")
	out := exif.Render(1000)
	if !bytes.Equal(out, minExifSetExpected) {
		t.Error("wrong output")
	}
	exif2 := Parse(out, 0)
	if exif2.ImageDescription() != "desc" {
		t.Error("wrong desc")
	}
	if exif2.ImageUniqueID() != "unique" {
		t.Errorf("wrong id %q", exif.ImageUniqueID())
	}
}

var minExifSetExpected = []byte{
	/* 0000 */ 0x49, 0x49, 0x2A, 0x00, // TIFF header
	/* 0004 */ 0x08, 0x00, 0x00, 0x00, // Pointer to IFD0
	/* 0008 */ 0x02, 0x00, // Count of two tags
	/* 000A */ 0x0E, 0x01, 0x02, 0x00, 0x05, 0x00, 0x00, 0x00, 0x26, 0x00, 0x00, 0x00, // ImageDescription
	/* 0016 */ 0x69, 0x87, 0x04, 0x00, 0x01, 0x00, 0x00, 0x00, 0x2B, 0x00, 0x00, 0x00, // ExifOffset
	/* 0022 */ 0x00, 0x00, 0x00, 0x00, // no next IFD
	/* 0026 */ 'd', 'e', 's', 'c', 0x00,
	/* 002B */ 0x01, 0x00, // Count of one tag
	/* 002D */ 0x20, 0xA4, 0x02, 0x00, 0x07, 0x00, 0x00, 0x00, 0x3D, 0x00, 0x00, 0x00, // ImageUniqueID
	/* 0039 */ 0x00, 0x00, 0x00, 0x00, // no next IFD
	/* 003D */ 'u', 'n', 'i', 'q', 'u', 'e', 0x00,
}

// bounded is a basis for testing of shrinking and growing.  It has a foreign
// (i.e., not managed by this app) IFD at the end, so that nothing can grow.
// And we're going back to big endian because it's easier to write tests for.
var bounded = []byte{
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // TIFF header
	/* 0004 */ 0x00, 0x00, 0x00, 0x08, // Pointer to IFD0
	/* 0008 */ 0x00, 0x02, // IFD0: 2 tags
	/* 000A */ 0x01, 0x0E, 0x00, 0x02, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x26, // ImageDescription
	/* 0016 */ 0x87, 0x69, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x2B, // ExifOffset
	/* 0022 */ 0x00, 0x00, 0x00, 0x44, // point to foreign IFD
	/* 0026 */ 'd', 'e', 's', 'c', 0x00,
	/* 002B */ 0x00, 0x01, // ExifIFD: 1 tag
	/* 002D */ 0xA4, 0x20, 0x00, 0x02, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x3D, // ImageUniqueID
	/* 0039 */ 0x00, 0x00, 0x00, 0x00, // no next IFD
	/* 003D */ 'u', 'n', 'i', 'q', 'u', 'e', 0x00,
	/* 0044 */ 0x00, 0x01, // Foreign IFD: 1 tag
	/* 0046 */ 0xEE, 0xEE, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0xEE, 0xEE, 0xEE, 0xEE, // fake tag
	/* 0052 */ 0x00, 0x00, 0x00, 0x00, // No next IFD
}

func TestBoundedRewrite(t *testing.T) {
	exif := Parse(bounded, 0)
	exif.ifd0.dirty = true
	exif.exifIFD.dirty = true
	out := exif.Render(1000)
	if &out[0] == &bounded[0] {
		t.Error("Render didn't rewrite bounded")
	}
	if !bytes.Equal(out, bounded) {
		t.Error("Render rewrite of bounded changed it")
	}
}

func TestBoundedShrink(t *testing.T) {
	exif := Parse(bounded, 0)
	exif.SetImageDescription("")
	out := exif.Render(1000)
	if !bytes.Equal(out, boundedShrinkExpected) {
		t.Error("wrong output")
	}
}

var boundedShrinkExpected = []byte{
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // TIFF header
	/* 0004 */ 0x00, 0x00, 0x00, 0x08, // Pointer to IFD0
	/* 0008 */ 0x00, 0x01, // IFD0: 1 tag
	/* 000A */ 0x87, 0x69, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1A, // ExifOffset
	/* 0016 */ 0x00, 0x00, 0x00, 0x44, // point to foreign IFD
	/* 001A */ 0x00, 0x01, // ExifIFD: 1 tag
	/* 001C */ 0xA4, 0x20, 0x00, 0x02, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x2C, // ImageUniqueID
	/* 0028 */ 0x00, 0x00, 0x00, 0x00, // no next IFD
	/* 002C */ 'u', 'n', 'i', 'q', 'u', 'e', 0x00,
	/* 0033 */ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // padding
	/* 0044 */ 0x00, 0x01, // Foreign IFD: 1 tag
	/* 0046 */ 0xEE, 0xEE, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0xEE, 0xEE, 0xEE, 0xEE, // fake tag
	/* 0052 */ 0x00, 0x00, 0x00, 0x00, // No next IFD
}

func TestBoundedGrow(t *testing.T) {
	exif := Parse(bounded, 0)
	exif.SetImageUniqueID("unique2")
	out := exif.Render(1000)
	if !bytes.Equal(out, boundedGrowExpected) {
		t.Error("wrong output")
	}
}

var boundedGrowExpected = []byte{
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // TIFF header
	/* 0004 */ 0x00, 0x00, 0x00, 0x08, // Pointer to IFD0
	/* 0008 */ 0x00, 0x02, // IFD0: 2 tags
	/* 000A */ 0x01, 0x0E, 0x00, 0x02, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x26, // ImageDescription
	/* 0016 */ 0x87, 0x69, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x56, // ExifOffset
	/* 0022 */ 0x00, 0x00, 0x00, 0x44, // point to foreign IFD
	/* 0026 */ 'd', 'e', 's', 'c', 0x00,
	/* 002B */ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // padding
	/* 0044 */ 0x00, 0x01, // Foreign IFD: 1 tag
	/* 0046 */ 0xEE, 0xEE, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0xEE, 0xEE, 0xEE, 0xEE, // fake tag
	/* 0052 */ 0x00, 0x00, 0x00, 0x00, // No next IFD
	/* 0056 */ 0x00, 0x01, // ExifIFD: 1 tag
	/* 0058 */ 0xA4, 0x20, 0x00, 0x02, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x68, // ImageUniqueID
	/* 0064 */ 0x00, 0x00, 0x00, 0x00, // no next IFD
	/* 0068 */ 'u', 'n', 'i', 'q', 'u', 'e', '2', 0x00,
}

func TestBoundedGrow2(t *testing.T) {
	exif := Parse(bounded, 0)
	exif.SetImageDescription("desc2")
	out := exif.Render(1000)
	if !bytes.Equal(out, boundedGrow2Expected) {
		t.Error("wrong output")
	}
}

var boundedGrow2Expected = []byte{
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // TIFF header
	/* 0004 */ 0x00, 0x00, 0x00, 0x08, // Pointer to IFD0
	/* 0008 */ 0x00, 0x02, // IFD0: 2 tags
	/* 000A */ 0x01, 0x0E, 0x00, 0x02, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x26, // ImageDescription
	/* 0016 */ 0x87, 0x69, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x56, // ExifOffset
	/* 0022 */ 0x00, 0x00, 0x00, 0x44, // point to foreign IFD
	/* 0026 */ 'd', 'e', 's', 'c', '2', 0x00,
	/* 002C */ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // padding
	/* 0044 */ 0x00, 0x01, // Foreign IFD: 1 tag
	/* 0046 */ 0xEE, 0xEE, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0xEE, 0xEE, 0xEE, 0xEE, // fake tag
	/* 0052 */ 0x00, 0x00, 0x00, 0x00, // No next IFD
	/* 0056 */ 0x00, 0x01, // ExifIFD: 1 tag
	/* 0058 */ 0xA4, 0x20, 0x00, 0x02, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x68, // ImageUniqueID
	/* 0064 */ 0x00, 0x00, 0x00, 0x00, // no next IFD
	/* 0068 */ 'u', 'n', 'i', 'q', 'u', 'e', 0x00,
}

func TestBoundedGrow3(t *testing.T) {
	exif := Parse(bounded, 0)
	// Make ExifIFD longer than IFD0.
	exif.SetImageUniqueID("something longer!")
	out := exif.Render(1000)
	if !bytes.Equal(out, boundedGrow3Expected) {
		t.Error("wrong output")
	}
	exif2 := Parse(out, 0)
	if exif2.ImageDescription() != "desc" {
		t.Error("wrong desc")
	}
	if exif2.ImageUniqueID() != "something longer!" {
		t.Error("wrong id")
	}
}

var boundedGrow3Expected = []byte{
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // TIFF header
	/* 0004 */ 0x00, 0x00, 0x00, 0x56, // Pointer to IFD0
	/* 0008 */ 0x00, 0x01, // ExifIFD: 1 tag
	/* 000A */ 0xA4, 0x20, 0x00, 0x02, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00, 0x1A, // ImageUniqueID
	/* 0016 */ 0x00, 0x00, 0x00, 0x00, // no next IFD
	/* 001A */ 's', 'o', 'm', 'e', 't', 'h', 'i', 'n', 'g', ' ', 'l', 'o', 'n', 'g', 'e', 'r', '!', 0x00,
	/* 002C */ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // padding
	/* 0044 */ 0x00, 0x01, // Foreign IFD: 1 tag
	/* 0046 */ 0xEE, 0xEE, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0xEE, 0xEE, 0xEE, 0xEE, // fake tag
	/* 0052 */ 0x00, 0x00, 0x00, 0x00, // No next IFD
	/* 0056 */ 0x00, 0x02, // IFD0: 2 tags
	/* 0058 */ 0x01, 0x0E, 0x00, 0x02, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x74, // ImageDescription
	/* 0064 */ 0x87, 0x69, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, // ExifOffset
	/* 0070 */ 0x00, 0x00, 0x00, 0x44, // point to foreign IFD
	/* 0074 */ 'd', 'e', 's', 'c', 0x00,
}
