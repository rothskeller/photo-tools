package operations

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newShowOp() Operation { return new(showOp) }

type showOp struct {
	fieldListOp
	hasRun bool
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *showOp) parseArgs(args []string) (remainingArgs []string, err error) {
	remainingArgs, _ = op.fieldListOp.parseArgs(args)
	if len(op.fields) == 0 {
		op.fields = []fields.Field{
			fields.TitleField,
			fields.DateTimeField,
			fields.ArtistField,
			fields.GPSField,
			fields.LocationField,
			fields.PlacesField,
			fields.PeopleField,
			fields.GroupsField,
			fields.TopicsField,
			fields.OtherKeywordsField,
			fields.CaptionField,
		}
	}
	return remainingArgs, nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *showOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *showOp) Run(files []MediaFile) error {
	if op.hasRun {
		fmt.Println()
	}
	var tw = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tFIELD\tVALUE")
	for _, file := range files {
		for _, field := range op.fields {
			values := field.GetValues(file.Handler)
			if len(values) == 0 {
				fmt.Fprintf(tw, "%s\t%s\t\n", file.Path, field.Label())
			} else {
				for _, value := range values {
					fmt.Fprintf(tw, "%s\t%s\t%s\n", file.Path, field.Label(), escapeString(field.RenderValue(value)))
				}
			}
		}
	}
	tw.Flush()
	op.hasRun = true
	return nil
}
