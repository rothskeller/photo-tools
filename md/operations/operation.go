// Package operations has definitions and code for each of the supported
// operations.
package operations

import (
	"os"
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
)

// MediaFile identifies, and provides the handler for, one media file named on
// the command line.
type MediaFile struct {
	Path     string
	File     *os.File
	Provider metadata.Provider
	Changed  bool
}

// escapeString escapes newlines in tabular output from the show and tabs
// operations.
func escapeString(s string) string {
	return strings.Replace(strings.Replace(s, "\\", "\\\\", -1), "\n", "\\n", -1)
}
