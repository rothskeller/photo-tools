package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

type placesField struct {
	baseField
}

// PlacesField is the field handler for place where media files were captured or
// places depicted in media files.
var PlacesField Field = &placesField{
	baseField{
		name:        "place",
		pluralName:  "places",
		label:       "Place",
		shortLabel:  "PL",
		multivalued: true,
	},
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *placesField) ParseValue(s string) (interface{}, error) {
	var v strmeta.Place
	if err := v.Parse(s); err != nil {
		return nil, err
	}
	return v, nil
}

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *placesField) RenderValue(v interface{}) string { return v.(strmeta.Place).String() }

// EqualValue compares two values for equality.
func (f *placesField) EqualValue(a interface{}, b interface{}) bool {
	return a.(strmeta.Place).Equal(b.(strmeta.Place))
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *placesField) GetValues(h filefmt.FileHandler) []interface{} {
	var places = strmeta.GetPlaces(h)
	var ifcs = make([]interface{}, len(places))
	for i := range places {
		ifcs[i] = places[i]
	}
	return ifcs
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *placesField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	var tags, values = strmeta.GetPlaceTags(h)
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// CheckValues returns whether the values of the field are tagged correctly.
func (f *placesField) CheckValues(h filefmt.FileHandler) (res strmeta.CheckResult) {
	if res = strmeta.CheckPlaces(h); res == strmeta.ChkPresent {
		res = strmeta.CheckResult(len(f.GetValues(h)))
	}
	return res
}

// SetValues sets all of the values of the field.
func (f *placesField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	var places = make([]strmeta.Place, len(v))
	for i := range v {
		places[i] = v[i].(strmeta.Place)
	}
	return strmeta.SetPlaces(h, places)
}
