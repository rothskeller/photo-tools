package tiffifd0

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
)

const tagOrientation uint16 = 0x112

// getOrientation reads the value of the Orientation field from the IFD.
func (p *Provider) getOrientation() (err error) {
	if tag := p.ifd.Tag(tagOrientation); tag != nil {
		var v int
		if v, err = tag.AsShort(); err != nil {
			return fmt.Errorf("ImageDescription: %s", err)
		}
		p.orientation = metadata.Orientation(v)
	}
	return nil
}

// Orientation returns the value of the Orientation field.
func (p *Provider) Orientation() (value metadata.Orientation) { return p.orientation }

// OrientationTags returns a list of tag names for the Orientation field, and a
// parallel list of values held by those tags.
func (p *Provider) OrientationTags() (tags []string, values [][]metadata.Orientation) {
	return []string{"IFD0 Orientation"}, [][]metadata.Orientation{{p.orientation}}
}

// SetOrientation sets the value of the Orientation field.
func (p *Provider) SetOrientation(value metadata.Orientation) error {
	if value == 0 || value == metadata.Rotate0 {
		p.orientation = 0
		p.ifd.DeleteTag(tagOrientation)
		return nil
	}
	if value == p.orientation {
		return nil
	}
	p.orientation = value
	p.ifd.AddTag(tagOrientation, 3).SetShort(int(value))
	return nil
}
