# Metadata Library Design

## Package `filefmt`

Provides a file-type-independent interface for working with media files, and a
function to return the correct implementation of that interface for any given
media file. The interface knows how to extract metadata blocks from the file,
and how to update the file with new or revised metadata blocks.

## Packages `filefmt/jpeg` and `filefmt/xmp`

Provide implementations of the media file interface for JPEG and XMP files,
respectively.

## Package `metadata`

Provides data types for the various metadata fields of interest, that are
independent of how those data types are encoded in different metadata blocks.
These data types are designed to support the union of possible values for those
fields as stored in any of the supported metadata block formats.

## Packages `metadata/exif`, `metadata/iptc`, and `metadata/xmp`

Provide structures and methods for decoding and encoding EXIF, IPTC-IIM, and XMP
metadata blocks, respectively, and for translating the data in them to and from
the types defined in package `metadata`.

## Package `strmeta`

Provides types and methods for reconciling the data from multiple metadata
blocks in a media file, and digesting them down to the simplified metadata view
used in my photo tools.
