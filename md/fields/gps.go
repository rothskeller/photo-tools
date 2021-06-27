package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/strmeta"
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
	},
}

// ParseValue parses a string and returns a value for the field.  It
// returns an error if the string is invalid.
func (f *gpsField) ParseValue(s string) (interface{}, error) {
	var gps metadata.GPSCoords
	if err := gps.Parse(s); err != nil {
		return nil, err
	}
	return &gps, nil
}

// RenderValue takes a value for the field and renders it in string form
// for display.
func (f *gpsField) RenderValue(v interface{}) string {
	return v.(*metadata.GPSCoords).String()
}

// EqualValue compares two values for equality.  (This is only called for
// multivalued fields.)
func (f *gpsField) EqualValue(a interface{}, b interface{}) bool {
	panic("should not be called")
}

// GetValues returns all of the values of the field.  (For single-valued
// fields, the return slice will have at most one entry.)  Empty values
// should not be included.
func (f *gpsField) GetValues(h filefmt.FileHandler) []interface{} {
	if gps := strmeta.GetGPSCoords(h); !gps.Empty() {
		return []interface{}{&gps}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond
// to the field in its first return slice, and a parallel slice of the
// values of those tags (which may be zero values).
func (f *gpsField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	if tags, values := strmeta.GetGPSCoordsTags(h); len(tags) != 0 {
		var ivals = make([]interface{}, len(values))
		for i := range values {
			ivals[i] = &values[i]
		}
		return tags, ivals
	}
	return nil, nil
}

// SetValues sets all of the values of the field.
func (f *gpsField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	switch len(v) {
	case 0:
		return strmeta.SetGPSCoords(h, metadata.GPSCoords{})
	case 1:
		return strmeta.SetGPSCoords(h, *v[0].(*metadata.GPSCoords))
	default:
		return errors.New("gps cannot have multiple values")
	}
}

// CheckValues returns whether the values of the field in the target are
// tagged correctly, and are consistent with the values of the field in
// the reference.
func (f *gpsField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) strmeta.CheckResult {
	panic("not implemented") // TODO: Implement
}
