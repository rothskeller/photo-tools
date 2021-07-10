// Package iptc handles IPTC metadata blocks.
package iptc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

const (
	iptcPSIRID          uint16 = 0x404
	iptcPSIRHash        uint16 = 0x425
	iptcTagMarker       byte   = 0x1C
	idCodedCharacterSet uint16 = 0x015A
)

var psirResourceType = []byte("8BIM")
var utf8Escape1 = []byte{0x1B, 0x25, 0x47}
var utf8Escape2 = []byte{0x1B, 0x25, 0x2F, 0x49}

// Parse saves a portion of an IPTC block for later parsing.  offset is the
// offset of the IPTC block in the file, used for problem reporting.
func (p *IPTC) Parse(buf []byte, offset uint32) *IPTC {
	if p == nil {
		p = &IPTC{offset: offset, buf: buf}
	} else {
		p.buf = append(p.buf, buf...)
	}
	return p
}

// Check does the actual parsing of the IPTC block after all of its portions
// have been retrieved.
func (p *IPTC) Check() {
	if p == nil || !p.splitPSIRs() {
		return
	}
	for _, psir := range p.psir {
		if psir.id == iptcPSIRID {
			p.parseIPTC(psir)
		}
	}
	p.getBylines()
	p.getCaptionAbstract()
	p.getDateTimeCreated()
	p.getDigitalCreationDateTime()
	p.getKeywords()
	p.getLocation()
	p.getObjectName()
}

// splitPSIRs splits the block up into PSIRs.  It returns false if a format
// problem was found.
func (p *IPTC) splitPSIRs() bool {
	var poff uint32
	var bufend = uint32(len(p.buf))
	for poff < bufend {
		var psir psirt
		if poff+8 > bufend {
			p.log("[%x] incomplete PSIR", p.offset+poff)
			return false
		}
		if !bytes.Equal(p.buf[poff:poff+4], psirResourceType) {
			p.log("[%x] invalid PSIR resource type", p.offset+poff)
			return false
		}
		psir.id = binary.BigEndian.Uint16(p.buf[poff+4:])
		resNameLen := p.buf[poff+6]
		resNameSize := uint32(resNameLen + 1)
		if resNameSize%2 == 1 {
			resNameSize++
		}
		if poff+6+resNameSize+4 > bufend {
			p.log("[%x] incomplete PSIR", p.offset+poff)
			return false
		}
		psir.name = string(p.buf[poff+7 : poff+7+uint32(resNameLen)])
		resLenOff := poff + 6 + resNameSize
		resLen := binary.BigEndian.Uint32(p.buf[resLenOff:])
		resSize := resLen
		if resSize%2 == 1 {
			resSize++
		}
		if resLenOff+4+resSize > bufend {
			p.log("[%x] incomplete PSIR", p.offset+poff)
			println("return")
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
			p.log("[%x] incomplete DataSet", p.offset+offset)
			return
		}
		if buf[0] != iptcTagMarker {
			p.log("[%x] invalid DataSet tag marker", p.offset+offset)
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
				p.log("[%x] incomplete DataSet", p.offset+offset)
				return
			}
			if dlen > 4 {
				p.log("[%x] unsupported DataSet size %d", p.offset+offset, dlen)
				return
			}
			lenbuf := make([]byte, 8)
			copy(lenbuf[4-dlen:4], buf[:dlen])
			buf = buf[dlen:]
			offset += uint32(dlen)
			dlen = binary.BigEndian.Uint32(lenbuf)
		}
		if int(dlen) > len(buf) {
			p.log("[%x] incomplete DataSet", p.offset+offset)
			return
		}
		dset.data = buf[:dlen]
		buf = buf[dlen:]
		offset += uint32(dlen)
		if dset.id == idCodedCharacterSet {
			if !bytes.Equal(dset.data, utf8Escape1) && !bytes.Equal(dset.data, utf8Escape2) {
				p.log("[%x] block uses character set other than UTF-8", p.offset+offset)
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

func (p *IPTC) decodeString(by []byte, label string) string {
	if utf8.Valid(by) {
		return string(by)
	}
	if by2, err := charmap.ISO8859_1.NewDecoder().Bytes(by); err == nil {
		return string(by2)
	}
	p.log("cannot determine character set for %s", label)
	return ""
}

func (p *IPTC) log(f string, a ...interface{}) {
	problem := fmt.Sprintf(f, a...)
	p.Problems = append(p.Problems, "IPTC: "+problem)
}
