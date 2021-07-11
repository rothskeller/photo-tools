package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Faces returns the value of the Faces field.
func (p Provider) Faces() (value []string) {
	for _, sp := range p {
		if value = sp.Faces(); len(value) != 0 {
			return value
		}
	}
	return nil
}

// FacesTags returns a list of tag names for the Faces field, and a parallel
// list of values held by those tags.
func (p Provider) FacesTags() (tags []string, values [][]string) {
	for _, sp := range p {
		t, v := sp.FacesTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetFaces sets the value of the Faces field.
func (p Provider) SetFaces(value []string) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetFaces(value); err != nil && err != metadata.ErrNotSupported {
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
