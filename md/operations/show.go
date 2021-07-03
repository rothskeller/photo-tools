package operations

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/md/fields"
)

func newShowOp() Operation { return new(showOp) }

// showOp prints the canonical values of one or more fields in a table.
type showOp struct {
	fieldListOp
	hasRun   bool
	hasFaces bool
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
			fields.FacesField,
			fields.GroupsField,
			fields.TopicsField,
			fields.KeywordsField,
			fields.CaptionField,
		}
	}
	for _, field := range op.fields {
		if field == fields.FacesField {
			op.hasFaces = true
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
	fmt.Fprintln(tw, "FILE\t  FIELD\tVALUE")
	for _, file := range files {
		for _, field := range op.fields {
			var values []interface{}
			if field == fields.PeopleField && op.hasFaces {
				// Special case: don't include the same names as Person and Face.
				field := field.(interface {
					GetValuesNoFaces(filefmt.FileHandler) []interface{}
				})
				values = field.GetValuesNoFaces(file.Handler)
			} else {
				values = field.GetValues(file.Handler)
			}
			check := field.CheckValues(file.Handler, file.Handler)
			if len(values) == 0 && check < 0 {
				fmt.Fprintf(tw, "%s%s%s\t\n", file.Path, resultCodes[check], field.Label())
			} else {
				for _, value := range values {
					if check < 0 {
						fmt.Fprintf(tw, "%s%s%s\t%s\n", file.Path, resultCodes[check], field.Label(), escapeString(field.RenderValue(value)))
					} else {
						fmt.Fprintf(tw, "%s\t  %s\t%s\n", file.Path, field.Label(), escapeString(field.RenderValue(value)))
					}
				}
			}
		}
	}
	tw.Flush()
	op.hasRun = true
	return nil
}
