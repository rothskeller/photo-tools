package fields

import (
	"github.com/rothskeller/photo-tools/metadata"
)

type groupsField struct {
	hierValueField
}

// GroupsField is the field handler for group names, i.e., groups of people
// (organizations, teams, etc.) that are depicted in the media.
var GroupsField Field = &groupsField{
	hierValueField{
		baseField{
			name:        "group",
			pluralName:  "groups",
			label:       "Group",
			shortLabel:  "GR",
			multivalued: true,
		},
	},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *groupsField) GetValues(p metadata.Provider) []interface{} {
	var groups = p.Groups()
	var ifcs = make([]interface{}, len(groups))
	for i := range groups {
		ifcs[i] = groups[i]
	}
	return ifcs
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *groupsField) GetTags(p metadata.Provider) ([]string, []interface{}) {
	var tags, values = p.GroupsTags()
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// SetValues sets all of the values of the field.
func (f *groupsField) SetValues(p metadata.Provider, v []interface{}) error {
	var values = make([]metadata.HierValue, len(v))
	for i := range v {
		values[i] = v[i].(metadata.HierValue)
	}
	return p.SetGroups(values)
}
