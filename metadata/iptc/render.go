package iptc

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"sort"
)

// Dirty returns whether the IPTC data have changed and need to be saved.
func (p *IPTC) Dirty() bool {
	if p == nil || len(p.Problems) != 0 {
		return false
	}
	return p.dirty
}

// Render renders and returns the encoded IPTC block, reflecting the data that
// was read, as subsequently modified by any SetXXX calls.
func (p *IPTC) Render() []byte {
	if len(p.Problems) != 0 {
		panic("IPTC Render with parse problems")
	}
	var out bytes.Buffer
	for _, psir := range p.psir {
		if psir.id == iptcPSIRID && p.dirty {
			psirIPTC, psirHash := p.renderIPTC()
			p.renderPSIR(&out, psirIPTC)
			p.renderPSIR(&out, psirHash)
		} else if psir.id != iptcPSIRHash || !p.dirty {
			p.renderPSIR(&out, psir)
		}
	}
	return out.Bytes()
}

// renderPSIR renders a PSIR.
func (p *IPTC) renderPSIR(out *bytes.Buffer, psir *psirt) {
	var buf = make([]byte, 4)
	out.Write(psirResourceType)
	binary.BigEndian.PutUint16(buf, psir.id)
	out.Write(buf[:2])
	out.WriteByte(byte(len(psir.name)))
	out.WriteString(psir.name)
	if len(psir.name)%2 == 0 {
		out.WriteByte(0)
	}
	binary.BigEndian.PutUint32(buf, uint32(len(psir.buf)))
	out.Write(buf[:4])
	out.Write(psir.buf)
	if len(psir.buf)%2 == 1 {
		out.WriteByte(0)
	}
}

// renderIPTC generates new PSIRs for the IPTC data and its hash.
func (p *IPTC) renderIPTC() (psir, hash *psirt) {
	if dset := p.findDSet(idCodedCharacterSet); dset == nil {
		p.dsets = append(p.dsets, &dsett{0, idCodedCharacterSet, utf8Escape1})
	}
	sort.SliceStable(p.dsets, func(i, j int) bool {
		if p.dsets[i] == nil && p.dsets[j] != nil {
			return false
		}
		if p.dsets[j] == nil {
			return true
		}
		return p.dsets[i].id < p.dsets[j].id
	})
	psir = new(psirt)
	psir.id = iptcPSIRID
	var out bytes.Buffer
	var buf = make([]byte, 4)
	for _, dset := range p.dsets {
		if dset == nil {
			continue
		}
		out.WriteByte(iptcTagMarker)
		binary.BigEndian.PutUint16(buf, dset.id)
		out.Write(buf[:2])
		if len(dset.data) < 0x8000 {
			binary.BigEndian.PutUint16(buf, uint16(len(dset.data)))
			out.Write(buf[:2])
		} else {
			binary.BigEndian.PutUint16(buf, 0x8004)
			out.Write(buf[:2])
			binary.BigEndian.PutUint32(buf, uint32(len(dset.data)))
			out.Write(buf[:4])
		}
		out.Write(dset.data)
	}
	psir.buf = out.Bytes()
	hash = new(psirt)
	hash.id = iptcPSIRHash
	if h := md5.Sum(psir.buf); true {
		hash.buf = h[:]
	}
	return psir, hash
}

func (p *IPTC) deleteDSet(id uint16) {
	for i, dset := range p.dsets {
		if dset != nil && dset.id == id {
			p.dsets[i] = nil
			p.dirty = true
		}
	}
}

func (p *IPTC) setDSet(id uint16, val []byte) {
	dset := p.findDSet(id)
	if dset != nil {
		if !bytes.Equal(dset.data, val) {
			dset.data = val
			p.dirty = true
		}
	} else {
		p.dsets = append(p.dsets, &dsett{0, id, val})
		p.dirty = true
	}
}

func applyMax(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}
