package iptc

import (
	"strings"
	"unicode/utf8"
)

// MaxKeywordLen is the maximum length of one Keyword entry.
const MaxKeywordLen = 64

const idKeyword uint16 = 0x0219

// Keywords returns the IPTC Keyword tags, if any.
func (p *IPTC) Keywords() (vals []string) {
	if p == nil {
		return nil
	}
	for _, dset := range p.dsets {
		if dset != nil && dset.id == idKeyword {
			if utf8.Valid(dset.data) {
				vals = append(vals, strings.TrimSpace(string(dset.data)))
			} else {
				p.problems = append(p.problems, "ignoring non-UTF8 IPTC Keyword")
			}
		}
	}
	return vals
}

// SetKeywords sets the IPTC Keywords.
func (p *IPTC) SetKeywords(val []string) {
	if p == nil {
		return
	}
	if len(val) == 0 {
		p.deleteDSet(idKeyword)
		return
	}
	for i, dset := range p.dsets {
		if dset != nil && dset.id == idKeyword {
			if len(val) != 0 {
				next := applyMax(val[0], MaxKeywordLen)
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
		p.dsets = append(p.dsets, &dsett{0, idKeyword, []byte(applyMax(val[0], MaxKeywordLen))})
		p.dirty = true
		val = val[1:]
	}
}
