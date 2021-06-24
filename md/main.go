// md is a program for viewing and editing media file metadata, tailored to the
// metadata conventions in my library.
package main

import (
	"fmt"
	"os"

	"github.com/rothskeller/photo-tools/filefmt"
)

type fileHandler = filefmt.FileHandler // copied here to save typing

type mediafile struct {
	path    string
	handler fileHandler
}

type operation interface {
	check(batches [][]mediafile) error
	run(files []mediafile) error
}

func main() {
	var (
		ops      []operation
		batches  [][]mediafile
		sawError bool
	)
	ops, batches = parseCommandLine() // exits on error
	for _, batch := range batches {
		for _, op := range ops {
			if err := op.run(batch); err != nil {
				panic("not sure how to handle errors here; revisit once I know what's possible") // TODO
			}
		}
		for _, file := range batch {
			if err := file.handler.SaveMetadata(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file, err)
				sawError = true
			}
		}
	}
	if sawError {
		os.Exit(1)
	}
}
