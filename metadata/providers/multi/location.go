package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Location returns the value of the Location field.
func (p Provider) Location() (value metadata.Location) {
	for _, sp := range p {
		if value = sp.Location(); !value.Empty() {
			return value
		}
	}
	return metadata.Location{}
}

// LocationTags returns a list of tag names for the Location field, and a parallel
// list of values held by those tags.
func (p Provider) LocationTags() (tags []string, values []metadata.Location) {
	for _, sp := range p {
		t, v := sp.LocationTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetLocation sets the value of the Location field.
func (p Provider) SetLocation(value metadata.Location) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetLocation(value); err != nil && err != metadata.ErrNotSupported {
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
