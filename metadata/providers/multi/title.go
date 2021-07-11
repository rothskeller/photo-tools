package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Title returns the value of the Title field.
func (p Provider) Title() (value string) {
	for _, sp := range p {
		if value = sp.Title(); value != "" {
			return value
		}
	}
	return ""
}

// TitleTags returns a list of tag names for the Title field, and a parallel
// list of values held by those tags.
func (p Provider) TitleTags() (tags []string, values [][]string) {
	for _, sp := range p {
		t, v := sp.TitleTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetTitle sets the value of the Title field.
func (p Provider) SetTitle(value string) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetTitle(value); err != nil && err != metadata.ErrNotSupported {
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
