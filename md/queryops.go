package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func createShowOp(args []string) (op operation, remainingArgs []string, err error) {
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
		fields = parseField("all")
	}
	return showOp{fields: fields}, args, nil
}

type showOp struct {
	fields []*field
	hasRun bool
}

var _ operation = showOp{}

func (op showOp) check(batches [][]mediafile) error {
	return nil
}

func (op showOp) run(files []mediafile) error {
	if op.hasRun {
		fmt.Println()
	}
	var tw = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tFIELD\tVALUE")
	for _, file := range files {
		for _, field := range op.fields {
			values := field.getValues(file.handler)
			if len(values) == 0 {
				fmt.Fprintf(tw, "%s\t%s\t\n", file.path, field.label)
			} else {
				for _, value := range values {
					fmt.Fprintf(tw, "%s\t%s\t%s\n", file.path, field.label, escapeString(field.renderValue(value)))
				}
			}
		}
	}
	tw.Flush()
	op.hasRun = true
	return nil
}

func createTagsOp(args []string) (op operation, remainingArgs []string, err error) {
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
		fields = parseField("all")
	}
	return tagsOp{fields: fields}, args, nil
}

type tagsOp struct {
	fields []*field
	hasRun bool
}

var _ operation = tagsOp{}

func (op tagsOp) check(batches [][]mediafile) error {
	return nil
}

func (op tagsOp) run(files []mediafile) error {
	if op.hasRun {
		fmt.Println()
	}
	var tw = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tTAG\tVALUE")
	for _, file := range files {
		for _, field := range op.fields {
			tagNames, tagValues := field.getTags(file.handler)
			for i, tag := range tagNames {
				fmt.Fprintf(tw, "%s\t%s\t%s\n", file.path, tag, escapeString(field.renderValue(tagValues[i])))
			}
		}
	}
	tw.Flush()
	op.hasRun = true
	return nil
}

func escapeString(s string) string {
	return strings.Replace(strings.Replace(s, "\\", "\\\\", -1), "\n", "\\n", -1)
}

func createReadOp(args []string) (op operation, remainingArgs []string, err error) {
	var argfields []*field

	if len(args) != 0 {
		argfields = parseField(args[0])
	}
	if len(argfields) != 1 || argfields[0] != captionField {
		return nil, args, errors.New("read: field must be caption")
	}
	return readCaptionOp{}, args[1:], nil
}

type readCaptionOp struct{}

var _ operation = readCaptionOp{}

func (readCaptionOp) check(batches [][]mediafile) error {
	if len(batches) > 0 || len(batches[0]) > 0 {
		return errors.New("read caption: only one file allowed")
	}
	return nil
}

func (readCaptionOp) run(files []mediafile) error {
	caption := captionField.getValues(files[0].handler)
	if len(caption) != 0 {
		str := captionField.renderValue(caption[0])
		fmt.Print(str)
		if len(str) != 0 && str[len(str)-1] != '\n' {
			fmt.Println()
		}
	}
	return nil
}
