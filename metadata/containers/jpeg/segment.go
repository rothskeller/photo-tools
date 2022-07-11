// Package jpeg handles marshaling and unmarshaling of JPEG file segments.
package jpeg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers"
)

// A segmentGroup is one or more segments in a JPEG file with the same marker
// and namespace.
type segmentGroup struct {
	marker    byte
	namespace []byte
	reader    metadata.Reader
	container containers.Container
	size      int64
	csize     int64
}

var _ containers.Container = (*segmentGroup)(nil) // verify interface compliance

// Read reads a single segment from the JPEG file, and returns a segment group
// containing only one segment.
func (seg *segmentGroup) Read(r metadata.Reader) (err error) {
	var (
		buf    [64]byte
		size   int64
		offset int64
		count  int
	)
	if _, err = io.ReadFull(r, buf[0:2]); err != nil {
		return err
	}
	if buf[0] != 0xFF {
		return errors.New("invalid segment marker")
	}
	for buf[1] == 0xFF && tell(r) > 2 { // skip padding
		if _, err = r.Read(buf[1:2]); err != nil {
			return err
		}
	}
	seg.marker = buf[1]
	switch seg.marker {
	case 0x00, 0xFF:
		return errors.New("invalid segment marker")
	case 0x01, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9:
		return nil
	case 0xDA:
		offset = tell(r)
		seg.reader = io.NewSectionReader(r, offset, r.Size()-offset) // rest of the file
		return nil
	}
	if _, err = io.ReadFull(r, buf[0:2]); err != nil {
		return err
	}
	size = int64(binary.BigEndian.Uint16(buf[0:2])) - 2
	offset = tell(r)
	count, err = r.ReadAt(buf[:], offset)
	if err != nil && err != io.EOF {
		return err
	}
	switch {
	case seg.marker == markerJFIF && count >= len(nsJFIF) && bytes.Equal(buf[:len(nsJFIF)], nsJFIF):
		seg.namespace = nsJFIF
	case seg.marker == markerJFXX && count > len(nsJFXX) && bytes.Equal(buf[:len(nsJFXX)], nsJFXX):
		seg.namespace = nsJFXX
	case seg.marker == markerEXIF && count > len(nsEXIF) && bytes.Equal(buf[:len(nsEXIF)], nsEXIF):
		seg.namespace = nsEXIF
	case seg.marker == markerXMP && count > len(nsXMP) && bytes.Equal(buf[:len(nsXMP)], nsXMP):
		seg.namespace = nsXMP
	// Disabling the recognition of XMPext namespaces until we can better
	// implement handling multi-segment XMPext namespaces.
	// case seg.marker == markerXMPext && count > len(nsXMPext) && bytes.Equal(buf[:len(nsXMPext)], nsXMPext):
	// 	seg.namespace = nsXMPext
	case seg.marker == markerPSIR && count > len(nsPSIR) && bytes.Equal(buf[:len(nsPSIR)], nsPSIR):
		seg.namespace = nsPSIR
	}
	if _, err = r.Seek(size, io.SeekCurrent); err != nil {
		panic(err)
	}
	if seg.namespace != nil {
		offset += int64(len(seg.namespace))
		size -= int64(len(seg.namespace))
	}
	seg.reader = io.NewSectionReader(r, offset, size)
	return nil
}

// Empty returns whether the container is empty (and should therefore be omitted
// from the written file, along with whatever tag in the parent container points
// to it).
func (seg *segmentGroup) Empty() bool {
	if seg == nil {
		return true
	}
	if seg.container != nil {
		return seg.container.Empty()
	}
	return false
}

// Dirty returns whether the segment group has been changed.
func (seg *segmentGroup) Dirty() bool {
	if seg == nil {
		return false
	}
	if seg.container != nil {
		return seg.container.Dirty()
	}
	return false
}

// Layout computes the rendered layout of the container, i.e. prepares for a
// call to Write, and returns what the rendered size of the container will be.
func (seg *segmentGroup) Layout() int64 {
	if seg == nil {
		return 0
	}
	if seg.container == nil && seg.reader == nil {
		seg.size = 2
		return seg.size
	}
	if seg.marker == 0xDA {
		seg.size = 2 + seg.reader.Size()
		return seg.size
	}
	if seg.namespace == nil {
		seg.size = 4 + seg.reader.Size()
		return seg.size
	}
	if seg.container != nil {
		seg.size = seg.container.Layout()
		seg.csize = seg.size
	} else {
		seg.size = seg.reader.Size()
	}
	var chunk = 0xFFFF - int64(len(seg.namespace)) - 2
	var numChunks = (seg.size + chunk - 1) / chunk
	seg.size += numChunks * (int64(len(seg.namespace)) + 4)
	return seg.size
}

// Write writes a segment group to the specified writer.
func (seg *segmentGroup) Write(w io.Writer) (count int, err error) {
	var (
		buf   [64]byte
		chunk int64
		size  int64
		r     io.Reader
		n     int
		n64   int64
	)
	if seg == nil {
		return 0, nil
	}
	defer func() {
		if err == nil && seg.size != 0 && int(seg.size) != count {
			println(seg.marker, seg.namespace, seg.size, count)
			// TODO: problem with writing the same photo twice?
			panic("actual size different from predicted size")
		}
	}()
	// Make sure the reader (if any) is rewound.
	if seg.reader != nil {
		seg.reader.Seek(0, os.SEEK_SET)
	}
	// Handle the simple cases first.  Markers without data:
	if seg.container == nil && seg.reader == nil {
		buf[0] = 0xFF
		buf[1] = seg.marker
		return w.Write(buf[0:2])
	}
	// The marker for "remainder of the file":
	if seg.marker == 0xDA {
		buf[0] = 0xFF
		buf[1] = seg.marker
		n, err = w.Write(buf[0:2])
		count += n
		if err != nil {
			return count, err
		}
		n64, err = io.Copy(w, seg.reader)
		count += int(n64)
		return count, err
	}
	// A segment without a recognized namespace:
	if seg.namespace == nil {
		buf[0] = 0xFF
		buf[1] = seg.marker
		binary.BigEndian.PutUint16(buf[2:4], uint16(seg.reader.Size())+2)
		n, err = w.Write(buf[0:4])
		count += n
		if err != nil {
			return count, err
		}
		n64, err = io.Copy(w, seg.reader)
		count += int(n64)
		return count, err
	}
	// We're now looking at a segment group with a recognized namespace.
	// First, determine whether it can fit into one physical segment.
	chunk = 0xFFFF - int64(len(seg.namespace)) - 2
	if seg.container != nil {
		size = seg.csize
	} else {
		size = seg.reader.Size()
	}
	if size <= chunk {
		// Yes, it can, so write out that segment and we're done.
		buf[0] = 0xFF
		buf[1] = seg.marker
		binary.BigEndian.PutUint16(buf[2:4], uint16(len(seg.namespace))+uint16(size)+2)
		copy(buf[4:], seg.namespace)
		n, err = w.Write(buf[0 : len(seg.namespace)+4])
		count += n
		if err != nil {
			return count, err
		}
		if seg.container != nil {
			n, err = seg.container.Write(w)
		} else {
			n64, err = io.Copy(w, seg.reader)
			n = int(n64)
		}
		count += n
		return count, err
	}
	// We need to write multiple physical segments.  We'll need the source
	// in the form of an io.Reader for this, so if we need to render a
	// container, do so into a pipe in a goroutine.
	if seg.container != nil {
		var pr, pw = io.Pipe()
		go seg.container.Write(pw)
		r = pr
	} else {
		r = seg.reader
	}
	// Write out each physical segment.
	for size > 0 {
		if size < chunk {
			chunk = size
		}
		buf[0] = 0xFF
		buf[1] = seg.marker
		binary.BigEndian.PutUint16(buf[2:4], uint16(len(seg.namespace))+uint16(chunk)+2)
		copy(buf[4:], seg.namespace)
		n, err = w.Write(buf[0 : len(seg.namespace)+4])
		count += n
		if err != nil {
			return count, err
		}
		n64, err = io.CopyN(w, r, chunk)
		count += int(n64)
		if err != nil {
			return count, err
		}
		size -= chunk
	}
	return count, nil
}

// merge merges two segment groups (the "other" of which must have only one
// segment) into a single segment group.  They must have the same marker and
// namespace.  It returns the merged segment group (the receiver).  The receiver
// may be nil, in which case the argument is returned.  If skip is nonzero and
// receiver is nonnil, the first skip bytes of other are excluded from the
// merged segment group.
func (seg *segmentGroup) merge(other *segmentGroup, skip int64) *segmentGroup {
	if seg == nil {
		return other
	}
	if seg.marker != other.marker || !bytes.Equal(seg.namespace, other.namespace) {
		panic("can't merge unlike segments")
	}
	if mr, ok := seg.reader.(*multireader); ok {
		if skip != 0 {
			mr.rdrs = append(mr.rdrs, io.NewSectionReader(other.reader, skip, other.reader.Size()-skip))
		} else {
			mr.rdrs = append(mr.rdrs, other.reader.(*io.SectionReader))
		}
	} else {
		seg.reader = &multireader{rdrs: []*io.SectionReader{
			seg.reader.(*io.SectionReader),
			other.reader.(*io.SectionReader),
		}}
	}
	return seg
}

func tell(r metadata.Reader) (offset int64) {
	var err error

	if offset, err = r.Seek(0, io.SeekCurrent); err != nil {
		panic(err)
	}
	return offset
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
