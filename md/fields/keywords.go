package fields

import (
	"github.com/rothskeller/photo-tools/metadata"
)

type keywordsField struct {
	hierValueField
}

// KeywordsField is the field handler for keywords.
var KeywordsField Field = &keywordsField{
	hierValueField{
		baseField{
			name:        "keyword",
			pluralName:  "keywords",
			label:       "Keyword",
			shortLabel:  "KW",
			multivalued: true,
		},
	},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *keywordsField) GetValues(p metadata.Provider) []interface{} {
	var groups = p.Keywords()
	var ifcs = make([]interface{}, len(groups))
	for i := range groups {
		ifcs[i] = groups[i]
	}
	return ifcs
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *keywordsField) GetTags(p metadata.Provider) ([]string, []interface{}) {
	var tags, values = p.KeywordsTags()
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// SetValues sets all of the values of the field.
func (f *keywordsField) SetValues(p metadata.Provider, v []interface{}) error {
	var values = make([]metadata.HierValue, len(v))
	for i := range v {
		values[i] = v[i].(metadata.HierValue)
	}
	return p.SetKeywords(values)
}
