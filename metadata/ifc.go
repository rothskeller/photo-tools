package metadata

// Metadatum is the interface that all metadata types honor.
type Metadatum interface {
	// Parse sets the value from the input string.  It returns an error if
	// the input was invalid.
	Parse(string) error
	// String returns the value in string form, suitable for input to Parse.
	String() string
	// Empty returns true if the value contains no data.
	Empty() bool
	// Equal returns true if the receiver is equal to the argument.
	Equal(Metadatum) bool
	// Equivalent returns true if the receiver is equal to the argument, to
	// the precision of the least precise of the two.  If so, the second
	// return value is the more precise of the two.
	Equivalent(Metadatum) (bool, Metadatum)
}

// TaggedMetadatum associates a Metadatum with the name of the metadata tag from
// which it was read.
type TaggedMetadatum struct {
	Metadatum
	Tag string
}
