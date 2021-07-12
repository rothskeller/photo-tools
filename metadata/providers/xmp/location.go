package xmp

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	locationCreatedName = rdf.Name{Namespace: nsIPTC, Name: "LocationCreated"}
	locationShownName   = rdf.Name{Namespace: nsIPTC, Name: "LocationShown"}
	countryCodeName     = rdf.Name{Namespace: nsIPTC, Name: "CountryCode"}
	countryNameName     = rdf.Name{Namespace: nsIPTC, Name: "CountryName"}
	provinceStateName   = rdf.Name{Namespace: nsIPTC, Name: "ProvinceState"}
	cityName            = rdf.Name{Namespace: nsIPTC, Name: "City"}
	sublocationName     = rdf.Name{Namespace: nsIPTC, Name: "Sublocation"}
)

// getLocation reads the value of the Location field from the RDF.
func (p *Provider) getLocation() (err error) {
	if val := p.rdf.Property(locationCreatedName); val.Value != nil {
		switch val := val.Value.(type) {
		case rdf.Struct:
			if p.iptcLocationCreated, err = getLocationFromStruct(val); err != nil {
				return fmt.Errorf("Iptc4xmpExt:LocationCreated: %s", err)
			}
		default:
			return errors.New("Iptc4xmpExt:LocationCreated: wrong data type")
		}
	}
	if val := p.rdf.Property(locationShownName); val.Value != nil {
		switch val := val.Value.(type) {
		case rdf.Seq:
			p.iptcLocationsShown = make([]location, 0, len(val))
			for idx, loc := range val {
				switch loc := loc.Value.(type) {
				case rdf.Struct:
					if ls, err := getLocationFromStruct(loc); err == nil {
						p.iptcLocationsShown = append(p.iptcLocationsShown, ls)
					} else {
						return fmt.Errorf("Iptc4xmpExt:LocationShown[%d]: %s", idx, err)
					}
				default:
					return fmt.Errorf("Iptc4xmpExt:LocationShown[%d]: wrong data type", idx)
				}
			}
		default:
			return errors.New("Iptc4xmpExt:LocationShown: wrong data type")
		}
	}
	return nil
}
func getLocationFromStruct(str rdf.Struct) (loc location, err error) {
	if loc.CountryCode, err = getString(str[countryCodeName]); err != nil {
		return location{}, fmt.Errorf("CountryCode: %s", err)
	}
	if loc.CountryName, err = getAlt(str[countryNameName]); err != nil {
		return location{}, fmt.Errorf("CountryName: %s", err)
	}
	if loc.State, err = getAlt(str[provinceStateName]); err != nil {
		return location{}, fmt.Errorf("ProvinceState: %s", err)
	}
	if loc.City, err = getAlt(str[cityName]); err != nil {
		return location{}, fmt.Errorf("City: %s", err)
	}
	if loc.Sublocation, err = getAlt(str[sublocationName]); err != nil {
		return location{}, fmt.Errorf("Sublocation: %s", err)
	}
	return loc, nil
}

// Location returns the value of the Location field.
func (p *Provider) Location() (value metadata.Location) {
	var loc location
	if !p.iptcLocationCreated.Empty() {
		loc = p.iptcLocationCreated
	} else if len(p.iptcLocationsShown) != 0 {
		loc = p.iptcLocationsShown[0]
	} else {
		return metadata.Location{}
	}
	return metadata.Location{
		CountryCode: loc.CountryCode,
		CountryName: loc.CountryName.Default(),
		State:       loc.State.Default(),
		City:        loc.City.Default(),
		Sublocation: loc.Sublocation.Default(),
	}
}

// LocationTags returns a list of tag names for the Location field, and a
// parallel list of values held by those tags.
func (p *Provider) LocationTags() (tags []string, values [][]metadata.Location) {
	tags = []string{"XMP  iptc:LocationCreated"}
	values = [][]metadata.Location{locationToValues(p.iptcLocationCreated)}
	if len(p.iptcLocationsShown) != 0 {
		tags = append(tags, "XMP  iptc:LocationShown")
		values = append(values, nil)
		for _, shown := range p.iptcLocationsShown {
			values[1] = append(values[1], locationToValues(shown)...)
		}
	}
	return tags, values
}
func locationToValues(loc location) (values []metadata.Location) {
	// What languages are used in the location?
	var langs []string
	for _, ai := range loc.CountryName {
		langs = addUnique(langs, ai.Lang)
	}
	for _, ai := range loc.State {
		langs = addUnique(langs, ai.Lang)
	}
	for _, ai := range loc.City {
		langs = addUnique(langs, ai.Lang)
	}
	for _, ai := range loc.Sublocation {
		langs = addUnique(langs, ai.Lang)
	}
	// Make a location for each language.
	for _, lang := range langs {
		var mdl metadata.Location
		mdl.CountryCode = loc.CountryCode
		if mdl.CountryName = loc.CountryName.Get(lang); mdl.CountryName == "" {
			mdl.CountryName = loc.CountryName.Default()
		}
		if mdl.State = loc.State.Get(lang); mdl.State == "" {
			mdl.State = loc.CountryName.Default()
		}
		if mdl.City = loc.City.Get(lang); mdl.City == "" {
			mdl.City = loc.CountryName.Default()
		}
		if mdl.Sublocation = loc.Sublocation.Get(lang); mdl.Sublocation == "" {
			mdl.Sublocation = loc.CountryName.Default()
		}
		if loc.Empty() {
			continue
		}
		values = append(values, mdl)
	}
	return values
}
func addUnique(list []string, val string) []string {
	for _, exist := range list {
		if exist == val {
			return list
		}
	}
	return append(list, val)
}

// SetLocation sets the value of the Location field.
func (p *Provider) SetLocation(value metadata.Location) error {
	p.iptcLocationsShown = nil
	p.rdf.RemoveProperty(locationShownName)
	if value.Empty() {
		p.iptcLocationCreated = location{}
		p.rdf.RemoveProperty(locationCreatedName)
		return nil
	}
	if value.CountryCode != p.iptcLocationCreated.CountryCode {
		goto DIFFERENT
	}
	switch len(p.iptcLocationCreated.CountryName) {
	case 0:
		if value.CountryName != "" {
			goto DIFFERENT
		}
	case 1:
		if value.CountryName != p.iptcLocationCreated.CountryName.Default() {
			goto DIFFERENT
		}
	default:
		goto DIFFERENT
	}
	switch len(p.iptcLocationCreated.State) {
	case 0:
		if value.State != "" {
			goto DIFFERENT
		}
	case 1:
		if value.State != p.iptcLocationCreated.State.Default() {
			goto DIFFERENT
		}
	default:
		goto DIFFERENT
	}
	switch len(p.iptcLocationCreated.City) {
	case 0:
		if value.City != "" {
			goto DIFFERENT
		}
	case 1:
		if value.City != p.iptcLocationCreated.City.Default() {
			goto DIFFERENT
		}
	default:
		goto DIFFERENT
	}
	switch len(p.iptcLocationCreated.Sublocation) {
	case 0:
		if value.Sublocation != "" {
			goto DIFFERENT
		}
	case 1:
		if value.Sublocation != p.iptcLocationCreated.Sublocation.Default() {
			goto DIFFERENT
		}
	default:
		goto DIFFERENT
	}
	return nil
DIFFERENT:
	p.iptcLocationCreated = location{
		CountryCode: value.CountryCode,
		CountryName: newAltString(value.CountryName),
		State:       newAltString(value.State),
		City:        newAltString(value.City),
		Sublocation: newAltString(value.Sublocation),
	}
	var str rdf.Struct
	str[countryCodeName] = makeString(p.iptcLocationCreated.CountryCode)
	str[countryNameName] = makeAlt(p.iptcLocationCreated.CountryName)
	str[provinceStateName] = makeAlt(p.iptcLocationCreated.State)
	str[cityName] = makeAlt(p.iptcLocationCreated.City)
	str[sublocationName] = makeAlt(p.iptcLocationCreated.Sublocation)
	p.rdf.SetProperty(locationCreatedName, rdf.Value{Value: str})
	return nil
}
