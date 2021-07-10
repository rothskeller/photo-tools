package xmpps

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

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
