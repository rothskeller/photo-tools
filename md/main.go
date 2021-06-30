// md is a program for viewing and editing media file metadata, tailored to the
// metadata conventions in my library.  See MANUAL.md for usage and details.
package main

import (
	"fmt"
	"os"

	"github.com/rothskeller/photo-tools/md/operations"
)

func main() {
	var (
		ops      []operations.Operation
		batches  [][]operations.MediaFile
		sawError bool
	)
	ops, batches = parseCommandLine() // exits on error
	for _, batch := range batches {
		for _, op := range ops {
			if err := op.Run(batch); err != nil {
				panic("not sure how to handle errors here; revisit once I know what's possible") // TODO
			}
		}
		for _, file := range batch {
			if file.Changed {
				if err := file.Handler.SaveMetadata(); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file.Path, err)
					sawError = true
				}
			}
		}
	}
	if sawError {
		os.Exit(1)
	}
}
