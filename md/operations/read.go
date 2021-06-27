package operations

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newReadOp() Operation { return new(readOp) }

type readOp struct{}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *readOp) parseArgs(args []string) (remainingArgs []string, err error) {
	if len(args) == 0 {
		return nil, errors.New("read: missing field name")
	}
	if field := fields.ParseField(args[0]); field == nil {
		return nil, errors.New("read: missing field name")
	} else if field != fields.CaptionField {
		return nil, fmt.Errorf("read: not supported for %q", field.Name())
	}
	return args[1:], nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *readOp) Check(batches [][]MediaFile) error {
	if len(batches) > 0 || len(batches[0]) > 0 {
		return errors.New("read caption: only one file allowed")
	}
	return nil
}

// Run executes the operation against the listed media files (one batch).
func (op *readOp) Run(files []MediaFile) error {
	caption := fields.CaptionField.GetValues(files[0].Handler)
	if len(caption) != 0 {
		str := fields.CaptionField.RenderValue(caption[0])
		fmt.Print(str)
		if len(str) != 0 && str[len(str)-1] != '\n' {
			fmt.Println()
		}
	}
	return nil
}
