package main

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata"
	strmeta "github.com/rothskeller/photo-tools/strmeta"
)

// A field represents a metadata field that can be read and/or changed by this
// program.  Fields are more abstract than tags: a field can represent a summary
// or merging of multiple tags, and multiple fields can be derived from the same
// tag.
type field struct {
	// name is the singular form of the name of the field as used on the
	// command line.
	name string
	// pluralName is the plural form of the name of the field as used on the
	// command line.  However, it is singular for single-valued fields.
	pluralName string
	// label is the label to identify the field in output tables.
	label string
	// multivalued is true if the field allows multiple values.
	multivalued bool
	// parseValue parses a string and returns a value for the field.  It
	// returns an error if the string is invalid.
	parseValue func(string) (interface{}, error)
	// renderValue takes a value for the field and renders it in string form
	// for display.
	renderValue func(interface{}) string
	// equalValue compares two values for equality.
	equalValue func(a, b interface{}) bool
	// getValues returns all of the values of the field.  (For single-valued
	// fields, the return slice will have at most one entry.)  Empty values
	// should not be included.
	getValues func(fileHandler) []interface{}
	// getTags returns the names of all of the metadata tags that correspond
	// to the field in its first return slice, and a parallel slice of the
	// values of those tags (which may be zero values).
	getTags func(fileHandler) ([]string, []interface{})
	// setValues sets all of the values of the field.
	setValues func(fileHandler, []interface{}) error
	// lang is the language tag for the field, if any.
	lang string
}

var artistField = &field{
	name: "artist", pluralName: "artist", label: "Artist",
	parseValue:  func(s string) (interface{}, error) { return s, nil },
	renderValue: func(v interface{}) string { return v.(string) },
	getValues: func(h fileHandler) []interface{} {
		if artist := strmeta.GetCreator(h); artist != "" {
			return []interface{}{artist}
		}
		return nil
	},
	getTags: func(h fileHandler) ([]string, []interface{}) {
		if tags, values := strmeta.GetCreatorTags(h); len(tags) != 0 {
			return tags, stringSliceToInterfaceSlice(values)
		}
		return nil, nil
	},
	setValues: func(h fileHandler, v []interface{}) error {
		switch len(v) {
		case 0:
			return strmeta.SetCreator(h, "")
		case 1:
			return strmeta.SetCreator(h, v[0].(string))
		default:
			return errors.New("artist cannot have multiple values")
		}
	},
}

var bothLocationsField = &field{ // only used internally by "copy all"
	name: "location", pluralName: "locations", label: "Locations",
	renderValue: func(v interface{}) string { return v.(*metadata.Location).String() },
	getValues: func(h fileHandler) []interface{} {
		if locations := strmeta.GetLocation(h); len(locations) != 0 {
			var ivals = make([]interface{}, len(locations))
			for i := range locations {
				ivals[i] = &locations[i]
			}
			return ivals
		}
		return nil
	},
	setValues: func(h fileHandler, v []interface{}) error {
		if len(v) != 0 {
			var locations = make([]metadata.Location, len(v))
			for i := range v {
				locations[i] = *v[i].(*metadata.Location)
			}
			return strmeta.SetLocation(h, locations)
		}
		if err := strmeta.SetLocation(h, nil); err != nil {
			return err
		}
		return setPlacesKeywords(h)
	},
}

var bothShownField = &field{ // only used internally by "copy all"
	name: "shown", pluralName: "shown", label: "Shown",
	renderValue: func(v interface{}) string { return v.(*metadata.Location).String() },
	getValues: func(h fileHandler) []interface{} {
		if shown := strmeta.GetShown(h); len(shown) != 0 {
			var ivals = make([]interface{}, len(shown))
			for i := range shown {
				ivals[i] = &shown[i]
			}
			return ivals
		}
		return nil
	},
	setValues: func(h fileHandler, v []interface{}) error {
		if len(v) != 0 {
			var shown = make([]metadata.Location, len(v))
			for i := range v {
				shown[i] = *v[i].(*metadata.Location)
			}
			return strmeta.SetShown(h, shown)
		}
		if err := strmeta.SetLocation(h, nil); err != nil {
			return err
		}
		return setPlacesKeywords(h)
	},
}

var captionField = &field{
	name: "caption", pluralName: "caption", label: "Caption",
	parseValue:  func(s string) (interface{}, error) { return s, nil },
	renderValue: func(v interface{}) string { return v.(string) },
	getValues: func(h fileHandler) []interface{} {
		if caption := strmeta.GetCaption(h); caption != "" {
			return []interface{}{caption}
		}
		return nil
	},
	getTags: func(h fileHandler) ([]string, []interface{}) {
		if tags, values := strmeta.GetCaptionTags(h); len(tags) != 0 {
			return tags, stringSliceToInterfaceSlice(values)
		}
		return nil, nil
	},
	setValues: func(h fileHandler, v []interface{}) error {
		switch len(v) {
		case 0:
			return strmeta.SetCaption(h, "")
		case 1:
			return strmeta.SetCaption(h, v[0].(string))
		default:
			return errors.New("caption cannot have multiple values")
		}
	},
}

var dateTimeField = &field{
	name: "datetime", pluralName: "datetime", label: "Date/Time",
	parseValue: func(s string) (interface{}, error) {
		var dt metadata.DateTime
		if err := dt.Parse(s); err != nil {
			return nil, err
		}
		return &dt, nil
	},
	renderValue: func(v interface{}) string { return v.(*metadata.DateTime).String() },
	getValues: func(h fileHandler) []interface{} {
		if datetime := strmeta.GetDateTime(h); !datetime.Empty() {
			return []interface{}{&datetime}
		}
		return nil
	},
	getTags: func(h fileHandler) ([]string, []interface{}) {
		if tags, values := strmeta.GetDateTimeTags(h); len(tags) != 0 {
			var ivals = make([]interface{}, len(values))
			for i := range values {
				ivals[i] = &values[i]
			}
			return tags, ivals
		}
		return nil, nil
	},
	setValues: func(h fileHandler, v []interface{}) error {
		switch len(v) {
		case 0:
			return strmeta.SetDateTime(h, metadata.DateTime{})
		case 1:
			return strmeta.SetDateTime(h, *v[0].(*metadata.DateTime))
		default:
			return errors.New("datetime cannot have multiple values")
		}
	},
}

var gpsField = &field{
	name: "gps", pluralName: "gps", label: "GPS Coords",
	parseValue: func(s string) (interface{}, error) {
		var gps metadata.GPSCoords
		if err := gps.Parse(s); err != nil {
			return nil, err
		}
		return &gps, nil
	},
	renderValue: func(v interface{}) string { return v.(*metadata.GPSCoords).String() },
	getValues: func(h fileHandler) []interface{} {
		if gps := strmeta.GetGPSCoords(h); !gps.Empty() {
			return []interface{}{&gps}
		}
		return nil
	},
	getTags: func(h fileHandler) ([]string, []interface{}) {
		if tags, values := strmeta.GetGPSCoordsTags(h); len(tags) != 0 {
			var ivals = make([]interface{}, len(values))
			for i := range values {
				ivals[i] = &values[i]
			}
			return tags, ivals
		}
		return nil, nil
	},
	setValues: func(h fileHandler, v []interface{}) error {
		switch len(v) {
		case 0:
			return strmeta.SetGPSCoords(h, metadata.GPSCoords{})
		case 1:
			return strmeta.SetGPSCoords(h, *v[0].(*metadata.GPSCoords))
		default:
			return errors.New("gps cannot have multiple values")
		}
	},
}

var groupsField = &field{
	name: "group", pluralName: "groups", label: "Group", multivalued: true,
	parseValue:  func(s string) (interface{}, error) { return metadata.ParseKeyword(s, "GROUPS") },
	renderValue: func(v interface{}) string { return v.(metadata.Keyword).StringWithoutPrefix("GROUPS") },
	equalValue:  func(a, b interface{}) bool { return a.(metadata.Keyword).Equal(b.(metadata.Keyword)) },
	getValues:   func(h fileHandler) []interface{} { return getKeywordValues(h, "GROUPS") },
	getTags:     func(h fileHandler) ([]string, []interface{}) { return getKeywordTags(h, "GROUPS") },
	setValues:   func(h fileHandler, v []interface{}) error { return setKeywordValues(h, "GROUPS", v) },
}

var keywordsField = &field{
	name: "keyword", pluralName: "keywords", label: "Keyword", multivalued: true,
	parseValue:  func(s string) (interface{}, error) { return metadata.ParseKeyword(s, "") },
	renderValue: func(v interface{}) string { return v.(metadata.Keyword).String() },
	equalValue:  func(a, b interface{}) bool { return a.(metadata.Keyword).Equal(b.(metadata.Keyword)) },
	getValues:   func(h fileHandler) []interface{} { return getKeywordValues(h, "") },
	getTags:     func(h fileHandler) ([]string, []interface{}) { return getKeywordTags(h, "") },
	setValues:   func(h fileHandler, v []interface{}) error { return setKeywordValues(h, "", v) },
}

var otherKeywordsField = &field{ // only used as part of "show all" and "tags all"
	name: "keyword", pluralName: "keywords", label: "Keyword", multivalued: true,
	renderValue: func(v interface{}) string { return v.(metadata.Keyword).String() },
	getValues: func(h fileHandler) []interface{} {
		if kws := strmeta.GetKeywords(h, ""); len(kws) != 0 {
			var ivals []interface{}
			for _, kw := range kws {
				if kw[0].Word != "GROUPS" && kw[0].Word != "PEOPLE" && kw[0].Word != "PLACES" && kw[0].Word != "TOPICS" {
					ivals = append(ivals, kw)
				}
			}
			if len(ivals) != 0 {
				return ivals
			}
		}
		return nil
	},
	getTags: func(h fileHandler) ([]string, []interface{}) {
		if tags, values := strmeta.GetKeywordsTags(h, ""); len(tags) != 0 {
			var itags []string
			var ivals []interface{}
			for i, kw := range values {
				if kw[0].Word != "GROUPS" && kw[0].Word != "PEOPLE" && kw[0].Word != "PLACES" && kw[0].Word != "TOPICS" {
					itags = append(itags, tags[i])
					ivals = append(ivals, kw)
				}
			}
			if len(ivals) != 0 {
				return itags, ivals
			}
		}
		return nil, nil
	},
}

var peopleField = &field{
	name: "person", pluralName: "people", label: "Person",
	parseValue:  func(s string) (interface{}, error) { return metadata.ParseKeyword(s, "PEOPLE") },
	renderValue: func(v interface{}) string { return v.(metadata.Keyword).StringWithoutPrefix("PEOPLE") },
	equalValue:  func(a, b interface{}) bool { return a.(metadata.Keyword).Equal(b.(metadata.Keyword)) },
	getValues:   func(h fileHandler) []interface{} { return getKeywordValues(h, "PEOPLE") },
	getTags:     func(h fileHandler) ([]string, []interface{}) { return getKeywordTags(h, "PEOPLE") },
	setValues:   func(h fileHandler, v []interface{}) error { return setKeywordValues(h, "PEOPLE", v) },
}

var placesField = &field{ // only used as part of "show all" and "tags all"
	name: "place", pluralName: "places", label: "Place", multivalued: true,
	renderValue: func(v interface{}) string { return v.(metadata.Keyword).StringWithoutPrefix("PLACES") },
	equalValue:  func(a, b interface{}) bool { return a.(metadata.Keyword).Equal(b.(metadata.Keyword)) },
	getValues:   func(h fileHandler) []interface{} { return getKeywordValues(h, "PLACES") },
	getTags:     func(h fileHandler) ([]string, []interface{}) { return getKeywordTags(h, "PLACES") },
}

var titleField = &field{
	name: "title", pluralName: "title", label: "Title",
	parseValue:  func(s string) (interface{}, error) { return s, nil },
	renderValue: func(v interface{}) string { return v.(string) },
	getValues: func(h fileHandler) []interface{} {
		if title := strmeta.GetTitle(h); title != "" {
			return []interface{}{title}
		}
		return nil
	},
	getTags: func(h fileHandler) ([]string, []interface{}) {
		if tags, values := strmeta.GetTitleTags(h); len(tags) != 0 {
			return tags, stringSliceToInterfaceSlice(values)
		}
		return nil, nil
	},
	setValues: func(h fileHandler, v []interface{}) error {
		switch len(v) {
		case 0:
			return strmeta.SetTitle(h, "")
		case 1:
			return strmeta.SetTitle(h, v[0].(string))
		default:
			return errors.New("title cannot have multiple values")
		}
	},
}

var topicsField = &field{
	name: "topic", pluralName: "topics", label: "Topic", multivalued: true,
	parseValue:  func(s string) (interface{}, error) { return metadata.ParseKeyword(s, "TOPICS") },
	renderValue: func(v interface{}) string { return v.(metadata.Keyword).StringWithoutPrefix("TOPICS") },
	equalValue:  func(a, b interface{}) bool { return a.(metadata.Keyword).Equal(b.(metadata.Keyword)) },
	getValues:   func(h fileHandler) []interface{} { return getKeywordValues(h, "TOPICS") },
	getTags:     func(h fileHandler) ([]string, []interface{}) { return getKeywordTags(h, "TOPICS") },
	setValues:   func(h fileHandler, v []interface{}) error { return setKeywordValues(h, "TOPICS", v) },
}

var untaggedLocationField = &field{
	name: "location", pluralName: "location", label: "Location",
	parseValue: func(s string) (interface{}, error) {
		var loc metadata.Location
		if err := loc.Parse(s); err != nil {
			return nil, err
		}
		return &loc, nil
	},
	renderValue: func(v interface{}) string { return v.(*metadata.Location).String() },
}

var untaggedShownField = &field{
	name: "shown", pluralName: "shown", label: "Shown",
	parseValue: func(s string) (interface{}, error) {
		var loc metadata.Location
		if err := loc.Parse(s); err != nil {
			return nil, err
		}
		return &loc, nil
	},
	renderValue: func(v interface{}) string { return v.(*metadata.Location).String() },
}

func stringSliceToInterfaceSlice(ss []string) (is []interface{}) {
	is = make([]interface{}, len(ss))
	for i, s := range ss {
		is[i] = s
	}
	return is
}

func getKeywordValues(h fileHandler, prefix string) []interface{} {
	if kws := strmeta.GetKeywords(h, prefix); len(kws) != 0 {
		var ivals = make([]interface{}, len(kws))
		for i := range kws {
			ivals[i] = kws[i]
		}
		return ivals
	}
	return nil
}

func getKeywordTags(h fileHandler, prefix string) ([]string, []interface{}) {
	if tags, values := strmeta.GetKeywordsTags(h, prefix); len(tags) != 0 {
		var ivals = make([]interface{}, len(values))
		for i := range values {
			ivals[i] = &values[i]
		}
		return tags, ivals
	}
	return nil, nil
}

func setKeywordValues(h fileHandler, prefix string, v []interface{}) error {
	if len(v) != 0 {
		var kws = make([]metadata.Keyword, len(v))
		for i := range v {
			kws[i] = v[i].(metadata.Keyword)
		}
		return strmeta.SetKeywords(h, prefix, kws)
	}
	return strmeta.SetKeywords(h, prefix, nil)
}

func setPlacesKeywords(h fileHandler) error {
	var kws []metadata.Keyword
LOCATIONS:
	for _, loc := range strmeta.GetLocation(h) {
		var kw metadata.Keyword
		if loc.CountryName != "" {
			kw = append(kw, metadata.KeywordComponent{Word: loc.CountryName})
		}
		if loc.State != "" {
			kw = append(kw, metadata.KeywordComponent{Word: loc.State})
		}
		if loc.City != "" {
			kw = append(kw, metadata.KeywordComponent{Word: loc.City})
		}
		if loc.Sublocation != "" {
			kw = append(kw, metadata.KeywordComponent{Word: loc.Sublocation})
		}
		if len(kw) == 0 {
			continue
		}
		for _, ex := range kws {
			if ex.Equal(kw) {
				continue LOCATIONS
			}
		}
		kws = append(kws, kw)
	}
SHOWN:
	for _, sh := range strmeta.GetShown(h) {
		var kw metadata.Keyword
		if sh.CountryName != "" {
			kw = append(kw, metadata.KeywordComponent{Word: sh.CountryName})
		}
		if sh.State != "" {
			kw = append(kw, metadata.KeywordComponent{Word: sh.State})
		}
		if sh.City != "" {
			kw = append(kw, metadata.KeywordComponent{Word: sh.City})
		}
		if sh.Sublocation != "" {
			kw = append(kw, metadata.KeywordComponent{Word: sh.Sublocation})
		}
		if len(kw) == 0 {
			continue
		}
		for _, ex := range kws {
			if ex.Equal(kw) {
				continue SHOWN
			}
		}
		kws = append(kws, kw)
	}
	return strmeta.SetKeywords(h, "PLACES", kws)
}
