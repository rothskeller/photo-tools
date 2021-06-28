package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
    clear fieldname            l(ocation)
    choose fieldname           t(itle)
    copy [fieldname...]        person (people)
    write caption              group(s), place(s), topic(s)
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
		// It's not "batch" or an operation, so it should be a file name.
		arg, args = args[0], args[1:]
		if _, err := os.Stat(arg); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			usage()
		}
		// Open the file and read its metadata.
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
	// Sort the files (they're all in one batch so far) so that the batching
	// operation works correctly.
	sort.Slice(batches[0], func(a, b int) bool {
		return batchSortFn(batches[0][a], batches[0][b])
	})
	if batch { // batch the files if requested
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
	// Give each operation the chance to check the batches for validity.
	for _, op := range ops {
		if err := op.Check(batches); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			usage()
		}
	}
	return ops, batches
}

// Files are sorted by basename first, then by variant (the part between the
// first dot and the extension, if any), then by directory, and finally by
// extension.  This ensures, for example, that "foo.jpg" comes before
// "foo.aaa.jpg".
func batchSortFn(a, b operations.MediaFile) bool {
	var adir, abase, avariant, aext string
	var bdir, bbase, bvariant, bext string
	adir = filepath.Dir(a.Path)
	bdir = filepath.Dir(b.Path)
	abase = filepath.Base(a.Path)
	bbase = filepath.Base(b.Path)
	if strings.HasSuffix(abase, ".xmp") {
		abase, aext = abase[:len(abase)-4], ".xmp"
	}
	if strings.HasSuffix(bbase, ".xmp") {
		bbase, bext = bbase[:len(bbase)-4], ".xmp"
	}
	if ext := filepath.Ext(abase); ext != "" {
		aext = ext + aext
		abase = abase[:len(abase)-len(ext)]
	}
	if ext := filepath.Ext(bbase); ext != "" {
		bext = ext + bext
		bbase = bbase[:len(bbase)-len(ext)]
	}
	if idx := strings.IndexByte(abase, '.'); idx >= 0 {
		abase, avariant = abase[:idx], abase[idx:]
	}
	if idx := strings.IndexByte(bbase, '.'); idx >= 0 {
		bbase, bvariant = bbase[:idx], bbase[idx:]
	}
	switch {
	case abase < bbase:
		return true
	case bbase < abase:
		return false
	case avariant < bvariant:
		return true
	case bvariant < avariant:
		return false
	case adir < bdir:
		return true
	case bdir < adir:
		return false
	default:
		return aext < bext
	}
}

// basename returns the path name, shorn of any directory information and
// anything after the first period (including the period itself).
func basename(path string) string {
	path = filepath.Base(path)
	if idx := strings.IndexByte(path, '.'); idx >= 0 {
		path = path[:idx]
	}
	return path
}
