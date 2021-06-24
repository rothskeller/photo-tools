package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/metadata"
)

func createAddOp(args []string) (op operation, remainingArgs []string, err error) {
	var (
		argfields []field
		value     metadata.Metadatum
	)
	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	switch len(argfields) {
	case 0:
		return nil, args[1:], errors.New("add: missing field name")
	case 1:
		if !argfields[0].multivalued() {
			return nil, args[1:], fmt.Errorf("add: not supported for %q", argfields[0].name())
		}
	default:
		return nil, args[1:], errors.New("add: not supported for \"all\"")
	}
	if len(args) < 2 {
		return nil, args[1:], errors.New("add: missing value")
	}
	value = argfields[1].newValue()
	if err = value.Parse(args[1]); err != nil {
		return nil, args[2:], fmt.Errorf("add %s: %s", args[0], err)
	}
	return &addOp{argfields[0], value}, args[2:], nil
}

type addOp struct {
	field field
	value metadata.Metadatum
}

var _ operation = addOp{}

func (op addOp) check(batches [][]mediafile) error {
	return nil
}

func (op addOp) run(files []mediafile) error {
	for _, file := range files {
		values := op.field.get(file.handler)
		found := false
		for _, v := range values {
			if v.Equal(op.value) {
				found = true
				break
			}
		}
		if !found {
			values = append(values, op.value)
			if err := op.field.set(file.handler, values); err != nil {
				return fmt.Errorf("%s: add %s: %s", file.path, op.field.name(), err)
			}
		}
	}
	return nil
}

func createChooseOp(args []string) (op operation, remainingArgs []string, err error) {
	var argfields []field

	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	switch len(argfields) {
	case 0:
		return nil, args[1:], errors.New("choose: missing field name")
	case 1:
		if argfields[0].multivalued() || argfields[0].langtagged() {
			return nil, args[1:], fmt.Errorf("choose: not supported for %q", argfields[0].name())
		}
	default:
		return nil, args[1:], errors.New("choose: not supported for \"all\"")
	}
	return &chooseOp{argfields[0]}, args[1:], nil
}

type chooseOp struct {
	field field
}

var _ operation = chooseOp{}

func (op chooseOp) check(batches [][]mediafile) error {
	return nil
}

func (op chooseOp) run(files []mediafile) error {
	var (
		values []metadata.Metadatum
		newv   metadata.Metadatum
		newvs  []metadata.Metadatum
		scan   *bufio.Scanner
		line   string
		tw     = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	)
	fmt.Fprintln(tw, "#\tFILE\tTAG\tVALUE")
	for _, file := range files {
		tags := op.field.tags(file.handler)
		for _, tag := range tags {
			values = append(values, tag.Metadatum)
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", len(values), file.path, tag.Tag, escapeString(tag.String()))
		}
		if len(tags) == 0 {
			fmt.Fprintf(tw, "\t%s\t(none)\t\n", file.path)
		}
	}
	tw.Flush()
	newv = op.field.newValue()
	scan = bufio.NewScanner(os.Stdin)
RETRY:
	if len(values) != 0 {
		fmt.Printf("Enter a new value for %s, or 0 to clear, 1-%d to copy, or nothing to skip.\n? ", op.field.name(), len(values))
	} else {
		fmt.Printf("Enter a new value for %s, 0 to clear, or nothing to skip.\n? ", op.field.name())
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
			newvs = []metadata.Metadatum{values[lnum-1]}
		} else {
			fmt.Printf("ERROR: no such line %s.\n", line)
			goto RETRY
		}
	} else {
		if err := newv.Parse(line); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			goto RETRY
		}
		newvs = []metadata.Metadatum{newv}
	}
	for _, file := range files {
		if err := op.field.set(file.handler, newvs); err != nil {
			return fmt.Errorf("%s: choose %s: %s", file.path, op.field.name(), err)
		}
	}
	return nil
}

func createClearOp(args []string) (op operation, remainingArgs []string, err error) {
	var argfields []field

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
		"clear " + argfields[0].pluralName(),
		argfields[0],
		[]metadata.Metadatum{},
	}, args[2:], nil
}

func createCopyOp(args []string) (op operation, remainingArgs []string, err error) {
	var fields []field

	for len(args) != 0 {
		argfields := parseField(args[0])
		if len(argfields) == 0 {
			break
		}
		for _, af := range argfields {
			var found bool
			for _, f := range fields {
				if f == af {
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
		fields = []field{}
	}
	return copyOp{fields}, args, nil
}

type copyOp struct {
	fields []field
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
	var values = make([][]metadata.Metadatum, len(op.fields))
	for idx, field := range op.fields {
		values[idx] = field.get(files[0].handler)
	}
	for _, file := range files[1:] {
		for idx, field := range op.fields {
			if err := field.set(file.handler, values[idx]); err != nil {
				return fmt.Errorf("%s: copy %s: %s", file.path, field.pluralName(), err)
			}
		}
	}
	return nil
}

func createRemoveOp(args []string) (op operation, remainingArgs []string, err error) {
	var (
		argfields []field
		value     metadata.Metadatum
	)
	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	switch len(argfields) {
	case 0:
		return nil, args[1:], errors.New("remove: missing field name")
	case 1:
		if !argfields[0].multivalued() {
			return nil, args[1:], fmt.Errorf("remove: not supported for %q", argfields[0].name())
		}
	default:
		return nil, args[1:], errors.New("remove: not supported for \"all\"")
	}
	if len(args) < 2 {
		return nil, args[1:], errors.New("remove: missing value")
	}
	value = argfields[1].newValue()
	if err = value.Parse(args[1]); err != nil {
		return nil, args[2:], fmt.Errorf("remove %s: %s", args[0], err)
	}
	return &removeOp{argfields[0], value}, args[2:], nil
}

type removeOp struct {
	field field
	value metadata.Metadatum
}

var _ operation = removeOp{}

func (op removeOp) check(batches [][]mediafile) error {
	return nil
}

func (op removeOp) run(files []mediafile) error {
	for _, file := range files {
		values := op.field.get(file.handler)
		j := 0
		for _, v := range values {
			if !v.Equal(op.value) {
				values[j] = v
				j++
			}
		}
		if j < len(values) {
			if err := op.field.set(file.handler, values[:j]); err != nil {
				return fmt.Errorf("%s: remove %s: %s", file.path, op.field.name(), err)
			}
		}
	}
	return nil
}

func createSetOp(args []string) (op operation, remainingArgs []string, err error) {
	var (
		argfields []field
		value     metadata.Metadatum
	)
	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	switch len(argfields) {
	case 0:
		return nil, args[1:], errors.New("set: missing field name")
	case 1:
		if argfields[0] == keywordField {
			return nil, args[1:], errors.New("set: not supported for keyword")
		}
	default:
		return nil, args[1:], errors.New("set: not supported for \"all\"")
	}
	if len(args) < 2 {
		return nil, args[1:], errors.New("set: missing value")
	}
	value = argfields[1].newValue()
	if err = value.Parse(args[1]); err != nil {
		return nil, args[2:], fmt.Errorf("set %s: %s", args[0], err)
	}
	return &setOp{
		"set " + argfields[0].pluralName(),
		argfields[0],
		[]metadata.Metadatum{value},
	}, args[2:], nil
}

type setOp struct {
	label  string
	field  field
	values []metadata.Metadatum
}

var _ operation = setOp{}

func (op setOp) check(batches [][]mediafile) error { return nil }
func (op setOp) run(files []mediafile) error {
	for _, file := range files {
		if err := op.field.set(file.handler, op.values); err != nil {
			return fmt.Errorf("%s: %s: %s", file.path, op.label, err)
		}
	}
	return nil
}

func createWriteOp(args []string) (op operation, remainingArgs []string, err error) {
	var (
		argfields []field
		by        []byte
		value     metadata.String
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
	value = metadata.String(string(by))
	return setOp{"write caption", captionField, []metadata.Metadatum{&value}}, args[1:], nil
}
