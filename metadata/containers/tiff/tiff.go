// Package tiff handles metadata blocks that use TIFF-style encoding.
package tiff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers"
)

// TIFF is a handler for a TIFF-like file (or portion of file).
type TIFF struct {
	r      metadata.Reader
	enc    binary.ByteOrder
	ifd0   *IFD
	ranges rangelist
	ifds   []*IFD
	end    uint32
}

var _ containers.Container = (*TIFF)(nil) // verify interface compliance

var tiffHeaderLE = []byte{0x49, 0x49, 0x2A, 0x00}
var tiffHeaderBE = []byte{0x4D, 0x4D, 0x00, 0x2A}

var zeros []byte

// Read reads and parses the container structure from the supplied Reader.  The
// reader will continue to be used after Read returns, and must remain open and
// usable as long as the Container is in scope.
func (t *TIFF) Read(r metadata.Reader) (err error) {
	var (
		buf [8]byte
	)
	t.r = r
	if _, err = r.Read(buf[:8]); err != nil {
		return errors.New("TIFF: invalid TIFF header")
	}
	if bytes.Equal(buf[:4], tiffHeaderBE) {
		t.enc = binary.BigEndian
	} else if bytes.Equal(buf[:4], tiffHeaderLE) {
		t.enc = binary.LittleEndian
	} else {
		return errors.New("TIFF: invalid TIFF header")
	}
	if _, err = r.Seek(int64(t.enc.Uint32(buf[4:8])), io.SeekStart); err != nil {
		panic(err)
	}
	t.ifd0 = &IFD{t: t}
	if err = t.ifd0.Read(r); err != nil {
		return fmt.Errorf("TIFF: IFD0: %s", err)
	}
	return nil
}

// Empty returns whether the container is empty (and should therefore be omitted
// from the written file, along with whatever tag in the parent container points
// to it).
func (t *TIFF) Empty() bool { return t.ifd0.Empty() }

// Dirty returns whether there have been any changes to the TIFF-like block.
func (t *TIFF) Dirty() bool {
	for _, ifd := range findAllIFDs(nil, t.ifd0) {
		if ifd.Dirty() {
			return true
		}
	}
	return false
}

// Layout computes the rendered layout of the container, i.e. prepares for a
// call to Write, and returns what the rendered size of the container will be.
func (t *TIFF) Layout() int64 {
	// Set the end pointer to the end of the file.  If there's a consumable
	// range that ends at the end of the file, drop it and set the end
	// pointer to the start of that range.  Round the end pointer up to an
	// even offset.
	t.end = uint32(t.r.Size())
	if len(t.ranges.r) != 0 && t.ranges.r[len(t.ranges.r)-1] == t.end {
		t.end = t.ranges.r[len(t.ranges.r)-2]
		t.ranges.r = t.ranges.r[:len(t.ranges.r)-2]
	}
	// Get a list of all of the IFDs to be rendered, in decreasing order by
	// size.
	removeEmptyIFDs(t.ifd0)
	t.ifds = findAllIFDs(nil, t.ifd0)
	for i := range t.ifds {
		t.ifds[i].size = t.ifds[i].Layout()
	}
	sort.Slice(t.ifds, func(i, j int) bool {
		return t.ifds[i].size > t.ifds[j].size
	})
	// Allocate space to each of the IFDs.
	for _, ifd := range t.ifds {
		if ifd.offset = t.ranges.consume(uint32(ifd.size)); ifd.offset == 0 {
			if t.end%2 == 1 {
				t.end++
			}
			ifd.offset = t.end
			t.end += uint32(ifd.size)
		}
	}
	// Resort the IFDs by file offset.
	sort.Slice(t.ifds, func(i, j int) bool {
		return t.ifds[i].offset < t.ifds[j].offset
	})
	return int64(t.end)
}

// Write writes the rendered container to the specified writer.
func (t *TIFF) Write(w io.Writer) (count int, err error) {
	var (
		n   int
		buf [8]byte
	)
	// Write the TIFF header.
	if t.enc == binary.BigEndian {
		copy(buf[:], tiffHeaderBE)
	} else {
		copy(buf[:], tiffHeaderLE)
	}
	t.enc.PutUint32(buf[4:], t.ifd0.offset)
	n, err = w.Write(buf[:])
	count += n
	if err != nil {
		return count, err
	}
	// Write the rest of the file, pulling from the IFDs, the input file,
	// and a buffer of zeros, as appropriate.
	for {
		// Which is the next IFD or unused range to write?
		var nextIFD, nextZero, nextAny uint32
		if len(t.ifds) != 0 {
			nextIFD = t.ifds[0].offset
			nextAny = nextIFD
		}
		if len(t.ranges.r) != 0 {
			nextZero = t.ranges.r[0]
			if nextAny == 0 || nextAny > nextZero {
				nextAny = nextZero
			}
		}
		if count == int(nextIFD) {
			// We're at the start of an IFD.  Render it.
			n, err = t.ifds[0].Write(w)
			count += n
			if err != nil {
				return count, err
			}
			t.ifds = t.ifds[1:]
		} else if count == int(nextZero) {
			// We're at the start of an unused range.  Write zeros.
			n, err = writeZeros(w, t.ranges.r[1]-t.ranges.r[0])
			count += n
			if err != nil {
				return count, err
			}
			t.ranges.r = t.ranges.r[2:]
		} else if nextAny != 0 {
			// We need to copy bytes from the input file up to the
			// next IFD or unused range.
			n, err = copyBytes(w, t.r, uint32(count), nextAny)
			count += n
			if err == io.EOF && count == int(nextAny-1) {
				// We could run out of input file with one byte
				// too few, because the input file ends on an
				// odd offset and the next IFD needs to be
				// aligned on an even offset.  If so, write one
				// zero.
				n, err = writeZeros(w, 1)
				count += n
			}
			if err != nil {
				return count, err
			}
		} else if count < int(t.end) {
			// There are no IFDs or unused ranges left, so we need
			// to copy the rest of the input file.
			n, err = copyBytes(w, t.r, uint32(count), t.end)
			count += n
			if err != nil {
				return count, err
			}
		} else {
			break
		}
	}
	if count != int(t.end) {
		panic("actual size different from predicted size")
	}
	return count, nil
}

// IFD0 returns the first IFD in the TIFF-like block.
func (t *TIFF) IFD0() *IFD {
	return t.ifd0
}

// Encoding returns the byte order of the TIFF-like block.  In rare cases
// (especially for tags of UNKNOWN type), calling code may need this in order to
// correctly interpret tag data.
func (t *TIFF) Encoding() binary.ByteOrder { return t.enc }

func removeEmptyIFDs(ifd *IFD) {
	for i := 0; i < len(ifd.tags); {
		if ifd.tags[i].toIFD != nil && len(ifd.tags[i].toIFD.tags) == 0 && ifd.tags[i].toIFD.next == 0 {
			ifd.DeleteTag(ifd.tags[i].tag)
			continue
		}
		if ifd.tags[i].toIFD != nil {
			removeEmptyIFDs(ifd.tags[i].toIFD)
		}
		i++
	}
	if ifd.nextIFD != nil && len(ifd.nextIFD.tags) == 0 && ifd.nextIFD.next == 0 {
		ifd.nextIFD = nil
		ifd.next = 0
	} else if ifd.nextIFD != nil {
		removeEmptyIFDs(ifd.nextIFD)
	}
}

func findAllIFDs(list []*IFD, ifd *IFD) []*IFD {
	list = append(list, ifd)
	for _, tag := range ifd.tags {
		if tag.toIFD != nil {
			list = findAllIFDs(list, tag.toIFD)
		}
	}
	if ifd.nextIFD != nil {
		list = findAllIFDs(list, ifd.nextIFD)
	}
	return list
}

func writeZeros(w io.Writer, size uint32) (count int, err error) {
	var n int

	if zeros == nil {
		zeros = make([]byte, 32768)
	}
	for size >= 32768 {
		n, err = w.Write(zeros)
		count += n
		if err != nil {
			return count, err
		}
		size -= 32768
	}
	if size != 0 {
		n, err = w.Write(zeros[:size])
		count += n
	}
	return count, err
}

func copyBytes(w io.Writer, r metadata.Reader, from, to uint32) (count int, err error) {
	if _, err := r.Seek(int64(from), io.SeekStart); err != nil {
		return 0, err
	}
	var size = to - from
	n, err := io.CopyN(w, r, int64(size))
	return int(n), err
}
