package tiff

import (
	"bytes"
	"sort"
)

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
	ifd.markResize()
	ifd.tags = append(ifd.tags, nil)
	copy(ifd.tags[idx+1:], ifd.tags[idx:])
	ifd.tags[idx] = &Tag{ifd: ifd, tag: id, ttype: 1}
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
	ifd.markResize()
	ifd.tags = append(ifd.tags[:idx], ifd.tags[idx+1:]...)
}

// SetBytes sets the tag value to the specified byte array.
func (tag *Tag) SetBytes(by []byte) {
	if old, err := tag.AsBytes(); err == nil && bytes.Equal(old, by) {
		return
	}
	newsize := len(by)
	if newsize%2 == 1 {
		newsize++
	}
	if newsize != int(tag.size()) {
		tag.ifd.markResize()
	} else {
		tag.ifd.dirty = true
	}
	tag.data = make([]byte, len(by))
	copy(tag.data, by)
	tag.count = uint32(len(by))
	tag.ttype = 1 // BYTE
}

// SetUnknown sets the tag value to the specified byte array, and type to
// UNKNOWN.
func (tag *Tag) SetUnknown(by []byte) {
	if old, err := tag.AsUnknown(); err == nil && bytes.Equal(old, by) {
		return
	}
	newsize := len(by)
	if newsize%2 == 1 {
		newsize++
	}
	if newsize != int(tag.size()) {
		tag.ifd.markResize()
	} else {
		tag.ifd.dirty = true
	}
	tag.data = make([]byte, len(by))
	copy(tag.data, by)
	tag.count = uint32(len(by))
	tag.ttype = 7 // UNKNOWN
}

// SetString sets the tag value to the specified string.
func (tag *Tag) SetString(s string) {
	if old, err := tag.AsString(); err == nil && old == s {
		return
	}
	newsize := len(s) + 1
	if newsize%2 == 1 {
		newsize++
	}
	if newsize != int(tag.size()) {
		tag.ifd.markResize()
	} else {
		tag.ifd.dirty = true
	}
	var encoded = make([]byte, len(s)+1)
	copy(encoded, s)
	tag.data = encoded
	tag.count = uint32(len(encoded))
	tag.ttype = 2 // ASCII
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
	if int(tag.size()) != 4*len(rat) {
		tag.ifd.markResize()
	} else {
		tag.ifd.dirty = true
	}
	var encoded = make([]byte, 4*len(rat))
	for i := range rat {
		tag.ifd.t.enc.PutUint32(encoded[4*i:], rat[i])
	}
	tag.data = encoded
	tag.count = uint32(len(rat) / 2)
	tag.ttype = 5 // RATIONAL
}

// AddNextIFD adds an IFD as the "next" IFD after the receiver.  If the receiver
// already has a "next" IFD, the existing one is returned.
func (ifd *IFD) AddNextIFD() (next *IFD, err error) {
	if next, err = ifd.NextIFD(); err == nil && next != nil {
		return next, nil
	} else if err != nil {
		return nil, err
	}
	ifd.nextIFD = &IFD{t: ifd.t, back: ifd, dirty: true, resize: true}
	ifd.dirty = true
	return ifd.nextIFD, nil
}

// AddIFD sets the tag value to be a pointer to a new IFD, and returns the new
// (empty) IFD.  If the tag value already is an IFD, the existing IFD is
// returned.
func (tag *Tag) AddIFD() *IFD {
	if ifd, err := tag.AsIFD(); err == nil && ifd != nil {
		return ifd
	}
	if tag.size() != 4 {
		tag.ifd.markResize()
	} else {
		tag.ifd.dirty = true
	}
	tag.toIFD = &IFD{t: tag.ifd.t, back: tag.ifd, dirty: true, resize: true}
	tag.count = 1
	tag.ttype = 4 // LONG
	tag.data = make([]byte, 4)
	return tag.toIFD
}

// markResize marks the IFD as having been changed, and therefore needing to be
// regenerated on Render.
func (ifd *IFD) markResize() {
	if ifd.resize {
		return
	}
	ifd.resize = true
	// Before we make any changes, we need to record all of the ranges used
	// by this IFD pre-change, so that we know they're available for reuse
	// when rendering.
	ifd.t.ranges.add(ifd.offset, ifd.offset+ifd.size)
	for _, tag := range ifd.tags {
		if tag.doff != 0 {
			ifd.t.ranges.add(tag.doff, tag.doff+tag.size())
		}
	}
	for ; ifd != nil; ifd = ifd.back {
		ifd.dirty = true
	}
}

// Dirty returns whether there have been any changes to the TIFF-like block.
func (t *TIFF) Dirty() bool {
	dirty := findDirtyIFDs(nil, t.ifd0)
	return len(dirty) != 0
}
