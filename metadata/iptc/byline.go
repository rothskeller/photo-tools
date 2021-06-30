package iptc

import (
	"strings"
)

// MaxBylineLen is the maximum length of one By-line entry.
const MaxBylineLen = 32

const idByline uint16 = 0x0250

// Bylines returns the values of the Byline tag.
func (p *IPTC) Bylines() []string { return p.bylines }

func (p *IPTC) getBylines() {
	for _, dset := range p.dsets {
		if dset != nil && dset.id == idByline {
			if byline := strings.TrimSpace(p.decodeString(dset.data, "Byline")); byline != "" {
				p.bylines = append(p.bylines, byline)
			}
		}
	}
}

// SetBylines sets the values of the Byline tag.
func (p *IPTC) SetBylines(v []string) error {
	if stringSliceEqualMax(v, p.bylines, MaxBylineLen) {
		return nil
	}
	p.bylines = make([]string, len(v))
	for i := range v {
		p.bylines[i] = applyMax(v[i], MaxBylineLen)
	}
	var idx int
	for i, dset := range p.dsets {
		if dset != nil && dset.id == idByline {
			if idx < len(p.bylines) {
				dset.data = []byte(p.bylines[idx])
				idx++
			} else {
				p.dsets[i] = nil
			}
		}
	}
	for idx < len(p.bylines) {
		p.dsets = append(p.dsets, &dsett{0, idByline, []byte(p.bylines[idx])})
		idx++
	}
	p.dirty = true
	return nil
}
