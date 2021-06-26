// Package exif handles EXIF metadata blocks.
package exif

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

const (
	tagExifIFDOffset uint16 = 0x8769
	tagGPSIFDOffset  uint16 = 0x8825
)

var tiffHeaderLE = []byte{0x49, 0x49, 0x2A, 0x00}
var tiffHeaderBE = []byte{0x4D, 0x4D, 0x00, 0x2A}

// Parse parses an EXIF block and returns the parse results.  offset is the
// starting offset of the EXIF block within the image file; it is used in
// problem messages.
func Parse(buf []byte, offset uint32) (exif *EXIF) {
	exif = &EXIF{offset: offset, buf: buf}
	if len(buf) < 8 || !bytes.HasPrefix(buf, tiffHeaderBE) && !bytes.HasPrefix(buf, tiffHeaderLE) {
		exif.Problems = append(exif.Problems, "[0] EXIF has invalid TIFF header")
		return exif
	}
	if bytes.HasPrefix(buf, tiffHeaderBE) {
		exif.enc = binary.BigEndian
	} else {
		exif.enc = binary.LittleEndian
	}
	exif.ranges = append(exif.ranges, []uint32{0, 8})
	if exif.ifd0 = exif.parseIFD(exif.enc.Uint32(buf[4:])); exif.ifd0 == nil {
		return exif
	}
	if eifdt := exif.ifd0.findTag(tagExifIFDOffset); eifdt != nil {
		exif.exifIFD = exif.parseIFD(exif.enc.Uint32(eifdt.data))
	}
	if gifdt := exif.ifd0.findTag(tagGPSIFDOffset); gifdt != nil {
		exif.gpsIFD = exif.parseIFD(exif.enc.Uint32(gifdt.data))
	}
	if !exif.verifyNoOverlaps() {
		return exif
	}
	exif.ranges = nil // free the space
	exif.getArtist()
	exif.getDateTime()
	exif.getImageDescription()
	if exif.exifIFD != nil {
		exif.getUserComment()
		exif.getDateTimeDigitized()
		exif.getDateTimeOriginal()
	}
	if exif.gpsIFD != nil {
		exif.getGPSCoords()
	}
	return exif
}

func (p *EXIF) parseIFD(offset uint32) (ifd *ifdt) {
	var dirsize uint32

	ifd = &ifdt{offset: offset}
	if int(ifd.offset+6) > len(p.buf) {
		p.log(ifd.offset, "IFD directory reaches outside EXIF block")
		return nil
	}
	count := p.enc.Uint16(p.buf[ifd.offset:])
	if count >= 0x8000 {
		p.log(ifd.offset, "invalid IFD tag count")
		return nil
	}
	ifd.size = 12*uint32(count) + 6
	dirsize = ifd.size
	if int(ifd.offset)+int(ifd.size) > len(p.buf) {
		p.log(ifd.offset, "IFD directory reaches outside EXIF block")
		return nil
	}
	nextoff := ifd.offset + 12*uint32(count) + 2
	ifd.next = p.enc.Uint32(p.buf[nextoff:])
	for i := 0; i < int(count); i++ {
		iop := ifd.offset + 2 + 12*uint32(i)
		tag, end := p.parseTag(iop)
		if tag == nil {
			return nil
		}
		if tag.doff == nextoff {
			// Some JPEGs are malformed in that they don't have a
			// next IFD pointer at the end of the Exif IFD
			// directory.  This would appear to be one, since we
			// found tag data where that pointer should be.
			ifd.next = 0
			dirsize -= 4
		}
		ifd.tags = append(ifd.tags, tag)
		if size := end - ifd.offset; size > ifd.size {
			ifd.size = size
		}
	}
	p.ranges = append(p.ranges, []uint32{ifd.offset, ifd.offset + dirsize})
	return ifd
}

func (p *EXIF) parseTag(offset uint32) (tag *tagt, end uint32) {
	tag = new(tagt)
	tag.offset = offset
	tag.tag = p.enc.Uint16(p.buf[offset:])
	tag.ttype = p.enc.Uint16(p.buf[offset+2:])
	tag.count = p.enc.Uint32(p.buf[offset+4:])
	tag.doff = p.enc.Uint32(p.buf[offset+8:])
	end = offset + 12
	var size uint32
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
		p.log(offset, "unknown IFD tag type %d", tag.ttype)
		return nil, 0
	}
	size *= tag.count
	tag.data = make([]byte, size)
	if size <= 4 {
		copy(tag.data, p.buf[offset+8:])
		tag.doff = 0
	} else {
		end = tag.doff + size
		p.ranges = append(p.ranges, []uint32{tag.doff, end})
		if int(end) > len(p.buf) {
			p.log(offset, "IFD tag reaches outside EXIF block")
			return nil, 0
		}
		copy(tag.data, p.buf[tag.doff:])
	}
	return tag, end
}

func (p *EXIF) verifyNoOverlaps() bool {
	sort.Slice(p.ranges, func(i, j int) bool {
		return p.ranges[i][0] < p.ranges[j][0]
	})
	for i := 1; i < len(p.ranges); i++ {
		if p.ranges[i][0] < p.ranges[i-1][1] {
			p.log(p.ranges[i][0], "structure error: overlapping ranges")
			p.ifd0 = nil
			p.exifIFD = nil
			p.gpsIFD = nil
			return false
		}
	}
	return true
}

func (ifd *ifdt) findTag(tag uint16) *tagt {
	for _, t := range ifd.tags {
		if t.tag == tag {
			return t
		}
	}
	return nil
}

func (p *EXIF) asciiAt(tag *tagt, label string) string {
	if tag.ttype != 2 {
		p.log(tag.offset, "%s is not ASCII", label)
		return ""
	}
	by := tag.data
	// The spec requires a trailing NUL, but sometimes it's not there, and
	// sometimes there are more than one.
	by = bytes.TrimRightFunc(by, func(r rune) bool { return r == 0 })
	s := string(by)
	if utf8.ValidString(s) {
		return strings.TrimSpace(s)
	}
	if s2, err := charmap.ISO8859_1.NewDecoder().String(s); err == nil {
		// p.log(tag.offset, "%s is in unknown character set, assuming ISO-8859-1", label)
		// Don't log this; I don't want to block file access based on this assumption.
		return strings.TrimSpace(s2)
	}
	p.log(tag.offset, "%s is in unknown character set, removing non-ASCII characters", label)
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if r < 32 || r > 126 {
			return -1
		}
		return r
	}, s))
}

func (p *EXIF) log(offset uint32, f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	p.Problems = append(p.Problems, fmt.Sprintf("EXIF[%x] %s", p.offset+offset, s))
}
