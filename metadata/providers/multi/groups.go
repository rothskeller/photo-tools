package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Groups returns the value of the Groups field.
func (p Provider) Groups() (value []metadata.HierValue) {
	for _, sp := range p {
		if value = sp.Groups(); len(value) != 0 {
			return value
		}
	}
	return nil
}

// GroupsTags returns a list of tag names for the Groups field, and a parallel
// list of values held by those tags.
func (p Provider) GroupsTags() (tags []string, values []metadata.HierValue) {
	for _, sp := range p {
		t, v := sp.GroupsTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetGroups sets the value of the Groups field.
func (p Provider) SetGroups(value []metadata.HierValue) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetGroups(value); err != nil && err != metadata.ErrNotSupported {
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
