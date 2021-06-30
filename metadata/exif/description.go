package exif

const tagImageDescription uint16 = 0x10E

// ImageDescription returns the value of the ImageDescription tag.
func (p *EXIF) ImageDescription() string { return p.imageDescription }

func (p *EXIF) getImageDescription() {
	if idt := p.ifd0.findTag(tagImageDescription); idt != nil {
		p.imageDescription = p.asciiAt(idt, "ImageDescription")
	}
}

// SetImageDescription sets the value of the ImageDescription tag.
func (p *EXIF) SetImageDescription(v string) error {
	if v == p.imageDescription {
		return nil
	}
	p.imageDescription = v
	if p.imageDescription == "" {
		p.deleteTag(p.ifd0, tagImageDescription)
	} else {
		p.setASCIITag(p.ifd0, tagImageDescription, p.imageDescription)
	}
	return nil
}
