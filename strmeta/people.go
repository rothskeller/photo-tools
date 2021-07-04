package strmeta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
)

// GetPeople returns the highest priority person values.
func GetPeople(h filefmt.FileHandler) (people []string) {
	kws := getFilteredKeywords(h, personPredicate, false)
	pmap := make(map[string]bool)
	for _, kw := range kws {
		name := kw[1]
		if name != "" && !pmap[name] {
			people = append(people, name)
			pmap[name] = true
		}
	}
	for _, face := range GetFaces(h) {
		if !pmap[face] {
			people = append(people, face)
			pmap[face] = true
		}
	}
	return people
}

// GetPersonTags returns all of the person tags and their values.
func GetPersonTags(h filefmt.FileHandler) (tags []string, values []string) {
	tags, kws := getFilteredKeywordTags(h, personPredicate)
	values = make([]string, len(kws))
	for i := range kws {
		values[i] = kws[i][1]
	}
	return tags, values
}

// CheckPeople determines whether the people are tagged correctly.
func CheckPeople(h filefmt.FileHandler) (res CheckResult) {
	if res = checkFilteredKeywords(h, personPredicate); res == ChkConflictingValues {
		return res
	}
	// Also check whether all face regions have corresponding people tags.
	// If not, it's reported as IncorrectlyTagged.
	var faces = GetFaces(h)
	var people = getFilteredKeywords(h, personPredicate, false)
	for _, face := range faces {
		var found = false
		for _, person := range people {
			if person[1] == face {
				found = true
				break
			}
		}
		if !found {
			res = ChkIncorrectlyTagged
			break
		}
	}
	return res
}

// SetPeople sets the person tags.
func SetPeople(h filefmt.FileHandler, v []string) error {
	var (
		facelist []string
		kws      = make([]metadata.Keyword, len(v))
		names    = make(map[string]bool)
		faces    = make(map[string]bool)
	)
	for _, face := range GetFaces(h) {
		faces[face] = false
	}
	for i, g := range v {
		if strings.TrimSpace(g) == "" {
			return errors.New("empty person not allowed")
		}
		if names[g] {
			return fmt.Errorf("cannot list name %q twice", g)
		}
		names[g] = true
		if _, ok := faces[g]; ok {
			faces[g] = true
		}
		kws[i] = append(metadata.Keyword{"People"}, v[i])
	}
	for face, seen := range faces {
		if seen {
			facelist = append(facelist, face)
		}
	}
	if err := setFilteredKeywords(h, kws, personPredicate); err != nil {
		return err
	}
	if err := SetFaces(h, facelist); err != nil {
		return err
	}
	return nil
}

// personPredicate is the predicate satisfied by keyword tags that encode person
// names.
func personPredicate(kw metadata.Keyword) bool {
	return len(kw) == 2 && kw[0] == "People"
}
