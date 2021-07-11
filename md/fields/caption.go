package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata"
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
func (f *captionField) GetValues(p metadata.Provider) []interface{} {
	if value := p.Caption(); value != "" {
		return []interface{}{value}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *captionField) GetTags(p metadata.Provider) ([]string, [][]interface{}) {
	tags, values := p.CaptionTags()
	ilist := make([][]interface{}, len(values))
	for i := range values {
		ilist[i] = stringSliceToInterfaceSlice(values[i])
	}
	return tags, ilist
}

// SetValues sets all of the values of the field.
func (f *captionField) SetValues(p metadata.Provider, v []interface{}) error {
	switch len(v) {
	case 0:
		return p.SetCaption("")
	case 1:
		return p.SetCaption(v[0].(string))
	default:
		return errors.New("caption cannot have multiple values")
	}
}
