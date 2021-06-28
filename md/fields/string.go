package fields

// stringField provides methods that all string fields have in common.
type stringField struct {
	baseField
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *stringField) ParseValue(s string) (interface{}, error) { return s, nil }

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *stringField) RenderValue(v interface{}) string { return v.(string) }

// EqualValue compares two values for equality.
func (f *stringField) EqualValue(a interface{}, b interface{}) bool { return a.(string) == b.(string) }

// stringSliceToInterfaceSlice is a utility function used by string-valued
// fields.
func stringSliceToInterfaceSlice(ss []string) (is []interface{}) {
	is = make([]interface{}, len(ss))
	for i, s := range ss {
		is[i] = s
	}
	return is
}
