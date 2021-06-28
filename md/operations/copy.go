package operations

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newCopyOp() Operation { return new(copyOp) }

// copyOp copies values of the specified fields from the first file to all other
// files.
type copyOp struct {
	fieldListOp
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *copyOp) parseArgs(args []string) (remainingArgs []string, err error) {
	remainingArgs, _ = op.fieldListOp.parseArgs(args)
	if len(op.fields) == 0 {
		op.fields = []fields.Field{
			fields.ArtistField,
			fields.CaptionField,
			fields.DateTimeField,
			fields.GPSField,
			fields.KeywordsField,
			fields.LocationField,
			fields.TitleField,
		}
	}
	return remainingArgs, nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *copyOp) Check(batches [][]MediaFile) error {
	for _, batch := range batches {
		if len(batch) == 1 {
			if len(batches) != 1 {
				return errors.New("copy: must list at least two files in each batch")
			}
			return errors.New("copy: must list at least two files")
		}
	}
	return nil
}

// Run executes the operation against the listed media files (one batch).
func (op *copyOp) Run(files []MediaFile) error {
	var values = make([][]interface{}, len(op.fields))
	for idx, field := range op.fields {
		values[idx] = field.GetValues(files[0].Handler)
	}
	for _, file := range files[1:] {
		for idx, field := range op.fields {
			if err := field.SetValues(file.Handler, values[idx]); err != nil {
				return fmt.Errorf("%s: copy %s: %s", file.Path, field.PluralName(), err)
			}
		}
		file.Changed = true
	}
	return nil
}
