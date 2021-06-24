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
			p.CaptionAbstract.SetMaxLength(MaxCaptionAbstractLen)
			p.CaptionAbstract.Parse(strings.TrimSpace(string(dset.data)))
		} else {
			p.log("ignoring non-UTF8 Caption/Abstract")
		}
		return
	}
}

func (p *IPTC) setCaptionAbstract() {
	if p.CaptionAbstract.Empty() {
		p.deleteDSet(idCaptionAbstract)
		return
	}
	p.CaptionAbstract.SetMaxLength(MaxCaptionAbstractLen)
	encoded := []byte(p.CaptionAbstract.String())
	p.setDSet(idCaptionAbstract, encoded)
}
