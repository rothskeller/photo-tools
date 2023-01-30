package exififd

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
)

const tagOrientation uint16 = 0x112

// getOrientation reads the value of the Orientation field from the RDF.
func (p *Provider) getOrientation() (err error) {
	tag := p.ifd.Tag(tagOrientation)
	if tag == nil {
		return nil
	}
	data, err := tag.AsShort()
	if err != nil {
		return fmt.Errorf("Orientation: %s", err)
	}
	if data < 1 || data > 8 {
		return fmt.Errorf("Orientation: illegal value %d", data)
	}
	p.orientation = metadata.Orientation(data)
	return nil
}

// Orientation returns the value of the Orientation field.
func (p *Provider) Orientation() (value metadata.Orientation) { return p.orientation }

// OrientationTags returns a list of tag names for the Orientation field, and a
// parallel list of values held by those tags.
func (p *Provider) OrientationTags() (tags []string, values [][]metadata.Orientation) {
	if p.orientation == 0 || p.orientation == metadata.Rotate0 {
		return nil, nil
	}
	return []string{"EXIF Orientation"}, [][]metadata.Orientation{{p.orientation}}
}

// SetOrientation sets the value of the Orientation field.
func (p *Provider) SetOrientation(value metadata.Orientation) error {
	if value == 0 || value == metadata.Rotate0 {
		p.orientation = 0
		p.ifd.DeleteTag(tagOrientation)
	} else {
		p.orientation = value
		p.ifd.AddTag(tagOrientation, 3).SetShort(int(value))
	}
	return nil
}
