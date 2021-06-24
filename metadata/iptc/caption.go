package iptc

import (
	"strings"
	"unicode/utf8"
)

// MaxCaptionAbstractLen is the maximum length of the Caption/Abstract entry.
const MaxCaptionAbstractLen = 2000

const idCaptionAbstract uint16 = 0x0278

func (p *IPTC) getCaptionAbstract() {
	if dset := p.findDSet(idCaptionAbstract); dset != nil {
		if utf8.Valid(dset.data) {
			p.CaptionAbstract = strings.TrimSpace(string(dset.data))
		} else {
			p.log("ignoring non-UTF8 Caption/Abstract")
		}
		return
	}
}

func (p *IPTC) setCaptionAbstract() {
	if p.CaptionAbstract == "" {
		p.deleteDSet(idCaptionAbstract)
		return
	}
	encoded := []byte(applyMax(p.CaptionAbstract, MaxCaptionAbstractLen))
	p.setDSet(idCaptionAbstract, encoded)
}
