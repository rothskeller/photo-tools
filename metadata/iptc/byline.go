package iptc

import (
	"strings"
)

// MaxBylineLen is the maximum length of one By-line entry.
const MaxBylineLen = 32

const idByline uint16 = 0x0250

func (p *IPTC) getBylines() {
	for _, dset := range p.dsets {
		if dset != nil && dset.id == idByline {
			if byline := strings.TrimSpace(p.decodeString(dset.data, "Byline")); byline != "" {
				p.Bylines = append(p.Bylines, byline)
			}
		}
	}
	p.saveBylines = append([]string{}, p.Bylines...)
}

func (p *IPTC) setBylines() {
	if stringSliceEqualMax(p.Bylines, p.saveBylines, MaxBylineLen) {
		return
	}
	var idx int
	for i, dset := range p.dsets {
		if dset != nil && dset.id == idByline {
			if idx < len(p.Bylines) {
				dset.data = []byte(applyMax(p.Bylines[idx], MaxBylineLen))
				idx++
			} else {
				p.dsets[i] = nil
			}
		}
	}
	for idx < len(p.Bylines) {
		p.dsets = append(p.dsets, &dsett{0, idByline, []byte(applyMax(p.Bylines[idx], MaxBylineLen))})
		idx++
	}
	p.dirty = true
}
