package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Places returns the value of the Places field.
func (p Provider) Places() (value []metadata.HierValue) {
	for _, sp := range p {
		if value = sp.Places(); len(value) != 0 {
			return value
		}
	}
	return nil
}

// PlacesTags returns a list of tag names for the Places field, and a parallel
// list of values held by those tags.
func (p Provider) PlacesTags() (tags []string, values []metadata.HierValue) {
	for _, sp := range p {
		t, v := sp.PlacesTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetPlaces sets the value of the Places field.
func (p Provider) SetPlaces(value []metadata.HierValue) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetPlaces(value); err != nil && err != metadata.ErrNotSupported {
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
