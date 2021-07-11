// Package metadata contains the definition of a MediaFile and the methods for
// working with one.
package metadata

import "errors"

// MediaFile is the interface honored by all media file handlers; there are
// different concrete implementations for different file types.
type MediaFile interface {
	// Filename returns the media file name.
	Filename() string
	// Providers returns a list of metadata providers for the media file.
	// Which providers are returned, and in what order, depends on the media
	// file type and on which metadata sources are present in the file.
	Providers() []Provider
	// Dirty returns a flag indicating whether any of the metadata for the
	// media file have been changed, and therefore the file needs to be
	// saved.
	Dirty() bool
	// Save saves the updated metadata to the file.
	Save() error
}

// A Provider is a source for metadata information and/or a receiver for
// metadata changes.  The set of metadata providers for a file, and their order,
// depends on the file type and the presence of particular metadata blocks in
// that file.
type Provider interface {
	// ProviderName is the name for the provider, for debug purposes.
	ProviderName() string

	// Caption returns the value of the Caption field.
	Caption() (value string)
	// CaptionTags returns a list of tag names for the Caption field, and a
	// parallel list of values held by those tags.
	CaptionTags() (tags []string, values [][]string)
	// SetCaption sets the value of the Caption field.
	SetCaption(value string) error

	// Creator returns the value of the Creator field.
	Creator() (value string)
	// CreatorTags returns a list of tag names for the Creator field, and a
	// parallel list of values held by those tags.
	CreatorTags() (tags []string, values [][]string)
	// SetCreator sets the value of the Creator field.
	SetCreator(value string) error

	// DateTime returns the value of the DateTime field.
	DateTime() (value DateTime)
	// DateTimeTags returns a list of tag names for the DateTime field, and
	// a parallel list of values held by those tags.
	DateTimeTags() (tags []string, values []DateTime)
	// SetDateTime sets the value of the DateTime field.
	SetDateTime(value DateTime) error

	// Faces returns the values of the Faces field.
	Faces() (values []string)
	// FacesTags returns a list of tag names for the Faces field, and a
	// parallel list of values held by those tags.
	FacesTags() (tags []string, values [][]string)
	// SetFaces sets the values of the Faces field.
	SetFaces(values []string) error

	// GPS returns the values of the GPS field.
	GPS() (value GPSCoords)
	// GPSTags returns a list of tag names for the GPS field, and a parallel
	// list of values held by those tags.
	GPSTags() (tags []string, values []GPSCoords)
	// SetGPS sets the values of the GPS field.
	SetGPS(value GPSCoords) error

	// Groups returns the values of the Groups field.
	Groups() (values []HierValue)
	// GroupsTags returns a list of tag names for the Groups field, and a
	// parallel list of values held by those tags.
	GroupsTags() (tags []string, values [][]HierValue)
	// SetGroups sets the values of the Groups field.
	SetGroups(values []HierValue) error

	// Keywords returns the values of the Keywords field.
	Keywords() (values []HierValue)
	// KeywordsTags returns a list of tag names for the Keywords field, and
	// a parallel list of values held by those tags.
	KeywordsTags() (tags []string, values [][]HierValue)
	// SetKeywords sets the values of the Keywords field.
	SetKeywords(values []HierValue) error

	// Location returns the value of the Location field.
	Location() (value Location)
	// LocationTags returns a list of tag names for the Location field, and
	// a parallel list of values held by those tags.
	LocationTags() (tags []string, values [][]Location)
	// SetLocation sets the value of the Location field.
	SetLocation(values Location) error

	// People returns the values of the People field.
	People() (values []string)
	// PeopleTags returns a list of tag names for the People field, and a
	// parallel list of values held by those tags.
	PeopleTags() (tags []string, values [][]string)
	// SetPeople sets the values of the People field.
	SetPeople(values []string) error

	// Places returns the values of the Places field.
	Places() (values []HierValue)
	// PlacesTags returns a list of tag names for the Places field, and a
	// parallel list of values held by those tags.
	PlacesTags() (tags []string, values [][]HierValue)
	// SetPlaces sets the values of the Places field.
	SetPlaces(values []HierValue) error

	// Title returns the value of the Title field.
	Title() (value string)
	// TitleTags returns a list of tag names for the Title field, and a
	// parallel list of values held by those tags.
	TitleTags() (tags []string, values [][]string)
	// SetTitle sets the values of the Title field.
	SetTitle(value string) error

	// Topics returns the values of the Topics field.
	Topics() (values []HierValue)
	// TopicsTags returns a list of tag names for the Topics field, and a
	// parallel list of values held by those tags.
	TopicsTags() (tags []string, values [][]HierValue)
	// SetTopics sets the values of the Topics field.
	SetTopics(values []HierValue) error
}

// ErrNotSupported is the error returned from a SetXXX function when the
// underlying metadata does not support the XXX field.
var ErrNotSupported = errors.New("unsupported metadata field")

// BaseProvider provides a base implementation for all providers.  The base
// implementation returns empty values for Get, empty lists for Tags, and
// ErrNotSupported for Set, for all metadata fields that aren't overridden by
// the concrete implementation.
type BaseProvider struct{}

// Caption returns the value of the Caption field.
func (p BaseProvider) Caption() string { return "" }

// CaptionTags returns a list of tag names for the Caption field, and a
// parallel list of values held by those tags.
func (p BaseProvider) CaptionTags() ([]string, [][]string) { return nil, nil }

// SetCaption sets the value of the Caption field.
func (p BaseProvider) SetCaption(value string) error { return ErrNotSupported }

// Creator returns the value of the Creator field.
func (p BaseProvider) Creator() string { return "" }

// CreatorTags returns a list of tag names for the Creator field, and a
// parallel list of values held by those tags.
func (p BaseProvider) CreatorTags() ([]string, [][]string) { return nil, nil }

// SetCreator sets the value of the Creator field.
func (p BaseProvider) SetCreator(value string) error { return ErrNotSupported }

// DateTime returns the value of the DateTime field.
func (p BaseProvider) DateTime() DateTime { return DateTime{} }

// DateTimeTags returns a list of tag names for the DateTime field, and
// a parallel list of values held by those tags.
func (p BaseProvider) DateTimeTags() ([]string, []DateTime) { return nil, nil }

// SetDateTime sets the value of the DateTime field.
func (p BaseProvider) SetDateTime(value DateTime) error { return ErrNotSupported }

// Faces returns the values of the Faces field.
func (p BaseProvider) Faces() []string { return nil }

// FacesTags returns a list of tag names for the Faces field, and a
// parallel list of values held by those tags.
func (p BaseProvider) FacesTags() ([]string, [][]string) { return nil, nil }

// SetFaces sets the values of the Faces field.
func (p BaseProvider) SetFaces(values []string) error { return ErrNotSupported }

// GPS returns the values of the GPS field.
func (p BaseProvider) GPS() GPSCoords { return GPSCoords{} }

// GPSTags returns a list of tag names for the GPS field, and a parallel
// list of values held by those tags.
func (p BaseProvider) GPSTags() ([]string, []GPSCoords) { return nil, nil }

// SetGPS sets the values of the GPS field.
func (p BaseProvider) SetGPS(value GPSCoords) error { return ErrNotSupported }

// Groups returns the values of the Groups field.
func (p BaseProvider) Groups() []HierValue { return nil }

// GroupsTags returns a list of tag names for the Groups field, and a
// parallel list of values held by those tags.
func (p BaseProvider) GroupsTags() ([]string, [][]HierValue) { return nil, nil }

// SetGroups sets the values of the Groups field.
func (p BaseProvider) SetGroups(values []HierValue) error { return ErrNotSupported }

// Keywords returns the values of the Keywords field.
func (p BaseProvider) Keywords() []HierValue { return nil }

// KeywordsTags returns a list of tag names for the Keywords field, and
// a parallel list of values held by those tags.
func (p BaseProvider) KeywordsTags() ([]string, [][]HierValue) { return nil, nil }

// SetKeywords sets the values of the Keywords field.
func (p BaseProvider) SetKeywords(values []HierValue) error { return ErrNotSupported }

// Location returns the value of the Location field.
func (p BaseProvider) Location() Location { return Location{} }

// LocationTags returns a list of tag names for the Location field, and
// a parallel list of values held by those tags.
func (p BaseProvider) LocationTags() ([]string, [][]Location) { return nil, nil }

// SetLocation sets the value of the Location field.
func (p BaseProvider) SetLocation(values Location) error { return ErrNotSupported }

// People returns the values of the People field.
func (p BaseProvider) People() []string { return nil }

// PeopleTags returns a list of tag names for the People field, and a
// parallel list of values held by those tags.
func (p BaseProvider) PeopleTags() ([]string, [][]string) { return nil, nil }

// SetPeople sets the values of the People field.
func (p BaseProvider) SetPeople(values []string) error { return ErrNotSupported }

// Places returns the values of the Places field.
func (p BaseProvider) Places() []HierValue { return nil }

// PlacesTags returns a list of tag names for the Places field, and a
// parallel list of values held by those tags.
func (p BaseProvider) PlacesTags() ([]string, [][]HierValue) { return nil, nil }

// SetPlaces sets the values of the Places field.
func (p BaseProvider) SetPlaces(values []HierValue) error { return ErrNotSupported }

// Title returns the value of the Title field.
func (p BaseProvider) Title() string { return "" }

// TitleTags returns a list of tag names for the Title field, and a
// parallel list of values held by those tags.
func (p BaseProvider) TitleTags() ([]string, [][]string) { return nil, nil }

// SetTitle sets the values of the Title field.
func (p BaseProvider) SetTitle(value string) error { return ErrNotSupported }

// Topics returns the values of the Topics field.
func (p BaseProvider) Topics() []HierValue { return nil }

// TopicsTags returns a list of tag names for the Topics field, and a
// parallel list of values held by those tags.
func (p BaseProvider) TopicsTags() ([]string, [][]HierValue) { return nil, nil }

// SetTopics sets the values of the Topics field.
func (p BaseProvider) SetTopics(values []HierValue) error { return ErrNotSupported }
