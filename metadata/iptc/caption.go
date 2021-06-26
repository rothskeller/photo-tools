package iptc

import (
	"strings"
)

// MaxCaptionAbstractLen is the maximum length of the Caption/Abstract entry.
const MaxCaptionAbstractLen = 2000

const idCaptionAbstract uint16 = 0x0278

func (p *IPTC) getCaptionAbstract() {
	if dset := p.findDSet(idCaptionAbstract); dset != nil {
		p.CaptionAbstract = strings.TrimSpace(p.decodeString(dset.data, "CaptionAbstract"))
		p.saveCaptionAbstract = p.CaptionAbstract
		return
	}
}

func (p *IPTC) setCaptionAbstract() {
	if stringEqualMax(p.CaptionAbstract, p.saveCaptionAbstract, MaxCaptionAbstractLen) {
		return
	}
	if p.CaptionAbstract == "" {
		p.deleteDSet(idCaptionAbstract)
		return
	}
	p.setDSet(idCaptionAbstract, []byte(applyMax(p.CaptionAbstract, MaxCaptionAbstractLen)))
}
