package operations

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/md/fields"
	"github.com/rothskeller/photo-tools/strmeta"
)

func newCheckOp() Operation { return new(checkOp) }

// checkOp displays a table giving the tagging correctness and consistency of
// each field.
type checkOp struct {
	fieldListOp
	hasRun bool
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *checkOp) parseArgs(args []string) (remainingArgs []string, err error) {
	remainingArgs, _ = op.fieldListOp.parseArgs(args)
	if len(op.fields) == 0 {
		op.fields = []fields.Field{
			fields.ArtistField,
			fields.DateTimeField,
			fields.GPSField,
			fields.PlacesField,
			fields.PeopleField,
			fields.GroupsField,
			fields.TopicsField,
			fields.TitleField,
			fields.CaptionField,
			fields.KeywordsField,
			fields.LocationField,
		}
	}
	return remainingArgs, nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *checkOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *checkOp) Run(files []MediaFile) error {
	if op.hasRun { // put a newline between batches for readability
		fmt.Println()
	}
	out := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprint(out, "FILE")
	for _, field := range op.fields {
		fmt.Fprintf(out, "\t%s", field.ShortLabel())
	}
	fmt.Fprintln(out)
	for _, file := range files {
		fmt.Fprint(out, file.Path)
		for _, field := range op.fields {
			result := field.CheckValues(files[0].Handler, file.Handler)
			if result <= 0 || (result == strmeta.ChkPresent && !field.Multivalued()) {
				fmt.Fprint(out, resultCodes[result])
			} else {
				fmt.Fprintf(out, "\t%2d", result)
			}
		}
		fmt.Fprintln(out)
	}
	out.Flush()
	op.hasRun = true
	return nil
}

var resultCodes = map[strmeta.CheckResult]string{
	strmeta.ChkConflictingValues: "\t!=",
	strmeta.ChkExpectedAbsent:    "\t--",
	strmeta.ChkIncorrectlyTagged: "\t[]",
	strmeta.ChkOptionalAbsent:    "\t  ",
	strmeta.ChkPresent:           "\t ✓",
}
