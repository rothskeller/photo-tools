// Package operations has definitions and code for each of the supported
// operations.
package operations

import (
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
)

// Operation is the interface satisfied by all operation handlers.
type Operation interface {
	// parseArgs parses the arguments for the operation, returning the
	// remaining argument list or an error.
	parseArgs(args []string) (remainingArgs []string, err error)
	// Check verifies that the operation is valid for the listed batches of
	// media files.  (Some operations require certain numbers of batches,
	// certain numbers of files per batch, certain media types, etc.).
	Check(batches [][]MediaFile) error
	// Run executes the operation against the listed media files (one
	// batch).  It returns whether any changes were made, or an error if one
	// occurred.
	Run(files []MediaFile) error
}

// MediaFile identifies, and provides the handler for, one media file named on
// the command line.
type MediaFile struct {
	Path    string
	Handler filefmt.FileHandler
	Changed bool
}

func escapeString(s string) string {
	return strings.Replace(strings.Replace(s, "\\", "\\\\", -1), "\n", "\\n", -1)
}
