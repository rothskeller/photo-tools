package tiff

import (
	"errors"
	"io"
	"sort"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers"
)

// IFD is a single Image File Directory contained within the TIFF-like block.
type IFD struct {
	t       *TIFF
	back    *IFD
	tags    []*Tag
	offset  uint32 // of start of IFD, relative to start of TIFF
	size    int64  // of entire rendered IFD including data
	next    uint32 // offset of next IFD relative to start of TIFF
	nextIFD *IFD
	dirty   bool
}

var _ containers.Container = (*IFD)(nil) // verify interface compliance

// Read reads and parses an IFD from the supplied Reader.  It expects that
// seek offsets in the supplied Reader are relative to the beginning of the TIFF
// container containing the IFD, and the current seek offset of the Reader is
// the beginning of the IFD.
func (ifd *IFD) Read(r metadata.Reader) (err error) {
	var (
		buf       [12]byte
		count     uint16
		noNextIFD bool
		dirsize   uint32
	)
	if start, err := r.Seek(0, io.SeekCurrent); err != nil {
		panic(err)
	} else {
		ifd.offset = uint32(start)
	}
	if n, _ := r.Read(buf[0:2]); n < 2 {
		return errors.New("can't read IFD size")
	}
	count = ifd.t.enc.Uint16(buf[0:2])
	ifd.tags = make([]*Tag, count)
	dirsize = 12*uint32(count) + 6
	for i := int64(0); i < int64(count); i++ {
		ifd.tags[i] = &Tag{ifd: ifd}
		if err = ifd.tags[i].Read(ifd.t.r); err != nil {
			return err
		}
		if ifd.tags[i].doff == ifd.offset+12*uint32(count)+2 {
			// Some JPEGs are malformed in that they don't have a
			// next IFD pointer at the end of the Exif IFD
			// directory.  This would appear to be one, since we
			// found tag data where that pointer should be.
			noNextIFD = true
			dirsize -= 4
		}
	}
	ifd.t.ranges.add(ifd.offset, ifd.offset+dirsize)
	sort.Slice(ifd.tags, func(i, j int) bool {
		return ifd.tags[i].tag < ifd.tags[j].tag
	})
	if !noNextIFD {
		if n, _ := ifd.t.r.Read(buf[0:4]); n < 4 {
			return errors.New("can't read IFD next pointer")
		}
		ifd.next = ifd.t.enc.Uint32(buf[0:4])
	}
	return nil
}

// Dirty returns whether the contents of the IFD have been changed.
func (ifd *IFD) Dirty() bool {
	for _, tag := range ifd.tags {
		if tag.toIFD != nil && tag.toIFD.Dirty() {
			return true
		}
		if tag.container != nil && tag.container.Dirty() {
			return true
		}
	}
	return ifd.dirty
}

// Size returns the rendered size of the IFD, in bytes.
func (ifd *IFD) Size() int64 {
	ifd.size = 6 + 12*int64(len(ifd.tags))
	for _, tag := range ifd.tags {
		if tsz, _ := tag.size(); tsz > 4 {
			if ifd.size%2 == 1 {
				ifd.size++
			}
			ifd.size += int64(tsz)
		}
	}
	return ifd.size
}

// Write writes the rendered IFD to the specified writer.  Note that ifd.offset
// must be set to the offset of the start of the IFD relative to the start of
// the enclosing TIFF container.  This method does not write subsidiary IFDs.
func (ifd *IFD) Write(w io.Writer) (count int, err error) {
	var (
		buf    [4]byte
		offset uint32
		n      int
	)
	offset = ifd.offset + 12*uint32(len(ifd.tags)) + 6
	ifd.t.enc.PutUint16(buf[0:2], uint16(len(ifd.tags)))
	n, err = w.Write(buf[0:2])
	count += n
	if err != nil {
		return count, err
	}
	for _, tag := range ifd.tags {
		offset, n, err = tag.write(w, offset)
		count += n
		if err != nil {
			return count, err
		}
		if offset%2 == 1 {
			offset++
		}
	}
	if ifd.nextIFD != nil {
		ifd.t.enc.PutUint32(buf[0:4], ifd.nextIFD.offset)
	} else {
		ifd.t.enc.PutUint32(buf[0:4], ifd.next)
	}
	n, err = w.Write(buf[0:4])
	count += n
	if err != nil {
		return count, err
	}
	for _, tag := range ifd.tags {
		if tsz, _ := tag.size(); count%2 == 1 && tsz > 4 {
			n, err = w.Write([]byte{0})
			count += n
			if err != nil {
				return count, err
			}
		}
		n, err = tag.writeData(w)
		count += n
		if err != nil {
			return count, err
		}
	}
	if ifd.size != 0 && int(ifd.size) != count {
		panic("actual size different from predicted size")
	}
	return count, nil
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

// AddTag adds a tag to the IFD.  It does not set a value or type for the tag;
// that needs to be done separately.  If the tag already exists, the existing
// tag is returned.
func (ifd *IFD) AddTag(id uint16) *Tag {
	idx := sort.Search(len(ifd.tags), func(i int) bool {
		return ifd.tags[i].tag >= id
	})
	if idx < len(ifd.tags) && ifd.tags[idx].tag == id {
		return ifd.tags[idx]
	}
	ifd.dirty = true
	ifd.tags = append(ifd.tags, nil)
	copy(ifd.tags[idx+1:], ifd.tags[idx:])
	ifd.tags[idx] = &Tag{ifd: ifd, tag: id}
	return ifd.tags[idx]
}

// DeleteTag deletes a tag from the IFD.
func (ifd *IFD) DeleteTag(id uint16) {
	if ifd == nil {
		return
	}
	idx := sort.Search(len(ifd.tags), func(i int) bool {
		return ifd.tags[i].tag >= id
	})
	if idx >= len(ifd.tags) || ifd.tags[idx].tag != id {
		return
	}
	ifd.dirty = true
	ifd.tags = append(ifd.tags[:idx], ifd.tags[idx+1:]...)
}

// NextIFD returns the next IFD in the IFD chain, or nil if there is none.
func (ifd *IFD) NextIFD() (next *IFD, err error) {
	if ifd.nextIFD != nil {
		return ifd.nextIFD, nil
	}
	if ifd.next == 0 {
		return nil, nil
	}
	if _, err = ifd.t.r.Seek(int64(ifd.next), io.SeekStart); err != nil {
		panic(err)
	}
	ifd.nextIFD = &IFD{t: ifd.t, back: ifd}
	if err = ifd.nextIFD.Read(ifd.t.r); err != nil {
		return nil, err
	}
	return ifd.nextIFD, nil
}

// AddNextIFD adds an IFD as the "next" IFD after the receiver.  If the receiver
// already has a "next" IFD, the existing one is returned.
func (ifd *IFD) AddNextIFD() (next *IFD, err error) {
	if next, err = ifd.NextIFD(); err == nil && next != nil {
		return next, nil
	} else if err != nil {
		return nil, err
	}
	ifd.nextIFD = &IFD{t: ifd.t, back: ifd, dirty: true}
	ifd.dirty = true
	return ifd.nextIFD, nil
}
