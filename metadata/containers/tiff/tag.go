package tiff

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers"
	"golang.org/x/text/encoding/charmap"
)

// Tag is a single tag in an IFD.
type Tag struct {
	ifd   *IFD
	tag   uint16
	ttype uint16
	doff  uint32 // offset of data relative to start of TIFF; zero if data is embedded in IFD entry
	// Exactly one of the following four fields is non-nil:
	data      []byte
	reader    metadata.Reader
	container containers.Container
	toIFD     *IFD
}

// Read reads a single tag from the reader.  On entry, the file pointer should
// be positioned at the IFD entry for the tag; on exit, it will be positioned
// just past the IFD entry.
func (tag *Tag) Read(r metadata.Reader) (err error) {
	var (
		buf   [12]byte
		size  uint32
		count uint32
	)
	if _, err = io.ReadFull(r, buf[0:12]); err != nil {
		return fmt.Errorf("can't read IFD entry: %s", err)
	}
	tag.tag = tag.ifd.t.enc.Uint16(buf[0:2])
	tag.ttype = tag.ifd.t.enc.Uint16(buf[2:4])
	count = tag.ifd.t.enc.Uint32(buf[4:8])
	tag.doff = tag.ifd.t.enc.Uint32(buf[8:12])
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
		return fmt.Errorf("unknown IFD tag type %d", tag.ttype)
	}
	size *= count
	if size <= 4 {
		tag.data = make([]byte, size)
		copy(tag.data, buf[8:])
		tag.doff = 0
		return nil
	}
	if size%2 == 1 {
		tag.ifd.t.ranges.add(tag.doff, tag.doff+size+1)
	} else {
		tag.ifd.t.ranges.add(tag.doff, tag.doff+size)
	}
	if size > 1024 {
		tag.reader = io.NewSectionReader(r, int64(tag.doff), int64(size))
		return nil
	}
	tag.data = make([]byte, size)
	if n, _ := r.ReadAt(tag.data, int64(tag.doff)); n < int(size) {
		return fmt.Errorf("can't read tag %x data at offset %x", tag.tag, tag.ttype)
	}
	return nil
}

// SetContainer sets the container handler for the tag.  If the tag has a
// container handler, we defer to it for Dirty, Size, and Write.
func (tag *Tag) SetContainer(c containers.Container) {
	tag.container = c
	tag.data = nil
	tag.reader = nil
	tag.toIFD = nil
}

// AsBytes decodes the byte array in the tag.  It returns an error if the tag
// has the wrong type.
func (tag *Tag) AsBytes() (by []byte, err error) {
	if tag.ttype != 1 {
		return nil, errors.New("tag type is not BYTE")
	}
	if tag.data == nil {
		by = make([]byte, tag.reader.Size())
		if _, err = tag.reader.ReadAt(by, 0); err != nil {
			return nil, err
		}
	} else {
		by = make([]byte, len(tag.data))
		copy(by, tag.data)
	}
	return by, nil
}

// AsBytesReader decodes the byte array in the tag.  It returns an error if the
// tag has the wrong type.
func (tag *Tag) AsBytesReader() (metadata.Reader, error) {
	if tag.ttype != 1 {
		return nil, errors.New("tag type is not BYTE")
	}
	if tag.data == nil {
		return tag.reader, nil
	}
	return bytes.NewReader(tag.data), nil
}

// SetBytes sets the tag value to the specified byte array.
func (tag *Tag) SetBytes(by []byte) {
	if old, err := tag.AsBytes(); err == nil && bytes.Equal(old, by) {
		return
	}
	tag.ttype = 1 // BYTE
	tag.data = make([]byte, len(by))
	copy(tag.data, by)
	tag.container = nil
	tag.reader = nil
	tag.toIFD = nil
	tag.ifd.dirty = true
}

// AsUnknown decodes the byte array in the "unknown type" tag.  It returns an
// error if the tag has the wrong type.
func (tag *Tag) AsUnknown() (by []byte, err error) {
	if tag.ttype != 7 {
		return nil, errors.New("tag type is not UNKNOWN")
	}
	if tag.data == nil {
		by = make([]byte, tag.reader.Size())
		if _, err = tag.reader.ReadAt(by, 0); err != nil {
			return nil, err
		}
	} else {
		by = make([]byte, len(tag.data))
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

// SetUnknown sets the tag value to the specified byte array, and type to
// UNKNOWN.
func (tag *Tag) SetUnknown(by []byte) {
	if old, err := tag.AsUnknown(); err == nil && bytes.Equal(old, by) {
		return
	}
	tag.ttype = 7 // UNKNOWN
	tag.data = make([]byte, len(by))
	copy(tag.data, by)
	tag.container = nil
	tag.reader = nil
	tag.toIFD = nil
	tag.ifd.dirty = true
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

// SetString sets the tag value to the specified string.
func (tag *Tag) SetString(s string) {
	if old, err := tag.AsString(); err == nil && old == s {
		return
	}
	var encoded = make([]byte, len(s)+1)
	tag.ttype = 2 // ASCII
	copy(encoded, s)
	tag.data = encoded
	tag.container = nil
	tag.reader = nil
	tag.toIFD = nil
	tag.ifd.dirty = true
}

// AsRationals decodes the rationals in the tag.  It returns them as an even
// number of uint32s, alternating numerators and denominators.  It returns an
// error if the tag has the wrong type.
func (tag *Tag) AsRationals() (rat []uint32, err error) {
	if tag.ttype != 5 {
		return nil, errors.New("tag type is not RATIONAL")
	}
	rat = make([]uint32, len(tag.data)/4)
	for i := 0; i < len(tag.data)/8; i++ {
		rat[2*i] = tag.ifd.t.enc.Uint32(tag.data[8*i:])
		rat[2*i+1] = tag.ifd.t.enc.Uint32(tag.data[8*i+4:])
	}
	return rat, nil
}

// SetRationals sets the tag value to the specified list of rational values
// (passed as an even number of uint32s, alternating numerator and denominator).
func (tag *Tag) SetRationals(rat []uint32) {
	if len(rat)%2 != 0 {
		panic("SetRationals with odd slice length")
	}
	if len(rat) == 0 {
		panic("SetRationals with empty slice")
	}
	if old, err := tag.AsRationals(); err == nil {
		if len(old) == len(rat) {
			var mismatch = false
			for i := range old {
				if old[i] != rat[i] {
					mismatch = true
					break
				}
			}
			if !mismatch {
				return
			}
		}
	}
	var encoded = make([]byte, 4*len(rat))
	for i := range rat {
		tag.ifd.t.enc.PutUint32(encoded[4*i:], rat[i])
	}
	tag.ttype = 5 // RATIONAL
	tag.data = encoded
	tag.container = nil
	tag.reader = nil
	tag.toIFD = nil
	tag.ifd.dirty = true
}

// AsIFD returns the IFD that the tag points to.
func (tag *Tag) AsIFD() (ifd *IFD, err error) {
	if tag.toIFD != nil {
		return tag.toIFD, nil
	}
	if tag.ttype != 4 {
		return nil, errors.New("tag type is not LONG")
	}
	if len(tag.data) != 4 {
		return nil, errors.New("tag count is not 1")
	}
	if _, err = tag.ifd.t.r.Seek(int64(tag.ifd.t.enc.Uint32(tag.data)), io.SeekStart); err != nil {
		panic(err)
	}
	ifd = &IFD{t: tag.ifd.t, back: tag.ifd}
	if err = ifd.Read(ifd.t.r); err != nil {
		return nil, err
	}
	tag.toIFD = ifd
	tag.reader = nil
	tag.container = nil
	tag.data = nil
	return ifd, nil
}

// AddIFD sets the tag value to be a pointer to a new IFD, and returns the new
// (empty) IFD.  If the tag value already is an IFD, the existing IFD is
// returned.
func (tag *Tag) AddIFD() (ifd *IFD, err error) {
	if tag.ttype != 0 {
		if ifd, err = tag.AsIFD(); err != nil || ifd != nil {
			return ifd, err
		}
	}
	ifd = &IFD{t: tag.ifd.t, back: tag.ifd, dirty: true}
	tag.ttype = 4 // LONG
	tag.toIFD = ifd
	tag.data = nil
	tag.reader = nil
	tag.container = nil
	tag.ifd.dirty = true
	return ifd, nil
}

// write writes the IFD entry for a single tag.  offset is the next available
// offset for writing data (i.e., the pointer that should be put in the IFD
// entry if a pointer is needed).  newoff is the resulting offset, adjusted for
// the space needed to write this tag's data.  count is the number of bytes
// actually written (always 12, unless there's an error).
func (tag *Tag) write(w io.Writer, offset uint32) (newoff uint32, count int, err error) {
	var (
		buf    [12]byte
		size   uint32
		tcount uint32
	)
	size, tcount = tag.size()
	tag.ifd.t.enc.PutUint16(buf[0:2], tag.tag)
	tag.ifd.t.enc.PutUint16(buf[2:4], tag.ttype)
	tag.ifd.t.enc.PutUint32(buf[4:8], tcount)
	switch {
	case tag.toIFD != nil:
		tag.ifd.t.enc.PutUint32(buf[8:12], tag.toIFD.offset)
	case size <= 4:
		if tag.container != nil {
			panic("not expecting container of size <= 4")
		}
		copy(buf[8:12], tag.data)
	default:
		tag.ifd.t.enc.PutUint32(buf[8:12], offset)
		offset += size
	}
	count, err = w.Write(buf[0:12])
	return offset, count, err
}

// writeData writes the data associated with the tag.  It's a no-op if the data
// were included in the IFD entry.
func (tag *Tag) writeData(w io.Writer) (count int, err error) {
	size, _ := tag.size()
	if size <= 4 { // data was embedded in IFD entry, nothing to write
		return 0, nil
	}
	if tag.container != nil {
		if count, err = tag.container.Write(w); err != nil {
			return count, err
		}
		if count < int(size) {
			// size got rounded up to the nearest multiple of the
			// underlying data type.  We need to emit zeros to
			// match that.
			var (
				zeros [8]byte
				n     int
			)
			n, err = w.Write(zeros[:int(size)-count])
			count += n
		}
		return count, err
	}
	if tag.reader != nil {
		n64, err := io.Copy(w, tag.reader)
		return int(n64), err
	}
	return w.Write(tag.data)
}

// size computes the size of the encoded tag's data.  It does not include any
// alignment padding at the end, but it does ensure that the result is a
// multiple of the underlying data type of the tag.  It returns both the size
// in bytes, and the count of instances of the underlying data type.
func (tag *Tag) size() (size, count uint32) {
	var unit uint32

	switch tag.ttype {
	case 1, 2, 7:
		unit = 1
	case 3, 8:
		unit = 2
	case 4, 9:
		unit = 4
	case 5, 10:
		unit = 8
	default:
		panic("unknown IFD tag type")
	}
	switch {
	case tag.toIFD != nil:
		size = 4
	case tag.container != nil:
		size = uint32(tag.container.Size())
	case tag.reader != nil:
		size = uint32(tag.reader.Size())
	default:
		size = uint32(len(tag.data))
	}
	if size%unit != 0 {
		size += unit - (size % unit)
	}
	return size, size / unit
}
