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
	out       *tabwriter.Writer
	lastCount int
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
			fields.FacesField,
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
	if op.out == nil {
		op.out = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprint(op.out, "FILE")
		for _, field := range op.fields {
			fmt.Fprintf(op.out, "\t%s", field.ShortLabel())
		}
		fmt.Fprintln(op.out)
	} else if op.lastCount > 1 || len(files) > 1 {
		for range op.fields {
			fmt.Fprint(op.out, "\t")
		}
		fmt.Fprintln(op.out)
	}
	op.lastCount = len(files)
	for _, file := range files {
		fmt.Fprint(op.out, file.Path)
		for _, field := range op.fields {
			result := field.CheckValues(files[0].Handler, file.Handler)
			if result <= 0 || (result == strmeta.ChkPresent && !field.Multivalued()) {
				fmt.Fprint(op.out, resultCodes[result])
			} else {
				fmt.Fprintf(op.out, "\t%2d", result)
			}
		}
		fmt.Fprintln(op.out)
	}
	return nil
}

// Finish finishes the operation after all batches have been processed.
func (op *checkOp) Finish() error {
	op.out.Flush()
	return nil
}

var resultCodes = map[strmeta.CheckResult]string{
	strmeta.ChkConflictingValues: "\t!=",
	strmeta.ChkExpectedAbsent:    "\t--",
	strmeta.ChkIncorrectlyTagged: "\t[]",
	strmeta.ChkOptionalAbsent:    "\t  ",
	strmeta.ChkPresent:           "\t ✓",
}
