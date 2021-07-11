package fields

import "github.com/rothskeller/photo-tools/metadata"

// hierValueField provides methods that all hierarchical value fields have in
// common.
type hierValueField struct {
	baseField
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *hierValueField) ParseValue(s string) (interface{}, error) {
	hv, err := metadata.ParseHierValue(s)
	if err != nil {
		return nil, err
	}
	return hv, nil
}

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *hierValueField) RenderValue(v interface{}) string { return v.(metadata.HierValue).String() }

// EmptyValue returns whether a value for the field is empty.
func (f *hierValueField) EmptyValue(v interface{}) bool { return len(v.(metadata.HierValue)) == 0 }

// EqualValue compares two values for equality.
func (f *hierValueField) EqualValue(a interface{}, b interface{}) bool {
	return a.(metadata.HierValue).Equal(b.(metadata.HierValue))
}
