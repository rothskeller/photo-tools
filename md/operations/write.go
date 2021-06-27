package operations

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newWriteOp() Operation { return new(writeOp) }

type writeOp struct {
	value string
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *writeOp) parseArgs(args []string) (remainingArgs []string, err error) {
	if len(args) == 0 {
		return nil, errors.New("write: missing field name")
	}
	if field := fields.ParseField(args[0]); field == nil {
		return nil, errors.New("write: missing field name")
	} else if field != fields.CaptionField {
		return nil, fmt.Errorf("write: not supported for %q", field.Name())
	}
	if by, err := io.ReadAll(os.Stdin); err == nil {
		op.value = string(by)
	} else {
		return args[1:], fmt.Errorf("write: standard input: %s", err)
	}
	return args[1:], nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *writeOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *writeOp) Run(files []MediaFile) error {
	for _, file := range files {
		if err := fields.CaptionField.SetValues(file.Handler, []interface{}{op.value}); err != nil {
			return fmt.Errorf("%s: write caption: %s", file.Path, err)
		}
	}
	return nil
}
