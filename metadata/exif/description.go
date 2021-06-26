package exif

const tagImageDescription uint16 = 0x10E

func (p *EXIF) getImageDescription() {
	if idt := p.ifd0.findTag(tagImageDescription); idt != nil {
		p.ImageDescription = p.asciiAt(idt, "ImageDescription")
		p.saveImageDescription = p.ImageDescription
	}
}

func (p *EXIF) setImageDescription() {
	if p.ImageDescription == p.saveImageDescription {
		return
	}
	if p.ImageDescription == "" {
		p.deleteTag(p.ifd0, tagImageDescription)
	} else {
		p.setASCIITag(p.ifd0, tagImageDescription, p.ImageDescription)
	}
}
