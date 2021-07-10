package xmpiptc

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var xmlLang = rdf.Name{Namespace: rdf.NSxml, Name: "lang"}

// An altString is a set of language alternatives for a single conceptual
// string.  The first alternative is the default language.
type altString []altItem

// An altItem is a single language variant of an altString.
type altItem struct {
	Value string
	Lang  string
}

// newAltString creates a new altString, with a single default alternative.
func newAltString(s string) altString {
	return altString{{s, "x-default"}}
}

// Empty returns true if the AltString contains no non-empty values.
func (as altString) Empty() bool {
	for _, ai := range as {
		if ai.Value != "" {
			return false
		}
	}
	return true
}

// Default returns the default string from the altString.
func (as altString) Default() string {
	if len(as) == 0 {
		return ""
	}
	return as[0].Value
}

// Get returns the value of the AltString for the specified language.
func (as altString) Get(lang string) string {
	for _, alt := range as {
		if alt.Lang == lang {
			return alt.Value
		}
	}
	return ""
}

// A Location is the textual description of a location, with language
// alternative strings.
type location struct {
	CountryCode string
	CountryName altString
	State       altString
	City        altString
	Sublocation altString
}

// Empty returns true if the value contains no data.
func (loc location) Empty() bool {
	return loc.CountryCode == "" &&
		loc.CountryName.Empty() &&
		loc.State.Empty() &&
		loc.City.Empty() &&
		loc.Sublocation.Empty()
}

// getAlt returns the value of a Language Alternative value from the XMP.
func getAlt(from rdf.Struct, name rdf.Name) (as altString, err error) {
	if val, ok := from[name]; ok {
		switch val := val.Value.(type) {
		case rdf.Alt:
			as = make(altString, 0, len(val))
			for _, str := range val {
				var lang string
				if lt, ok := str.Qualifiers[xmlLang]; ok {
					switch lt := lt.Value.(type) {
					case string:
						lang = lt
					}
				}
				switch str := str.Value.(type) {
				case string:
					as = append(as, altItem{Value: str, Lang: lang})
				default:
					return nil, errors.New("wrong data type")
				}
			}
		default:
			return nil, errors.New("wrong data type")
		}
	}
	return as, nil
}

// setAlt sets a language alternative value in the XMP.
func setAlt(in rdf.Struct, name rdf.Name, as altString) {
	if len(as) == 0 {
		delete(in, name)
	} else {
		var values = make([]rdf.Value, len(as))
		for i := range as {
			values[i] = rdf.Value{
				Qualifiers: rdf.Struct{xmlLang: rdf.Value{Value: as[i].Lang}},
				Value:      as[i].Value,
			}
		}
		in[name] = rdf.Value{Value: rdf.Alt(values)}
	}
}

// getString returns the value of a simple string from the XMP.
func getString(from rdf.Struct, name rdf.Name) (str string, err error) {
	if val, ok := from[name]; ok {
		switch val := val.Value.(type) {
		case string:
			str = val
		default:
			return "", errors.New("wrong data type")
		}
	}
	return str, nil
}

// setString sets a string value in the XMP.
func setString(in rdf.Struct, name rdf.Name, str string) {
	if str == "" {
		delete(in, name)
	} else {
		in[name] = rdf.Value{Value: str}
	}
}
