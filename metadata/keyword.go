package metadata

import (
	"errors"
	"strings"
)

// A Keyword represents a hierarchical keyword.
type Keyword []string

// ParseKeyword parses a keyword.  If prefix is non-nil, it is prepended to the
// keyword.
func ParseKeyword(s, prefix string) (kw Keyword, err error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	var words = strings.Split(s, "/")
	if prefix != "" {
		kw = append(kw, prefix)
	}
	for _, w := range words {
		if strings.IndexByte(w, '|') >= 0 {
			return nil, errors.New("keywords cannot contain | characters")
		}
		if w := strings.TrimSpace(w); w != "" {
			kw = append(kw, w)
		} else {
			return nil, errors.New("keywords cannot have empty components")
		}
	}
	return kw, nil
}

func (kw Keyword) String() string {
	return strings.Join(kw, " / ")
}

// StringWithoutPrefix returns the string form of the keyword with the expected
// initial prefix removed.
func (kw Keyword) StringWithoutPrefix(prefix string) string {
	if len(kw) == 0 {
		return ""
	}
	if kw[0] != prefix {
		panic("StringWithoutPrefix called on keyword with wrong prefix")
	}
	return kw[1:].String()
}

// Equal returns whether two keywords are the same.
func (kw Keyword) Equal(other Keyword) bool {
	if len(kw) != len(other) {
		return false
	}
	for i := range kw {
		if kw[i] != other[i] {
			return false
		}
	}
	return true
}
