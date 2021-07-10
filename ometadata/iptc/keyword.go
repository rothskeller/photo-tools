package iptc

import (
	"strings"
)

// MaxKeywordLen is the maximum length of one Keyword entry.
const MaxKeywordLen = 64

const idKeyword uint16 = 0x0219

// Keywords returns the values of the Keyword tag.
func (p *IPTC) Keywords() []string { return p.keywords }

func (p *IPTC) getKeywords() {
	for _, dset := range p.dsets {
		if dset != nil && dset.id == idKeyword {
			if kw := strings.TrimSpace(p.decodeString(dset.data, "Keyword")); kw != "" {
				p.keywords = append(p.keywords, kw)
			}
		}
	}
}

// SetKeywords sets the values of the Keyword tag.
func (p *IPTC) SetKeywords(v []string) error {
	if stringSliceEqualMax(v, p.keywords, MaxKeywordLen) {
		return nil
	}
	p.keywords = make([]string, len(v))
	for i := range v {
		p.keywords[i] = applyMax(v[i], MaxKeywordLen)
	}
	if len(p.keywords) == 0 {
		p.deleteDSet(idKeyword)
		return nil
	}
	var idx int
	for i, dset := range p.dsets {
		if dset != nil && dset.id == idKeyword {
			if idx < len(p.keywords) {
				dset.data = []byte(p.keywords[idx])
				idx++
			} else {
				p.dsets[i] = nil
			}
		}
	}
	for idx < len(p.keywords) {
		p.dsets = append(p.dsets, &dsett{0, idKeyword, []byte(p.keywords[idx])})
		idx++
	}
	p.dirty = true
	return nil
}
