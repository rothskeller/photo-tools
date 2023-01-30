package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Orientation returns the value of the Orientation field.
func (p Provider) Orientation() (value metadata.Orientation) {
	for _, sp := range p {
		if value = sp.Orientation(); value != metadata.Rotate0 {
			return value
		}
	}
	return metadata.Rotate0
}

// OrientationTags returns a list of tag names for the Orientation field, and a
// parallel list of values held by those tags.
func (p Provider) OrientationTags() (tags []string, values [][]metadata.Orientation) {
	for _, sp := range p {
		t, v := sp.OrientationTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetOrientation sets the value of the Orientation field.
func (p Provider) SetOrientation(value metadata.Orientation) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetOrientation(value); err != nil && err != metadata.ErrNotSupported {
			return err
		} else if err == nil {
			set = true
		}
	}
	if !set {
		return metadata.ErrNotSupported
	}
	return nil
}
