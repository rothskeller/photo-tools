package strmeta

// removeEmpty returns a new slice containing all non-empty strings from the
// argument slice.
func removeEmpty(sl []string) (out []string) {
	out = make([]string, 0, len(sl))
	for _, s := range sl {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// removeDuplicates returns a new slice containing the unique strings from the
// argument slice, in the order they first appear in the argument slice.
func removeDuplicates(sl []string) (out []string) {
	out = make([]string, 0, len(sl))
	sm := make(map[string]bool)
	for _, s := range sl {
		if !sm[s] {
			out = append(out, s)
			sm[s] = true
		}
	}
	return out
}

// equalMaxLen returns true if the two argument strings are equal, or if the
// second one is equal to the first one, truncated to maxlen.
func equalMaxLen(a, b string, maxlen int) bool {
	if a == b {
		return true
	}
	if len(a) > maxlen && len(b) == maxlen && a[:maxlen] == b {
		return true
	}
	return false
}

// listMatchItem returns whether the list matches the value.  It matches if the
// list contains only one item that is equal to the value, or if the value and
// the list are both empty.  Equality comparison is done using equalMaxLen if
// maxlen is nonzero.
func listMatchItem(list []string, val string, maxlen int) bool {
	if len(list) == 0 {
		return val == ""
	}
	if len(list) > 1 {
		return false
	}
	if list[0] == val {
		return true
	}
	if maxlen > 0 && len(val) > maxlen && list[0] == val[:maxlen] {
		return true
	}
	return false
}
