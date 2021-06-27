package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
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
		label:      "Locations",
		shortLabel: " L",
	},
}

// ParseValue parses a string and returns a value for the field.  It
// returns an error if the string is invalid.
func (f *locationField) ParseValue(s string) (interface{}, error) {
	var loc strmeta.Location
	if err := loc.Parse(s); err != nil {
		return nil, err
	}
	return &loc, nil
}

// RenderValue takes a value for the field and renders it in string form
// for display.
func (f *locationField) RenderValue(v interface{}) string {
	return v.(*strmeta.Location).String()
}

// EqualValue compares two values for equality.  (This is only called for
// multivalued fields.)
func (f *locationField) EqualValue(a interface{}, b interface{}) bool {
	panic("should not be called")
}

// GetValues returns all of the values of the field.  (For single-valued
// fields, the return slice will have at most one entry.)  Empty values
// should not be included.
func (f *locationField) GetValues(h filefmt.FileHandler) []interface{} {
	if location := strmeta.GetLocation(h); !location.Empty() {
		return []interface{}{&location}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond
// to the field in its first return slice, and a parallel slice of the
// values of those tags (which may be zero values).
func (f *locationField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	if tags, values := strmeta.GetLocationTags(h); len(tags) != 0 {
		var ivals = make([]interface{}, len(values))
		for i := range values {
			ivals[i] = &values[i]
		}
		return tags, ivals
	}
	return nil, nil
}

// SetValues sets all of the values of the field.
func (f *locationField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	switch len(v) {
	case 0:
		return strmeta.SetLocation(h, strmeta.Location{})
	case 1:
		return strmeta.SetLocation(h, *v[0].(*strmeta.Location))
	default:
		return errors.New("location cannot have multiple values")
	}
}

// CheckValues returns whether the values of the field in the target are
// tagged correctly, and are consistent with the values of the field in
// the reference.
func (f *locationField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) strmeta.CheckResult {
	return strmeta.CheckLocation(ref, tgt)
}
