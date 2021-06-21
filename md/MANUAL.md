# md Manual

The `md` command reads and manipulates media file metadata, following the
conventions I use in my media library.

    usage: md [fieldname][op][value]... file...

Run md with a list of zero or more queries or changes, followed by one or more
names of files to query or change. If no queries or changes are listed, md will
show all of the metadata of all of the listed files.

Queries are one or more field names without assignments:

    a(rtist)
    d(atetime)
    t(itle)
    c(aption)
    k(eywords)
    g(ps)
    l(ocation)
    s(hown)

These can be fully written out, or written as single letters. In single letter
form, they can be combined without spaces, so `da` means `datetime artist`.

The following keywords are also queries, when used without an assignment:

    places
    people
    groups
    topics

They must be fully spelled out, but they are case-insensitive. They are
equivalent to a keyword query for hierarchical keywords whose first component is
the all-caps version of that word, e.g. `people` looks for "PEOPLE/..."
keywords.

Non-keyword assignments are one of the non-keyword field names:

    a(rtist)
    d(atetime)
    t(itle)
    c(aption)
    g(ps)
    l(ocation)
    s(hown)

followed by an `=` sign, optionally followed by a value. (Omitting the value
deletes the corresponding metadata tags.) As above, the field names can be
fully written out, or written as single letters, but the single letter forms
cannot be combined. See below for value formats.

Keyword assignments have one of the following forms:

    k(eywords)=         removes all keywords
    k(eywords)+=value   adds a keyword
    k(eywords)-=value   removes a keyword

Arguments beginning with a + or - are taken to be keyword assignments:

    +value              adds a keyword
    -value              removes a keyword

In addition, the four top-level keywords (`places`, `people`, `groups`, and
`topics`) get special handling:

    word=             removes all WORD/... keywords
    word=value        removes all WORD/... keywords, then adds one
    word+=value       adds a WORD/... keyword
    word-=value       removes a WORD/... keyword

Again, the word must be fully spelled out, but is case-insensitive.

Artist values are plain strings; while the metadata definitions support multiple
values, I don't use them. They should be first-name-first, e.g. "Steven Roth",
possibly followed by comma, space, company name if relevant.

Datetime values are `YYYY-MM-DDTHH:MM:SS.sss±HH:MM`, with the `T` and everything
after it optional. If the time is omitted, it is recorded as midnight. If the
time zone is omitted, no time zone is recorded.

Title is a plain, single-line string; language variants are not used.

Caption is a plain, multi-line string. A single dash (`-`) means to read the
caption from standard input. (This works for captions only.)

GPS coordinates are a string containing two or three floating point numbers,
separated by commas and optional spaces. The first two numbers are the latitude
and longitude in degrees. The third number, if present, is the altitude in
meters. GPS coordinates correspond to the location captured (`l(ocation)`), not
the location shown (`s(hown)`).

Location captured (`l(ocation)`) and location shown (`s(hown)`) have the same
format: a set of between zero and four components separated by slashes. (The
slashes can have spaces around them, which are ignored.) These components are
interpreted in order as follows:

1. A three-letter ISO-3166 alpha-3 country code. It is not case sensitive. The
   corresponding country name is looked up by the tool.
2. A state or province name, or, if the first component is USA, a two-letter
   state code (not case sensitive).
3. A name of a city within the state.
4. A name of a sublocation within the city or, if the city is blank, within the
   state.

Hierarchical keywords are specified as a sequence of levels separated by
slashes. A dot at the beginning of a level says that level is for grouping only
and is not an actual keyword itself. The four well-known top-level keywords are
always for grouping only. Thus

    TOPICS/.Activities/Hockey

would be recorded as one hierarchical keyword

    TOPICS/Activities/Hockey

and one flat keyword

    Hockey

while

    GROUPS/Hewlett-Packard/SSHA Team

would be recorded as two hierarchical keywords

    GROUPS/Hewlett-Packard/SSHA Team
    GROUPS/Hewlett-Packard

and two flat keywords

    Hewlett-Packard
    SSHA Team

== Output formats ==

If a single file is being queried, the output a line for each value of each
requested field, containing the field name, possible change marker, and field
value. Embedded newlines and backslashes in the value are escaped. The field
name is omitted if only one field was requested.

As a special case, if a single file is being queried, and the only field being
queried is the caption, the caption value is emitted to standard output by
itself, unescaped, with a trailing newline.

If multiple files are being queried, but only a single field, each line starts
with the filename, followed by the possible change marker and field value.

If multiple files and multiple fields are being queried, each file's report
starts with a `=== FILENAME ===` line, followed by the output as if it was the
only file.

The change marker `CHG` appears on a value if setting the field to the value
shown would make any change to the metadata tags. This could happen if:

- The tags for a field have conflicting values.
- The value for a field does not confirm to our metadata model.
- The value for a field is not expressed in all of the canonical tags.
- The value for a field is expressed in any non-canonical tags.
