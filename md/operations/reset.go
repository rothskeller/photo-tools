package operations

import (
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
	"github.com/rothskeller/photo-tools/metadata"
)

// Reset resets the values of one or more fields to their primary value(s), thus
// clearing up any inconsistencies or tagging errors.
func Reset(args []string, files []MediaFile) (err error) {
	var fieldlist []fields.Field

	if fieldlist, err = parseFieldList("reset", args); err != nil {
		return err
	}
	if len(fieldlist) == 0 {
		fieldlist = []fields.Field{
			fields.ArtistField,
			fields.CaptionField,
			fields.DateTimeField,
			fields.FacesField,
			fields.GPSField,
			fields.GroupsField,
			fields.KeywordsField,
			fields.LocationField,
			fields.PeopleField,
			fields.PlacesField,
			fields.TitleField,
			fields.TopicsField,
		}
	}
	for i, file := range files {
		for _, field := range fieldlist {
			values := field.GetValues(file.Provider)
			if err := field.SetValues(file.Provider, values); err == metadata.ErrNotSupported {
				continue
			} else if err != nil {
				return fmt.Errorf("%s: reset %s: %s", file.Path, field.PluralName(), err)
			}
			files[i].Changed = true
		}
	}
	return nil
	// NOTE: this doesn't actually correct all tagging errors.  If the wrong
	// tags are set, or tags are set to the wrong value, that will be fixed.
	// But if the right tags are set, but their encoding wasn't correct, the
	// encoding won't be fixed because the underlying file handler won't see
	// a change to the field value and won't rewrite the tag.  I can live
	// with that.
}
