package operations

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newChooseOp() Operation { return new(chooseOp) }

// chooseOp displays all of the tagged values for a field, lets the user choose
// one (or type a new value), and then sets the field to that value across all
// files in the batch.
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
		newvs  []interface{}
		scan   *bufio.Scanner
		line   string
		tw     = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	)
	// Print the tag table.
	fmt.Fprintln(tw, "#\tFILE\tTAG\tVALUE")
	for _, file := range files {
		tagNames, tagValues := op.field.GetTags(file.Handler)
		for i, tag := range tagNames {
			values = append(values, tagValues[i])
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", len(values), file.Path, tag, escapeString(op.field.RenderValue(tagValues[i])))
		}
		if len(tagNames) == 0 {
			fmt.Fprintf(tw, "\t%s\t(none)\t\n", file.Path)
		}
	}
	tw.Flush()
	// Repeat reading lines from stdin until we get a valid answer.
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
	// Parse the response.  Is it a line number list?
	if nums, showedErr := parseLineNumberSet(line, len(values)); showedErr {
		goto RETRY
	} else if len(nums) == 1 && nums[0] == 0 {
		newvs = nil
	} else if len(nums) != 0 {
		newvs = make([]interface{}, len(nums))
		for i := range nums {
			newvs[i] = values[nums[i]-1]
		}
	} else { // Not a line number list; is it a valid value for the field?
		if newv, err := op.field.ParseValue(line); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			goto RETRY
		} else {
			newvs = []interface{}{newv}
		}
	}
	// Set these value(s) on all files in the batch.
	for i, file := range files {
		if err := op.field.SetValues(file.Handler, newvs); err != nil {
			return fmt.Errorf("%s: choose %s: %s", file.Path, op.field.Name(), err)
		}
		files[i].Changed = true
	}
	return nil
}

func parseLineNumberSet(s string, max int) (nums []int, showedError bool) {
	var (
		err   error
		seen  = make(map[int]bool)
		parts = strings.Fields(s)
	)
	nums = make([]int, len(parts))
	for i := range parts {
		if nums[i], err = strconv.Atoi(parts[i]); err != nil {
			return nil, false
		}
		if nums[i] < 0 {
			return nil, false
		}
		if nums[i] == 0 && i != 0 {
			fmt.Println("ERROR: 0 must appear alone, not with other line numbers")
			return nil, true
		}
		if nums[i] > max {
			fmt.Printf("ERROR: no such line number %d\n", nums[i])
			return nil, true
		}
		if seen[nums[i]] {
			fmt.Printf("ERROR: line number %d repeated\n", nums[i])
			return nil, true
		}
		seen[nums[i]] = true
	}
	return nums, false
}
