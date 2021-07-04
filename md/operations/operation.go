// Package operations has definitions and code for each of the supported
// operations.
package operations

import (
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
)

// MediaFile identifies, and provides the handler for, one media file named on
// the command line.
type MediaFile struct {
	Path    string
	Handler filefmt.FileHandler
	Changed bool
}

// escapeString escapes newlines in tabular output from the show and tabs
// operations.
func escapeString(s string) string {
	return strings.Replace(strings.Replace(s, "\\", "\\\\", -1), "\n", "\\n", -1)
}
