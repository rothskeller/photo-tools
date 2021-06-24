package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func createShowOp(args []string) (op operation, remainingArgs []string, err error) {
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
	return showOp(fields), args, nil
}

type showOp []field

var _ operation = showOp{}

func (op showOp) check(batches [][]mediafile) error {
	return nil
}

func (op showOp) run(files []mediafile) error {
	var tw = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tFIELD\tVALUE")
	for _, file := range files {
		for _, field := range op {
			values := field.get(file.handler)
			if len(values) == 0 {
				fmt.Fprintf(tw, "%s\t%s\t\n", file.path, field.label())
			} else {
				for _, value := range values {
					fmt.Fprintf(tw, "%s\t%s\t%s\n", file.path, field.label(), escapeString(value.String()))
				}
			}
		}
	}
	tw.Flush()
	return nil
}

func createTagsOp(args []string) (op operation, remainingArgs []string, err error) {
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
	return showOp(fields), args, nil
}

type tagsOp []field

var _ operation = tagsOp{}

func (op tagsOp) check(batches [][]mediafile) error {
	return nil
}

func (op tagsOp) run(files []mediafile) error {
	var tw = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tTAG\tVALUE")
	for _, file := range files {
		for _, field := range op {
			tags := field.tags(file.handler)
			for _, tag := range tags {
				fmt.Fprintf(tw, "%s\t%s\t%s\n", file.path, tag.Tag, escapeString(tag.String()))
			}
		}
	}
	tw.Flush()
	return nil
}

func escapeString(s string) string {
	return strings.Replace(strings.Replace(s, "\\", "\\\\", -1), "\n", "\\n", -1)
}

func createReadOp(args []string) (op operation, remainingArgs []string, err error) {
	var argfields []field

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
	caption := captionField.get(files[0].handler)
	if !caption.Empty() {
		str := caption.String()
		fmt.Print(str)
		if len(str) != 0 && str[len(str)-1] != '\n' {
			fmt.Println()
		}
	}
	return nil
}
