package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/strmeta"
)

// baseKWField is an abstract base class for keyword fields, providing routines
// that they all share.
type baseKWField struct {
	baseField
	prefix   string
	expected bool
}

// ParseValue parses a string and returns a value for the field.  It
// returns an error if the string is invalid.
func (f *baseKWField) ParseValue(s string) (interface{}, error) {
	return metadata.ParseKeyword(s, f.prefix)
}

// RenderValue takes a value for the field and renders it in string form
// for display.
func (f *baseKWField) RenderValue(v interface{}) string {
	if f.prefix != "" {
		return v.(metadata.Keyword).StringWithoutPrefix(f.prefix)
	}
	return v.(metadata.Keyword).String()
}

// EqualValue compares two values for equality.
func (f *baseKWField) EqualValue(a interface{}, b interface{}) bool {
	return a.(metadata.Keyword).Equal(b.(metadata.Keyword))
}

// GetValues returns all of the values of the field.  (For single-valued
// fields, the return slice will have at most one entry.)  Empty values
// should not be included.
func (f *baseKWField) GetValues(h filefmt.FileHandler) []interface{} {
	var kws = strmeta.GetKeywords(h)
	if f.prefix != "" {
		var filtered []metadata.Keyword
		for _, kw := range kws {
			if len(kw) != 0 && kw[0] == f.prefix {
				filtered = append(filtered, kw)
			}
		}
		kws = filtered
	}
	if len(kws) == 0 {
		return nil
	}
	var ifcs = make([]interface{}, len(kws))
	for i := range kws {
		ifcs[i] = kws[i]
	}
	return ifcs
}

// GetTags returns the names of all of the metadata tags that correspond
// to the field in its first return slice, and a parallel slice of the
// values of those tags (which may be zero values).
func (f *baseKWField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	var tags, values = strmeta.GetKeywordsTags(h)
	if f.prefix != "" {
		j := 0
		for i, kw := range values {
			if len(kw) != 0 && kw[0] == f.prefix {
				tags[j] = tags[i]
				values[j] = kw
				j++
			}
		}
		tags = tags[:j]
		values = values[:j]
	}
	if len(values) == 0 {
		return nil, nil
	}
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// SetValues sets all of the values of the field.
func (f *baseKWField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	var all []metadata.Keyword
	if f.prefix == "" {
		all = make([]metadata.Keyword, len(v))
		for i := range v {
			all[i] = v[i].(metadata.Keyword)
		}
	} else {
		j := 0
		for _, kw := range all {
			if len(kw) != 0 && kw[0] != f.prefix {
				all[j] = kw
				j++
			}
		}
		all = all[:j]
		for _, kw := range v {
			all = append(all, kw.(metadata.Keyword))
		}
	}
	return strmeta.SetKeywords(h, all)
}

// CheckValues returns whether the values of the field in the target are
// tagged correctly, and are consistent with the values of the field in
// the reference.
func (f *baseKWField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) strmeta.CheckResult {
	if result := strmeta.CheckKeywords(ref, tgt); result < 0 {
		return result
	}
	if count := len(f.GetValues(tgt)); count > 0 {
		return strmeta.CheckResult(count)
	} else if f.expected {
		return strmeta.ChkExpectedAbsent
	}
	return strmeta.ChkOptionalAbsent
}
