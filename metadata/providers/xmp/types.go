package xmp

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
func getAlt(val rdf.Value) (as altString, err error) {
	if val.Value == nil {
		return nil, nil
	}
	if val, ok := val.Value.(rdf.Alt); ok {
		as = make(altString, 0, len(val))
		for _, str := range val {
			var lang string
			if lt, ok := str.Qualifiers[xmlLang]; ok {
				if lt, ok := lt.Value.(string); ok {
					lang = lt
				} else {
					return nil, errors.New("wrong data type")
				}
			}
			if str, ok := str.Value.(string); ok {
				as = append(as, altItem{Value: str, Lang: lang})
			} else {
				return nil, errors.New("wrong data type")
			}
		}
		return as, nil
	}
	return nil, errors.New("wrong data type")
}

// makeAlt makes an rdf.Value with a language alternative string.
func makeAlt(as altString) rdf.Value {
	var values = make([]rdf.Value, len(as))
	for i := range as {
		values[i] = rdf.Value{
			Qualifiers: rdf.Struct{xmlLang: rdf.Value{Value: as[i].Lang}},
			Value:      as[i].Value,
		}
	}
	return rdf.Value{Value: rdf.Alt(values)}
}

// makeBag creates an unordered array of text values for the XMP.
func makeBag(bag []string) rdf.Value {
	var values = make([]rdf.Value, len(bag))
	for i := range bag {
		values[i] = rdf.Value{Value: bag[i]}
	}
	return rdf.Value{Value: rdf.Bag(values)}
}

// makeSeq creates an ordered array of text values for the XMP.
func makeSeq(seq []string) rdf.Value {
	var values = make([]rdf.Value, len(seq))
	for i := range seq {
		values[i] = rdf.Value{Value: seq[i]}
	}
	return rdf.Value{Value: rdf.Seq(values)}
}

// getString returns the value of a simple string from the XMP.
func getString(val rdf.Value) (str string, err error) {
	if val.Value == nil {
		return "", nil
	}
	if val, ok := val.Value.(string); ok {
		return val, nil
	}
	return "", errors.New("wrong data type")
}

// makeString creates a string value for the XMP.
func makeString(str string) rdf.Value { return rdf.Value{Value: str} }

// getStrings returns the value of an array of text values from the XMP.  It
// accepts either Bag or Seq, or a single string value.
func getStrings(val rdf.Value) (list []string, err error) {
	var vals []rdf.Value
	switch val := val.Value.(type) {
	case nil:
		return nil, nil
	case rdf.Seq:
		vals = val
	case rdf.Bag:
		vals = val
	case string:
		return []string{val}, nil
	default:
		return nil, errors.New("wrong data type")
	}
	list = make([]string, 0, len(vals))
	for _, str := range vals {
		switch str := str.Value.(type) {
		case string:
			list = append(list, str)
		default:
			return nil, errors.New("wrong data type")
		}
	}
	return list, nil
}
