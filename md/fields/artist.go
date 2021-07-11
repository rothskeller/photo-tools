package fields

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata"
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
			expected:   true,
		},
	},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *artistField) GetValues(p metadata.Provider) []interface{} {
	if value := p.Creator(); value != "" {
		return []interface{}{value}
	}
	return nil
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *artistField) GetTags(p metadata.Provider) ([]string, [][]interface{}) {
	tags, values := p.CreatorTags()
	ilist := make([][]interface{}, len(values))
	for i := range values {
		ilist[i] = stringSliceToInterfaceSlice(values[i])
	}
	return tags, ilist
}

// SetValues sets all of the values of the field.
func (f *artistField) SetValues(p metadata.Provider, v []interface{}) error {
	switch len(v) {
	case 0:
		return p.SetCreator("")
	case 1:
		return p.SetCreator(v[0].(string))
	default:
		return errors.New("creator cannot have multiple values")
	}
}
