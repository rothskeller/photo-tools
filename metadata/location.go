package metadata

// A Location is the textual description of a location.
type Location struct {
	CountryCode string
	CountryName AltString
	State       AltString
	City        AltString
	Sublocation AltString
}

// Empty returns true if the value contains no data.
func (loc Location) Empty() bool {
	return loc.CountryCode == "" &&
		EmptyAltString(loc.CountryName) &&
		EmptyAltString(loc.State) &&
		EmptyAltString(loc.City) &&
		EmptyAltString(loc.Sublocation)
}

// Equal returns true if the receiver is equal to the argument.
func (loc Location) Equal(other Location) bool {
	return loc.CountryCode == other.CountryCode &&
		EqualAltStrings(loc.CountryName, other.CountryName) &&
		EqualAltStrings(loc.State, other.State) &&
		EqualAltStrings(loc.City, other.City) &&
		EqualAltStrings(loc.Sublocation, other.Sublocation)
}

// Copy copies the argument Location into the receiver Location.
func (loc Location) Copy(other Location) {
	loc.CountryCode = other.CountryCode
	loc.CountryName = CopyAltString(other.CountryName)
	loc.State = CopyAltString(other.State)
	loc.City = CopyAltString(other.City)
	loc.Sublocation = CopyAltString(other.Sublocation)
}
