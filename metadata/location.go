package metadata

import (
	"errors"
	"regexp"
	"strings"
)

// Location contains a location where the media was captured, or which is
// depicted in the media.  It's a simplified form of metadata.Location, with all
// language alternatives removed.
type Location struct {
	CountryCode string
	CountryName string
	State       string
	City        string
	Sublocation string
}

var countryCodeRE = regexp.MustCompile(`^[A-Z]{2,3}$`)

// Parse sets the value from the input string.
func (loc *Location) Parse(val string) (err error) {
	*loc = Location{}
	parts := strings.Split(val, "/")
	if len(parts) > 5 {
		return errors.New("too many components in location")
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	for len(parts) < 5 {
		parts = append(parts, "")
	}
	seenEmpty := false
	for _, part := range parts {
		if part == "" {
			seenEmpty = true
		} else if seenEmpty {
			return errors.New("missing component in location")
		}
	}
	loc.CountryCode = strings.ToUpper(parts[0])
	loc.CountryName = parts[1]
	loc.State = parts[2]
	loc.City = parts[3]
	loc.Sublocation = parts[4]
	if loc.CountryCode != "" && !countryCodeRE.MatchString(loc.CountryCode) {
		return errors.New("invalid country code in location")
	}
	return nil
}

// String returns the value in string form, suitable for input to Parse.
func (loc Location) String() string {
	var (
		sb strings.Builder
	)
	if loc.CountryCode != "" {
		sb.WriteString(loc.CountryCode)
		sb.WriteByte(' ')
	}
	sb.WriteByte('/')
	if loc.CountryName != "" {
		sb.WriteString(loc.CountryName)
		sb.WriteByte(' ')
	}
	sb.WriteByte('/')
	if loc.State != "" {
		sb.WriteString(loc.State)
		sb.WriteByte(' ')
	}
	sb.WriteByte('/')
	if loc.City != "" {
		sb.WriteString(loc.City)
		sb.WriteByte(' ')
	}
	sb.WriteByte('/')
	sb.WriteString(loc.Sublocation)
	return strings.TrimRight(sb.String(), "/ ")
}

// Empty returns true if the location has no data.
func (loc Location) Empty() bool {
	return loc.CountryCode == "" && loc.CountryName == "" && loc.State == "" && loc.City == "" && loc.Sublocation == ""
}

// Equal returns true if the two locations are the same.
func (loc Location) Equal(other Location) bool {
	return loc.CountryCode == other.CountryCode &&
		loc.CountryName == other.CountryName &&
		loc.State == other.State &&
		loc.City == other.City &&
		loc.Sublocation == other.Sublocation
}
