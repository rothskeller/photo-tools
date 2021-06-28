package operations

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newTagsOp() Operation { return new(tagsOp) }

// tagsOp prints all of the tagged value of one or more fields in a table.
type tagsOp struct {
	fieldListOp
	hasRun bool
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *tagsOp) parseArgs(args []string) (remainingArgs []string, err error) {
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
			fields.KeywordsField,
			fields.CaptionField,
		}
	}
	return remainingArgs, nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *tagsOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *tagsOp) Run(files []MediaFile) error {
	if op.hasRun {
		fmt.Println()
	}
	var tw = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tTAG\tVALUE")
	for _, file := range files {
		for _, field := range op.fields {
			tagNames, tagValues := field.GetTags(file.Handler)
			for i, tag := range tagNames {
				fmt.Fprintf(tw, "%s\t%s\t%s\n", file.Path, tag, escapeString(field.RenderValue(tagValues[i])))
			}
		}
	}
	tw.Flush()
	op.hasRun = true
	return nil
}
