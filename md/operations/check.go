package operations

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/md/fields"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/strmeta"
)

var checkFields = []fields.Field{
	fields.ArtistField,
	fields.DateTimeField,
	fields.GPSField,
	fields.PlacesField,
	fields.PeopleField,
	fields.FacesField,
	fields.GroupsField,
	fields.TopicsField,
	fields.TitleField,
	fields.CaptionField,
	fields.KeywordsField,
	fields.LocationField,
}

// Check displays a table giving the tagging correctness of each field.
func Check(args []string, files []MediaFile) (err error) {
	var out *tabwriter.Writer

	if len(args) != 0 {
		return errors.New("check: excess arguments")
	}
	out = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprint(out, "FILE")
	for _, field := range checkFields {
		fmt.Fprintf(out, "\t%s", field.ShortLabel())
	}
	fmt.Fprintln(out)
	for _, file := range files {
		fmt.Fprint(out, file.Path)
		for _, field := range checkFields {
			result := checkField(file.Provider, field)
			fmt.Fprintf(out, "\t%s", result)
		}
		fmt.Fprintln(out)
	}
	out.Flush()
	return nil
}

func checkField(p metadata.Provider, field fields.Field) string {
	var canon = field.GetValues(p)
	// TODO PROBLEM comparing the values return by Tags works fine for
	// single-valued fields, but it isn't going to detect all problems with
	// multi-valued fields.  How do we implement check in a provider model?
}

var resultCodes = map[strmeta.CheckResult]string{
	strmeta.ChkConflictingValues: "\t!=",
	strmeta.ChkExpectedAbsent:    "\t--",
	strmeta.ChkIncorrectlyTagged: "\t[]",
	strmeta.ChkOptionalAbsent:    "\t  ",
	strmeta.ChkPresent:           "\t ✓",
}
