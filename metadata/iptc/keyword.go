package iptc

import (
	"strings"
)

// MaxKeywordLen is the maximum length of one Keyword entry.
const MaxKeywordLen = 64

const idKeyword uint16 = 0x0219

func (p *IPTC) getKeywords() {
	for _, dset := range p.dsets {
		if dset != nil && dset.id == idKeyword {
			if kw := strings.TrimSpace(p.decodeString(dset.data, "Keyword")); kw != "" {
				p.Keywords = append(p.Keywords, kw)
			}
		}
	}
	p.saveKeywords = append([]string{}, p.Keywords...)
}

func (p *IPTC) setKeywords() {
	if stringSliceEqualMax(p.Keywords, p.saveKeywords, MaxKeywordLen) {
		return
	}
	if len(p.Keywords) == 0 {
		p.deleteDSet(idKeyword)
		return
	}
	var idx int
	for i, dset := range p.dsets {
		if dset != nil && dset.id == idKeyword {
			if idx < len(p.Keywords) {
				dset.data = []byte(applyMax(p.Keywords[idx], MaxKeywordLen))
				idx++
			} else {
				p.dsets[i] = nil
			}
		}
	}
	for idx < len(p.Keywords) {
		p.dsets = append(p.dsets, &dsett{0, idKeyword, []byte(applyMax(p.Keywords[idx], MaxKeywordLen))})
		idx++
	}
	p.dirty = true
}
