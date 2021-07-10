package fields

import (
	"github.com/rothskeller/photo-tools/metadata"
)

type topicsField struct {
	hierValueField
}

// TopicsField is the field handler for topic names, i.e., topics of people
// (organizations, teams, etc.) that are depicted in the media.
var TopicsField Field = &topicsField{
	hierValueField{
		baseField{
			name:        "topic",
			pluralName:  "topics",
			label:       "Topic",
			shortLabel:  "TP",
			multivalued: true,
		},
	},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *topicsField) GetValues(p metadata.Provider) []interface{} {
	var groups = p.Topics()
	var ifcs = make([]interface{}, len(groups))
	for i := range groups {
		ifcs[i] = groups[i]
	}
	return ifcs
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *topicsField) GetTags(p metadata.Provider) ([]string, []interface{}) {
	var tags, values = p.TopicsTags()
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// SetValues sets all of the values of the field.
func (f *topicsField) SetValues(p metadata.Provider, v []interface{}) error {
	var values = make([]metadata.HierValue, len(v))
	for i := range v {
		values[i] = v[i].(metadata.HierValue)
	}
	return p.SetTopics(values)
}
