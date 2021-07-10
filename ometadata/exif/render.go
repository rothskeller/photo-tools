package exif

import (
	"bytes"
)

// Dirty returns whether the EXIF data have changed and need to be saved.
func (p *EXIF) Dirty() bool {
	if p == nil || len(p.Problems) != 0 {
		return false
	}
	return p.tl.Dirty()
}

// Render renders and returns the encoded EXIF block, applying any changes made
// to the metadata fields of the EXIF structure.  maxSize is the maximum allowed
// size of the block.
func (p *EXIF) Render(max uint64) (out []byte) {
	var buf bytes.Buffer

	if len(p.Problems) != 0 {
		panic("EXIF Render with parse problems")
	}
	if err := p.tl.Render(&buf); err != nil {
		panic("EXIF render error: " + err.Error())
	}
	if buf.Len() > int(max) {
		panic("EXIF block doesn't fit within maximum size")
	}
	return buf.Bytes()
}

func (p *EXIF) addEXIFIFD() {
	tag := p.ifd0.AddTag(tagExifIFDOffset)
	p.exifIFD = tag.AddIFD()
}

func (p *EXIF) addGPSIFD() {
	tag := p.ifd0.AddTag(tagGPSIFDOffset)
	p.gpsIFD = tag.AddIFD()
	tag = p.gpsIFD.AddTag(0)
	tag.SetBytes([]byte{2, 3, 0, 0})
}
