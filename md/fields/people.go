package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

// peopleField is the field handler for people depicted in the media.
type peopleField struct {
	stringField
}

// PeopleField is the field handler for people depicted in the media.
var PeopleField Field = &peopleField{
	stringField{baseField{
		name:        "person",
		pluralName:  "people",
		label:       "Person",
		shortLabel:  "PE",
		multivalued: true,
	}},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *peopleField) GetValues(h filefmt.FileHandler) []interface{} {
	return stringSliceToInterfaceSlice(strmeta.GetPeople(h))
}

// GetValuesNoFaces is a special case: it returns only those people values that
// don't also have face values.  It is used by "show" when showing both the
// people and face fields, to avoid redundancy.
func (f *peopleField) GetValuesNoFaces(h filefmt.FileHandler) []interface{} {
	var people = strmeta.GetPeople(h)
	var faces = strmeta.GetFaces(h)
	var facemap = make(map[string]bool, len(faces))
	for _, face := range faces {
		facemap[face] = true
	}
	var ifcs = make([]interface{}, 0, len(people))
	for _, person := range people {
		if !facemap[person] {
			ifcs = append(ifcs, person)
		}
	}
	return ifcs
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *peopleField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	var tags, values = strmeta.GetPersonTags(h)
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// CheckValues returns whether the values of the field in the target are tagged
// correctly, and are consistent with the values of the field in the reference.
func (f *peopleField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) (res strmeta.CheckResult) {
	if res = strmeta.CheckPeople(ref, tgt); res == strmeta.ChkPresent {
		res = strmeta.CheckResult(len(f.GetValues(ref)))
	}
	return res
}

// SetValues sets all of the values of the field.
func (f *peopleField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	var people = make([]string, len(v))
	for i := range v {
		people[i] = v[i].(string)
	}
	return strmeta.SetPeople(h, people)
}
