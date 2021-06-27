package operations

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newChooseOp() Operation { return new(chooseOp) }

type chooseOp struct {
	field fields.Field
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *chooseOp) parseArgs(args []string) (remainingArgs []string, err error) {
	if len(args) == 0 {
		return nil, errors.New("choose: missing field name")
	}
	if op.field = fields.ParseField(args[0]); op.field == nil {
		return nil, errors.New("choose: missing field name")
	}
	if op.field.Multivalued() {
		return nil, fmt.Errorf("choose: not supported for %q", op.field.Name())
	}
	return args[1:], nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *chooseOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *chooseOp) Run(files []MediaFile) error {
	var (
		values []interface{}
		newv   interface{}
		newvs  []interface{}
		scan   *bufio.Scanner
		line   string
		tw     = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	)
	fmt.Fprintln(tw, "#\tFILE\tTAG\tVALUE")
	for _, file := range files {
		tagNames, tagValues := op.field.GetTags(file.Handler)
		values = append(values, tagValues...)
		for i, tag := range tagNames {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", len(values), file.Path, tag, escapeString(op.field.RenderValue(tagValues[i])))
		}
		if len(tagNames) == 0 {
			fmt.Fprintf(tw, "\t%s\t(none)\t\n", file.Path)
		}
	}
	tw.Flush()
	scan = bufio.NewScanner(os.Stdin)
RETRY:
	if len(values) != 0 {
		fmt.Printf("Enter a new value for %s, or 0 to clear, 1-%d to copy, or nothing to skip.\n? ", op.field.Name(), len(values))
	} else {
		fmt.Printf("Enter a new value for %s, 0 to clear, or nothing to skip.\n? ", op.field.Name())
	}
	if !scan.Scan() {
		return scan.Err()
	}
	if line = scan.Text(); line == "" {
		return nil
	}
	if lnum, err := strconv.Atoi(line); err == nil {
		if lnum == 0 {
			newvs = nil
		} else if lnum > 0 && lnum <= len(values) {
			newvs = []interface{}{values[lnum-1]}
		} else {
			fmt.Printf("ERROR: no such line %s.\n", line)
			goto RETRY
		}
	} else {
		if newv, err = op.field.ParseValue(line); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			goto RETRY
		}
		newvs = []interface{}{newv}
	}
	for _, file := range files {
		if err := op.field.SetValues(file.Handler, newvs); err != nil {
			return fmt.Errorf("%s: choose %s: %s", file.Path, op.field.Name(), err)
		}
	}
	return nil
}
