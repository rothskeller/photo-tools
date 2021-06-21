// Package iptc handles IPTC metadata blocks.
package iptc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	iptcPSIRID          uint16 = 0x404
	iptcPSIRHash        uint16 = 0x425
	iptcTagMarker       byte   = 0x1C
	idCodedCharacterSet uint16 = 0x015A
)

// IPTC is a an IPTC parser and generator.
type IPTC struct {
	offset   uint32
	buf      []byte
	psir     []*psirt
	dsets    []*dsett
	dirty    bool
	problems []string
}
type psirt struct {
	offset uint32
	id     uint16
	name   string
	buf    []byte
}
type dsett struct {
	offset uint32
	id     uint16
	data   []byte
}

var psirResourceType = []byte("8BIM")
var utf8Escape1 = []byte{0x1B, 0x25, 0x47}
var utf8Escape2 = []byte{0x1B, 0x25, 0x2F, 0x49}

// Parse parses an IPTC block and returns the parse results.  offset is the
// offset of the IPTC block in the file, used for problem reporting.
func Parse(buf []byte, offset uint32) (iptc *IPTC) {
	iptc = &IPTC{offset: offset, buf: buf}
	if !iptc.splitPSIRs() {
		return iptc
	}
	for _, psir := range iptc.psir {
		if psir.id == iptcPSIRID {
			iptc.parseIPTC(psir)
		}
	}
	return iptc
}

// Problems returns the accumulated set of problems found.
func (p *IPTC) Problems() []string {
	if p == nil {
		return nil
	}
	return p.problems
}

// CanUpdate returns whether the IPTC block can be safely rewritten.
func (p *IPTC) CanUpdate() bool {
	return len(p.problems) == 0
}

// splitPSIRs splits the block up into PSIRs.  It returns false if a format
// problem was found.
func (p *IPTC) splitPSIRs() bool {
	var poff uint32
	var bufend = uint32(len(p.buf))
	for poff < bufend {
		var psir psirt
		if poff+8 > bufend {
			p.problems = append(p.problems, fmt.Sprintf("[%x] incomplete PSIR", p.offset+poff))
			return false
		}
		if !bytes.Equal(p.buf[poff:poff+4], psirResourceType) {
			p.problems = append(p.problems, fmt.Sprintf("[%x] invalid PSIR resource type", p.offset+poff))
			return false
		}
		psir.id = binary.BigEndian.Uint16(p.buf[poff+4:])
		resNameLen := p.buf[poff+6]
		resNameSize := uint32(resNameLen + 1)
		if resNameSize%2 == 1 {
			resNameSize++
		}
		if poff+6+resNameSize+4 > bufend {
			p.problems = append(p.problems, fmt.Sprintf("[%x] incomplete PSIR", p.offset+poff))
			return false
		}
		psir.name = string(p.buf[poff+7 : poff+7+uint32(resNameLen)])
		resLenOff := poff + 6 + resNameSize
		resLen := binary.BigEndian.Uint32(p.buf[resLenOff:])
		resSize := resLen
		if resSize%2 == 1 {
			resSize++
		}
		if resLenOff+resSize > bufend {
			p.problems = append(p.problems, fmt.Sprintf("[%x] incomplete PSIR", p.offset+poff))
			return false
		}
		psir.buf = p.buf[resLenOff+4 : resLenOff+4+resLen]
		psir.offset = resLenOff + 4
		p.psir = append(p.psir, &psir)
		poff = resLenOff + 4 + resSize
	}
	return true
}

// parseIPTC parses the PSIR that contains the IPTC data.
func (p *IPTC) parseIPTC(psir *psirt) {
	var (
		buf    = psir.buf
		offset = psir.offset
	)

	for len(buf) != 0 {
		var dset dsett

		if len(buf) < 5 {
			p.problems = append(p.problems, fmt.Sprintf("[%x] incomplete IPTC DataSet", p.offset+offset))
			return
		}
		if buf[0] != iptcTagMarker {
			p.problems = append(p.problems, fmt.Sprintf("[%x] invalid IPTC DataSet tag marker", p.offset+offset))
			return
		}
		dset.offset = offset
		dset.id = binary.BigEndian.Uint16(buf[1:3])
		dlen := uint32(binary.BigEndian.Uint16(buf[3:5]))
		buf = buf[5:]
		offset += 5
		if dlen&0x8000 != 0 {
			dlen &^= 0x8000
			if int(dlen) > len(buf) {
				p.problems = append(p.problems, fmt.Sprintf("[%x] incomplete IPTC DataSet", p.offset+offset))
				return
			}
			if dlen > 4 {
				p.problems = append(p.problems, fmt.Sprintf("[%x] unsupported IPTC DataSet size %d", p.offset+offset, dlen))
				return
			}
			lenbuf := make([]byte, 8)
			copy(lenbuf[4-dlen:4], buf[:dlen])
			buf = buf[dlen:]
			offset += uint32(dlen)
			dlen = binary.BigEndian.Uint32(lenbuf)
		}
		if int(dlen) > len(buf) {
			p.problems = append(p.problems, fmt.Sprintf("[%x] incomplete IPTC DataSet", p.offset+offset))
			return
		}
		dset.data = buf[:dlen]
		buf = buf[dlen:]
		offset += uint32(dlen)
		if dset.id == idCodedCharacterSet {
			if !bytes.Equal(dset.data, utf8Escape1) && !bytes.Equal(dset.data, utf8Escape2) {
				p.problems = append(p.problems, fmt.Sprintf("[%x] IPTC block uses character set other than UTF-8", p.offset+offset))
				return
			}
		}
		p.dsets = append(p.dsets, &dset)
	}
}

func (p *IPTC) findDSet(id uint16) *dsett {
	if p == nil {
		return nil
	}
	for _, dset := range p.dsets {
		if dset != nil && dset.id == id {
			return dset
		}
	}
	return nil
}
