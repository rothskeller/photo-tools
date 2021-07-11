package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata"
)

type locationField struct {
	baseField
}

// LocationField is the field handler for the location field, which gives a
// textual description of the location where the media was captured.
var LocationField Field = &locationField{
	baseField{
		name:       "location",
		pluralName: "location",
		label:      "Location",
		shortLabel: " L",
	},
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *locationField) ParseValue(s string) (interface{}, error) {
	var loc metadata.Location
	if err := loc.Parse(s); err != nil {
		return nil, err
	}
	return loc, nil
}

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *locationField) RenderValue(v interface{}) string {
	return v.(metadata.Location).String()
}

// EmptyValue returns whether a value for the field is empty.
func (f *locationField) EmptyValue(v interface{}) bool { return v.(metadata.Location).Empty() }

// EqualValue compares two values for equality.
func (f *locationField) EqualValue(a interface{}, b interface{}) bool {
	return a.(metadata.Location).Equal(b.(metadata.Location))
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *locationField) GetValues(p metadata.Provider) []interface{} {
	if location := p.Location(); !location.Empty() {
		return []interface{}{location}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *locationField) GetTags(p metadata.Provider) ([]string, [][]interface{}) {
	tags, values := p.LocationTags()
	var ivals = make([][]interface{}, len(values))
	for i := range values {
		ivals[i] = make([]interface{}, len(values[i]))
		for j := range values[i] {
			ivals[i][j] = values[i][j]
		}
	}
	return tags, ivals
}

// SetValues sets all of the values of the field.
func (f *locationField) SetValues(p metadata.Provider, v []interface{}) error {
	switch len(v) {
	case 0:
		return p.SetLocation(metadata.Location{})
	case 1:
		return p.SetLocation(v[0].(metadata.Location))
	default:
		return errors.New("location cannot have multiple values")
	}
}
