// Package strmeta contains the translation between actual metadata tags and the
// simplified metadata model used in my library.
package strmeta

import (
	"fmt"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
)

type fileHandler = filefmt.FileHandler // copied to save typing

// A CheckResult gives the result of a Check operation on a field of a file.
// A positive integer is a count of values in a multi-valued field that is
// tagged correctly.  All other values have defined constants.
type CheckResult int

// Values for CheckResult
const (
	ChkPresent           CheckResult = 1
	ChkOptionalAbsent    CheckResult = 0
	ChkExpectedAbsent    CheckResult = -1
	ChkIncorrectlyTagged CheckResult = -2
	ChkConflictingValues CheckResult = -3
)

// stringEqualMax compares two strings, with the second one subject to the
// specified maximum length, and returns whether they should be considered
// equal for check purposes.
func stringEqualMax(a, b string, bmax int) bool {
	if a == b {
		return true
	}
	if len(a) > bmax && a[:bmax] == b {
		return true
	}
	return false
}

func tagsForStringList(tags, values []string, label string, ss []string) (newt, newv []string) {
	if len(ss) == 0 {
		tags = append(tags, label)
		values = append(values, "")
	}
	for _, s := range ss {
		tags = append(tags, label)
		values = append(values, s)
	}
	return tags, values
}

func tagsForAltString(tags, values []string, label string, as metadata.AltString) (newt, newv []string) {
	if len(as) == 0 {
		tags = append(tags, label)
		values = append(values, "")
	}
	for i, ls := range as {
		if i == 0 && ls.Lang == "" {
			tags = append(tags, label)
		} else {
			tags = append(tags, fmt.Sprintf("%s[%s]", label, ls.Lang))
		}
		values = append(values, ls.Value)
	}
	return tags, values
}
