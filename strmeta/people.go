package strmeta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp"
)

// A Person represents a person depicted in a media file.
type Person struct {
	// Name is the person's full name, as I would normally address them.
	Name string
	// FaceRegion indicates whether the media file has an area marked as
	// being this person's face.
	FaceRegion bool
}

// Parse parses a person.
func (g *Person) Parse(s string) error {
	g.FaceRegion = false
	if strings.HasSuffix(s, " [F]") {
		// We shouldn't see this on input, but since we want Parse and
		// String to be a faithful round trip, we should handle it.
		s = s[:len(s)-4]
		g.FaceRegion = true
	}
	g.Name = strings.TrimSpace(s)
	return nil
}

// String returns the formatted string form of the person name, suitable for
// input to Parse.
func (g Person) String() string {
	if g.Name != "" && g.FaceRegion {
		return g.Name + " [F]"
	}
	return g.Name
}

// Empty returns whether the person name is empty.
func (g Person) Empty() bool { return g.Name == "" }

// Equal returns whether two person names are equal.
func (g Person) Equal(other Person) bool {
	return g.Name == other.Name
	// Disregard the FaceRegion flag.
}

// GetPeople returns the highest priority person values.
func GetPeople(h filefmt.FileHandler) []Person {
	kws := getFilteredKeywords(h, personPredicate, false)
	people := make([]Person, 0, len(kws))
	pmap := make(map[string]*Person)
	for i := range kws {
		name := kws[i][1]
		if name != "" && pmap[name] == nil {
			people = append(people, Person{Name: name})
			pmap[name] = &people[len(people)-1]
		}
	}
	if xmp := h.XMP(false); xmp != nil {
		var faces []string
		if len(xmp.MPRegPersonDisplayNames()) != 0 {
			faces = xmp.MPRegPersonDisplayNames()
		} else {
			faces = xmp.MWGRSNames()
		}
		for _, face := range faces {
			if pmap[face] != nil {
				pmap[face].FaceRegion = true
			} else {
				people = append(people, Person{Name: face, FaceRegion: true})
				pmap[face] = &people[len(people)-1]
			}
		}
	}
	return people
}

// GetPersonTags returns all of the person tags and their values.
func GetPersonTags(h filefmt.FileHandler) (tags []string, values []Person) {
	tags, kws := getFilteredKeywordTags(h, personPredicate)
	values = make([]Person, len(kws))
	for i := range kws {
		values[i].Name = kws[i][1]
	}
	if xmp := h.XMP(false); xmp != nil {
		for _, face := range xmp.MPRegPersonDisplayNames() {
			tags = append(tags, "XMP  MP:Regions")
			values = append(values, Person{Name: face})
			// Note, we do not set FaceRegion to true, because that
			// would result in printing a [F] suffix, which is
			// rendundant given that the tag name is being shown.
		}
		for _, face := range xmp.MWGRSNames() {
			tags = append(tags, "XMP  mwg-rs:RegionInfo")
			values = append(values, Person{Name: face})
		}
	}
	return tags, values
}

// CheckPeople determines whether the people are tagged correctly, and are
// consistent with the reference.
func CheckPeople(ref, h filefmt.FileHandler) (res CheckResult) {
	var (
		xmp  *xmp.XMP
		pmap map[string]bool
	)
	if res = checkFilteredKeywords(ref, h, personPredicate); res == ChkConflictingValues {
		return res
	}
	xmp = h.XMP(false)
	if xmp == nil || (len(xmp.MPRegPersonDisplayNames()) == 0 && len(xmp.MWGRSNames()) == 0) {
		return res
	}
	if len(xmp.MPRegPersonDisplayNames()) != 0 && len(xmp.MWGRSNames()) != 0 && len(xmp.MPRegPersonDisplayNames()) != len(xmp.MWGRSNames()) {
		return ChkConflictingValues
	}
	pmap = make(map[string]bool)
	for _, p := range getFilteredKeywords(ref, personPredicate, false) {
		pmap[p[1]] = true
	}
	for _, f := range xmp.MPRegPersonDisplayNames() {
		if !pmap[f] {
			return ChkConflictingValues
		}
	}
	for _, f := range xmp.MWGRSNames() {
		if !pmap[f] {
			return ChkConflictingValues
		}
	}
	if len(xmp.MPRegPersonDisplayNames()) != len(xmp.MWGRSNames()) { // i.e., one of them is zero and the other isn't
		return ChkIncorrectlyTagged
	}
	return res // the result from checking the keywords
}

// SetPeople sets the person tags.
func SetPeople(h filefmt.FileHandler, v []Person) error {
	var (
		kws   = make([]metadata.Keyword, len(v))
		names = make(map[string]bool)
		faces = make(map[string]bool)
	)
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.MPRegPersonDisplayNames()) != 0 {
			for _, f := range xmp.MPRegPersonDisplayNames() {
				faces[f] = false
			}
		} else {
			for _, f := range xmp.MWGRSNames() {
				faces[f] = false
			}
		}
	}
	for i, g := range v {
		if g.Empty() {
			return errors.New("empty person not allowed")
		}
		if names[g.Name] {
			return fmt.Errorf("cannot list name %q twice", g.Name)
		}
		names[g.Name] = true
		if _, ok := faces[g.Name]; ok {
			faces[g.Name] = true
		} else if g.FaceRegion {
			return fmt.Errorf("cannot add face region for %q", g.Name)
		}
		kws[i] = append(metadata.Keyword{"People"}, v[i].Name)
	}
	for face, seen := range faces {
		if !seen {
			return fmt.Errorf("cannot remove %q: face regions are not removable with this tool", face)
		}
	}
	if err := setFilteredKeywords(h, kws, personPredicate); err != nil {
		return err
	}
	if xmp := h.XMP(false); xmp != nil {
		if list := xmp.MPRegPersonDisplayNames(); len(list) != 0 {
			j := 0
			for _, name := range list {
				if faces[name] {
					list[j] = name
					j++
				}
			}
			if j < len(list) {
				if err := xmp.SetMPRegPersonDisplayNames(list[:j]); err != nil {
					return err
				}
			}
		}
		if list := xmp.MWGRSNames(); len(list) != 0 {
			j := 0
			for _, name := range list {
				if faces[name] {
					list[j] = name
					j++
				}
			}
			if j < len(list) {
				if err := xmp.SetMWGRSNames(list[:j]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// personPredicate is the predicate satisfied by keyword tags that encode person
// names.
func personPredicate(kw metadata.Keyword) bool {
	return len(kw) == 2 && kw[0] == "People"
}
