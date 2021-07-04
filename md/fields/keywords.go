package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

type keywordsField struct {
	baseField
}

// KeywordsField is the field handler for keywords.
var KeywordsField Field = &keywordsField{
	baseField{
		name:        "keyword",
		pluralName:  "keywords",
		label:       "Keyword",
		shortLabel:  "KW",
		multivalued: true,
	},
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *keywordsField) ParseValue(s string) (interface{}, error) {
	var v strmeta.Keyword
	if err := v.Parse(s); err != nil {
		return nil, err
	}
	return v, nil
}

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *keywordsField) RenderValue(v interface{}) string { return v.(strmeta.Keyword).String() }

// EqualValue compares two values for equality.
func (f *keywordsField) EqualValue(a interface{}, b interface{}) bool {
	return a.(strmeta.Keyword).Equal(b.(strmeta.Keyword))
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *keywordsField) GetValues(h filefmt.FileHandler) []interface{} {
	var keywords = strmeta.GetKeywords(h)
	var ifcs = make([]interface{}, len(keywords))
	for i := range keywords {
		ifcs[i] = keywords[i]
	}
	return ifcs
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *keywordsField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	var tags, values = strmeta.GetKeywordTags(h)
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// CheckValues returns whether the values of the field are tagged correctly.
func (f *keywordsField) CheckValues(h filefmt.FileHandler) (res strmeta.CheckResult) {
	if res = strmeta.CheckKeywords(h); res == strmeta.ChkPresent {
		res = strmeta.CheckResult(len(f.GetValues(h)))
	}
	return res
}

// SetValues sets all of the values of the field.
func (f *keywordsField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	var keywords = make([]strmeta.Keyword, len(v))
	for i := range v {
		keywords[i] = v[i].(strmeta.Keyword)
	}
	return strmeta.SetKeywords(h, keywords)
}
