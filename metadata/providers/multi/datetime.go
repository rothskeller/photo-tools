package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// DateTime returns the value of the DateTime field.
func (p Provider) DateTime() (value metadata.DateTime) {
	for _, sp := range p {
		if value = sp.DateTime(); !value.Empty() {
			return value
		}
	}
	return metadata.DateTime{}
}

// DateTimeTags returns a list of tag names for the DateTime field, and a parallel
// list of values held by those tags.
func (p Provider) DateTimeTags() (tags []string, values []metadata.DateTime) {
	for _, sp := range p {
		t, v := sp.DateTimeTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetDateTime sets the value of the DateTime field.
func (p Provider) SetDateTime(value metadata.DateTime) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetDateTime(value); err != nil && err != metadata.ErrNotSupported {
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
