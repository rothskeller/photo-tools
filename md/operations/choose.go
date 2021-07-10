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

// Choose displays all of the tagged values for a field, lets the user choose
// one or more (or type new values), and then sets the field to those values
// across all target files.
func Choose(args []string, files []MediaFile) (err error) {
	var (
		field  fields.Field
		values []interface{}
		newvs  []interface{}
		scan   *bufio.Scanner
		line   string
		tw     = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	)
	switch len(args) {
	case 0:
		return errors.New("choose: missing field name")
	case 1:
		break
	default:
		return errors.New("choose: excess arguments")
	}
	if field = fields.ParseField(args[0]); field == nil {
		return fmt.Errorf("choose: %q is not a recognized field name", args[0])
	}
	// Print the tag table.
	fmt.Fprintln(tw, "#\tFILE\tTAG\tVALUE")
	for _, file := range files {
		tagNames, tagValues := field.GetTags(file.Provider)
		for i, tag := range tagNames {
			values = append(values, tagValues[i])
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", len(values), file.Path, tag, escapeString(field.RenderValue(tagValues[i])))
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
		fmt.Printf("Enter a new value for %s, or 0 to clear, 1-%d to copy, or nothing to skip.\n? ", field.Name(), len(values))
	} else {
		fmt.Printf("Enter a new value for %s, 0 to clear, or nothing to skip.\n? ", field.Name())
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
	} else { // Not a line number list; is it a set of valid values for the field?
		list := strings.Split(line, ";")
		for _, item := range list {
			if newv, err := field.ParseValue(strings.TrimSpace(item)); err != nil {
				fmt.Printf("ERROR: %s\n", err)
				goto RETRY
			} else {
				newvs = append(newvs, newv)
			}
		}
	}
	// Set these value(s) on all files in the batch.
	for i, file := range files {
		if err := field.SetValues(file.Provider, newvs); err != nil {
			return fmt.Errorf("%s: choose %s: %s", file.Path, field.Name(), err)
		}
		files[i].Changed = true
	}
	return nil
}

func parseLineNumberSet(s string, max int) (nums []int, showedError bool) {
	var seen = make(map[int]bool)

	nums = ParseNumberList(s)
	switch len(nums) {
	case 0:
		return nil, false
	case 1:
		if nums[0] == 0 {
			return nums, false // otherwise it would be rejected below
		}
	}
	for _, num := range nums {
		if num == 0 || num >= max {
			fmt.Printf("ERROR: no such line number %d\n", num)
			return nil, true
		}
		if seen[num] {
			fmt.Printf("ERROR: line number %d repeated\n", num)
			return nil, true
		}
		seen[num] = true
	}
	return nums, false
}

// ParseNumberList parses a string containing a space-separated list of
// non-negative integers and/or ranges of non-negative integers.  It returns nil
// if the input does not look like such a list.
func ParseNumberList(s string) (list []int) {
	var parts = strings.Fields(s)
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		idx := strings.IndexByte(part, '-')
		if idx < 0 {
			if num, err := strconv.Atoi(part); err == nil {
				list = append(list, num)
			} else {
				return nil
			}
		} else {
			num1, err := strconv.Atoi(part[:idx])
			if err != nil {
				return nil
			}
			num2, err := strconv.Atoi(part[idx+1:])
			if err != nil || num2 < num1 {
				return nil
			}
			for i := num1; i <= num2; i++ {
				list = append(list, i)
			}
		}
	}
	return list
}
