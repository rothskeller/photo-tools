// md is a program for viewing and editing media file metadata, tailored to the
// metadata conventions in my library.  See MANUAL.md for usage and details.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"sort"

	"github.com/rothskeller/photo-tools/md/operations"
	"github.com/rothskeller/photo-tools/metadata/filefmts"
)

func main() {
	var (
		args            []string
		fnames          []string
		files           []operations.MediaFile
		sawError        bool
		ignoreNoHandler bool
		disallowWrites  bool
		saveMetadata    bool
		saveSet         bool
		err             error
	)
	// First, check for files given on the command line.
	args = os.Args[1:]
	for len(args) != 0 {
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			break
		}
		fnames = append(fnames, args[0])
		args = args[1:]
		saveSet = true
	}
	// If no files on command line, check for file selection keyword.
	if len(fnames) == 0 && len(args) != 0 {
		switch args[0] {
		case "all", "al":
			args = args[1:]
			fnames, err = restoreFullSet()
		case "batch", "b", "ba", "bat", "batc":
			args = args[1:]
			fnames, err = getFirstBatch()
		case "next", "n", "ne", "nex":
			args = args[1:]
			fnames, err = getNextBatch()
		case "prev", "p", "pr", "pre":
			args = args[1:]
			fnames, err = getPrevBatch()
		case "select", "sel", "sele", "selec":
			args = args[1:]
			fnames, err = selectSubset()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	// If no files or selection keyword on command line, get targeted
	// subset of remembered set.
	if len(fnames) == 0 {
		fnames = getTargetedFiles()
	}
	// If no remembered set, read the current directory.
	if len(fnames) == 0 {
		var dirents []fs.DirEntry

		if dirents, err = os.ReadDir("."); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		for _, dirent := range dirents {
			fnames = append(fnames, dirent.Name())
		}
		sort.Strings(fnames)
		ignoreNoHandler, disallowWrites, saveSet = true, true, true
	}
	// Get a handler and read the metadata for each identified file.
	for _, fname := range fnames {
		var (
			fh      *os.File
			handler filefmts.FileFormat
		)
		if fh, err = os.Open(fname); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			continue
		}
		if handler, err = filefmts.HandlerFor(fh); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			fh.Close()
			sawError = true
			continue
		}
		if handler == nil {
			if !ignoreNoHandler {
				fmt.Fprintf(os.Stderr, "ERROR: %s: not a supported file type\n", fname)
				sawError = true
			}
			fh.Close()
			continue
		}
		files = append(files, operations.MediaFile{Path: fname, File: fh, Provider: handler.Provider()})
	}
	// If no successfully read files, exit.
	if len(files) == 0 {
		if !sawError {
			fmt.Fprintln(os.Stderr, "ERROR: no files to act on")
		}
		os.Exit(1)
	}
	// If we have a new set (from files on command line or reading current
	// directory), save it.
	if saveSet {
		fnames = make([]string, len(files))
		for i := range files {
			fnames[i] = files[i].Path
		}
		writeMDFile(fnames)
	}
	// Choose an operation.
	if len(args) == 0 {
		err = operations.Check(args, files)
	} else {
		switch args[0] {
		case "add", "ad",
			"choose", "cho", "choo", "choos",
			"clear", "cl", "cle", "clea", "clr",
			"copy", "co", "cop", "cp",
			"remove", "rem", "remo", "remov", "rm",
			"reset", "res", "rese",
			"set", "se",
			"write", "w", "wr", "wri", "writ":
			if disallowWrites {
				fmt.Fprintf(os.Stderr, "ERROR: %q operation not allowed when defaulting to all files in directory\n", args[0])
				os.Exit(2)
			}
			saveMetadata = true
		}
		switch args[0] {
		case "add", "ad":
			err = operations.Add(args[1:], files)
		case "check", "che", "chec", "chk":
			err = operations.Check(args[1:], files)
		case "choose", "cho", "choo", "choos":
			err = operations.Choose(args[1:], files)
		case "clear", "cl", "cle", "clea", "clr":
			err = operations.Clear(args[1:], files)
		case "copy", "co", "cop", "cp":
			err = operations.Copy(args[1:], files)
		case "read", "rea", "rd":
			err = operations.Read(args[1:], files)
		case "remove", "rem", "remo", "remov", "rm":
			err = operations.Remove(args[1:], files)
		case "reset", "res", "rese":
			err = operations.Reset(args[1:], files)
		case "set", "se":
			err = operations.Set(args[1:], files)
		case "show", "sh":
			err = operations.Show(args[1:], files)
		case "tags", "t", "ta", "tag":
			err = operations.Tags(args[1:], files)
		case "write", "w", "wr", "wri", "writ":
			err = operations.Write(args[1:], files)
		default:
			fmt.Fprintf(os.Stderr, "ERROR: %q is not a recognized operation\n", args[0])
			os.Exit(1)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if saveMetadata {
		for _, file := range files {
			if file.Changed {
				// if err := file.Handler.SaveMetadata(); err != nil {
				// 	fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file.Path, err)
				// 	sawError = true
				// }
			}
		}
	}
	if sawError {
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `
usage: md [file...] [operation]
       md [file-selection] [operation]
Selections: all batch next prev select
Operations: add check choose clear copy read remove reset set show tags write
Fields: artist caption datetime faces gps groups keywords location people
        places title topics
See MANUAL.md for more details.
`)
	os.Exit(2)
}
