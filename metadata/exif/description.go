package exif

const tagImageDescription uint16 = 0x10E

// ImageDescription returns the value of the ImageDescription tag.
func (p *EXIF) ImageDescription() string { return p.imageDescription }

func (p *EXIF) getImageDescription() {
	var err error

	if tag := p.ifd0.Tag(tagImageDescription); tag != nil {
		if p.imageDescription, err = tag.AsString(); err != nil {
			p.log("ImageDescription: %s", err)
		}
	}
}

// SetImageDescription sets the value of the ImageDescription tag.
func (p *EXIF) SetImageDescription(v string) error {
	if v == p.imageDescription {
		return nil
	}
	p.imageDescription = v
	if p.imageDescription == "" {
		p.ifd0.DeleteTag(tagImageDescription)
	} else {
		p.ifd0.AddTag(tagImageDescription).SetString(p.imageDescription)
	}
	return nil
}
