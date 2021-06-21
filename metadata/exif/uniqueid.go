package exif

import (
	"bytes"
)

const tagImageUniqueID uint16 = 0xA420

// ImageUniqueID returns the ImageUniqueID tag, if provided.
func (p *EXIF) ImageUniqueID() string {
	if p != nil && p.exifIFD != nil {
		if idt := p.exifIFD.findTag(tagImageUniqueID); idt != nil {
			return p.asciiAt(idt, "ImageUniqueID")
		}
	}
	return ""
}

// SetImageUniqueID sets the ImageUniqueID tag.
func (p *EXIF) SetImageUniqueID(id string) {
	if p == nil || p.exifIFD == nil {
		return
	}
	if id == "" {
		p.deleteTag(p.exifIFD, tagImageUniqueID)
		return
	}
	tag := p.exifIFD.findTag(tagImageUniqueID)
	if tag == nil {
		tag = &tagt{tag: tagImageUniqueID, ttype: 2, count: 1, data: []byte{0}}
		p.addTag(p.exifIFD, tag)
	}
	encbytes := []byte(id + "\000")
	if !bytes.Equal(encbytes, tag.data) {
		tag.data = encbytes
		tag.count = uint32(len(encbytes))
		p.exifIFD.dirty = true
	}
}
