package tiff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/rothskeller/photo-tools/metadata"
	"golang.org/x/text/encoding/charmap"
)

var tiffHeaderLE = []byte{0x49, 0x49, 0x2A, 0x00}
var tiffHeaderBE = []byte{0x4D, 0x4D, 0x00, 0x2A}

// Read returns a handler for a TIFF-like file (or portion of file) covered by
// the specified reader.  It returns an error if the data read from the reader
// does not conform to TIFF layout.
func Read(r metadata.Reader) (t *TIFF, err error) {
	var (
		buf [8]byte
	)
	if n, _ := r.ReadAt(buf[:8], 0); n < 4 {
		return nil, errors.New("TIFF: invalid TIFF header")
	}
	t = &TIFF{r: r}
	if bytes.Equal(buf[:4], tiffHeaderBE) {
		t.enc = binary.BigEndian
	} else if bytes.Equal(buf[:4], tiffHeaderLE) {
		t.enc = binary.LittleEndian
	} else {
		return nil, errors.New("TIFF: invalid TIFF header")
	}
	t.ifd0 = &IFD{t: t, offset: t.enc.Uint32(buf[4:])}
	if err := t.ifd0.read(); err != nil {
		return nil, fmt.Errorf("TIFF: IFD0: %s", err)
	}
	return t, nil
}

// IFD0 returns the first IFD in the TIFF-like block.
func (t *TIFF) IFD0() *IFD {
	return t.ifd0
}

// NextIFD returns the next IFD in the IFD chain, or nil if there is none.
func (ifd *IFD) NextIFD() (next *IFD, err error) {
	if ifd.nextIFD != nil {
		return ifd.nextIFD, nil
	}
	if ifd.next == 0 {
		return nil, nil
	}
	next = &IFD{t: ifd.t, back: ifd, offset: ifd.next}
	if err = next.read(); err != nil {
		return nil, err
	}
	ifd.nextIFD = next
	return next, nil
}

// read reads the IFD at the offset specified in the IFD structure.
func (ifd *IFD) read() (err error) {
	var (
		buf       [12]byte
		count     uint16
		noNextIFD bool
	)
	if n, _ := ifd.t.r.ReadAt(buf[:2], int64(ifd.offset)); n < 2 {
		return errors.New("can't read IFD size")
	}
	count = ifd.t.enc.Uint16(buf[:2])
	ifd.tags = make([]*Tag, count)
	ifd.size = 12*uint32(count) + 6
	for i := int64(0); i < int64(count); i++ {
		if n, _ := ifd.t.r.ReadAt(buf[:12], int64(ifd.offset)+12*i+2); n < 12 {
			return errors.New("can't read IFD entry")
		}
		if ifd.tags[i], err = ifd.readTag(buf[:]); err != nil {
			return err
		}
		if ifd.tags[i].doff == ifd.offset+12*uint32(count)+2 {
			// Some JPEGs are malformed in that they don't have a
			// next IFD pointer at the end of the Exif IFD
			// directory.  This would appear to be one, since we
			// found tag data where that pointer should be.
			noNextIFD = true
			ifd.size -= 4
		}
	}
	sort.Slice(ifd.tags, func(i, j int) bool {
		return ifd.tags[i].tag < ifd.tags[j].tag
	})
	if !noNextIFD {
		if n, _ := ifd.t.r.ReadAt(buf[:4], int64(ifd.offset)+12*int64(count)+2); n < 4 {
			return errors.New("can't read IFD next pointer")
		}
		ifd.next = ifd.t.enc.Uint32(buf[:4])
	}
	return nil
}

func (ifd *IFD) readTag(buf []byte) (tag *Tag, err error) {
	var size uint32

	tag = &Tag{
		ifd:   ifd,
		tag:   ifd.t.enc.Uint16(buf[:2]),
		ttype: ifd.t.enc.Uint16(buf[2:4]),
		count: ifd.t.enc.Uint32(buf[4:8]),
		doff:  ifd.t.enc.Uint32(buf[8:12]),
	}
	switch tag.ttype {
	case 1, 2, 7:
		size = 1
	case 3, 8:
		size = 2
	case 4, 9:
		size = 4
	case 5, 10:
		size = 8
	default:
		return nil, fmt.Errorf("unknown IFD tag type %d", tag.ttype)
	}
	size *= tag.count
	if size > 1024 {
		tag.reader = io.NewSectionReader(ifd.t.r, int64(tag.doff), int64(size))
		return tag, nil
	}
	tag.data = make([]byte, size)
	if size <= 4 {
		copy(tag.data, buf[8:])
		tag.doff = 0
	} else {
		if n, _ := ifd.t.r.ReadAt(tag.data, int64(tag.doff)); n < int(size) {
			return nil, fmt.Errorf("can't read tag %x data at offset %x", tag.tag, tag.ttype)
		}
	}
	return tag, nil
}

// Tag returns the specified tag from the IFD, or nil if it doesn't exist in the
// IFD.
func (ifd *IFD) Tag(id uint16) *Tag {
	idx := sort.Search(len(ifd.tags), func(i int) bool {
		return ifd.tags[i].tag >= id
	})
	if idx < len(ifd.tags) && ifd.tags[idx].tag == id {
		return ifd.tags[idx]
	}
	return nil
}

// AsBytes decodes the byte array in the tag.  It returns an error if the tag
// has the wrong type.
func (tag *Tag) AsBytes() (by []byte, err error) {
	if tag.ttype != 1 {
		return nil, errors.New("tag type is not BYTE")
	}
	by = make([]byte, tag.count)
	if tag.data == nil {
		if _, err = tag.reader.ReadAt(by, 0); err != nil {
			return nil, err
		}
	} else {
		copy(by, tag.data)
	}
	return by, nil
}

// AsBytesReader decodes the byte array in the tag.  It returns an error if the tag
// has the wrong type.
func (tag *Tag) AsBytesReader() (metadata.Reader, error) {
	if tag.ttype != 1 {
		return nil, errors.New("tag type is not BYTE")
	}
	if tag.data == nil {
		return tag.reader, nil
	}
	return bytes.NewReader(tag.data), nil
}

// AsUnknown decodes the byte array in the "unknown type" tag.  It returns an
// error if the tag has the wrong type.
func (tag *Tag) AsUnknown() (by []byte, err error) {
	if tag.ttype != 7 {
		return nil, errors.New("tag type is not UNKNOWN")
	}
	by = make([]byte, tag.count)
	if tag.data == nil {
		if _, err = tag.reader.ReadAt(by, 0); err != nil {
			return nil, err
		}
	} else {
		copy(by, tag.data)
	}
	return by, nil
}

// AsUnknownReader decodes the byte array in the tag.  It returns an error if
// the tag has the wrong type.
func (tag *Tag) AsUnknownReader() (metadata.Reader, error) {
	if tag.ttype != 7 {
		return nil, errors.New("tag type is not UNKNOWN")
	}
	if tag.data == nil {
		return tag.reader, nil
	}
	return bytes.NewReader(tag.data), nil
}

// AsLongReader decodes the byte array in the tag.  It returns an error if the
// tag has the wrong type.
func (tag *Tag) AsLongReader() (metadata.Reader, error) {
	if tag.ttype != 4 {
		return nil, errors.New("tag type is not LONG")
	}
	if tag.data == nil {
		return tag.reader, nil
	}
	return bytes.NewReader(tag.data), nil
}

// AsString decodes the string in the tag.  It returns an error if the tag has
// the wrong type or the character encoding can't be guessed.
func (tag *Tag) AsString() (string, error) {
	if tag.ttype != 2 {
		return "", errors.New("tag type is not ASCII")
	}
	by := tag.data
	// The spec requires a trailing NUL, but sometimes it's not there, and
	// sometimes there are more than one.
	by = bytes.TrimRightFunc(by, func(r rune) bool { return r == 0 })
	s := string(by)
	if utf8.ValidString(s) {
		return strings.TrimSpace(s), nil
	}
	if s2, err := charmap.ISO8859_1.NewDecoder().String(s); err == nil {
		return strings.TrimSpace(s2), nil
	}
	return "", errors.New("tag data uses unknown character set")
}

// AsIFD returns the IFD that the tag points to.
func (tag *Tag) AsIFD() (ifd *IFD, err error) {
	if tag.toIFD != nil {
		return tag.toIFD, nil
	}
	if tag.ttype != 4 {
		return nil, errors.New("tag type is not LONG")
	}
	if tag.count != 1 {
		return nil, errors.New("tag count is not 1")
	}
	ifd = &IFD{t: tag.ifd.t, back: tag.ifd, offset: tag.ifd.t.enc.Uint32(tag.data)}
	if err = ifd.read(); err != nil {
		return nil, err
	}
	tag.toIFD = ifd
	return ifd, nil
}

// AsRationals decodes the rationals in the tag.  It returns them as an even
// number of uint32s, alternating numerators and denominators.  It returns an
// error if the tag has the wrong type.
func (tag *Tag) AsRationals() (rat []uint32, err error) {
	if tag.ttype != 5 {
		return nil, errors.New("tag type is not RATIONAL")
	}
	rat = make([]uint32, 2*tag.count)
	for i := 0; i < int(tag.count); i++ {
		rat[2*i] = tag.ifd.t.enc.Uint32(tag.data[8*i:])
		rat[2*i+1] = tag.ifd.t.enc.Uint32(tag.data[8*i+4:])
	}
	return rat, nil
}

// NextTag returns the next tag in the IFD with the specified ID or a higher ID.
// It returns nil if there are no such tags.
func (ifd *IFD) NextTag(after uint16) *Tag {
	idx := sort.Search(len(ifd.tags), func(i int) bool {
		return ifd.tags[i].tag >= after
	})
	if idx < len(ifd.tags) {
		return ifd.tags[idx]
	}
	return nil
}

// Encoding returns the byte order of the TIFF-like block.  In rare cases
// (especially for tags of UNKNOWN type), calling code may need this in order to
// correctly interpret tag data.
func (t *TIFF) Encoding() binary.ByteOrder { return t.enc }
