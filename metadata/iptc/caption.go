package iptc

import (
	"strings"
)

// MaxCaptionAbstractLen is the maximum length of the Caption/Abstract entry.
const MaxCaptionAbstractLen = 2000

const idCaptionAbstract uint16 = 0x0278

// CaptionAbstract returns the value of the Caption/Abstract tag.
func (p *IPTC) CaptionAbstract() string { return p.captionAbstract }

func (p *IPTC) getCaptionAbstract() {
	if dset := p.findDSet(idCaptionAbstract); dset != nil {
		p.captionAbstract = strings.TrimSpace(p.decodeString(dset.data, "CaptionAbstract"))
	}
}

// SetCaptionAbstract sets the value of the Caption/Abstract tag.
func (p *IPTC) SetCaptionAbstract(v string) error {
	if stringEqualMax(v, p.captionAbstract, MaxCaptionAbstractLen) {
		return nil
	}
	p.captionAbstract = applyMax(v, MaxCaptionAbstractLen)
	if p.captionAbstract == "" {
		p.deleteDSet(idCaptionAbstract)
		return nil
	}
	p.setDSet(idCaptionAbstract, []byte(p.captionAbstract))
	return nil
}
