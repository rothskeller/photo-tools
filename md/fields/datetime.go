package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/strmeta"
)

type datetimeField struct {
	baseField
}

// DateTimeField is the field handler for the date/time field, which records
// the time at which the media was originally captured.
var DateTimeField Field = &datetimeField{
	baseField{
		name:       "datetime",
		pluralName: "datetime",
		label:      "Date/Time",
		shortLabel: "DT",
	},
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *datetimeField) ParseValue(s string) (interface{}, error) {
	var dt metadata.DateTime
	if err := dt.Parse(s); err != nil {
		return nil, err
	}
	return &dt, nil
}

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *datetimeField) RenderValue(v interface{}) string {
	return v.(*metadata.DateTime).String()
}

// EqualValue compares two values for equality.  (This is only called for
// multivalued fields.)
func (f *datetimeField) EqualValue(a interface{}, b interface{}) bool {
	panic("should not be called")
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *datetimeField) GetValues(h filefmt.FileHandler) []interface{} {
	if datetime := strmeta.GetDateTime(h); !datetime.Empty() {
		return []interface{}{&datetime}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *datetimeField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	if tags, values := strmeta.GetDateTimeTags(h); len(tags) != 0 {
		var ivals = make([]interface{}, len(values))
		for i := range values {
			ivals[i] = &values[i]
		}
		return tags, ivals
	}
	return nil, nil
}

// SetValues sets all of the values of the field.
func (f *datetimeField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	switch len(v) {
	case 0:
		return strmeta.SetDateTime(h, metadata.DateTime{})
	case 1:
		return strmeta.SetDateTime(h, *v[0].(*metadata.DateTime))
	default:
		return errors.New("datetime cannot have multiple values")
	}
}

// CheckValues returns whether the values of the field in the target are tagged
// correctly, and are consistent with the values of the field in the reference.
func (f *datetimeField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) strmeta.CheckResult {
	return strmeta.CheckDateTime(ref, tgt)
}
