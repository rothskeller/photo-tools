package metadata

import (
	"errors"
	"strings"
)

// A HierValue represents a hierarchical value.
type HierValue []string

// ParseHierValue parses a hierarchical value.
func ParseHierValue(s string) (hv HierValue, err error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	var words = strings.Split(s, "/")
	for _, w := range words {
		if strings.IndexByte(w, '|') >= 0 {
			return nil, errors.New("hierarchical values cannot contain | characters")
		}
		if w := strings.TrimSpace(w); w != "" {
			hv = append(hv, w)
		} else {
			return nil, errors.New("hierarchical values cannot have empty components")
		}
	}
	return hv, nil
}

func (hv HierValue) String() string {
	return strings.Join(hv, " / ")
}

// Equal returns whether two hierarchical values are the same.
func (hv HierValue) Equal(other HierValue) bool {
	if len(hv) != len(other) {
		return false
	}
	for i := range hv {
		if hv[i] != other[i] {
			return false
		}
	}
	return true
}
