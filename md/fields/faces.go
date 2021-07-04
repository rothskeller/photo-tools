package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
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
func (f *facesField) GetValues(h filefmt.FileHandler) []interface{} {
	return stringSliceToInterfaceSlice(strmeta.GetFaces(h))
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *facesField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	var tags, values = strmeta.GetFaceTags(h)
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// CheckValues returns whether the values of the field are tagged correctly.
func (f *facesField) CheckValues(h filefmt.FileHandler) (res strmeta.CheckResult) {
	if res = strmeta.CheckFaces(h); res == strmeta.ChkPresent {
		res = strmeta.CheckResult(len(f.GetValues(h)))
	}
	return res
}

// SetValues sets all of the values of the field.
func (f *facesField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	var faces = make([]string, len(v))
	for i := range v {
		faces[i] = v[i].(string)
	}
	return strmeta.SetFaces(h, faces)
}
