package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
)

func usage() {
	fmt.Fprint(os.Stderr, `usage: md [batch] [operation...] file...
  Operations:                Fields:
    show [fieldname...]        a(rtist)
    tags [fieldname...]        c(aption)
    set fieldname value        d(atetime)
    add fieldname value        g(ps)
    remove fieldname value     k(eywords) (kw)
    clear fieldname            l(ocation)[:lang]
    choose fieldname           s(hown)[:lang]
    copy [fieldname...]        t(itle)
    write caption              group(s), person (people), topic(s)
    read caption               all
See MANUAL.md for more details.
`)
	os.Exit(2)
}

// parseCommandLine parses the command line and returns the list of operations
// and the list of files.  It aborts the program if there are any errors.
func parseCommandLine() (ops []operation, batches [][]mediafile) {
	var (
		args      []string
		arg       string
		batch     bool
		op        operation
		handler   fileHandler
		fileError bool
		err       error
	)
	args = os.Args[1:]
	batches = append(batches, []mediafile{})
ARGS:
	for len(args) != 0 {
		arg, args = args[0], args[1:]
		if arg == "batch" {
			if batch {
				fmt.Fprintf(os.Stderr, "ERROR: \"batch\" specified repeatedly\n")
				usage()
			}
			batch = true
		}
		if factory := optypes[arg]; factory != nil {
			if op, args, err = factory(args); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
				usage()
			}
			ops = append(ops, op)
			continue ARGS
		}
		if _, err := os.Stat(arg); os.IsNotExist(err) {
			usage()
		}
		if handler = filefmt.HandlerFor(arg); handler == nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s: unsupported file format\n", arg)
			fileError = true
			continue ARGS
		}
		handler.ReadMetadata()
		if problems := handler.Problems(); len(problems) != 0 {
			for _, problem := range problems {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", arg, problem)
			}
			fileError = true
			continue ARGS
		}
		batches[0] = append(batches[0], mediafile{arg, handler})
	}
	if len(batches[0]) == 0 && !fileError {
		fmt.Fprintf(os.Stderr, "ERROR: no files specified\n")
		usage()
	}
	if len(batches[0]) == 0 {
		os.Exit(1)
	}
	if len(ops) == 0 {
		op, _, _ := createShowOp([]string{})
		ops = append(ops, op)
	}
	if batch {
		var (
			bnum = 0
			fnum = 1
			base = basename(batches[0][0].path)
		)
		for fnum < len(batches[bnum]) {
			if nb := basename(batches[bnum][fnum].path); nb == base {
				fnum++
			} else {
				batches = append(batches, batches[bnum][fnum:])
				batches[bnum] = batches[bnum][:fnum]
				base, bnum, fnum = nb, bnum+1, 1
			}
		}
	}
	for _, op := range ops {
		if err := op.check(batches); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			usage()
		}
	}
	return ops, batches
}

func basename(path string) string {
	path = filepath.Base(path)
	if idx := strings.IndexByte(path, '.'); idx >= 0 {
		path = path[:idx]
	}
	return path
}

type operationFactory func(args []string) (op operation, remainingArgs []string, err error)

var optypes = map[string]operationFactory{
	"add":    createAddOp,
	"choose": createChooseOp,
	"clear":  createClearOp,
	"copy":   createCopyOp,
	"read":   createReadOp,
	"remove": createRemoveOp,
	"set":    createSetOp,
	"show":   createShowOp,
	"tags":   createTagsOp,
	"write":  createWriteOp,
}

func parseField(arg string) []*field {
	switch arg {
	case "artist", "a":
		return []*field{artistField}
	case "caption", "c":
		return []*field{captionField}
	case "datetime", "date", "time", "d":
		return []*field{dateTimeField}
	case "gps", "g":
		return []*field{gpsField}
	case "keyword", "keywords", "kw", "k":
		return []*field{keywordsField}
	case "location", "loc", "l":
		return []*field{locationField}
	case "title", "t":
		return []*field{titleField}
	case "group", "groups":
		return []*field{groupsField}
	case "person", "people":
		return []*field{peopleField}
	case "place", "places":
		return []*field{placesField}
	case "topic", "topics":
		return []*field{topicsField}
	case "all":
		return []*field{titleField, dateTimeField, artistField, gpsField, locationField,
			placesField, peopleField, groupsField, topicsField, otherKeywordsField, captionField}
	}
	return nil
}
