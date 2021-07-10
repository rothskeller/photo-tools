package metadata

// An AltString is a set of language alternatives for a single conceptual
// string.  The first alternative is the default language.
type AltString []AltItem

// An AltItem is a single language variant of an AltString.
type AltItem struct {
	Value string
	Lang  string
}

// NewAltString creates a new AltString, with a single default alternative.
func NewAltString(s string) AltString {
	return AltString{{s, "x-default"}}
}

// Copy returns a deep copy of the provide AltString.
func (as AltString) Copy() (nas AltString) {
	nas = make(AltString, len(as))
	for i := range as {
		nas[i] = as[i]
	}
	return nas
}

// Empty returns true if the AltString contains no non-empty values.
func (as AltString) Empty() bool {
	for _, ai := range as {
		if ai.Value != "" {
			return false
		}
	}
	return true
}

// Equal returns whether twe AltStrings are equal.
func (as AltString) Equal(other AltString) bool {
	if len(as) != len(other) {
		return false
	}
	for i := range as {
		if as[i].Lang != other[i].Lang ||
			as[i].Value != other[i].Value {
			return false
		}
	}
	return true
}

// Default returns the default string from the AltString.
func (as AltString) Default() string {
	if len(as) == 0 {
		return ""
	}
	return as[0].Value
}

// Get returns the value of the AltString for the specified language.
func (as AltString) Get(lang string) string {
	for _, alt := range as {
		if alt.Lang == lang {
			return alt.Value
		}
	}
	return ""
}
