package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata"
)

type gpsField struct {
	baseField
}

// GPSField is the field handler for the GPS coordinates field, which expresses
// the location where the media was captured (or, failing that, the location
// where it should be shown on a map).
var GPSField Field = &gpsField{
	baseField{
		name:       "gps",
		pluralName: "gps",
		label:      "GPS Coords",
		shortLabel: " G",
		expected:   true,
	},
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *gpsField) ParseValue(s string) (interface{}, error) {
	var gps metadata.GPSCoords
	if err := gps.Parse(s); err != nil {
		return nil, err
	}
	return gps, nil
}

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *gpsField) RenderValue(v interface{}) string {
	return v.(metadata.GPSCoords).String()
}

// EmptyValue returns whether a value for the field is empty.
func (f *gpsField) EmptyValue(v interface{}) bool { return v.(metadata.GPSCoords).Empty() }

// EqualValue compares two values for equality.
func (f *gpsField) EqualValue(a interface{}, b interface{}) bool {
	return a.(metadata.GPSCoords).Equivalent(b.(metadata.GPSCoords))
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *gpsField) GetValues(p metadata.Provider) []interface{} {
	if gps := p.GPS(); !gps.Empty() {
		return []interface{}{gps}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *gpsField) GetTags(p metadata.Provider) ([]string, [][]interface{}) {
	if tags, values := p.GPSTags(); len(tags) != 0 {
		var ivals = make([][]interface{}, len(values))
		for i := range values {
			ivals[i] = []interface{}{values[i]}
		}
		return tags, ivals
	}
	return nil, nil
}

// SetValues sets all of the values of the field.
func (f *gpsField) SetValues(p metadata.Provider, v []interface{}) error {
	switch len(v) {
	case 0:
		return p.SetGPS(metadata.GPSCoords{})
	case 1:
		return p.SetGPS(v[0].(metadata.GPSCoords))
	default:
		return errors.New("gps cannot have multiple values")
	}
}
