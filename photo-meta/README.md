# photo-meta

The photo-meta command is used to set photo metadata to match my timeline media
metadata standard (see media/text/setup/timeline.txt).  Its usage is:

```x
photo-meta [options] image-file...
    -read-xmp
    -adjust-time offset
    -tz offset
    -date date-time-specification
    -gps lat,long[,alt]
    -loc countrycode:state:{city|city:sublocation|sublocation}
    -creator creator
    -mine
    -person [+-=]person,...
    -title title
    -desc description|@filename
    -kw [+-=]keyword,...
```

For each listed image file, it will validate and update the date and time,
location, creator, people list, title, description, and keywords, as follows.

## XMP Sidecar Files

When handling images in formats that support embedded metadata, photo-meta
generally ignores sidecar .xmp files (issuing a warning if one exists).  When
the -read-xmp flag is specified, the metadata embedded in the image is merged
with that from the XMP file, with the XMP file taking precedence, prior to
following the algorithms below.  Results are always written to the embedded
metadata; the XMP file is never touched.

When handling images in formats that don't support embedded metadata, photo-meta
reads and writes the XMP sidecar file, regardless of whether -read-xmp was
specified.

## Date and Time

Pseudocode for date and time handling:

```x
If there is an EXIF DateTimeOriginal value
    If it is consistent with the year and month in the directory name
        If -date was specified
            Issue warning
        If -set-tz was specified
            Set OffsetTimeOriginal according to -set-tz value
        If OffsetTimeOriginal is not present
            Issue warning and set OffsetTimeOriginal to PST or PDT as appropriate.
        Image time is taken from DateTimeOriginal, SubSecTimeOriginal (if present), and OffsetTimeOriginal.
        If -adjust-time was specified
            Adjust image time as requested.
            Set DateTimeOriginal to image time.
            If SubSecTimeOriginal is present or new time has fractional seconds
                Set SubSecTimeOriginal to image time.
    Else
        Remove EXIF DateTimeOriginal, SubSecTimeOriginal, and OffsetTimeOriginal.
If no image time found by above
    If -set-tz or -adjust-time were specified
        Issue warning
    If -date was specified
        Image time is taken from -date.
    Else if XMP-dc:Date is present
        Image time is taken from XMP-dc:Date
    Else if IPTC DateCreated is present
        Image time is taken from IPTC DateCreated and (if present) TimeCreated
    Else
        Prompt and read image time from stdin.  (User can decline to provide.)
If image time found by above
    Set XMP-dc:Date to image time.
    If IPTC DateCreated is present
        Set IPTC DateCreated and TimeCreated to image time.
```

## Location

Pseudocode for location handling:

```x
If EXIF has GPS data:
    Copy it to XMP-iptcExt:LocationCreated.
    If -gps was specified
        Issue warning
Else if -gps was specified:
    Save it to XMP-iptcExt:LocationCreated.
Else if there is no GPS data in XMP:iptcExt:LocationCreated:
    Prompt and read GPS data from stdin.
    If provided:
        Save it to XMP-iptcExt:LocationCreated.
If -loc was specified:
    Save it to XMP-iptcExt:LocationCreated.
Else if there is no text data in XMP:iptcExt:LocationCreated:
    Prompt and read text location data from stdin.
    If provided:
        Save it to XMP-iptcExt:LocationCreated.
```

## Creator

The creator of the image is taken from the -creator flag, the -mine flag
(equivalent to -creator 'Steve Roth'), the XMP-dc:Creator tag, the IPTC By-line
tag, the EXIF:Artist tag, or stdin, in order of preference.  If provided, it is
set in XMP-dc:Creator, and updated in IPTC:By-line and/or EXIF:Artist if those
tags exist.

## People

If -person is specified, the named people are added to, removed from, or set
into the XMP-iptcExt:PersonInImage tag.  Otherwise, if no such tag exists,
named people are read from stdin and added to that tag.

## Title, Description, and Keywords

If -title is specified, it is stored in XMP-dc:Title.  Otherwise, if no such tag
exists, the value of it is read from IPTC Headline or from stdin.  If the IPTC
Headline tag exists, the title is written to it.

If -description is specified, it is stored in XMP-dc:Description.  Otherwise, if
no such tag exists, the value of it is read from IPTC Caption-Abstract, EXIF
ImageDescription, or stdin.  If the IPTC Caption-Abstract and/or EXIF
ImageDescription tags exist, the description is written to them.

If -kw is specified, the named keywords are added to, removed from, or set into
the XMP-dc:Subject tag.  Otherwise, if no such tag exists, the values of it are
read from IPTC Keywords or from stdin.  If the IPTC Keywords tag exists, the
keywords written to it as well.
