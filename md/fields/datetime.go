package fields

import (
	"errors"
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
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
		expected:   true,
	},
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *datetimeField) ParseValue(s string) (interface{}, error) {
	var dt metadata.DateTime
	if err := dt.Parse(s); err != nil {
		return nil, err
	}
	return dt, nil
}

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *datetimeField) RenderValue(v interface{}) string {
	var str = v.(metadata.DateTime).String()
	if str == "" {
		return str
	}
	var date = str[:10]
	var time = str[11:]
	time = strings.Replace(time, "+", " +", -1)
	time = strings.Replace(time, "-", " -", -1)
	time = strings.Replace(time, "Z", " Z", -1)
	return date + " " + time
}

// EmptyValue returns whether a value for the field is empty.
func (f *datetimeField) EmptyValue(v interface{}) bool { return v.(metadata.DateTime).Empty() }

// EqualValue compares two values for equality.
func (f *datetimeField) EqualValue(a interface{}, b interface{}) bool {
	return a.(metadata.DateTime).Equivalent(b.(metadata.DateTime))
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *datetimeField) GetValues(p metadata.Provider) []interface{} {
	if datetime := p.DateTime(); !datetime.Empty() {
		return []interface{}{datetime}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *datetimeField) GetTags(p metadata.Provider) ([]string, [][]interface{}) {
	if tags, values := p.DateTimeTags(); len(tags) != 0 {
		var ivals = make([][]interface{}, len(values))
		for i := range values {
			ivals[i] = []interface{}{values[i]}
		}
		return tags, ivals
	}
	return nil, nil
}

// SetValues sets all of the values of the field.
func (f *datetimeField) SetValues(p metadata.Provider, v []interface{}) error {
	switch len(v) {
	case 0:
		return p.SetDateTime(metadata.DateTime{})
	case 1:
		return p.SetDateTime(v[0].(metadata.DateTime))
	default:
		return errors.New("datetime cannot have multiple values")
	}
}
