package strmeta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/iptc"
)

// A Keyword is a free-form hierarchical keyword.
type Keyword []string

// Parse parses a keyword, as a hierarchical string with levels separated by
// slashes and optional whitespace.  Pipe symbols are disallowed due to
// underlying storage formats, and empty levels are disallowed (although a
// completely empty string is allowed).  Keywords that would be interpreted as
// keywords, people, places, or topics are disallowed.
func (g *Keyword) Parse(s string) error {
	kw, err := metadata.ParseKeyword(s, "")
	if err == nil {
		if !keywordPredicate(kw) {
			return fmt.Errorf("invalid keyword (reserved for %s field)", strings.ToLower(kw[0]))
		}
		*g = Keyword(kw)
	}
	return err
}

// String returns the formatted string form of the keyword, suitable for
// input to Parse.
func (g Keyword) String() string { return metadata.Keyword(g).String() }

// Empty returns whether the keyword is empty.
func (g Keyword) Empty() bool { return len(g) == 0 }

// Equal returns whether two keywords are equal.
func (g Keyword) Equal(other Keyword) bool {
	return metadata.Keyword(g).Equal(metadata.Keyword(other))
}

// GetKeywords returns the highest priority keyword values.
func GetKeywords(h filefmt.FileHandler) []Keyword {
	kws := getFilteredKeywords(h, keywordPredicate, true)
	keywords := make([]Keyword, len(kws))
	for i := range kws {
		keywords[i] = Keyword(kws[i])
	}
	return keywords
}

// GetKeywordTags returns all of the keyword tags and their values.
func GetKeywordTags(h filefmt.FileHandler) (tags []string, values []Keyword) {
	tags, kws := getFilteredKeywordTags(h, keywordPredicate)
	values = make([]Keyword, len(kws))
	for i := range kws {
		values[i] = Keyword(kws[i])
	}
	return tags, values
}

// CheckKeywords determines whether the keywords are tagged correctly.
func CheckKeywords(h filefmt.FileHandler) (res CheckResult) {
	if res = checkFilteredKeywords(h, keywordPredicate); res == ChkConflictingValues {
		return res
	}
	// Check on the "keywords" field also checks the consistency of the flat
	// keyword tags.  For this purpose, we look at all keywords, not just
	// the ones matching keywordPredicate.
	var (
		values []metadata.Keyword
		refmap = make(map[string]bool)
	)
	values = getFilteredKeywords(h, allKeywordsFilter, true)
	for _, kw := range values {
		if len(kw) != 0 {
			refmap[kw[len(kw)-1]] = true
		}
	}
	if xmp := h.XMP(false); xmp != nil {
		var tgtmap = make(map[string]bool)
		for _, kw := range xmp.DCSubject() {
			tgtmap[kw] = true
		}
		if r := checkMaps(refmap, tgtmap); r == ChkConflictingValues {
			return r
		} else if r < res {
			res = r
		}
	}
	if i := h.IPTC(); i != nil {
		for kw := range refmap {
			if len(kw) > iptc.MaxKeywordLen {
				refmap[kw[:iptc.MaxKeywordLen]] = true
				delete(refmap, kw)
			}
		}
		var tgtmap = make(map[string]bool)
		for _, kw := range i.Keywords() {
			if len(kw) > iptc.MaxKeywordLen {
				res = ChkIncorrectlyTagged
				tgtmap[kw[:iptc.MaxKeywordLen]] = true
			} else {
				tgtmap[kw] = true
			}
		}
		if r := checkMaps(refmap, tgtmap); r == ChkConflictingValues {
			return r
		} else if r < res {
			res = r
		}
	}
	return res
}

// SetKeywords sets the keyword tags.
func SetKeywords(h filefmt.FileHandler, v []Keyword) error {
	var kws = make([]metadata.Keyword, len(v))
	for i, g := range v {
		if g.Empty() {
			return errors.New("empty keyword not allowed")
		}
		kws[i] = metadata.Keyword(v[i])
	}
	return setFilteredKeywords(h, kws, keywordPredicate)
}

// keywordPredicate is the predicate satisfied by keyword tags that encode keyword
// names.
func keywordPredicate(kw metadata.Keyword) bool {
	return !groupPredicate(kw) && !personPredicate(kw) && !placePredicate(kw) && !topicPredicate(kw)
}
