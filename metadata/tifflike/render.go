package tifflike

import (
	"encoding/binary"
	"io"
	"sort"
)

var zeros []byte

// Render writes the TIFF-like block to the provided writer.  Portions of the
// input reader that belong to unread or unmodified IFDs are left unchanged and
// unmoved.  The modified IFDs are written in the spaces they previously used,
// or at the end of the block.  Unused space is filled with zeros.
func (t *TIFFLike) Render(w io.Writer) (err error) {
	var (
		inputend int64
		offset   uint32
		dirty    []*IFD
		buf      [8]byte
	)
	// We need to know where the end of the input file is.
	if inputend, err = t.r.Seek(0, io.SeekEnd); err == nil {
		offset = uint32(inputend)
	} else {
		return err
	}
	// If we have a reusable range up against the end of the input file,
	// treat it as being past the end of the file instead.  This avoids
	// unnecessary fragmentation.
	offset = t.ranges.removeTrailer(offset)
	// The file might end on an odd byte boundary.
	if offset%2 == 1 {
		offset++
	}
	// Find all of the dirty IFDs and sort them by size, descending.
	dirty = findDirtyIFDs(nil, t.ifd0)
	for _, ifd := range dirty {
		ifd.computeSize()
	}
	sort.Slice(dirty, func(i, j int) bool { return dirty[i].size > dirty[j].size })
	// Assign space to each of them, and then re-sort them by offset.
	for _, ifd := range dirty {
		if ifd.resize {
			if ifd.offset = t.ranges.consume(ifd.size); ifd.offset == 0 {
				ifd.offset, offset = offset, offset+ifd.size
			}
		}
	}
	sort.Slice(dirty, func(i, j int) bool { return dirty[i].offset < dirty[j].offset })
	// Write the TIFF header.
	if t.enc == binary.BigEndian {
		copy(buf[:], tiffHeaderBE)
	} else {
		copy(buf[:], tiffHeaderLE)
	}
	t.enc.PutUint32(buf[4:], t.ifd0.offset)
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	offset = 8
	// Write the rest of the file, pulling from the dirty IFDs, the input
	// file, and a buffer of zeros, as appropriate.
	for {
		// Which is the next dirty IFD or unused range to write?
		var nextDirty, nextZero, nextAny uint32
		if len(dirty) != 0 {
			nextDirty = dirty[0].offset
			nextAny = nextDirty
		}
		if len(t.ranges.r) != 0 {
			nextZero = t.ranges.r[0]
			if nextAny == 0 || nextAny > nextZero {
				nextAny = nextZero
			}
		}
		if offset == nextDirty {
			// We're at the start of a dirty IFD.  Render it.
			if err := dirty[0].render(w); err != nil {
				return err
			}
			offset += dirty[0].size
			dirty = dirty[1:]
		} else if offset == nextZero {
			// We're at the start of an unused range.  Write zeros.
			if err := writeZeros(w, t.ranges.r[1]-t.ranges.r[0]); err != nil {
				return err
			}
			offset += t.ranges.r[1] - t.ranges.r[0]
			t.ranges.r = t.ranges.r[2:]
		} else if nextAny != 0 && nextAny <= uint32(inputend) {
			// We need to copy bytes from the input file up to the
			// next dirty IFD or unused range.
			if err := copyBytes(w, t.r, offset, nextAny); err != nil {
				return err
			}
			offset = nextAny
		} else if nextAny != 0 && nextAny == uint32(inputend)+1 {
			// We need to copy the rest of the bytes from the input
			// file, plus we need to add a zero byte to align the
			// file offset on an even-numbered boundary.
			if err := copyBytes(w, t.r, offset, uint32(inputend)); err != nil {
				return err
			}
			if err := writeZeros(w, 1); err != nil {
				return err
			}
			offset = nextAny
		} else if nextAny != 0 {
			panic("input file ended early")
		} else if offset < uint32(inputend) {
			// There are no dirty IFDs or unused ranges left, so we
			// need to copy the rest of the input file.
			if err := copyBytes(w, t.r, offset, uint32(inputend)); err != nil {
				return err
			}
			offset = uint32(inputend)
		} else {
			break
		}
	}
	return nil
}
func findDirtyIFDs(dirty []*IFD, ifd *IFD) []*IFD {
	if ifd.dirty && len(ifd.tags) != 0 {
		dirty = append(dirty, ifd)
	}
	for _, tag := range ifd.tags {
		if tag.toIFD != nil {
			dirty = findDirtyIFDs(dirty, tag.toIFD)
		}
	}
	if ifd.nextIFD != nil {
		dirty = findDirtyIFDs(dirty, ifd.nextIFD)
	}
	return dirty
}

// computeSize computes the size of the encoded IFD.
func (ifd *IFD) computeSize() {
	ifd.size = 6
	for _, tag := range ifd.tags {
		ifd.size += 12
		if dsize := tag.size(); dsize > 4 {
			ifd.size += dsize
		}
	}
}

// size computes the size of the encoded tag's data.
func (tag *Tag) size() (size uint32) {
	switch tag.ttype {
	case 1, 2, 7:
		size = tag.count
	case 3, 8:
		size = 2 * tag.count
	case 4, 9:
		size = 4 * tag.count
	case 5, 10:
		size = 8 * tag.count
	default:
		panic("unknown IFD tag type")
	}
	if size > 4 && size%2 == 1 {
		size++
	}
	return size
}

func (ifd *IFD) render(w io.Writer) (err error) {
	var (
		buf    [4]byte
		offset uint32
	)
	offset = ifd.offset + 12*uint32(len(ifd.tags)) + 6
	ifd.t.enc.PutUint16(buf[0:2], uint16(len(ifd.tags)))
	if _, err := w.Write(buf[0:2]); err != nil {
		return err
	}
	for _, tag := range ifd.tags {
		if offset, err = tag.render(w, offset); err != nil {
			return err
		}
	}
	if ifd.nextIFD != nil {
		ifd.t.enc.PutUint32(buf[0:4], ifd.nextIFD.offset)
	} else {
		ifd.t.enc.PutUint32(buf[0:4], ifd.next)
	}
	if _, err := w.Write(buf[0:4]); err != nil {
		return err
	}
	for _, tag := range ifd.tags {
		if err := tag.renderData(w); err != nil {
			return err
		}
	}
	return nil
}
func (tag *Tag) render(w io.Writer, offset uint32) (newoff uint32, err error) {
	var buf [12]byte

	tag.ifd.t.enc.PutUint16(buf[0:2], tag.tag)
	tag.ifd.t.enc.PutUint16(buf[2:4], tag.ttype)
	tag.ifd.t.enc.PutUint32(buf[4:8], tag.count)
	if tag.toIFD != nil {
		tag.ifd.t.enc.PutUint32(buf[8:12], tag.toIFD.offset)
	} else if dsize := tag.size(); dsize <= 4 {
		copy(buf[8:12], tag.data)
	} else {
		tag.ifd.t.enc.PutUint32(buf[8:12], offset)
		offset += dsize
	}
	if _, err := w.Write(buf[0:12]); err != nil {
		return 0, err
	}
	return offset, nil
}
func (tag *Tag) renderData(w io.Writer) (err error) {
	dsize := tag.size()
	if dsize <= 4 {
		return nil
	}
	if _, err := w.Write(tag.data); err != nil {
		return err
	}
	if uint32(len(tag.data)) == dsize-1 {
		if err := writeZeros(w, 1); err != nil {
			return err
		}
	}
	return nil
}

func writeZeros(w io.Writer, size uint32) (err error) {
	// NextIFD thing to write is a block of zeros covering an
	// unused range.
	if zeros == nil {
		zeros = make([]byte, 32768)
	}
	for size >= 32768 {
		if _, err := w.Write(zeros); err != nil {
			return err
		}
		size -= 32768
	}
	if size != 0 {
		if _, err := w.Write(zeros[:size]); err != nil {
			return err
		}
	}
	return nil
}

func copyBytes(w io.Writer, r tiffReader, from, to uint32) (err error) {
	if _, err := r.Seek(int64(from), io.SeekStart); err != nil {
		return err
	}
	var size = to - from
	if _, err := io.CopyN(w, r, int64(size)); err != nil {
		return err
	}
	return nil
}
