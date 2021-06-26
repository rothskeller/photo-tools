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
    set fieldname value
    add fieldname value
    remove fieldname value
    clear fieldname
    choose fieldname
    copy [fieldname...]
    write caption
    read caption

The `show` operation displays the value of each named field in each named file.
They are shown in a table with file name, field name, and field value columns.
Where a file's metadata has conflicting values for a field, only the value(s)
from the highest priority metadata tag are shown. If no fields are named, `all`
is assumed.

The `tags` operation displays the values of each named field in each named file.
They are shown in a table with file name, metadata tag name, and metadata tag
value columns. All values of all metadata tags for the requested fields are
shown. If no fields are named, `all` is assumed.

The `set` operation removes all values of the named field, and then adds the
specified value, in each of the named files. As a safety precaution,
`set keyword` is not allowed.

The `add` operation adds the specified value to the list of values for the named
field, in each of the named files. It is valid only for fields that can have
multiple values.

The `remove` operation removes the specified value from the list of values for
the named field, in each of the named files. It is valid only for fields that
can have multiple values.

The `clear` operation removes all values of the specified field, and all
corresponding metadata tags, from each of the named files.

The `choose` operation displays all values the named field in the named files,
just like the `tags` operation. It then allows the user to choose one of those
values (or manually enter some other value), which it applies to each of the
named files just like the `set` operation. The `choose` operation is not valid
for `group`, `keyword`, `person`, and `place`, and `topic` fields.

The `copy` operation requires at least two named files. It copies the values of
the named field(s) from the first named file to all of the other named files.
If no fields are named, all fields are copied.

The `write caption` operation is like `show caption`, except that the caption is
written to standard output without any table formatting.

The `read caption` operation is like `set caption`, except that the value is
read from standard input rather than taken on the command line.

If the `batch` prefix is given, the named files are batched by basename, and the
operations are run against each batch separately. This is really only useful
for the `choose` operation, although it also produces cosmetic differences in
the output of the `show` and `tags` operations.

## Fields

The possible field names (and allowed abbreviations) are:

    all
    artist    (a)
    caption   (c)
    datetime  (d, date, time)
    gps       (g)
    group     (groups)
    keyword   (k, kw, keywords)
    location  (l, loc)
    person    (people)
    place     (places)
    title     (t)
    topic     (topics)

The `all` word is the same as listing all field names. It is allowed only for
the `show` and `tags` operations.

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

The `gps` field is the GPS coordinates of the location where the media was
captured (or, if not known exactly, the place where they should be shown on a
map). It is represented as two or three signed floating point numbers separated
by commas. The first two are the latitude and longitude in degrees. If a third
one is present, it is the altitude, and must be followed by a suffix of `m`
(meters) or `ft` (feet). (On output, altitude is always reported in feet.)

The `keyword` field contains a list of keywords associated with the media.
Keywords are hierarchical, with components separated by slashes. The `group`,
`person`, `place`, and `topic` fields are shorthand for `keyword` with the first
component of the value assumed to be `GROUPS`, `PEOPLE`, `PLACES`, and `TOPICS`,
respectively.

The `location` field contains a textual description of the location of the
media. It has the form

    countrycode / countryname / state / city / sublocation

Parts that are unused can be left blank, and trailing slashes can be omitted.
Spaces around the slashes are optional and insignificant.

The `title` field contains a one-line, short title for the media, expressed in
title case. While some metadata tags allow for titles to be provided in
multiple languages, `md` treats it as a language-invariant, single-value field
(presumed to be English).
