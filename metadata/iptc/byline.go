package iptc

import (
	"strings"
	"unicode/utf8"

	"github.com/rothskeller/photo-tools/metadata"
)

// MaxBylineLen is the maximum length of one By-line entry.
const MaxBylineLen = 32

const idByline uint16 = 0x0250

func (p *IPTC) getBylines() {
	for _, dset := range p.dsets {
		if dset != nil && dset.id == idByline {
			if utf8.Valid(dset.data) {
				if byline := strings.TrimSpace(string(dset.data)); byline != "" {
					var s metadata.String
					s.SetMaxLength(MaxBylineLen)
					s.Parse(byline)
					p.Bylines = append(p.Bylines, &s)
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
				p.Bylines[idx].SetMaxLength(MaxBylineLen)
				if next := p.Bylines[idx].String(); next != string(dset.data) {
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
		p.Bylines[idx].SetMaxLength(MaxBylineLen)
		p.dsets = append(p.dsets, &dsett{0, idByline, []byte(p.Bylines[idx].String())})
		p.dirty = true
		idx++
	}
}
