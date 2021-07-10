// Package jpeg handles marshaling and unmarshaling of JPEG file segments.
package jpeg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	markerJFIF   byte = 0xE0
	markerJFXX   byte = 0xE0
	markerEXIF   byte = 0xE1
	markerXMP    byte = 0xE1
	markerXMPext byte = 0xE1
	markerPSIR   byte = 0xED
)

var (
	nsJFIF   = []byte("JFIF\000")
	nsJFXX   = []byte("JFXX\000")
	nsEXIF   = []byte("Exif\000\000")
	nsXMP    = []byte("http://ns.adobe.com/xap/1.0/\000")
	nsXMPext = []byte("http://ns.adobe.com/xmp/extension/\000")
	nsPSIR   = []byte("Photoshop 3.0\000")
)

// A JPEG is a container of Segments.
type JPEG struct {
	start  segment
	jfif   []segment
	exif   []segment
	xmp    []segment
	xmpext []segment
	psir   []segment
	others []segment
	end    segment
}

// A segment is one segment in a JPEG file.
type segment struct {
	marker    byte
	namespace []byte
	reader    *io.SectionReader
}

// SegmentReader is the interface that must be honored by any Reader in a
// segment.
type SegmentReader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	Size() int64
}

// Read creates a new JPEG container handler, reading the specified reader.  It
// returns an error if the container is ill-formed or unreadable.
func Read(r SegmentReader) (jpeg *JPEG, err error) {
	var (
		offset int64
		seg    segment
	)
	jpeg = new(JPEG)
	if offset, seg, err = jpeg.readSegment(r, offset); err == io.EOF || seg.marker != 0xD8 {
		return nil, errors.New("JPEG: not a jpeg file")
	} else if err != nil {
		return nil, fmt.Errorf("JPEG: %s", err)
	}
	for seg.marker != 0xDA {
		if offset, seg, err = jpeg.readSegment(r, offset); err != nil {
			return nil, fmt.Errorf("JPEG: %s", err)
		}
	}
	return jpeg, nil
}
func (jpeg *JPEG) readSegment(r SegmentReader, offset int64) (newoff int64, seg segment, err error) {
	var (
		buf   [64]byte
		size  int64
		count int
		list  *[]segment
	)
	if _, err = r.ReadAt(buf[0:2], offset); err != nil {
		return 0, segment{}, err
	}
	if buf[0] != 0xFF {
		return 0, segment{}, errors.New("invalid segment marker")
	}
	offset += 2
	for buf[1] == 0xFF && offset > 2 { // skip padding
		if _, err = r.ReadAt(buf[1:2], offset); err != nil {
			return 0, segment{}, err
		}
		offset++
	}
	seg.marker = buf[1]
	switch seg.marker {
	case 0x00, 0xFF:
		return 0, segment{}, errors.New("invalid segment marker")
	case 0xD8:
		jpeg.start = seg
		return offset, seg, nil
	case 0x01, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD9:
		jpeg.others = append(jpeg.others, seg)
		return offset, seg, nil
	case 0xDA:
		offset -= 2
		seg.reader = io.NewSectionReader(r, offset, r.Size()-offset)
		jpeg.end = seg
		return offset, seg, nil
	}
	if _, err = r.ReadAt(buf[0:2], offset); err != nil {
		return 0, segment{}, err
	}
	offset += 2
	size = int64(binary.BigEndian.Uint16(buf[0:2])) - 2
	list = &jpeg.others
	if seg.marker == 0xE1 || seg.marker == 0xED {
		count, err = r.ReadAt(buf[:], offset)
		if err != nil && err != io.EOF {
			return 0, segment{}, err
		}
		switch {
		case seg.marker == markerJFIF && count >= len(nsJFIF) && bytes.Equal(buf[:len(nsJFIF)], nsJFIF):
			seg.namespace, list = nsJFIF, &jpeg.jfif
		case seg.marker == markerJFXX && count > len(nsJFXX) && bytes.Equal(buf[:len(nsJFXX)], nsJFXX):
			seg.namespace, list = nsJFXX, &jpeg.jfif
		case seg.marker == markerEXIF && count > len(nsEXIF) && bytes.Equal(buf[:len(nsEXIF)], nsEXIF):
			seg.namespace, list = nsEXIF, &jpeg.exif
		case seg.marker == markerXMP && count > len(nsXMP) && bytes.Equal(buf[:len(nsXMP)], nsXMP):
			seg.namespace, list = nsXMP, &jpeg.xmp
		case seg.marker == markerXMPext && count > len(nsXMPext) && bytes.Equal(buf[:len(nsXMPext)], nsXMPext):
			seg.namespace, list = nsXMPext, &jpeg.xmpext
		case seg.marker == markerPSIR && count > len(nsPSIR) && bytes.Equal(buf[:len(nsPSIR)], nsPSIR):
			seg.namespace, list = nsPSIR, &jpeg.psir
		}
		if seg.namespace != nil {
			offset += int64(len(seg.namespace))
			size -= int64(len(seg.namespace))
		}
	}
	seg.reader = io.NewSectionReader(r, offset, size)
	*list = append(*list, seg)
	offset += size
	return offset, seg, nil
}

// Render renders a JPEG file to the specified writer.
func (jpeg *JPEG) Render(w io.Writer) (err error) {
	if err = writeSegment(w, jpeg.start); err != nil {
		return err
	}
	for _, seg := range jpeg.jfif {
		if err = writeSegment(w, seg); err != nil {
			return err
		}
	}
	for _, seg := range jpeg.exif {
		if err = writeSegment(w, seg); err != nil {
			return err
		}
	}
	for _, seg := range jpeg.xmp {
		if err = writeSegment(w, seg); err != nil {
			return err
		}
	}
	for _, seg := range jpeg.xmpext {
		if err = writeSegment(w, seg); err != nil {
			return err
		}
	}
	for _, seg := range jpeg.psir {
		if err = writeSegment(w, seg); err != nil {
			return err
		}
	}
	for _, seg := range jpeg.others {
		if err = writeSegment(w, seg); err != nil {
			return err
		}
	}
	return writeSegment(w, jpeg.end)
}
func writeSegment(w io.Writer, seg segment) (err error) {
	var buf [64]byte

	switch seg.marker {
	case 0x01, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9:
		buf[0] = 0xFF
		buf[1] = seg.marker
		_, err = w.Write(buf[0:2])
	case 0xDA:
		buf[0] = 0xFF
		buf[1] = seg.marker
		if _, err = w.Write(buf[0:2]); err != nil {
			return err
		}
		_, err = io.Copy(w, seg.reader)
	case 0xE1, 0xED:
		buf[0] = 0xFF
		buf[1] = seg.marker
		binary.BigEndian.PutUint16(buf[2:4], uint16(len(seg.namespace))+uint16(seg.reader.Size())+2)
		copy(buf[4:], seg.namespace)
		if _, err = w.Write(buf[0 : len(seg.namespace)+4]); err != nil {
			return err
		}
		_, err = io.Copy(w, seg.reader)
	default:
		buf[0] = 0xFF
		buf[1] = seg.marker
		binary.BigEndian.PutUint16(buf[2:4], uint16(seg.reader.Size())+2)
		if _, err = w.Write(buf[0:4]); err != nil {
			return err
		}
		_, err = io.Copy(w, seg.reader)
	}
	return err
}

// EXIF returns the contents of the EXIF segment, if any.
func (jpeg *JPEG) EXIF() SegmentReader { return getSegment(jpeg.exif) }

// XMP returns the contents of the XMP segment, if any.
func (jpeg *JPEG) XMP() SegmentReader {
	var (
		buf  [1]byte
		size int64
	)
	if len(jpeg.xmp) == 0 {
		return nil
	}
	// Many XMP segments in my library have extraneous null bytes at the
	// end, which the RDF parser can't handle.  Detect and remove them.
	size = jpeg.xmp[0].reader.Size()
	for size > 0 {
		jpeg.xmp[0].reader.ReadAt(buf[:], size-1)
		if buf[0] == 0 {
			size--
		} else {
			break
		}
	}
	if size < jpeg.xmp[0].reader.Size() {
		return io.NewSectionReader(jpeg.xmp[0].reader, 0, size)
	}
	return jpeg.xmp[0].reader
	// Note also that we're ignoring the possibility of multiple XMP
	// segments here.  That isn't allowed.  If it did happen, the single
	// segment that we're returning would be incomplete XML and would fail
	// to parse.
}

// XMPext returns the contents of the XMP extension segment, if any.
func (jpeg *JPEG) XMPext() SegmentReader {
	// Each XMPext segment has a header at the front of it that we need to
	// skip past, so we can't use the plain getSegment logic here.
	switch len(jpeg.xmpext) {
	case 0:
		return nil
	case 1:
		return io.NewSectionReader(jpeg.xmpext[0].reader, 40, jpeg.xmpext[0].reader.Size()-40)
	}
	mr := new(multireader)
	for _, seg := range jpeg.xmpext {
		mr.rdrs = append(mr.rdrs, io.NewSectionReader(seg.reader, 40, seg.reader.Size()-40))
	}
	return mr
}

// PSIR returns the contents of the PSIR segment, if any.
func (jpeg *JPEG) PSIR() SegmentReader { return getSegment(jpeg.psir) }

func getSegment(segs []segment) SegmentReader {
	switch len(segs) {
	case 0:
		return nil
	case 1:
		return segs[0].reader
	}
	mr := new(multireader)
	for _, seg := range segs {
		mr.rdrs = append(mr.rdrs, seg.reader)
	}
	return mr
}

// SetEXIF sets the contents of the EXIF segment to those provided by the
// supplied reader.
func (jpeg *JPEG) SetEXIF(r SegmentReader) error {
	return setSegment(&jpeg.exif, markerEXIF, nsEXIF, r)
}

// SetXMP sets the contents of the XMP segment to those provided by the
// supplied reader.
func (jpeg *JPEG) SetXMP(r SegmentReader) error {
	if err := setSegment(&jpeg.xmp, markerXMP, nsXMP, r); err != nil {
		return err
	}
	if len(jpeg.xmp) != 1 {
		return errors.New("XMP exceeds JPEG segment size")
	}
	return nil
}

// SetPSIR sets the contents of the PSIR segment to those provided by the
// supplied reader.
func (jpeg *JPEG) SetPSIR(r SegmentReader) error {
	return setSegment(&jpeg.psir, markerPSIR, nsPSIR, r)
}

func setSegment(list *[]segment, marker byte, namespace []byte, reader SegmentReader) error {
	var (
		chunksize int64
		offset    int64
		size      int64
	)
	*list = (*list)[:0]
	chunksize = 0xFFFF - int64(len(namespace)) - 2
	offset = 0
	size = reader.Size()
	for size > chunksize {
		*list = append(*list, segment{marker, namespace, io.NewSectionReader(reader, offset, chunksize)})
		offset += chunksize
		size -= chunksize
	}
	*list = append(*list, segment{marker, namespace, io.NewSectionReader(reader, offset, size)})
	return nil
}

type multireader struct {
	rdrs   []*io.SectionReader
	offset int64
}

func (mr *multireader) Read(buf []byte) (n int, err error) {
	rnum, rdroff := mr.split(mr.offset)
	n, err = mr.rdrs[rnum].ReadAt(buf, rdroff)
	mr.offset += int64(n)
	if err == io.EOF && rnum < len(mr.rdrs)-1 {
		err = nil
	}
	return n, err
}

func (mr *multireader) ReadAt(buf []byte, offset int64) (n int, err error) {
	for len(buf) != 0 {
		var pn int

		rnum, rdroff := mr.split(offset)
		pn, err = mr.rdrs[rnum].ReadAt(buf, rdroff)
		n += pn
		if err == io.EOF && rnum < len(mr.rdrs)-1 {
			err = nil
		}
		if err != nil {
			return n, err
		}
		buf = buf[pn:]
		offset += int64(pn)
	}
	return n, nil
}

func (mr *multireader) Seek(offset int64, from int) (newoff int64, err error) {
	switch from {
	case io.SeekStart:
		mr.offset = offset
	case io.SeekCurrent:
		mr.offset += offset
	case io.SeekEnd:
		mr.offset = mr.Size() + offset
	}
	return mr.offset, nil
}

func (mr *multireader) Size() (size int64) {
	for _, rdr := range mr.rdrs {
		size += rdr.Size()
	}
	return size
}

func (mr *multireader) split(offset int64) (rnum int, rdroff int64) {
	for i, rdr := range mr.rdrs {
		if offset < rdr.Size() || i == len(mr.rdrs)-1 {
			return i, offset
		}
		offset -= rdr.Size()
	}
	panic("not reachable")
}
