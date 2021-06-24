package iptc

import (
	"strings"
	"unicode/utf8"

	"github.com/rothskeller/photo-tools/metadata"
)

// MaxKeywordLen is the maximum length of one Keyword entry.
const MaxKeywordLen = 64

const idKeyword uint16 = 0x0219

func (p *IPTC) getKeywords() {
	for _, dset := range p.dsets {
		if dset != nil && dset.id == idKeyword {
			if utf8.Valid(dset.data) {
				if kw := strings.TrimSpace(string(dset.data)); kw != "" {
					var s metadata.String
					s.SetMaxLength(MaxKeywordLen)
					s.Parse(kw)
					p.Keywords = append(p.Keywords, &s)
				}
			} else {
				p.log("ignoring non-UTF8 Keyword")
			}
		}
	}
}

func (p *IPTC) setKeywords() {
	if len(p.Keywords) == 0 {
		p.deleteDSet(idKeyword)
		return
	}
	var idx int
	for i, dset := range p.dsets {
		if dset != nil && dset.id == idKeyword {
			if idx < len(p.Keywords) {
				p.Keywords[idx].SetMaxLength(MaxKeywordLen)
				if next := p.Keywords[idx].String(); next != string(dset.data) {
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
	for idx < len(p.Keywords) {
		p.Keywords[idx].SetMaxLength(MaxKeywordLen)
		p.dsets = append(p.dsets, &dsett{0, idKeyword, []byte(p.Keywords[idx].String())})
		p.dirty = true
		idx++
	}
}
