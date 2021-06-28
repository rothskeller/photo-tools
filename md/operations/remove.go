package operations

import "fmt"

func newRemoveOp() Operation { return &removeOp{fieldValueOp{name: "remove"}} }

// removeOp removes a value from a multi-valued field.
type removeOp struct {
	fieldValueOp
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *removeOp) parseArgs(args []string) (remainingArgs []string, err error) {
	if remainingArgs, err = op.fieldValueOp.parseArgs(args); err != nil {
		return nil, err
	}
	if !op.field.Multivalued() {
		return nil, fmt.Errorf("remove: not supported for %q", op.field.Name())
	}
	return remainingArgs, nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *removeOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *removeOp) Run(files []MediaFile) error {
	for _, file := range files {
		// Get the current values.
		values := op.field.GetValues(file.Handler)
		// Remove the one we were asked to remove.
		j := 0
		for _, v := range values {
			if !op.field.EqualValue(v, op.value) {
				values[j] = v
				j++
			}
		}
		// If we found it, set the new value list that leaves it out.
		if j < len(values) {
			if err := op.field.SetValues(file.Handler, values[:j]); err != nil {
				return fmt.Errorf("%s: remove %s: %s", file.Path, op.field.Name(), err)
			}
			file.Changed = true
		}
	}
	return nil
}
