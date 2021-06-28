package strmeta

import (
	"errors"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
)

// A Place represents a place related to a media (i.e., where it was captured,
// or depicted in it).
type Place []string

// Parse parses a place name, as a hierarchical string with levels separated by
// slashes and optional whitespace.  Pipe symbols are disallowed due to
// underlying storage formats, and empty levels are disallowed (although a
// completely empty string is allowed).
func (g *Place) Parse(s string) error {
	kw, err := metadata.ParseKeyword(s, "")
	if err == nil {
		*g = Place(kw)
	}
	return err
}

// String returns the formatted string form of the place name, suitable for
// input to Parse.
func (g Place) String() string { return metadata.Keyword(g).String() }

// Empty returns whether the place name is empty.
func (g Place) Empty() bool { return len(g) == 0 }

// Equal returns whether two place names are equal.
func (g Place) Equal(other Place) bool {
	return metadata.Keyword(g).Equal(metadata.Keyword(other))
}

// GetPlaces returns the highest priority place values.
func GetPlaces(h filefmt.FileHandler) []Place {
	kws := getFilteredKeywords(h, placePredicate, false)
	places := make([]Place, len(kws))
	for i := range kws {
		places[i] = Place(kws[i][1:])
	}
	return places
}

// GetPlaceTags returns all of the place tags and their values.
func GetPlaceTags(h filefmt.FileHandler) (tags []string, values []Place) {
	tags, kws := getFilteredKeywordTags(h, placePredicate)
	values = make([]Place, len(kws))
	for i := range kws {
		values[i] = Place(kws[i][1:])
	}
	return tags, values
}

// CheckPlaces determines whether the places are tagged correctly, and are
// consistent with the reference.
func CheckPlaces(ref, h filefmt.FileHandler) (res CheckResult) {
	res = checkFilteredKeywords(ref, h, placePredicate)
	if res == ChkOptionalAbsent {
		return ChkExpectedAbsent
	}
	return res
}

// SetPlaces sets the place tags.
func SetPlaces(h filefmt.FileHandler, v []Place) error {
	var kws = make([]metadata.Keyword, len(v))
	for i, g := range v {
		if g.Empty() {
			return errors.New("empty place name not allowed")
		}
		kws[i] = append(metadata.Keyword{"Places"}, v[i]...)
	}
	if err := setFilteredKeywords(h, kws, placePredicate); err != nil {
		return err
	}
	// SetPlaces will clear Location if it's not congruent with any places.
	var location = GetLocation(h)
	var parts []string
	if location.CountryName != "" {
		parts = append(parts, location.CountryName)
	}
	if location.State != "" {
		parts = append(parts, location.State)
	}
	if location.City != "" {
		parts = append(parts, location.City)
	}
	if location.Sublocation != "" {
		parts = append(parts, location.Sublocation)
	}
	if len(parts) == 0 {
		return nil
	}
	for _, p := range v {
		if locationCongruentToPlace(parts, p) {
			return nil
		}
	}
	return SetLocation(h, Location{})
}
func locationCongruentToPlace(loc []string, place Place) bool {
	for len(loc) != 0 && len(place) != 0 {
		if loc[0] == place[0] {
			loc = loc[1:]
		}
		place = place[1:]
	}
	return len(loc) == 0
}

// placePredicate is the predicate satisfied by keyword tags that encode place
// names.
func placePredicate(kw metadata.Keyword) bool {
	return len(kw) >= 2 && kw[0] == "Places"
}
