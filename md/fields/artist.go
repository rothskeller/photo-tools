package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

type artistField struct {
	stringField
}

// ArtistField is the field handler for the Artist field, i.e., the person who
// originally captured the media.  (It is referred to as "Creator" in package
// strmeta, since that is a better description, but md uses "artist" so that it
// has a unique initial letter.)
var ArtistField Field = &artistField{
	stringField{
		baseField{
			name:       "artist",
			pluralName: "artist", // single-valued, so singular
			label:      "Artist",
			shortLabel: " A",
		},
	},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *artistField) GetValues(h filefmt.FileHandler) []interface{} {
	if artist := strmeta.GetCreator(h); artist != "" {
		return []interface{}{artist}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *artistField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	if tags, values := strmeta.GetCreatorTags(h); len(tags) != 0 {
		return tags, stringSliceToInterfaceSlice(values)
	}
	return nil, nil
}

// SetValues sets all of the values of the field.
func (f *artistField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	switch len(v) {
	case 0:
		return strmeta.SetCreator(h, "")
	case 1:
		return strmeta.SetCreator(h, v[0].(string))
	default:
		return errors.New("artist cannot have multiple values")
	}
}

// CheckValues returns whether the values of the field are tagged correctly.
func (f *artistField) CheckValues(h filefmt.FileHandler) strmeta.CheckResult {
	return strmeta.CheckCreator(h)
}
