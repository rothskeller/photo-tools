# md Manual

The `md` command reads and manipulates media file metadata, following the
conventions I use in my media library.

    usage: md [batch] [operation...] file...

Files and operations may occur in any order, but there must be at least one
file. If no operation is given, `show all` is assumed. Operations are applied
in sequence, left to right.

## Operations

The possible operations are:

    show [fieldname...]
    tags [fieldname...]
    check [fieldname...]
    set fieldname value
    add fieldname value
    remove fieldname value
    reset [fieldname...]
    clear fieldname
    choose fieldname
    copy [fieldname...]
    write caption
    read caption

If no fields are named For the `show`, `tags`, `check`, `reset`, or `copy`
operations, or if they are given a field name of `all`, they act on all known
fields.

The `show` operation displays the value of each named field (or all fields) in
each named file. They are shown in a table with file name, field name, and
field value columns. Where a file's metadata has conflicting values for a
field, only the value(s) from the highest priority metadata tag are shown.

The `tags` operation displays the values of each named field (or all fields) in
each named file. They are shown in a table with file name, metadata tag name,
and metadata tag value columns. All values of all metadata tags for the
requested fields are shown.

The `check` operation verifies that all of the named files are correctly tagged,
and have the same values, for the named fields (or all fields). It displays a
table with one row per named file and one column per named field (plus the
leftmost column for the file name). In each cell of this table, it displays one
of the following:

    '  ' for an optional field that is not set
    '--' for an expected field that is not set
    ' ✓' for a single-valued field that is set, and tagged correctly
    ' 3' value count for a multi-valued field that is set, and tagged correctly
    '!=' for a field whose tags don't agree with each other, or don't agree with
         the field in the first-named file
    '[]' for a field whose value isn't tagged correctly

The `set` operation removes all values of the named field, and then adds the
specified value, in each of the named files. As a safety precaution,
`set keyword` is not allowed.

The `add` operation adds the specified value to the list of values for the named
field, in each of the named files. It is valid only for fields that can have
multiple values.

The `remove` operation removes the specified value from the list of values for
the named field, in each of the named files. It is valid only for fields that
can have multiple values.

The `reset` operation corrects the tagging of the named fields (or all fields)
in all named files, using the value(s) from the highest priority metadata tag
for those fields (i.e., the same one shown by `show`, generally the first one
listed by `tags`).

The `clear` operation removes all values of the specified field, and all
corresponding metadata tags, from each of the named files.

The `choose` operation displays all values of the named field in the named
files, just like the `tags` operation. It then allows the user to choose one of
those values (or manually enter some other value), which it applies to each of
the named files just like the `set` operation. The `choose` operation is not
valid for fields that can have multiple values.

The `copy` operation requires at least two named files. It copies the values of
the named fields (or all fields) from the first named file to all of the other
named files.

The `read caption` operation is like `show caption`, except that the caption is
written to standard output without any table formatting.

The `write caption` operation is like `set caption`, except that the value is
read from standard input rather than taken on the command line.

If the `batch` prefix is given, the named files are batched by basename, and the
operations are run against each batch separately. This is really only useful for
the `check` and `choose` operations, although it also produces cosmetic
differences in the output of the `show` and `tags` operations.

## Fields

The possible field names (and allowed abbreviations) are:

    artist    (a)
    caption   (c)
    datetime  (d, date, time)
    face      (f, faces)
    gps       (g)
    group     (groups)
    keyword   (k, kw, keywords)
    location  (l, loc)
    person    (people)
    place     (places)
    title     (t)
    topic     (topics)

The `artist` field is the name of the person who captured the original media,
e.g. "Steven Roth". When appropriate, it could be a company name, or a person's
name and company name separated by a comma and space. While many metadata tags
support multiple artist values, `md` treats it as a single-value field.

The `caption` field is a prose description of the media. While some metadata
tags allow for captions to be provided in multiple languages, `md` treats it as
a language-invariant, single-value field (presumed to be English).

The `datetime` field is the date and time at which the original media was
captured, as precisely as is known. It is represented in RFC 3339 format, i.e.,
YYYY-MM-DDTHH:MM:SS.sss±HH:MM. On input, the THH:MM:SS.sss can be omitted, in
which case midnight is assumed (and will be subsequently reported). Fractional
seconds can be omitted. The time zone can be omitted, indicating that it is
unknown. `Z` can be used on input, and is always used on output, in place of
`+00:00` or `-00:00` to represent UTC.

The `face` field applies only to images; it is the list of people names
associated with face regions in the image. See Special Behaviors, below, for
the relationship between the `face` and `person` fields. Note that this tool
cannot add face regions; it can only recognize them and remove them.

The `gps` field is the GPS coordinates of the location where the media was
captured (or, if not known exactly, the place where they should be shown on a
map). It is represented as two or three signed floating point numbers separated
by commas. The first two are the latitude and longitude in degrees. If a third
one is present, it is the altitude, and must be followed by a suffix of `m`
(meters) or `ft` (feet). (On output, altitude is always reported in feet.)

The `group` field contains a list of groups (teams, organizations, etc.) that
are depicted in the media. Group names are hierarchical, with components
separated by slashes.

The `keyword` field contains a list of keywords associated with the media.
Keywords are hierarchical, with components separated by slashes. Note that
while values of the `group`, `person`, `place`, and `topic` are stored in
underlying metadata as keywords, they are not reported or managed by the
`keyword` field. The `keyword` field only reports and acts on other keywords.

The `location` field contains a textual description of the location of the
media. It has the form

    countrycode / countryname / state / city / sublocation

Parts that are unused can be left blank, and trailing slashes can be omitted.
Spaces around the slashes are optional and insignificant. English-language names
should be used when they exist. See Special Behaviors, below for the
relationship between the `location` and `place` fields.

The `person` field contains a list of names of people (or in a few cases, pets)
who are depicted in the media. People should be listed by full name, as they
are informally addressed. See Special Behaviors, below, for the relationship
between the `face` and `person` fields.

The `place` field contains a list of places relating to the media: the place
where it was captured and/or the place(s) depicted in it. Places are
hierarchical, with components separated by slashes. See Special Behaviors,
below, for the relationship between the `location` and `place` fields. The
components of a place typically include:

- Country name (spelled out unless it is "USA")
- State, province, or similar region of the country (spelled out) (if applicable)
- Region within the state (currently only used with CA and HI)
- City name, or name of similarly prominent location (park, monument)
- Refinements within the city or similar location (landmorks)

If the name of a place is different in English from its name as spoken by the
people who live there, the `place` field should contain two values for that
place, one with each name.

The `title` field contains a one-line, short title for the media, expressed in
title case. While some metadata tags allow for titles to be provided in
multiple languages, `md` treats it as a language-invariant, single-value field
(presumed to be English).

The `topic` field contains a list of topics of the media (activities, events,
etc.). Topic names are hierarchical, with components separated by slashes.

## Special Behaviors

The `face` and `person` fields are related; for every face region, there should
be a like-named person keyword. The following special behaviors apply:

- The `show` operation does not show values of the `person` field for which
  corresponding values of the `face` field are also being shown.
- Face regions without a corresponding person keyword are flagged by `show` and
  `check` as inconsistencies.
- Removing a person value also removes any corresponding face region.

The `location` and `place` fields are related; if a media has a location, it
should also have a congruent place value, such that the countryname, state,
city, and sublocation components of the location appear as components of the
place value, in the same order but possibly interspersed with other components
of the place value. The following special behaviors apply:

- Locations without a congruent place value are flagged by `show` and `check` as
  inconsistencies.
- Changing the place values will clear the location unless it is congruent with
  one of the resulting place values.
