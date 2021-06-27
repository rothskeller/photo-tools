package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/md/operations"
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
func parseCommandLine() (ops []operations.Operation, batches [][]operations.MediaFile) {
	var (
		args      []string
		arg       string
		batch     bool
		op        operations.Operation
		handler   filefmt.FileHandler
		fileError bool
		err       error
	)
	args = os.Args[1:]
	batches = append(batches, []operations.MediaFile{})
	for len(args) != 0 {
		if args[0] == "batch" {
			if batch {
				fmt.Fprintf(os.Stderr, "ERROR: \"batch\" specified repeatedly\n")
				usage()
			}
			args = args[1:]
			batch = true
			continue
		}
		if op, args, err = operations.ParseOperation(args); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			usage()
		} else if op != nil {
			ops = append(ops, op)
			continue
		}
		arg, args = args[0], args[1:]
		if _, err := os.Stat(arg); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			usage()
		}
		if handler = filefmt.HandlerFor(arg); handler == nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s: unsupported file format\n", arg)
			fileError = true
			continue
		}
		handler.ReadMetadata()
		if problems := handler.Problems(); len(problems) != 0 {
			for _, problem := range problems {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", arg, problem)
			}
			fileError = true
			continue
		}
		batches[0] = append(batches[0], operations.MediaFile{Path: arg, Handler: handler})
	}
	if len(batches[0]) == 0 && !fileError {
		fmt.Fprintf(os.Stderr, "ERROR: no files specified\n")
		usage()
	}
	if len(batches[0]) == 0 {
		os.Exit(1)
	}
	if len(ops) == 0 {
		op, _, _ = operations.ParseOperation([]string{"show"})
		ops = append(ops, op)
	}
	if batch {
		var (
			bnum = 0
			fnum = 1
			base = basename(batches[0][0].Path)
		)
		for fnum < len(batches[bnum]) {
			if nb := basename(batches[bnum][fnum].Path); nb == base {
				fnum++
			} else {
				batches = append(batches, batches[bnum][fnum:])
				batches[bnum] = batches[bnum][:fnum]
				base, bnum, fnum = nb, bnum+1, 1
			}
		}
	}
	for _, op := range ops {
		if err := op.Check(batches); err != nil {
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
