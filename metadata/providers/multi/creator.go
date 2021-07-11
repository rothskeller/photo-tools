package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Creator returns the value of the Creator field.
func (p Provider) Creator() (value string) {
	for _, sp := range p {
		if value = sp.Creator(); value != "" {
			return value
		}
	}
	return ""
}

// CreatorTags returns a list of tag names for the Creator field, and a
// parallel list of values held by those tags.
func (p Provider) CreatorTags() (tags []string, values [][]string) {
	for _, sp := range p {
		t, v := sp.CreatorTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetCreator sets the value of the Creator field.
func (p Provider) SetCreator(value string) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetCreator(value); err != nil && err != metadata.ErrNotSupported {
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
