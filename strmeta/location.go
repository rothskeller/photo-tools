package strmeta

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/iptc"
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
func (loc *Location) String() string {
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
func (loc *Location) Empty() bool {
	if loc == nil {
		return true
	}
	return loc.CountryCode == "" && loc.CountryName == "" && loc.State == "" && loc.City == "" && loc.Sublocation == ""
}

// GetLocation returns the highest priority location value.
func GetLocation(h fileHandler) Location {
	xmp := h.XMP(false)
	if xmp != nil {
		if !xmp.IPTCLocationCreated().Empty() {
			return mdLocationToSTRLocation(xmp.IPTCLocationCreated())
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		var loc Location
		loc.CountryCode = iptc.CountryPLCode()
		loc.CountryName = iptc.CountryPLName()
		loc.State = iptc.ProvinceState()
		loc.City = iptc.City()
		loc.Sublocation = iptc.Sublocation()
		if !loc.Empty() {
			return loc
		}
	}
	if xmp != nil {
		for _, shown := range xmp.IPTCLocationsShown() {
			if !shown.Empty() {
				return mdLocationToSTRLocation(shown)
			}
		}
	}
	return Location{}
}
func mdLocationToSTRLocation(md metadata.Location) (str Location) {
	str.CountryCode = md.CountryCode
	str.CountryName = md.CountryName.Default()
	str.State = md.State.Default()
	str.City = md.City.Default()
	str.Sublocation = md.Sublocation.Default()
	return str
}

// GetLocationTags returns all of the location tags and their values.
func GetLocationTags(h fileHandler) (tags []string, values []Location) {
	if xmp := h.XMP(false); xmp != nil {
		tags, values = mdLocationToTags(tags, values, "XMP  iptc:LocationCreated", xmp.IPTCLocationCreated(), true)
		for _, shown := range xmp.IPTCLocationsShown() {
			tags, values = mdLocationToTags(tags, values, "XMP  iptc:LocationShown", shown, false)
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		var loc Location
		loc.CountryCode = iptc.CountryPLCode()
		loc.CountryName = iptc.CountryPLName()
		loc.State = iptc.ProvinceState()
		loc.City = iptc.City()
		loc.Sublocation = iptc.Sublocation()
		tags = append(tags, "IPTC Location")
		values = append(values, loc)
	}
	return tags, values
}
func mdLocationToTags(tags []string, values []Location, label string, md metadata.Location, addEmpty bool) ([]string, []Location) {
	// What languages are used in the location?
	var langs []string
	for _, ai := range md.CountryName {
		langs = addUnique(langs, ai.Lang)
	}
	for _, ai := range md.State {
		langs = addUnique(langs, ai.Lang)
	}
	for _, ai := range md.City {
		langs = addUnique(langs, ai.Lang)
	}
	for _, ai := range md.Sublocation {
		langs = addUnique(langs, ai.Lang)
	}
	// Make a location for each language.
	var added = false
	for _, lang := range langs {
		var loc Location
		loc.CountryCode = md.CountryCode
		if loc.CountryName = md.CountryName.Get(lang); loc.CountryName == "" {
			loc.CountryName = md.CountryName.Default()
		}
		if loc.State = md.State.Get(lang); loc.State == "" {
			loc.State = md.CountryName.Default()
		}
		if loc.City = md.City.Get(lang); loc.City == "" {
			loc.City = md.CountryName.Default()
		}
		if loc.Sublocation = md.Sublocation.Get(lang); loc.Sublocation == "" {
			loc.Sublocation = md.CountryName.Default()
		}
		if loc.Empty() {
			continue
		}
		if lang == "" {
			tags = append(tags, label)
		} else {
			tags = append(tags, fmt.Sprintf("%s[%s]", label, lang))
		}
		values = append(values, loc)
		added = true
	}
	if !added && addEmpty {
		tags = append(tags, label)
		values = append(values, Location{})
	}
	return tags, values
}
func addUnique(list []string, val string) []string {
	for _, exist := range list {
		if exist == val {
			return list
		}
	}
	return append(list, val)
}

// CheckLocation determines whether the location is tagged correctly.
func CheckLocation(h fileHandler) (res CheckResult) {
	var value = GetLocation(h)

	xmp := h.XMP(false)
	if xmp != nil {
		if !xmp.IPTCLocationCreated().Empty() {
			if xmp.IPTCLocationCreated().CountryCode != value.CountryCode {
				return ChkConflictingValues
			}
			switch len(xmp.IPTCLocationCreated().CountryName) {
			case 0:
				if value.CountryName != "" {
					return ChkConflictingValues
				}
			case 1:
				if value.CountryName != xmp.IPTCLocationCreated().CountryName[0].Value {
					return ChkConflictingValues
				}
			default:
				return ChkConflictingValues
			}
			switch len(xmp.IPTCLocationCreated().State) {
			case 0:
				if value.State != "" {
					return ChkConflictingValues
				}
			case 1:
				if value.State != xmp.IPTCLocationCreated().State[0].Value {
					return ChkConflictingValues
				}
			default:
				return ChkConflictingValues
			}
			switch len(xmp.IPTCLocationCreated().City) {
			case 0:
				if value.City != "" {
					return ChkConflictingValues
				}
			case 1:
				if value.City != xmp.IPTCLocationCreated().City[0].Value {
					return ChkConflictingValues
				}
			default:
				return ChkConflictingValues
			}
			switch len(xmp.IPTCLocationCreated().Sublocation) {
			case 0:
				if value.Sublocation != "" {
					return ChkConflictingValues
				}
			case 1:
				if value.Sublocation != xmp.IPTCLocationCreated().Sublocation[0].Value {
					return ChkConflictingValues
				}
			default:
				return ChkConflictingValues
			}
		}
		switch len(xmp.IPTCLocationsShown()) {
		case 0:
			break
		case 1:
			if !xmp.IPTCLocationCreated().Equal(xmp.IPTCLocationsShown()[0]) {
				return ChkConflictingValues
			}
			res = ChkIncorrectlyTagged
		default:
			return ChkConflictingValues
		}
	}
	if i := h.IPTC(); i != nil {
		if !stringEqualMax(value.CountryCode, i.CountryPLCode(), iptc.MaxCountryPLCodeLen) ||
			!stringEqualMax(value.CountryName, i.CountryPLName(), iptc.MaxCountryPLNameLen) ||
			!stringEqualMax(value.State, i.ProvinceState(), iptc.MaxProvinceStateLen) ||
			!stringEqualMax(value.City, i.City(), iptc.MaxCityLen) ||
			!stringEqualMax(value.Sublocation, i.Sublocation(), iptc.MaxSublocationLen) {
			if i.CountryPLCode() != "" || i.CountryPLName() != "" || i.ProvinceState() != "" || i.City() != "" ||
				i.Sublocation() != "" {
				return ChkConflictingValues
			}
			res = ChkIncorrectlyTagged
		}
	}
	if !value.Empty() && res == 0 {
		return ChkPresent
	}
	return res
}

// SetLocation sets the location tags.
func SetLocation(h fileHandler, v Location) error {
	if xmp := h.XMP(true); xmp != nil {
		if err := xmp.SetIPTCLocationCreated(strLocationToMDLocation(v)); err != nil {
			return err
		}
		if err := xmp.SetIPTCLocationsShown(nil); err != nil { // Always remove unwanted tag
			return err
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if err := iptc.SetCountryPLCode(v.CountryCode); err != nil {
			return err
		}
		if err := iptc.SetCountryPLName(v.CountryName); err != nil {
			return err
		}
		if err := iptc.SetProvinceState(v.State); err != nil {
			return err
		}
		if err := iptc.SetCity(v.City); err != nil {
			return err
		}
		if err := iptc.SetSublocation(v.Sublocation); err != nil {
			return err
		}
	}
	return nil
}
func strLocationToMDLocation(str Location) (md metadata.Location) {
	md.CountryCode = str.CountryCode
	if str.CountryName != "" {
		md.CountryName = metadata.NewAltString(str.CountryName)
	}
	if str.State != "" {
		md.State = metadata.NewAltString(str.State)
	}
	if str.City != "" {
		md.City = metadata.NewAltString(str.City)
	}
	if str.Sublocation != "" {
		md.Sublocation = metadata.NewAltString(str.Sublocation)
	}
	return md
}
