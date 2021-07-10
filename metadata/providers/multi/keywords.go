package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Keywords returns the value of the Keywords field.
func (p Provider) Keywords() (value []metadata.HierValue) {
	for _, sp := range p {
		if value = sp.Keywords(); len(value) != 0 {
			return value
		}
	}
	return nil
}

// KeywordsTags returns a list of tag names for the Keywords field, and a parallel
// list of values held by those tags.
func (p Provider) KeywordsTags() (tags []string, values []metadata.HierValue) {
	for _, sp := range p {
		t, v := sp.KeywordsTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetKeywords sets the value of the Keywords field.
func (p Provider) SetKeywords(value []metadata.HierValue) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetKeywords(value); err != nil && err != metadata.ErrNotSupported {
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
