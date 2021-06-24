package metadata

import (
	"errors"
	"strings"
)

type Keyword []KeywordComponent

type KeywordComponent struct {
	Word              string
	OmitWhenFlattened bool
}

// ParseKeyword parses a keyword.  If prefix is non-nil, it is prepended to the
// keyword as an omit-when-flattened.
func ParseKeyword(s, prefix string) (kw Keyword, err error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	var words = strings.Split(s, "/")
	if prefix != "" {
		kw = append(kw, KeywordComponent{prefix, true})
	}
	for _, w := range words {
		if strings.IndexByte(w, '|') >= 0 {
			return nil, errors.New("keywords cannot contain | characters")
		}
		if w := strings.TrimSpace(w); w != "" {
			kw = append(kw, KeywordComponent{Word: w})
		} else {
			return nil, errors.New("keywords cannot have empty components")
		}
	}
	return kw, nil
}

func (kw Keyword) String() string {
	var parts = make([]string, len(kw))
	for i := range kw {
		parts[i] = kw[i].Word
	}
	return strings.Join(parts, " / ")
}

// StringWithoutPrefix returns the string form of the keyword with the expected
// initial prefix removed.
func (kw Keyword) StringWithoutPrefix(prefix string) string {
	if len(kw) == 0 {
		return ""
	}
	if kw[0].Word != prefix {
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
		if kw[i].Word != other[i].Word {
			return false
		}
	}
	return true
}
