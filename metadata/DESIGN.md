# Metadata Library

This library manages metadata in my media archive, presenting it as a simple and
regular metadata model and hiding the complexities of various file types and
encoding schemes.

The simple metadata model into which all media are squeezed is:

```Go
Caption  string
Creator  string
DateTime metadata.DateTime
Faces    []string
GPS      metadata.GPSCoords
Groups   []metadata.HierValue
Keywords []metadata.HierValue
Location metadata.Location
People   []string
Places   []metadata.HierValue
Title    string
Topics   []metadata.HierValue
```

## The `metadata` Package

The top-level `metadata` package defines the data types used in this model:
`DateTime`, `GPSCoords`, `HierValue`, and `Location`. It also defines the
interface for a metadata `Provider` that allows reading those values from, and
changing them in, an arbitrary metadata source.

For each metadata field `XXX` of type `T`, the `Provider` interface contains
three functions:

```Go
XXX() (value T)
XXXTags() (tags []string, values [][]T)
SetXXX(value T) (err error)
```

The first function, with the plain name of the field, returns the "best" value
for the field. Since a media file may have multiple sources of metadata for a
field, they are prioritized, and the "best" value is the one that comes from the
highest priority source that has a value for that field. (For multi-valued
fields, "best" means all of the values from the highest priority source that has
any values.)

The second function, `XXXTags`, returns all of the values of the field from all
sources in the media file. Its return values are based on this algorithm,
applied to each source in priority order:

1. If the source has multiple values for the field, an element is added to
   `tags` with the name of the source, and an element is added to `values` with
   the values from that source. The expectation is that none of these values is
   empty, but that is generally not checked.
2. If the source has a single, non-empty value for the field, an element is
   added to `tags` with the name of the source, and an element is added to
   `values` with a single-element slice containing the value.
3. If the source does not exist in the file, has no value, or has a single empty
   value, but calling `SetXXX` with a non-empty value would create the source
   and give it that non-empty value, then an element is added to `tags` with the
   name of the source, and an element is added to `values` with a nil slice.
4. If the source does not exist in the file, has no value, or has a single empty
   value, and calling `SetXXX` with a non-empty value would not create the
   source or give it that non-empty value, then nothing is added to `tags` or
   `values` for that source.

The third function, `SetXXX`, sets the value(s) for the field in all of the
sources where my metadata model says it should be recorded. It also removes old
values of the field from any other sources that are not part of my metadata
model.

The assertion between these three functions is that, if `SetXXX` is called with
(possibly empty) value(s) for a field, a subsequent `XXX` call will return those
same values, and a subsequent `XXXTags` call will return one or more sources for
that field, all of which have those same values, and no other sources.

Note, however, that "same value" is used guardedly. Underlying metadata do not
always store all values precisely. For example, EXIF stores latitude values as
rational numbers, while XMP stores them as floating point numbers; the
conversions between the two are not exact. Similarly, some DateTime sources can
store subsecond granularity and others cannot. And some sources truncate string
values to a maximum length, while other sources do not. To the extent possible,
such variances are ignored.

## The `filefmts` Package

The `filefmts` package contains an interface that all file format handlers must
honor, and a factory function that returns the file format handler appropriate
for a particular file.

Each file format has a subpackage under `filefmts` that defines a handler for
that format. Each file format subpackage contains a `Read` function which,
given a handle to a file, determines whether the file is of the type it knows
how to handle, and if so, returns a handler for it. (Determination is made
based on file contents, not on file name.)

The file format handlers each have a `Provider` function, which returns the
metadata provider for that file. The provider can be used to query and set
metadata within the file, as described above. In many cases, the provider that
is returned is actually a merger of multiple providers from different metadata
structures in the file, but this is transparent to the caller.

## The `containers` Packages

Media files are containers, using various encoding schemes to contain a variety
of identified blocks of information. In many cases, those containers nest: the
top-level container may include identified blocks that are themselves
containers. Real-world media files often have as many as three levels of nested
containers.

The subpackages of `containers` each define a container type. They each have a
`Read` function that knows how to parse the container, and a `Render` function
that knows how to save it. They also provide container-specific functions that
providers based on that container can call to query and set values in the
container. Note that any call to a set function on a container marks that
container dirty; it is the responsibility of each provider to avoid calling the
container's set function if the underlying data is not changing.

The defined containers are:

- `iim` is the container for the IPTC Information Interchange Model, an obsolete
  but still widely used container format.
- `jpeg` is the container format for a JPEG file. The JFIF and EXIF standards
  describe conflicting requirements for JPEG files; this container format (like
  most modern software working with JPEGs) reads both, and writes something that
  doesn't technically comply with either standard, but that virtually all
  JPEG-reading software can handle.
- `photoshop` is the container format for Photoshop Information Resources
  (PSIRs).
- `rdf` is the container format for Extensible Metadata Platform (XMP) metadata,
  the newest and most complete metadata format.
- `tiff` is the container format defined by the Tagged Image File Format (TIFF)
  specification. In addition to being the top-level container format for TIFF
  and DNG files, this container format is also used in the EXIF segments of JPEG
  files.

## The `providers` Packages

The subpackages of `providers` each offer a `metadata.Provider` that allows
reading and writing the metadata in a particular container. In some cases,
multiple providers handle different subsets of a container, to allow those
subsets to be placed differently in the provider priority order. These
providers handle all of the details of translating between the metadata format
in the container and my canonical metadata model.

Each `provider` subpackage has a `New` function, which takes the appropriate
container instance as input, and returns a `metadata.Provider` for it.

The defined provider subpackages are:

- `exififd`: Provider for the EXIF IFD in a TIFF container.
- `gpsifd`: Provider for the GPS IFD in a TIFF container.
- `iptc`: Provider for the IPTC data in an IIM container.
- `jpegifd0`: Provider for the root IFD in the EXIF TIFF container of a JPEG
  file.
- `multi`: Provider that merges the results of a list of other providers.
- `tiffifd0`: Provider for the root IFD in a TIFF file.
- `xmp`: Provider for the native XMP metadata in an XMP/RDF container.
- `xmpexif`: Provider for the mirror of EXIF metadata in an XMP/RDF container.
- `xmpext`: Pseudo-provider for the XMP metadata in an XMP extension segment of
  a JPEG file. This package doesn't actually provide anything; it just does
  error checking.
- `xmpiptc`: Provider for the mirror of IPTC metadata in an XMP/RDF container.
- `xmpps`: Provider for the mirror of Photoshop metadata in an XMP/RDF
  container.
- `xmptiff`: Provider for the mirror of TIFF metadata in an XMP/RDF container.

## File Structures

A JPEG file will have some or all of the following structure:

```x
JPEG container
    EXIF segment
        TIFF container
            IFD0
                jpegifd0 provider
            EXIF IFD
                exififd provider
            GPS IFD
                gpsifd provider
    Photoshop segment
        Photoshop container
            IPTC PSIR
                iim container
                    iptc provider
    XMP segment
        rdf container
            xmp provider
            xmpexif provider
            expiptc provider
            xmpps provider
            xmptiff provider
    XMP extension segment
        rdf container
            xmpext provider
```

A TIFF file will have some or all of the following structure:

```x
TIFF container
    IFD0
        tiffifd0 provider
        IPTC tag
            iim container
                iptc provider
        Photoshop tag
            Photoshop container
                IPTC PSIR
                    iim container
                        iptc provider
        XMP tag
            rdf container
                xmp provider
                xmpexif provider
                expiptc provider
                xmpps provider
                xmptiff provider

    EXIF IFD
        exififd provider
    GPS IFD
        gpsifd provider
```
