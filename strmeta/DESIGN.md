# STR Media Metadata

This package is the heart of the metadata management for my media library. It
translates between the plethora of metadata tags that exist in various media and
the simplified metadata model that I use in my library.

This file describes my metadata model and how it relates to the underlying
metadata. For clarity, "field" always refers to a datum in the simplified
model; "tag" always refers to the raw underlying metadata.

My metadata model consists of these fields, each of which is discussed below:

- Caption
- Creator
- DateTime
- GPSCoords
- Groups
- Keywords
- Location
- People
- Places
- Title
- Topics

## Library API

For each of the fields in my metadata, there are four functions supplied by this
library.

`GetXXX(h filefmt.FileHandler) valuetype` gets the value of the field from the
highest priority underlying tag that contains any data. For a multivalued
field, the return is `[]valuetype`.

`GetXXXTags(h filefmt.FileHandler) (tags []string, values []valuetype)` gets all
of the values of the field from all underlying tags, and returns the tag names
and corresponding values as two parallel slices. Tags that have multiple values
get multiple entries in each of the returned slices. Tags that have no value
get one entry in each of the returned slices, with the tag name and an empty
value. Exceptions: empty tags are not returned if the entire metadata block
containing them is absent from the media file, or if the tag is marked `[rm]` in
the field descriptions below.

`CheckXXX(ref, tgt filefmt.FileHandler) CheckResult` returns an indicator of
whether the field value(s) are tagged correctly in tgt, and whether its value(s)
are consistent with the value(s) of the tag in ref.

`SetXXX(h filefmt.FileHandler, value valuetype) error` sets the value of the
field, in all of the underlying tags that support it (except those marked `[rm]`
in the descriptions below, which are cleared instead). For a multivalued field,
the second parameter is `[]valuetype`, and the method replaces all existing
values with the values supplied. This method returns an error if the value(s)
are invalid or some other error occurs.

All of these methods take, as their first parameter, the handle to the file on
which they are supposed to act.

For each field, the `valuetype` is given in the field descriptions below. For
several types, it is a primitive Go `string`. In all other cases, it is a type
that supports this interface:

    type ValueType[T any] interface {
        Parse(s string) (err error)
        String() string
        Empty() bool
        Equal(other T) bool
    }

## Caption Field

The caption of a media is an optional block of text in prose form, describing
the contents of the media. Its intent is not to be redundant with other
metadata, but to provide back-story needed to appreciate the media. The value
type for this field is `string`.

The caption comes from the following tags, in priority order. (See the bottom of
this file for definitions of the notes.)

    XMP  dc:description        [alt]
    XMP  exif:UserComments     [mult] [rm]
    XMP  tiff:ImageDescription [alt]
    EXIF UserComment           [rm]
    EXIF ImageDescription
    IPTC CaptionAbstract       [max]

## Creator Field

The creator of a media is the person who originally captured or created it. It
is a `string` containing the person's full name in natural reading order,
optionally followed by a comma, a space, and a company name if that is relevant.
In some cases it may be a company name only. It should be considered a required
field in all media, unless the creator is not known.

The creator comes from the following tags, in priority order:

    XMP  dc:creator   [mult]
    XMP  tiff:Artist  [rm]
    IPTC By-line      [mult] [max]
    EXIF Artist       [mult]

## DateTime Field

The date/time of a media is the date and time at which it was originally
captured, as precisely as is known. It is stored as a `metadata.DateTime`
structure. It is a required field in all media.

The date and time come from the following tags, in priority order:

    EXIF DateTimeOriginal         [set]
    EXIF DateTimeDigitized        [set] [rm]
    EXIF DateTime                 [set] [rm]
    XMP  exif:DateTimeOriginal
    XMP  exif:DateTimeDigitized   [rm]
    XMP  ps:DateCreated
    XMP  xmp:CreateDate
    XMP  tiff:DateTime            [rm]
    XMP  xmp:ModifyDate           [rm]
    XMP  xmp:MetadataDate         [rm]
    IPTC DateTimeCreated          [set]
    IPTC DigitalCreationDateTime  [set] [rm]

## GPSCoords Field

The GPS coordinates of a media specify the place where the media was captured,
or failing that, the place where it should be shown on a map. It should be
considered a required field for all media that have locality, but it can be
omitted for media that aren't connected with a place (e.g. memes or songs). It
is stored as a `metadata.GPSCoords` structure, which contains a latitude,
longitude, and optional altitude.

The GPS coordinates come from the following tags, in priority order:

    XMP  exif:GPSCoords  [set]
    EXIF GPSCoords       [set]

## Groups Field

Optional groups associated with a media identify the groups of people (teams,
organizations, etc.) with which the media is associated. Group names form a
hierarchy. In most cases, group names will have only one component (e.g. "Schola
Cantorum"). In some cases, two levels of hierarchy are useful (e.g.
"Hewlett-Packard / KDT Team"). Use of more than two levels of hierarchy is not
expected. Group names are stored as a `strmeta.Group`, which is a string slice
with one element per level of the hierarchy.

Group names are stored internally as hierarchical keywords with an initial
component of `Groups`. As such, they come from the following tags, in priority
order:

    XMP  digiKam:TagsList
    XMP  lr:HierarchicalSubject
    XMP  dc:Subject              [flat]
    IPTC Keywords                [flat] [max]

NOTE: although group names are stored in keyword metadata tags, they are not
considered to be keywords in my metadata model, and they are not returned or
changed by the Keywords functions.

## Keywords Field

The keywords on a media give optional, unstructured additional information about
it that may be useful when trying to search for media with those keywords.
Keywords form a hierarchy; each keyword is stored as a `metadata.Keyword`, which
is a string slice with one element per level of the hierarchy.

Keywords come from the following tags, in priority order:

    XMP  digiKam:TagsList
    XMP  lr:HierarchicalSubject
    XMP  dc:Subject              [flat]
    IPTC Keywords                [flat] [max]

NOTE: in addition to storing keywords as defined here, these tags are also used
to store values for the Groups, People, Places, and Topics fields (which see).
This is an implementation detail hidden by this library, so (for example)
GetKeywords does not return such values, and SetKeywords does not change them.

## Location Field

The location of a media is a textual description of a location related to the
media: usually the place where it was captured, but sometimes a place depicted
in it. It is stored as a `metadata.Location` structure, which contains a country
code, country name, state name, city name, and sublocation name.

I don't use location any more, but much of my library has location values, so
this library handles them. Location values need not be added to new media, and
they should be removed from existing media if they are not accurate. To that
end, whenever the values of the Places field are changed, this library will
automatically remove the location value unless it is consistent with one of the
resulting Places values. (Consistent, in this context, means that the country
name, state name, city name, and sublocation name in the location appear in
order as components of a Place value, although that Place value may have
additional components.)

The location comes from the following tags, in priority order:

    XMP  iptc:LocationCreated  [alt]
    IPTC Location              [set] [max]
    XMP  iptc:LocationShown    [alt] [mult] [rm]

## People Field

The people field identifies people (or pets) depicted in a media. Each person
value is stored in a `strmeta.Person` structure, which contains the person's
name (see below) and also a flag indicating whether there is a face region in
the media associated with that person.

The metadata for a media should strive to include a person value for each person
prominently depicted in it. However, there is no need to try to name every
person who appears as a tiny speck in a crowd photo. The guideline should be,
"if I search for media depicting this person, would I want this media to be
found?"

The name that should be used for each person is the name that I normally call
them (i.e., Ken Roth, not Kenneth Roth). If someone uses a different name now
than they did when the media was captured, it's best to have a person value for
each name to aid in searching. (Exception: person values reflecting names the
person finds offensive or hurtful, such as deadnames, should be removed.)
Either way, when someone changes name, all media in the library should be
retagged to reflect that.

People values are stored internally as hierarchical keywords with an initial
component of `People`. As such, they come from the following tags, in priority
order:

    XMP  digiKam:TagsList
    XMP  lr:HierarchicalSubject
    XMP  dc:Subject              [flat]
    IPTC Keywords                [flat] [max]

In addition, people values may also be stored in tags describing face regions:

    XMP  MP:RegionInfo/MPRI:Regions/MPReg:PersonDisplayName  [ro]
    XMP  mwg-rs:RegionInfo/Regions/Name                      [ro]

People values that are stored as face regions cannot be changed by this library.
Any attempt to remove one will fail. People values added using this library
will not have face regions.

NOTE: although people values are stored in keyword metadata tags, they are not
considered to be keywords in my metadata model, and they are not returned or
changed by the Keywords functions.

## Places Field

Places associated with a media identify the location where it was captured, and
optionally the location(s) depicted in it. Places are considered a required
field for all media that have locality. Places form a hierarchy: country name,
state or province name (when applicable), region name (when helpful), city or
equivalent major area name, and refinements within that as needed. Each place is
stored as a `strmeta.Place`, which is a string slice with one element per level
of the hierarchy. Examples include:

    USA / California / Sierra Nevada / Yosemite National Park / Mist Trail
    USA / California / SF Bay Area / Sunnyvale / Baylands Park
    Austria / Vienna / St. Stephen's Cathedral
    Japan / Gifu / Takayama

When a place has a different name in English than in the language spoken in that
place, it's best to add two place values, covering both. For example, the last
two examples above would be seen alongside their native language equivalents:

    Österreich / Wien / Stephansdom
    日本 / 岐阜県 / 高山市

Note that names should be fully spelled out. The two defined exceptions to that
rule are both shown above: USA and SF Bay Area.

Places are stored internally as hierarchical keywords with an initial component
of `Places`. As such, they come from the following tags, in priority order:

    XMP  digiKam:TagsList
    XMP  lr:HierarchicalSubject
    XMP  dc:Subject              [flat]
    IPTC Keywords                [flat] [max]

NOTE: although places are stored in keyword metadata tags, they are not
considered to be keywords in my metadata model, and they are not returned or
changed by the Keywords functions.

## Title Field

The title of a media is an optional short, one-line title `string`, formatted in
title case.

The title comes from the following tags, in priority order:

    XMP  dc:title    [alt]
    IPTC ObjectName  [max]

## Topics Field

Optional topics associated with a media identify the topic of the media, e.g.,
the event, the activity taking place, or any other similar attribute that might
be helpful to people searching for similar media. Topic names form a hierarchy.
Topic names are stored as a `strmeta.Topic`, which is a string slice with one
element per level of the hierarchy.

Topic names are stored internally as hierarchical keywords with an initial
component of `Topics`. As such, they come from the following tags, in priority
order:

    XMP  digiKam:TagsList
    XMP  lr:HierarchicalSubject
    XMP  dc:Subject              [flat]
    IPTC Keywords                [flat] [max]

NOTE: although topic names are stored in keyword metadata tags, they are not
considered to be keywords in my metadata model, and they are not returned or
changed by the Keywords functions.

## Tag Notes

In the field descriptions above, notes indicate peculiarities and limitations of
each tag that contributes to the field. Here are definitions of these notes:

`[alt]` This underlying metadata tag allows for translations of each value into
multiple languages. My metadata model has only one language.  
`[flat]` This underlying metadata tag does not support hierarchical keywords;
it has a flat namespace. When a hierarchical keyword is stored in this tag,
only the final (leaf) component of it is stored.  
`[mult]` This underlying metadata tag allows for multiple values, but my
metadata model only allows one.  
`[max]` This underlying metadata tag has a length limitation. When a value is
stored in this tag, it is truncated to fit the limitation.  
`[rm]` This underlying metadata tag is read, if present, but is removed whenever
the field is rewritten.  
`[ro]` This underlying metadata tag is read-only. Attempts to change it will
result in an error.  
`[set]` This entry refers to a related set of underlying metadata tags.
