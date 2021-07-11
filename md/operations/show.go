package operations

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/md/fields"
	"github.com/rothskeller/photo-tools/metadata"
)

// Show prints the canonical values of one or more fields in a table.
func Show(args []string, files []MediaFile) (err error) {
	var (
		fieldlist []fields.Field
		hasFaces  bool
		tw        *tabwriter.Writer
	)
	if fieldlist, err = parseFieldList("show", args); err != nil {
		return err
	}
	if len(fieldlist) == 0 {
		fieldlist = []fields.Field{
			fields.TitleField,
			fields.DateTimeField,
			fields.ArtistField,
			fields.GPSField,
			fields.LocationField,
			fields.PlacesField,
			fields.PeopleField,
			fields.FacesField,
			fields.GroupsField,
			fields.TopicsField,
			fields.KeywordsField,
			fields.CaptionField,
		}
	}
	for _, field := range fieldlist {
		if field == fields.FacesField {
			hasFaces = true
		}
	}
	tw = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\t  FIELD\tVALUE")
	for _, file := range files {
		for _, field := range fieldlist {
			var values []interface{}
			if field == fields.PeopleField && hasFaces {
				// Special case: don't include the same names as Person and Face.
				field := field.(interface {
					GetValuesNoFaces(metadata.Provider) []interface{}
				})
				values = field.GetValuesNoFaces(file.Provider)
			} else {
				values = field.GetValues(file.Provider)
			}
			check := checkField(file.Provider, field, false)
			if len(values) == 0 && check != "  " {
				fmt.Fprintf(tw, "%s\t%s%s\t\n", file.Path, check, field.Label())
			} else {
				for _, value := range values {
					fmt.Fprintf(tw, "%s\t%s%s\t%s\n", file.Path, check, field.Label(), escapeString(field.RenderValue(value)))
				}
			}
		}
	}
	tw.Flush()
	return nil
}
