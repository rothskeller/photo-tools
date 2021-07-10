package xmptiff

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

// Default returns the default string from the altString.
func (as altString) Default() string {
	if len(as) == 0 {
		return ""
	}
	return as[0].Value
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

// getStrings returns the value of an array of text values from the XMP.  It
// accepts either Bag or Seq, or a single string value.
func getStrings(from rdf.Struct, name rdf.Name) (list []string, err error) {
	var vals []rdf.Value
	if val, ok := from[name]; ok {
		switch val := val.Value.(type) {
		case rdf.Seq:
			vals = val
		case rdf.Bag:
			vals = val
		case string:
			return []string{val}, nil
		default:
			return nil, errors.New("wrong data type")
		}
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
