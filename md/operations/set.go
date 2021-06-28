package operations

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newSetOp() Operation { return &setOp{fieldValueOp{name: "set"}} }

// setOp sets the value of a field.
type setOp struct {
	fieldValueOp
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *setOp) parseArgs(args []string) (remainingArgs []string, err error) {
	if remainingArgs, err = op.fieldValueOp.parseArgs(args); err != nil {
		return nil, err
	}
	if op.field == fields.KeywordsField {
		return nil, errors.New(`set: not supported for "keyword" (you probably want "add" instead)`)
	}
	return remainingArgs, nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *setOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *setOp) Run(files []MediaFile) error {
	for _, file := range files {
		if err := op.field.SetValues(file.Handler, []interface{}{op.value}); err != nil {
			return fmt.Errorf("%s: clear %s: %s", file.Path, op.field.PluralName(), err)
		}
		file.Changed = true
	}
	return nil
}
