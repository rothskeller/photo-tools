package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"
)

func createAddOp(args []string) (op operation, remainingArgs []string, err error) {
	var (
		argfields []*field
		value     interface{}
	)
	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	switch len(argfields) {
	case 0:
		return nil, args[1:], errors.New("add: missing field name")
	case 1:
		if !argfields[0].multivalued {
			return nil, args[1:], fmt.Errorf("add: not supported for %q", argfields[0].name)
		}
	default:
		return nil, args[1:], errors.New("add: not supported for \"all\"")
	}
	if len(args) < 2 {
		return nil, args[1:], errors.New("add: missing value")
	}
	if value, err = argfields[0].parseValue(args[1]); err != nil {
		return nil, args[2:], fmt.Errorf("add %s: %s", args[0], err)
	}
	return &addOp{argfields[0], value}, args[2:], nil
}

type addOp struct {
	field *field
	value interface{}
}

var _ operation = addOp{}

func (op addOp) check(batches [][]mediafile) error {
	return nil
}

func (op addOp) run(files []mediafile) error {
	for _, file := range files {
		values := op.field.getValues(file.handler)
		found := false
		for _, v := range values {
			if op.field.equalValue(v, op.value) {
				found = true
				break
			}
		}
		if !found {
			values = append(values, op.value)
			if err := op.field.setValues(file.handler, values); err != nil {
				return fmt.Errorf("%s: add %s: %s", file.path, op.field.name, err)
			}
		}
	}
	return nil
}

func createChooseOp(args []string) (op operation, remainingArgs []string, err error) {
	var argfields []*field

	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	switch len(argfields) {
	case 0:
		return nil, args[1:], errors.New("choose: missing field name")
	case 1:
		if argfields[0].multivalued {
			return nil, args[1:], fmt.Errorf("choose: not supported for %q", argfields[0].name)
		}
	default:
		return nil, args[1:], errors.New("choose: not supported for \"all\"")
	}
	return &chooseOp{argfields[0]}, args[1:], nil
}

type chooseOp struct {
	field *field
}

var _ operation = chooseOp{}

func (op chooseOp) check(batches [][]mediafile) error {
	return nil
}

func (op chooseOp) run(files []mediafile) error {
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
		tagNames, tagValues := op.field.getTags(file.handler)
		values = append(values, tagValues...)
		for i, tag := range tagNames {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", len(values), file.path, tag, escapeString(op.field.renderValue(tagValues[i])))
		}
		if len(tagNames) == 0 {
			fmt.Fprintf(tw, "\t%s\t(none)\t\n", file.path)
		}
	}
	tw.Flush()
	scan = bufio.NewScanner(os.Stdin)
RETRY:
	if len(values) != 0 {
		fmt.Printf("Enter a new value for %s, or 0 to clear, 1-%d to copy, or nothing to skip.\n? ", op.field.name, len(values))
	} else {
		fmt.Printf("Enter a new value for %s, 0 to clear, or nothing to skip.\n? ", op.field.name)
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
		if newv, err = op.field.parseValue(line); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			goto RETRY
		}
		newvs = []interface{}{newv}
	}
	for _, file := range files {
		if err := op.field.setValues(file.handler, newvs); err != nil {
			return fmt.Errorf("%s: choose %s: %s", file.path, op.field.name, err)
		}
	}
	return nil
}

func createClearOp(args []string) (op operation, remainingArgs []string, err error) {
	var argfields []*field

	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	switch len(argfields) {
	case 0:
		return nil, args[1:], errors.New("clear: missing field name")
	case 1:
		break
	default:
		return nil, args[1:], errors.New("clear: not supported for \"all\"")
	}
	return &setOp{
		"clear " + argfields[0].pluralName,
		argfields[0],
		nil,
	}, args[1:], nil
}

func createCopyOp(args []string) (op operation, remainingArgs []string, err error) {
	var fields []*field

	for len(args) != 0 {
		argfields := parseField(args[0])
		if len(argfields) == 0 {
			break
		}
		for _, af := range argfields {
			var found bool
			for _, f := range fields {
				if f.name == af.name {
					found = true
					break
				}
			}
			if !found {
				fields = append(fields, af)
			}
		}
		args = args[1:]
	}
	if len(fields) == 0 {
		fields = []*field{artistField, captionField, dateTimeField, gpsField, keywordsField, locationField, titleField}
	}
	return copyOp{fields}, args, nil
}

type copyOp struct {
	fields []*field
}

var _ operation = copyOp{}

func (op copyOp) check(batches [][]mediafile) error {
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

func (op copyOp) run(files []mediafile) error {
	var values = make([][]interface{}, len(op.fields))
	for idx, field := range op.fields {
		values[idx] = field.getValues(files[0].handler)
	}
	for _, file := range files[1:] {
		for idx, field := range op.fields {
			if err := field.setValues(file.handler, values[idx]); err != nil {
				return fmt.Errorf("%s: copy %s: %s", file.path, field.pluralName, err)
			}
		}
	}
	return nil
}

func createRemoveOp(args []string) (op operation, remainingArgs []string, err error) {
	var (
		argfields []*field
		value     interface{}
	)
	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	switch len(argfields) {
	case 0:
		return nil, args[1:], errors.New("remove: missing field name")
	case 1:
		if !argfields[0].multivalued {
			return nil, args[1:], fmt.Errorf("remove: not supported for %q", argfields[0].name)
		}
	default:
		return nil, args[1:], errors.New("remove: not supported for \"all\"")
	}
	if len(args) < 2 {
		return nil, args[1:], errors.New("remove: missing value")
	}
	if value, err = argfields[0].parseValue(args[1]); err != nil {
		return nil, args[2:], fmt.Errorf("remove %s: %s", args[0], err)
	}
	return &removeOp{argfields[0], value}, args[2:], nil
}

type removeOp struct {
	field *field
	value interface{}
}

var _ operation = removeOp{}

func (op removeOp) check(batches [][]mediafile) error {
	return nil
}

func (op removeOp) run(files []mediafile) error {
	for _, file := range files {
		values := op.field.getValues(file.handler)
		j := 0
		for _, v := range values {
			if !op.field.equalValue(v, op.value) {
				values[j] = v
				j++
			}
		}
		if j < len(values) {
			if err := op.field.setValues(file.handler, values[:j]); err != nil {
				return fmt.Errorf("%s: remove %s: %s", file.path, op.field.name, err)
			}
		}
	}
	return nil
}

func createSetOp(args []string) (op operation, remainingArgs []string, err error) {
	var (
		argfields []*field
		value     interface{}
	)
	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	switch len(argfields) {
	case 0:
		return nil, args[1:], errors.New("set: missing field name")
	case 1:
		if argfields[0] == keywordsField {
			return nil, args[1:], errors.New("set: not supported for keyword")
		}
	default:
		return nil, args[1:], errors.New("set: not supported for \"all\"")
	}
	if len(args) < 2 {
		return nil, args[1:], errors.New("set: missing value")
	}
	if value, err = argfields[0].parseValue(args[1]); err != nil {
		return nil, args[2:], fmt.Errorf("set %s: %s", args[0], err)
	}
	return &setOp{
		"set " + argfields[0].pluralName,
		argfields[0],
		[]interface{}{value},
	}, args[2:], nil
}

type setOp struct {
	label  string
	field  *field
	values []interface{}
}

var _ operation = setOp{}

func (op setOp) check(batches [][]mediafile) error { return nil }
func (op setOp) run(files []mediafile) error {
	for _, file := range files {
		if err := op.field.setValues(file.handler, op.values); err != nil {
			return fmt.Errorf("%s: %s: %s", file.path, op.label, err)
		}
	}
	return nil
}

func createWriteOp(args []string) (op operation, remainingArgs []string, err error) {
	var (
		argfields []*field
		by        []byte
		value     interface{}
	)
	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	if len(argfields) != 1 || argfields[0] != captionField {
		return nil, args, errors.New("write: field must be caption")
	}
	if by, err = io.ReadAll(os.Stdin); err != nil {
		return nil, args[1:], fmt.Errorf("write: standard input: %s", err)
	}
	if value, err = captionField.parseValue(string(by)); err != nil {
		return nil, args[1:], fmt.Errorf("write %s: %s", captionField.name, err)
	}
	return setOp{"write caption", captionField, []interface{}{value}}, args[1:], nil
}
