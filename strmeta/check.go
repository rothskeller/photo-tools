package strmeta

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
