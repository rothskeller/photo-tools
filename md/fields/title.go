package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
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

// GetValues returns all of the values of the field.  (For single-valued
// fields, the return slice will have at most one entry.)  Empty values
// should not be included.
func (f *titleField) GetValues(h filefmt.FileHandler) []interface{} {
	if title := strmeta.GetTitle(h); title != "" {
		return []interface{}{title}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond
// to the field in its first return slice, and a parallel slice of the
// values of those tags (which may be zero values).
func (f *titleField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	if tags, values := strmeta.GetTitleTags(h); len(tags) != 0 {
		return tags, stringSliceToInterfaceSlice(values)
	}
	return nil, nil
}

// SetValues sets all of the values of the field.
func (f *titleField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	switch len(v) {
	case 0:
		return strmeta.SetTitle(h, "")
	case 1:
		return strmeta.SetTitle(h, v[0].(string))
	default:
		return errors.New("title cannot have multiple values")
	}
}

// CheckValues returns whether the values of the field in the target are
// tagged correctly, and are consistent with the values of the field in
// the reference.
func (f *titleField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) strmeta.CheckResult {
	return strmeta.CheckTitle(ref, tgt)
}
