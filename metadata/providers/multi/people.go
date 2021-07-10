package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// People returns the value of the People field.
func (p Provider) People() (value []string) {
	for _, sp := range p {
		if value = sp.People(); len(value) != 0 {
			return value
		}
	}
	return nil
}

// PeopleTags returns a list of tag names for the People field, and a parallel
// list of values held by those tags.
func (p Provider) PeopleTags() (tags []string, values []string) {
	for _, sp := range p {
		t, v := sp.PeopleTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetPeople sets the value of the People field.
func (p Provider) SetPeople(value []string) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetPeople(value); err != nil && err != metadata.ErrNotSupported {
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
