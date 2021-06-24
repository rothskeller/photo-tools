package iptc

import (
	"strings"
	"unicode/utf8"
)

// MaxBylineLen is the maximum length of one By-line entry.
const MaxBylineLen = 32

const idByline uint16 = 0x0250

func (p *IPTC) getBylines() {
	for _, dset := range p.dsets {
		if dset != nil && dset.id == idByline {
			if utf8.Valid(dset.data) {
				if byline := strings.TrimSpace(string(dset.data)); byline != "" {
					p.Bylines = append(p.Bylines, byline)
				}
			} else {
				p.log("ignoring non-UTF8 By-line")
			}
		}
	}
}

func (p *IPTC) setBylines() {
	if len(p.Bylines) == 0 {
		p.deleteDSet(idByline)
		return
	}
	var idx int
	for i, dset := range p.dsets {
		if dset != nil && dset.id == idByline {
			if idx < len(p.Bylines) {
				if next := applyMax(p.Bylines[idx], MaxBylineLen); next != string(dset.data) {
					dset.data = []byte(next)
					p.dirty = true
				}
				idx++
			} else {
				p.dsets[i] = nil
				p.dirty = true
			}
		}
	}
	for idx < len(p.Bylines) {
		p.dsets = append(p.dsets, &dsett{0, idByline, []byte(applyMax(p.Bylines[idx], MaxBylineLen))})
		p.dirty = true
		idx++
	}
}
