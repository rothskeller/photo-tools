package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

type captionField struct {
	stringField
}

// CaptionField is the field handler for the caption field, which contains an
// optional prose description of the media.
var CaptionField Field = &captionField{
	stringField{
		baseField{
			name:       "caption",
			pluralName: "caption",
			label:      "Caption",
			shortLabel: " C",
		},
	},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *captionField) GetValues(h filefmt.FileHandler) []interface{} {
	if caption := strmeta.GetCaption(h); caption != "" {
		return []interface{}{caption}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *captionField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	if tags, values := strmeta.GetCaptionTags(h); len(tags) != 0 {
		return tags, stringSliceToInterfaceSlice(values)
	}
	return nil, nil
}

// SetValues sets all of the values of the field.
func (f *captionField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	switch len(v) {
	case 0:
		return strmeta.SetCaption(h, "")
	case 1:
		return strmeta.SetCaption(h, v[0].(string))
	default:
		return errors.New("caption cannot have multiple values")
	}
}

// CheckValues returns whether the values of the field in the target are tagged
// correctly, and are consistent with the values of the field in the reference.
func (f *captionField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) strmeta.CheckResult {
	return strmeta.CheckCaption(ref, tgt)
}
