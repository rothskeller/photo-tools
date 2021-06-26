package exif

import (
	"sort"
)

// Dirty returns whether the EXIF data have changed and need to be saved.
func (p *EXIF) Dirty() bool {
	if p == nil || len(p.Problems) != 0 {
		return false
	}
	p.setArtist()
	p.setDateTime()
	p.setImageDescription()
	if p.exifIFD != nil {
		p.setDateTimeDigitized()
		p.setDateTimeOriginal()
		p.setUserComment()
	}
	if p.gpsIFD != nil {
		p.setGPSCoords()
	}
	return p.ifd0.dirty || (p.exifIFD != nil && p.exifIFD.dirty) || (p.gpsIFD != nil && p.gpsIFD.dirty)
}

// Render renders and returns the encoded EXIF block, applying any changes made
// to the metadata fields of the EXIF structure.  maxSize is the maximum allowed
// size of the block.
func (p *EXIF) Render(max uint64) (out []byte) {
	if len(p.Problems) != 0 {
		panic("EXIF Render with parse problems")
	}
	if !p.Dirty() {
		return p.buf
	}
	out = p.render()
	if len(out) > int(max) {
		panic("EXIF block doesn't fit within maximum size")
	}
	return out
}

func (p *EXIF) render() (out []byte) {
	type ifddata struct {
		ifd     *ifdt
		data    []byte
		offset  uint32
		pointer uint32
	}
	var (
		newIFDs     []*ifddata
		newIFD0     *ifddata
		newEXIF     *ifddata
		newGPS      *ifddata
		exifIFDPtr  uint32
		gpsIFDPtr   uint32
		availRanges [][]uint32
	)
	// Ranges available for use in the rendering include the current IFD0,
	// the current ExifIFD and GPSIFD if there are any, and the unbounded
	// range starting at the end of the current buffer.
	availRanges = [][]uint32{
		{p.ifd0.offset, p.ifd0.offset + p.ifd0.size},
		{uint32(len(p.buf)), 0xFFFFFFFF},
	}
	if p.exifIFD != nil && p.exifIFD.offset != 0 {
		availRanges = append(availRanges, []uint32{p.exifIFD.offset, p.exifIFD.offset + p.exifIFD.size})
	}
	if p.gpsIFD != nil && p.gpsIFD.offset != 0 {
		availRanges = append(availRanges, []uint32{p.gpsIFD.offset, p.gpsIFD.offset + p.gpsIFD.size})
	}
	availRanges = coalesceAndSortAvailRanges(availRanges)
	// Render the IFDs with unknown offsets, just to get their sizes.
	newIFD0 = &ifddata{ifd: p.ifd0}
	newIFD0.data, exifIFDPtr, gpsIFDPtr = p.renderIFD(p.ifd0, 0)
	newIFDs = append(newIFDs, newIFD0)
	if p.exifIFD != nil {
		newEXIF = &ifddata{ifd: p.exifIFD}
		newEXIF.data, _, _ = p.renderIFD(p.exifIFD, 0)
		newIFDs = append(newIFDs, newEXIF)
	}
	if p.gpsIFD != nil {
		newGPS = &ifddata{ifd: p.gpsIFD}
		newGPS.data, _, _ = p.renderIFD(p.gpsIFD, 0)
		newIFDs = append(newIFDs, newGPS)
	}
	// Sort them smallest first.
	sort.Slice(newIFDs, func(i, j int) bool {
		return len(newIFDs[i].data) < len(newIFDs[j].data)
	})
	// Assign space to each of them, using the smallest available range.
	for _, ifd := range newIFDs {
		ifd.offset, availRanges = consumeRange(availRanges, uint32(len(ifd.data)))
	}
	// Re-render the IFDs with the correct offsets.
	for _, ifd := range newIFDs {
		if ifd == newIFD0 {
			ifd.data, exifIFDPtr, gpsIFDPtr = p.renderIFD(ifd.ifd, ifd.offset)
		} else {
			ifd.data, _, _ = p.renderIFD(ifd.ifd, ifd.offset)
		}
	}
	// Set the offsets in IFD0.
	if newEXIF != nil {
		p.enc.PutUint32(newIFD0.data[exifIFDPtr:], newEXIF.offset)
	}
	if newGPS != nil {
		p.enc.PutUint32(newIFD0.data[gpsIFDPtr:], newGPS.offset)
	}
	// Copy the input buffer, with all of the data we aren't touching, to
	// the output.
	out = make([]byte, availRanges[len(availRanges)-1][0])
	copy(out, p.buf)
	// Empty out all of the leftover available ranges (except for the
	// unbounded one at the end).  Not strictly necessary, but perhaps
	// helpful to someone examining the file later.
	if len(availRanges) > 1 {
		for _, r := range availRanges[:len(availRanges)-1] {
			for i := r[0]; i < r[1]; i++ {
				out[i] = 0
			}
		}
	}
	// Copy the newly rendered IFDs into the output buffer.
	for _, ifd := range newIFDs {
		copy(out[ifd.offset:], ifd.data)
	}
	p.enc.PutUint32(out[4:], newIFD0.offset)
	return out
}

func (p *EXIF) renderIFD(ifd *ifdt, base uint32) (out []byte, exifIFDPtr, gpsIFDPtr uint32) {
	out = make([]byte, 6+12*len(ifd.tags))
	p.enc.PutUint16(out, uint16(len(ifd.tags)))
	p.enc.PutUint32(out[2+12*len(ifd.tags):], ifd.next)
	var iop uint32 = 2
	sort.Slice(ifd.tags, func(i, j int) bool {
		return ifd.tags[i].tag < ifd.tags[j].tag
	})
	for _, tag := range ifd.tags {
		p.enc.PutUint16(out[iop:], tag.tag)
		p.enc.PutUint16(out[iop+2:], tag.ttype)
		p.enc.PutUint32(out[iop+4:], tag.count)
		if len(tag.data) <= 4 {
			copy(out[iop+8:], tag.data)
		} else {
			p.enc.PutUint32(out[iop+8:], base+uint32(len(out)))
			out = append(out, tag.data...)
		}
		if tag.tag == tagExifIFDOffset {
			exifIFDPtr = iop + 8
		}
		if tag.tag == tagGPSIFDOffset {
			gpsIFDPtr = iop + 8
		}
		iop += 12
	}
	return out, exifIFDPtr, exifIFDPtr
}

// consumeRange finds the smallest range that will contain the desired size, and
// returns the offset of it.  It revises the available ranges list to reflect
// the change.
func consumeRange(availRanges [][]uint32, size uint32) (offset uint32, newAvail [][]uint32) {
	for i, r := range availRanges {
		if r[1]-r[0] >= size {
			offset = r[0]
			availRanges[i][0] += size
			break
		}
	}
	newAvail = coalesceAndSortAvailRanges(availRanges)
	return offset, newAvail
}

// coalesceAndSortAvailRanges merges adjacent ranges, and then sorts the list of
// ranges from smallest to largest.
func coalesceAndSortAvailRanges(ranges [][]uint32) [][]uint32 {
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] < ranges[j][0]
	})
	j := 0
	for i := 1; i < len(ranges); i++ {
		if ranges[i][0] == ranges[i][1] {
			// skip
		} else if ranges[i][0] == ranges[j][1] {
			ranges[j][1] = ranges[i][1]
		} else {
			j++
			ranges[j] = ranges[i]
		}
	}
	ranges = ranges[:j+1]
	sort.Slice(ranges, func(i, j int) bool {
		si := ranges[i][1] - ranges[i][0]
		sj := ranges[j][1] - ranges[j][0]
		if si != sj {
			return si < sj
		}
		return ranges[i][0] < ranges[j][0]
	})
	return ranges
}

func (p *EXIF) addTag(ifd *ifdt, tag *tagt) {
	ifd.tags = append(ifd.tags, tag)
	ifd.dirty = true
}

// Most of the time, deleteTag should set the dirty flag.  For the rare cases
// when it shouldn't, there's a separate entrypoint, to keep from having to pass
// an extra argument in the common case.
func (p *EXIF) deleteTag(ifd *ifdt, tag uint16) {
	if ifd == nil {
		return
	}
	j := 0
	for _, t := range ifd.tags {
		if t.tag != tag {
			ifd.tags[j] = t
			j++
		}
	}
	if j < len(ifd.tags) {
		ifd.tags = ifd.tags[:j]
		ifd.dirty = true
	}
}

func (p *EXIF) addEXIFIFD() {
	p.exifIFD = new(ifdt)
	p.addTag(p.ifd0, &tagt{
		tag:   tagExifIFDOffset,
		ttype: 4, // LONG
		count: 1,
		data:  []byte{0, 0, 0, 0}, // filled in later
	})
}

func (p *EXIF) addGPSIFD() {
	p.gpsIFD = new(ifdt)
	p.addTag(p.gpsIFD, &tagt{
		tag:   0, // GPS Version
		ttype: 1, // BYTE
		count: 4,
		data:  []byte{2, 3, 0, 0},
	})
	p.addTag(p.ifd0, &tagt{
		tag:   tagGPSIFDOffset,
		ttype: 4, // LONG
		count: 1,
		data:  []byte{0, 0, 0, 0}, // filled in later
	})
}

func (p *EXIF) setASCIITag(ifd *ifdt, tnum uint16, val string) {
	tag := ifd.findTag(tnum)
	if tag == nil {
		tag = &tagt{tag: tnum, ttype: 2, count: 0, data: nil}
		p.addTag(ifd, tag)
	}
	tag.data = make([]byte, len(val)+1)
	copy(tag.data, val)
	tag.count = uint32(len(tag.data))
	ifd.dirty = true
}

func (p *EXIF) setRationalTag(ifd *ifdt, tnum uint16, val []uint32) {
	tag := ifd.findTag(tnum)
	if tag == nil {
		tag = &tagt{tag: tnum, ttype: 2, count: 0, data: nil}
		p.addTag(ifd, tag)
	}
	tag.data = make([]byte, len(val)*4)
	for i := range val {
		p.enc.PutUint32(tag.data[i*4:], val[i])
	}
	tag.count = uint32(len(val))
	ifd.dirty = true
}

func (p *EXIF) setByteTag(ifd *ifdt, tnum uint16, val byte) {
	tag := ifd.findTag(tnum)
	if tag == nil {
		tag = &tagt{tag: tnum, ttype: 1, count: 0, data: nil}
		p.addTag(ifd, tag)
	}
	tag.data = []byte{val}
	tag.count = 1
	ifd.dirty = true
}
