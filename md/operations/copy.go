package operations

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

// Copy copies values of the specified fields from the first file to all other
// target files.
func Copy(args []string, files []MediaFile) (err error) {
	var fieldlist []fields.Field

	if fieldlist, err = parseFieldList("copy", args); err != nil {
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
	if len(files) < 2 {
		return errors.New("copy: must list at least two files")
	}
	var values = make([][]interface{}, len(fieldlist))
	for idx, field := range fieldlist {
		values[idx] = field.GetValues(files[0].Handler)
	}
	for i, file := range files[1:] {
		for idx, field := range fieldlist {
			if err := field.SetValues(file.Handler, values[idx]); err != nil {
				return fmt.Errorf("%s: copy %s: %s", file.Path, field.PluralName(), err)
			}
		}
		files[i+1].Changed = true
	}
	return nil
}
