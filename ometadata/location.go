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
		loc.CountryName.Empty() &&
		loc.State.Empty() &&
		loc.City.Empty() &&
		loc.Sublocation.Empty()
}

// Equal returns true if the receiver is equal to the argument.
func (loc Location) Equal(other Location) bool {
	return loc.CountryCode == other.CountryCode &&
		loc.CountryName.Equal(other.CountryName) &&
		loc.State.Equal(other.State) &&
		loc.City.Equal(other.City) &&
		loc.Sublocation.Equal(other.Sublocation)
}

// Copy copies the argument Location into the receiver Location.
func (loc Location) Copy(other Location) {
	loc.CountryCode = other.CountryCode
	loc.CountryName = other.CountryName.Copy()
	loc.State = other.State.Copy()
	loc.City = other.City.Copy()
	loc.Sublocation = other.Sublocation.Copy()
}
