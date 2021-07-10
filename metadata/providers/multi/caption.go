package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Caption returns the value of the Caption field.
func (p Provider) Caption() (value string) {
	for _, sp := range p {
		if value = sp.Caption(); value != "" {
			return value
		}
	}
	return ""
}

// CaptionTags returns a list of tag names for the Caption field, and a
// parallel list of values held by those tags.
func (p Provider) CaptionTags() (tags []string, values []string) {
	for _, sp := range p {
		t, v := sp.CaptionTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetCaption sets the value of the Caption field.
func (p Provider) SetCaption(value string) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetCaption(value); err != nil && err != metadata.ErrNotSupported {
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
