package metadata

import "trimmer.io/go-xmp/xmp"

// AltStringArray is an array of AltStrings, i.e., multiple strings each of
// which may have language alternatives.
type AltStringArray = xmp.AltStringArray

// AltString is a string with language alternatives, one of which is a default.
type AltString = xmp.AltString

// AltItem is one language alternative within an AltString.
type AltItem = xmp.AltItem

// NewAltString creates a new AltString, with one alternative for each
// argument.  The alternatives can be strings, AltItems, or Stringers.
func NewAltString(items ...interface{}) AltString {
	return xmp.NewAltString(items...)
}

// CopyAltString returns a deep copy of the provide AltString.
func CopyAltString(as AltString) (nas AltString) {
	nas = make(AltString, len(as))
	for i := range as {
		nas[i] = as[i]
	}
	return nas
}

// EmptyAltString returns true if the AltString contains no non-empty values.
func EmptyAltString(as AltString) bool {
	for _, ai := range as {
		if ai.Value != "" {
			return false
		}
	}
	return true
}

// EqualAltStrings returns whether twe AltStrings are equal.
func EqualAltStrings(a, b AltString) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].IsDefault != b[i].IsDefault ||
			a[i].Lang != b[i].Lang ||
			a[i].Value != b[i].Value {
			return false
		}
	}
	return true
}
