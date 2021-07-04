# md Manual

The `md` command reads and manipulates media file metadata, following the
conventions I use in my media library.

    usage: md [file-selection] [operation]

## File Selection

`md` remembers a set of files, and a targeted subset of them, from one
invocation to the next. When `md` is invoked without any file selection
arguments, it acts on the targeted subset of the remembered set.

The remembered set is cleared whenever `md` is invoked from a new working
directory. If it is invoked from a new working directory without any file
selection arguments, its remembered set and targeted subset default to
containing all files in the new working directory that are of a supported type.
In this default case, operations that would modify files are disallowed to
prevent accidents.

File selection on the command line can be a list of files (not necessarily all
in the current directory), or one of the keywords `all`, `batch`, `next`,
`prev`, or `select`. If files are listed on the command line, they become the
new remembered set and targeted subset.

The `all` keyword sets the targeted subset to the entire remembered set.

The `batch` keyword sets the targeted subset to the first "batch" of files in
the remembered set that have the same basename (i.e., the filenames are the
same up to the first period).

The `next` and `prev` keywords change the targeted subset to the next or
previous batch, respectively, of files in the remembered set. They are valid
if the current targeted subset resulted from a previous `batch`, `next`, or
`prev` selection.

The `select` keyword prints a numbered list of the files in the remembered set
and allows the user to select, by number, which one(s) to include in the new
targeted subset.

## Operations

The possible operations are:

    add fieldname values
    check
    choose fieldname
    clear fieldname
    copy [fieldname...]
    read caption
    remove fieldname values
    reset [fieldname...]
    set fieldname values
    show [fieldname...]
    tags [fieldname...]
    write caption

Operation names can be abbreviated as long as they remain unique. If no
operation is given on the command line, `check` is assumed.

If no fields are named for the `copy`, `reset`, `show`, or `tags` operations, or
if they are given a field name of `all`, they act on all known fields.

All command line arguments after the field name for `add`, `remove`, and `set`
operations are joined together with a single space (to minimize the need for
quoting on the command line). The result is then split on semicolons into
individual values for the field. Whitespace around the values is ignored.

The `add` operation adds the specified value(s) to the list of values for the
named field, in each of the target files. It is valid only for fields that can
have multiple values.

The `check` operation verifies that all of the named files are correctly tagged.
It displays a table with one row pernamed file and one column per field (plus
the leftmost column for the file name). In each cell of this table, it displays
one of the following:

    '  ' for an optional field that is not set
    '--' for an expected field that is not set
    ' ✓' for a single-valued field that is set, and tagged correctly
    ' 3' value count for a multi-valued field that is set, and tagged correctly
    '!=' for a field whose tags don't agree with each other
    '[]' for a field whose value isn't tagged correctly

The `choose` operation displays all values of the named field in the target
files, just like the `tags` operation. It then allows the user to choose one of
those values (or manually enter some other value), which it applies to each of
the target files just like the `set` operation.

The `clear` operation removes all values of the specified field, and all
corresponding metadata tags, from each of the target files.

The `copy` operation requires at least two target files. It copies the values of
the named fields (or all fields) from the first target file to all of the other
target files.

The `read caption` operation is like `show caption`, except that the caption is
written to standard output without any table formatting.

The `remove` operation removes the specified value(s) from the list of values
for the named field, in each of the target files. It is valid only for fields
that can have multiple values.

The `reset` operation corrects the tagging of the named fields (or all fields)
in all target files, using the value(s) from the highest priority metadata tag
for those fields (i.e., the same one shown by `show`, generally the first one
listed by `tags`).

The `set` operation removes all values of the named field, and then adds the
specified value(s), in each of the target files.

The `show` operation displays the value of each named field (or all fields) in
each named file. They are shown in a table with file name, field name, and
field value columns. Where a file's metadata has conflicting values for a
field, only the value(s) from the highest priority metadata tag are shown.

The `tags` operation displays the values of each named field (or all fields) in
each named file. They are shown in a table with file name, metadata tag name,
and metadata tag value columns. All values of all metadata tags for the
requested fields are shown.

The `write caption` operation is like `set caption`, except that the value is
read from standard input rather than taken on the command line.

## Fields

The possible field names (and allowed variations) are:

    artist
    caption
    datetime  (time)
    face
    gps
    group
    keyword   (kw)
    location
    person
    place
    title
    topic

Field names can be abbreviated as long as they remain unique.

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
