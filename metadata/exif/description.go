package exif

import (
	"bytes"
)

const tagImageDescription uint16 = 0x10E

// ImageDescription returns the ImageDescription tag, if provided.
func (p *EXIF) ImageDescription() string {
	if p != nil && p.ifd0 != nil {
		if idt := p.ifd0.findTag(tagImageDescription); idt != nil {
			return p.asciiAt(idt, "ImageDescription")
		}
	}
	return ""
}

// SetImageDescription sets the ImageDescription tag.
func (p *EXIF) SetImageDescription(desc string) {
	if p == nil || p.ifd0 == nil {
		return
	}
	if desc == "" {
		p.deleteTag(p.ifd0, tagImageDescription)
		return
	}
	tag := p.ifd0.findTag(tagImageDescription)
	if tag == nil {
		tag = &tagt{tag: tagImageDescription, ttype: 2, count: 1, data: []byte{0}}
		p.addTag(p.ifd0, tag)
	}
	encbytes := []byte(desc + "\000")
	if !bytes.Equal(encbytes, tag.data) {
		tag.data = encbytes
		tag.count = uint32(len(encbytes))
		p.ifd0.dirty = true
	}
}
