package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata"
)

type titleField struct {
	stringField
}

// TitleField is the field handler for the title field, which contains a short
// one-liner title for the media, in title case.
var TitleField Field = &titleField{
	stringField{
		baseField{
			name:       "title",
			pluralName: "title",
			label:      "Title",
			shortLabel: " T",
		},
	},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *titleField) GetValues(p metadata.Provider) []interface{} {
	if value := p.Title(); value != "" {
		return []interface{}{value}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *titleField) GetTags(p metadata.Provider) ([]string, []interface{}) {
	tags, values := p.TitleTags()
	return tags, stringSliceToInterfaceSlice(values)
	return nil, nil
}

// SetValues sets all of the values of the field.
func (f *titleField) SetValues(p metadata.Provider, v []interface{}) error {
	switch len(v) {
	case 0:
		return p.SetTitle("")
	case 1:
		return p.SetTitle(v[0].(string))
	default:
		return errors.New("title cannot have multiple values")
	}
}
