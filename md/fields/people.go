package fields

import (
	"github.com/rothskeller/photo-tools/metadata"
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
func (f *peopleField) GetValues(p metadata.Provider) []interface{} {
	return stringSliceToInterfaceSlice(p.People())
}

// GetValuesNoFaces is a special case: it returns only those people values that
// don't also have face values.  It is used by "show" when showing both the
// people and face fields, to avoid redundancy.
func (f *peopleField) GetValuesNoFaces(p metadata.Provider) []interface{} {
	var people = p.People()
	var faces = p.Faces()
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
func (f *peopleField) GetTags(p metadata.Provider) ([]string, [][]interface{}) {
	tags, values := p.PeopleTags()
	ivals := make([][]interface{}, len(values))
	for i := range values {
		ivals[i] = stringSliceToInterfaceSlice(values[i])
	}
	return tags, ivals
}

// SetValues sets all of the values of the field.
func (f *peopleField) SetValues(p metadata.Provider, v []interface{}) error {
	values := make([]string, len(v))
	for i := range v {
		values[i] = v[i].(string)
	}
	return p.SetPeople(values)
}
