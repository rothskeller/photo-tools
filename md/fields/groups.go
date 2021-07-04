package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

type groupsField struct {
	baseField
}

// GroupsField is the field handler for group names, i.e., groups of people
// (organizations, teams, etc.) that are depicted in the media.
var GroupsField Field = &groupsField{
	baseField{
		name:        "group",
		pluralName:  "groups",
		label:       "Group",
		shortLabel:  "GR",
		multivalued: true,
	},
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *groupsField) ParseValue(s string) (interface{}, error) {
	var v strmeta.Group
	if err := v.Parse(s); err != nil {
		return nil, err
	}
	return v, nil
}

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *groupsField) RenderValue(v interface{}) string { return v.(strmeta.Group).String() }

// EqualValue compares two values for equality.
func (f *groupsField) EqualValue(a interface{}, b interface{}) bool {
	return a.(strmeta.Group).Equal(b.(strmeta.Group))
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *groupsField) GetValues(h filefmt.FileHandler) []interface{} {
	var groups = strmeta.GetGroups(h)
	var ifcs = make([]interface{}, len(groups))
	for i := range groups {
		ifcs[i] = groups[i]
	}
	return ifcs
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *groupsField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	var tags, values = strmeta.GetGroupTags(h)
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// CheckValues returns whether the values of the field are tagged correctly.
func (f *groupsField) CheckValues(h filefmt.FileHandler) (res strmeta.CheckResult) {
	if res = strmeta.CheckGroups(h); res == strmeta.ChkPresent {
		res = strmeta.CheckResult(len(f.GetValues(h)))
	}
	return res
}

// SetValues sets all of the values of the field.
func (f *groupsField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	var groups = make([]strmeta.Group, len(v))
	for i := range v {
		groups[i] = v[i].(strmeta.Group)
	}
	return strmeta.SetGroups(h, groups)
}
