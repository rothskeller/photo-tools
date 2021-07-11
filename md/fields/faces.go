package fields

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// facesField is the field handler for face regions in the media.
type facesField struct {
	stringField
}

// FacesField is the field handler for face regions in the media.
var FacesField Field = &facesField{
	stringField{
		baseField{
			name:        "face",
			pluralName:  "faces",
			label:       "Face",
			shortLabel:  " F",
			multivalued: true,
		},
	},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *facesField) GetValues(p metadata.Provider) []interface{} {
	return stringSliceToInterfaceSlice(p.Faces())
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *facesField) GetTags(p metadata.Provider) ([]string, [][]interface{}) {
	tags, values := p.FacesTags()
	ivals := make([][]interface{}, len(values))
	for i := range values {
		ivals[i] = stringSliceToInterfaceSlice(values[i])
	}
	return tags, ivals
}

// SetValues sets all of the values of the field.
func (f *facesField) SetValues(p metadata.Provider, v []interface{}) error {
	values := make([]string, len(v))
	for i := range v {
		values[i] = v[i].(string)
	}
	return p.SetFaces(values)
}
