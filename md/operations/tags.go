package operations

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/md/fields"
)

// Tags prints all of the tagged values of one or more fields in a table.
func Tags(args []string, files []MediaFile) (err error) {
	var (
		fieldlist []fields.Field
		tw        *tabwriter.Writer
	)
	if fieldlist, err = parseFieldList("tags", args); err != nil {
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
	tw = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tTAG\tVALUE")
	for _, file := range files {
		for _, field := range fieldlist {
			tagNames, tagValues := field.GetTags(file.Provider)
			for i, tag := range tagNames {
				fmt.Fprintf(tw, "%s\t%s\t%s\n", file.Path, tag, escapeString(field.RenderValue(tagValues[i])))
			}
		}
	}
	tw.Flush()
	return nil
}
