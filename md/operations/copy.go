package operations

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
	"github.com/rothskeller/photo-tools/metadata"
)

// Copy copies values of the specified fields from the first file to all other
// target files.
func Copy(args []string, files []MediaFile) (err error) {
	var (
		fieldlist []fields.Field
		allFields bool
	)
	if fieldlist, err = parseFieldList("copy", args); err != nil {
		return err
	}
	if len(fieldlist) == 0 {
		allFields = true
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
		values[idx] = field.GetValues(files[0].Provider)
	}
	for i, file := range files[1:] {
		for idx, field := range fieldlist {
			if err := field.SetValues(file.Provider, values[idx]); err != nil {
				if allFields && err == metadata.ErrNotSupported && len(values) == 0 {
					continue
				}
				return fmt.Errorf("%s: copy %s: %s", file.Path, field.PluralName(), err)
			}
			files[i+1].Changed = true
		}
	}
	return nil
}
