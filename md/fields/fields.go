// Package fields defines the fields that are used in the md tool.
package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

// Field is the interface honored by a field definition.
type Field interface {
	// Name returns the singular form of the name of the field as used on
	// the command line.
	Name() string
	// PluralName returns the plural form of the name of the field as used
	// on the command line.  However, it is singular for single-valued
	// fields.
	PluralName() string
	// Label returns the label to identify the field in output tables.
	Label() string
	// ShortLabel returns the two-character label that identifies the field
	// in tables produced by the check operation.
	ShortLabel() string
	// Multivalued returns true if the field allows multiple values.
	Multivalued() bool
	// ParseValue parses a string and returns a value for the field.  It
	// returns an error if the string is invalid.
	ParseValue(s string) (interface{}, error)
	// RenderValue takes a value for the field and renders it in string form
	// for display.
	RenderValue(v interface{}) string
	// EqualValue compares two values for equality.
	EqualValue(a, b interface{}) bool
	// GetValues returns all of the values of the field.  (For single-valued
	// fields, the return slice will have at most one entry.)  Empty values
	// should not be included.
	GetValues(h filefmt.FileHandler) []interface{}
	// GetTags returns the names of all of the metadata tags that correspond
	// to the field in its first return slice, and a parallel slice of the
	// values of those tags (which may be zero values).
	GetTags(h filefmt.FileHandler) ([]string, []interface{})
	// SetValues sets all of the values of the field.
	SetValues(h filefmt.FileHandler, v []interface{}) error
	// CheckValues returns whether the values of the field in the target are
	// tagged correctly, and are consistent with the values of the field in
	// the reference.
	CheckValues(ref, tgt filefmt.FileHandler) strmeta.CheckResult
}

// baseField is, in effect, an abstract base class for fields, providing methods
// that all fields use.
type baseField struct {
	name        string
	pluralName  string
	label       string
	shortLabel  string
	multivalued bool
}

// Name returns the singular form of the name of the field as used on
// the command line.
func (f *baseField) Name() string { return f.name }

// PluralName returns the plural form of the name of the field as used
// on the command line.  However, it is singular for single-valued
// fields.
func (f *baseField) PluralName() string { return f.pluralName }

// Label returns the label to identify the field in output tables.
func (f *baseField) Label() string { return f.label }

// ShortLabel returns the two-character label that identifies the field
// in tables produced by the check operation.
func (f *baseField) ShortLabel() string { return f.shortLabel }

// Multivalued returns true if the field allows multiple values.
func (f *baseField) Multivalued() bool { return f.multivalued }

// ParseField parses a string to see if it is a recognized field name or
// abbreviation.  If so, it returns the corresponding field handler.  Otherwise
// it returns nil.  This function does not handle "all".
func ParseField(arg string) Field {
	switch arg {
	case "artist", "a":
		return ArtistField
	case "caption", "c":
		return CaptionField
	case "datetime", "date", "time", "d":
		return DateTimeField
	case "gps", "g":
		return GPSField
	case "keyword", "keywords", "kw", "k":
		return KeywordsField
	case "location", "loc", "l":
		return LocationField
	case "title", "t":
		return TitleField
	case "group", "groups":
		return GroupsField
	case "person", "people":
		return PeopleField
	case "place", "places":
		return PlacesField
	case "topic", "topics":
		return TopicsField
	}
	return nil
}
