package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Topics returns the value of the Topics field.
func (p Provider) Topics() (value []metadata.HierValue) {
	for _, sp := range p {
		if value = sp.Topics(); len(value) != 0 {
			return value
		}
	}
	return nil
}

// TopicsTags returns a list of tag names for the Topics field, and a parallel
// list of values held by those tags.
func (p Provider) TopicsTags() (tags []string, values []metadata.HierValue) {
	for _, sp := range p {
		t, v := sp.TopicsTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetTopics sets the value of the Topics field.
func (p Provider) SetTopics(value []metadata.HierValue) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetTopics(value); err != nil && err != metadata.ErrNotSupported {
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
