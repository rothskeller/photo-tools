package metadata

import (
	"strings"
)

// String implements the Metadatum interface for a plain string.  The string may
// have a maximum length.
type String struct {
	val string
	max int
}

// NewString returns a new String with the specified value and unlimited length.
func NewString(s string) (ms *String) {
	if s != "" {
		return &String{s, 0}
	}
	return nil
}

// Parse sets the value from the input string.  It returns an error if the input
// was invalid.
func (s *String) Parse(v string) error {
	if s.max != 0 && len(v) > s.max {
		s.val = v[:s.max]
	} else {
		s.val = v
	}
	return nil
}

// String returns the value in string form, suitable for input to Parse.
func (s *String) String() string {
	if s == nil {
		return ""
	}
	return string(s.val)
}

// SetMaxLength sets the maximum length of the string.
func (s *String) SetMaxLength(max int) {
	s.max = max
	if len(s.val) > s.max {
		s.val = s.val[:s.max]
	}
}

// Empty returns true if the value contains no data.
func (s *String) Empty() bool {
	return s == nil || s.val == ""
}

// Equal returns true if the receiver is equal to the argument.
func (s *String) Equal(other Metadatum) bool {
	if s == nil && other == nil {
		return true
	}
	switch other := other.(type) {
	case *String:
		if (s == nil) != (other == nil) {
			return false
		}
		if s == nil {
			return true
		}
		return s.val == other.val
	default:
		return false
	}
}

// Equivalent returns true if the receiver is equal to the argument, to the
// precision of the least precise of the two.  If so, the second return value is
// the more precise of the two.
func (s *String) Equivalent(other Metadatum) (bool, Metadatum) {
	if s == nil && other == nil {
		return true, nil
	}
	switch other := other.(type) {
	case *String:
		if (s == nil) != (other == nil) {
			return false, nil
		}
		if s == nil {
			return true, nil
		}
		if s.val == other.val {
			return true, s
		}
		if other.max != 0 && len(other.val) == other.max && strings.HasPrefix(s.val, other.val) {
			return true, s
		}
		if s.max != 0 && len(s.val) == s.max && strings.HasPrefix(other.val, s.val) {
			return true, other
		}
		return false, nil
	default:
		return false, nil
	}
}

// Verify interface compliance.
var _ Metadatum = (*String)(nil)
