package operations

import (
	"fmt"
)

func newAddOp() Operation { return &addOp{fieldValueOp{name: "add"}} }

// addOp adds a value to a multivalued field.
type addOp struct {
	fieldValueOp
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *addOp) parseArgs(args []string) (remainingArgs []string, err error) {
	if remainingArgs, err = op.fieldValueOp.parseArgs(args); err != nil {
		return nil, err
	}
	if !op.field.Multivalued() {
		return nil, fmt.Errorf("add: not supported for %q", op.field.Name())
	}
	return remainingArgs, nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *addOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *addOp) Run(files []MediaFile) error {
	for _, file := range files {
		// Get the current values.
		values := op.field.GetValues(file.Handler)
		// Find out whether the value we're adding is already there.
		found := false
		for _, v := range values {
			if op.field.EqualValue(v, op.value) {
				found = true
				break
			}
		}
		// If not, add it.
		if !found {
			values = append(values, op.value)
			if err := op.field.SetValues(file.Handler, values); err != nil {
				return fmt.Errorf("%s: add %s: %s", file.Path, op.field.Name(), err)
			}
			file.Changed = true
		}
	}
	return nil
}
