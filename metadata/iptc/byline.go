package iptc

import (
	"strings"
	"unicode/utf8"
)

// MaxBylineLen is the maximum length of one By-line entry.
const MaxBylineLen = 32

const idByline uint16 = 0x0250

// Bylines returns the IPTC By-line tags, if any.
func (p *IPTC) Bylines() (vals []string) {
	if p == nil {
		return nil
	}
	for _, dset := range p.dsets {
		if dset != nil && dset.id == idByline {
			if utf8.Valid(dset.data) {
				vals = append(vals, strings.TrimSpace(string(dset.data)))
			} else {
				p.problems = append(p.problems, "ignoring non-UTF8 IPTC By-line")
			}
		}
	}
	return vals
}

// SetBylines sets the IPTC By-lines.
func (p *IPTC) SetBylines(val []string) {
	if p == nil {
		return
	}
	if len(val) == 0 {
		p.deleteDSet(idByline)
		return
	}
	for i, dset := range p.dsets {
		if dset != nil && dset.id == idByline {
			if len(val) != 0 {
				next := applyMax(val[0], MaxBylineLen)
				if next != string(dset.data) {
					dset.data = []byte(next)
					p.dirty = true
				}
				val = val[1:]
			} else {
				p.dsets[i] = nil
				p.dirty = true
			}
		}
	}
	for len(val) != 0 {
		p.dsets = append(p.dsets, &dsett{0, idByline, []byte(applyMax(val[0], MaxBylineLen))})
		p.dirty = true
		val = val[1:]
	}
}
