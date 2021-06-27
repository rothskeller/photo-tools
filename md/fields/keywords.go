package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

// KeywordsField is the field handler for all keywords regardless of prefix.  It
// is used when explicitly listed in any operation except "set", where it is
// prohibited for safety's sake.  It is also used when implied by "all" or an
// empty field list in any operations that accepts a field lists, except "show".
var KeywordsField Field = &baseKWField{
	baseField{
		name:        "keyword",
		pluralName:  "keywords",
		label:       "Keyword",
		shortLabel:  "KW",
		multivalued: true,
	},
	"",
}

// GroupsField is the field handler for the subset of keywords that start with
// GROUPS.  These represent groups of people (organizations, teams, etc.) that
// are depicted in the media.
var GroupsField Field = &baseKWField{
	baseField{
		name:        "group",
		pluralName:  "groups",
		label:       "Group",
		shortLabel:  "GR",
		multivalued: true,
	},
	"GROUPS",
}

// PeopleField is the field handler for the subset of keywords that start with
// GROUPS.  These represent groups of people (organizations, teams, etc.) that
// are depicted in the media.
var PeopleField Field = &baseKWField{
	baseField{
		name:        "person",
		pluralName:  "people",
		label:       "Person",
		shortLabel:  "PE",
		multivalued: true,
	},
	"PEOPLE",
}

// PlacesField is the field handler for the subset of keywords that start with
// GROUPS.  These represent groups of people (organizations, teams, etc.) that
// are depicted in the media.
var PlacesField Field = &baseKWField{
	baseField{
		name:        "place",
		pluralName:  "places",
		label:       "Place",
		shortLabel:  "PL",
		multivalued: true,
	},
	"PLACES",
}

// TopicsField is the field handler for the subset of keywords that start with
// GROUPS.  These represent groups of people (organizations, teams, etc.) that
// are depicted in the media.
var TopicsField Field = &baseKWField{
	baseField{
		name:        "topic",
		pluralName:  "topics",
		label:       "Topic",
		shortLabel:  "TO",
		multivalued: true,
	},
	"TOPICS",
}

type otherKeywordsField struct {
	baseKWField
}

// OtherKeywordsField is the field handler for a pseudo-field that is used only
// by the "show" operation when invoked with an "all" field list (or an empty
// field list implying "all").  The output of that show operation will include
// groups, people, places, and topics tags under their own headings, and this
// field will list all other keywords not covered by those other fields.
var OtherKeywordsField Field = &otherKeywordsField{
	baseKWField{
		baseField{
			name:        "keyword",
			pluralName:  "keywords",
			label:       "Keyword",
			shortLabel:  "KW",
			multivalued: true,
		},
		"",
	},
}

// GetValues returns all of the values of the field.  (For single-valued
// fields, the return slice will have at most one entry.)  Empty values
// should not be included.
func (f *otherKeywordsField) GetValues(h filefmt.FileHandler) []interface{} {
	var kws = strmeta.GetKeywords(h)
	j := 0
	for _, kw := range kws {
		if len(kw) != 0 && kw[0].Word != "GROUPS" && kw[0].Word != "PEOPLE" &&
			kw[0].Word != "PLACES" && kw[0].Word != "TOPICS" {
			kws[j] = kw
			j++
		}
	}
	kws = kws[:j]
	if len(kws) == 0 {
		return nil
	}
	var ifcs = make([]interface{}, len(kws))
	for i := range kws {
		ifcs[i] = kws[i]
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

// CheckValues should not be called for this pseudo-field.
func (f *otherKeywordsField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) strmeta.CheckResult {
	panic("should not be called")
}
