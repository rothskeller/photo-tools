package iptc

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

// MaxCaptionAbstractLen is the maximum length of the Caption/Abstract entry.
const MaxCaptionAbstractLen = 2000

const idCaptionAbstract uint16 = 0x0278

// CaptionAbstract returns the IPTC Caption/Abstract tag, if any.
func (p *IPTC) CaptionAbstract() string {
	if dset := p.findDSet(idCaptionAbstract); dset != nil {
		if utf8.Valid(dset.data) {
			return strings.TrimSpace(string(dset.data))
		}
		p.problems = append(p.problems, "ignoring non-UTF8 IPTC Caption/Abstract")
	}
	return ""
}

// SetCaptionAbstract sets the IPTC Caption/Abstract.
func (p *IPTC) SetCaptionAbstract(val string) {
	if p == nil {
		return
	}
	if val == "" {
		p.deleteDSet(idCaptionAbstract)
		return
	}
	dset := p.findDSet(idCaptionAbstract)
	encoded := []byte(applyMax(val, MaxCaptionAbstractLen))
	if dset != nil {
		if !bytes.Equal(encoded, dset.data) {
			dset.data = encoded
			p.dirty = true
		}
	} else {
		p.dsets = append(p.dsets, &dsett{0, idCaptionAbstract, encoded})
		p.dirty = true
	}
}
