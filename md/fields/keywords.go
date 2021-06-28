package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/strmeta"
)

// KeywordsField is the field handler for all keywords regardless of prefix.  It
// is used when explicitly listed in any operation except "set", where it is
// prohibited for safety's sake.  It is also used when implied by "all" or an
// empty field list in any operations that accepts a field list, except "check"
// and "show".
var KeywordsField Field = &baseKWField{
	baseField: baseField{
		name:        "keyword",
		pluralName:  "keywords",
		label:       "Keyword",
		shortLabel:  "KW",
		multivalued: true,
	},
	prefix:   "",
	expected: true,
}

// GroupsField is the field handler for the subset of keywords that start with
// Groups.  These represent groups of people (organizations, teams, etc.) that
// are depicted in the media.
var GroupsField Field = &baseKWField{
	baseField: baseField{
		name:        "group",
		pluralName:  "groups",
		label:       "Group",
		shortLabel:  "GR",
		multivalued: true,
	},
	prefix: "Groups",
}

// PlacesField is the field handler for the subset of keywords that start with
// Places.  These represent locations: either where the media was captured or
// a location depicted in it.
var PlacesField Field = &baseKWField{
	baseField: baseField{
		name:        "place",
		pluralName:  "places",
		label:       "Place",
		shortLabel:  "PL",
		multivalued: true,
	},
	prefix:   "Places",
	expected: true,
}

// TopicsField is the field handler for the subset of keywords that start with
// Topics.  These represent topics (activities, events, etc.) depicted in the
// media.
var TopicsField Field = &baseKWField{
	baseField: baseField{
		name:        "topic",
		pluralName:  "topics",
		label:       "Topic",
		shortLabel:  "TP",
		multivalued: true,
	},
	prefix: "Topics",
}

type otherKeywordsField struct {
	baseKWField
}

// OtherKeywordsField is the field handler for a pseudo-field that is used only
// by the "check" and "show" operations when invoked with an "all" field list
// (or an empty field list implying "all").  These operations list groups,
// people, places, and topics separately, and their "Keyword" heading is only
// for the remaining keywords.
var OtherKeywordsField Field = &otherKeywordsField{
	baseKWField{
		baseField: baseField{
			name:        "keyword",
			pluralName:  "keywords",
			label:       "Keyword",
			shortLabel:  "KW",
			multivalued: true,
		},
		prefix: "",
	},
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *otherKeywordsField) GetValues(h filefmt.FileHandler) []interface{} {
	var kws = strmeta.GetKeywords(h)
	var filtered []metadata.Keyword
	for _, kw := range kws {
		if len(kw) != 0 {
			continue
		}
		if kw[0] == "Groups" || kw[0] == "Places" || kw[0] == "Topics" {
			continue
		}
		if kw[0] == "People" && len(kw) == 2 {
			// Because People keywords with more than one component
			// can't be represented as People values, we treat them
			// as "other keywords" instead.
			continue
		}
		filtered = append(filtered, kw)
	}
	if len(filtered) == 0 {
		return nil
	}
	var ifcs = make([]interface{}, len(filtered))
	for i := range filtered {
		ifcs[i] = filtered[i]
	}
	return ifcs
}

// GetTags should not be called for this pseudo-field.
func (f *otherKeywordsField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	panic("should not be called")
}

// SetValues should not be called for this pseudo-field.
func (f *otherKeywordsField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	panic("should not be called")
}

// CheckValues returns whether the values of the field in the target are
// tagged correctly, and are consistent with the values of the field in
// the reference.
func (f *otherKeywordsField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) strmeta.CheckResult {
	if result := strmeta.CheckKeywords(ref, tgt); result <= 0 {
		return result
	}
	return strmeta.CheckResult(len(f.GetValues(tgt)))
}
