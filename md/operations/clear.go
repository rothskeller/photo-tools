package operations

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newClearOp() Operation { return new(clearOp) }

// clearOp removes all values of the specified field.
type clearOp struct {
	field fields.Field
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *clearOp) parseArgs(args []string) (remainingArgs []string, err error) {
	if len(args) == 0 {
		return nil, errors.New("clear: missing field name")
	}
	if op.field = fields.ParseField(args[0]); op.field == nil {
		return nil, errors.New("clear: missing field name")
	}
	return args[1:], nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *clearOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *clearOp) Run(files []MediaFile) error {
	for _, file := range files {
		if err := op.field.SetValues(file.Handler, nil); err != nil {
			return fmt.Errorf("%s: clear %s: %s", file.Path, op.field.PluralName(), err)
		}
		file.Changed = true
	}
	return nil
}
