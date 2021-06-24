// Package strmeta contains the definition of a Metadata structure reflecting my
// desired media metadata model, and the code to read it from and write it to
// real metadata sources.
package strmeta

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// The Metadata structure holds all of the metadata that I care about for a
// media file.
type Metadata struct {
	// Creator is the name of the person who captured the media, first name
	// first.  It may be suffixed with a comma, space, and company name if
	// relevant.  It may be a company name only if the person's name is not
	// known.
	Creator *metadata.String
	// DateTime is the date and time at which the media was originally
	// captured, in the form YYYY-MM-DDTHH:MM:SS.ss±HH:MM.  The fractional
	// seconds may have any number of digits or may be omitted.  The time
	// zone is omitted if not known.  UTC will have a time zone of "Z",
	// never "+00:00" or "-00:00".
	DateTime *metadata.DateTime
	// GPSCoords are the coordinates of the location where the media was
	// captured.
	GPSCoords *metadata.GPSCoords
	// LocationCaptured is the textual description of the location where the
	// media was captured.
	LocationCaptured metadata.Duolingual
	// LocationShown is the textual description of the location shown in the
	// media.  It is generally set only when different from
	// LocationCaptured.
	LocationShown metadata.Duolingual
	// Title is the title of the media: a short one-liner in title case.
	Title *metadata.String
	// Caption is a longer description of the media, using prose grammar and
	// capitalization.
	Caption *metadata.String
	// Keywords are the keywords associated with the media.  They reflect
	// the places, people, groups, and topics of the media, among other
	// things.
	Keywords []metadata.Keyword
}

/* Limitations:

Some of the metadata in the Metadata structure are more limited in their values
than the metadata formats of the underlying files can represent.  These include:
  * Some file formats can support multiple creators for a media.
  * Some file formats support more than one location shown in a media.
  * Some file formats support specifying the title and caption, and the place
    names in the locations, in multiple languages.
In addition, the underlying metadata formats support a plethora of other
metadata tags that I don't use at all.
*/
