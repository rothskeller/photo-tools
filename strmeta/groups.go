package strmeta

import (
	"errors"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
)

// A Group represents a group of people (team, organization, etc.) depicted in a
// media.
type Group []string

// Parse parses a group name, as a hierarchical string with levels separated by
// slashes and optional whitespace.  Pipe symbols are disallowed due to
// underlying storage formats, and empty levels are disallowed (although a
// completely empty string is allowed).
func (g *Group) Parse(s string) error {
	kw, err := metadata.ParseKeyword(s, "")
	if err == nil {
		*g = Group(kw)
	}
	return err
}

// String returns the formatted string form of the group name, suitable for
// input to Parse.
func (g Group) String() string { return metadata.Keyword(g).String() }

// Empty returns whether the group name is empty.
func (g Group) Empty() bool { return len(g) == 0 }

// Equal returns whether two group names are equal.
func (g Group) Equal(other Group) bool {
	return metadata.Keyword(g).Equal(metadata.Keyword(other))
}

// GetGroups returns the highest priority group values.
func GetGroups(h filefmt.FileHandler) []Group {
	kws := getFilteredKeywords(h, groupPredicate, false)
	groups := make([]Group, len(kws))
	for i := range kws {
		groups[i] = Group(kws[i][1:])
	}
	return groups
}

// GetGroupTags returns all of the group tags and their values.
func GetGroupTags(h filefmt.FileHandler) (tags []string, values []Group) {
	tags, kws := getFilteredKeywordTags(h, groupPredicate)
	values = make([]Group, len(kws))
	for i := range kws {
		values[i] = Group(kws[i][1:])
	}
	return tags, values
}

// CheckGroups determines whether the groups are tagged correctly.
func CheckGroups(h filefmt.FileHandler) CheckResult {
	return checkFilteredKeywords(h, groupPredicate)
}

// SetGroups sets the group tags.
func SetGroups(h filefmt.FileHandler, v []Group) error {
	var kws = make([]metadata.Keyword, len(v))
	for i, g := range v {
		if g.Empty() {
			return errors.New("empty group name not allowed")
		}
		kws[i] = append(metadata.Keyword{"Groups"}, v[i]...)
	}
	return setFilteredKeywords(h, kws, groupPredicate)
}

// groupPredicate is the predicate satisfied by keyword tags that encode group
// names.
func groupPredicate(kw metadata.Keyword) bool {
	return len(kw) >= 2 && kw[0] == "Groups"
}
