package exif

import (
	"github.com/rothskeller/photo-tools/metadata"
)

const tagImageDescription uint16 = 0x10E

func (p *EXIF) getImageDescription() {
	if idt := p.ifd0.findTag(tagImageDescription); idt != nil {
		p.ImageDescription = metadata.NewString(p.asciiAt(idt, "ImageDescription"))
	}
}

func (p *EXIF) setImageDescription() {
	if p.ImageDescription.Empty() {
		p.deleteTag(p.ifd0, tagImageDescription)
	} else {
		p.setASCIITag(p.ifd0, tagImageDescription, p.ImageDescription.String())
	}
}
