package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	artistName  = rdf.Name{Namespace: nsTIFF, Name: "Artist"}
	creatorName = rdf.Name{Namespace: nsDC, Name: "creator"}
)

// getCreator reads the value of the Creator field from the RDF.
func (p *Provider) getCreator() (err error) {
	var artist string

	if artist, err = getString(p.rdf.Property(artistName)); err != nil {
		// tiff:Artist is supposed to be a single string, but sometimes
		// I see it as a sequence of strings.  We'll accept that.
		if p.tiffArtist, err = getStrings(p.rdf.Property(artistName)); err != nil {
			return fmt.Errorf("tiff:Artist: %s", err)
		}
	} else if artist != "" {
		p.tiffArtist = []string{artist}
	}
	if p.dcCreator, err = getStrings(p.rdf.Property(creatorName)); err != nil {
		return fmt.Errorf("dc:creator: %s", err)
	}
	return nil
}

// Creator returns the value of the Creator field.
func (p *Provider) Creator() (value string) {
	if len(p.dcCreator) != 0 {
		return p.dcCreator[0]
	}
	if len(p.tiffArtist) != 0 {
		return p.tiffArtist[0]
	}
	return ""
}

// CreatorTags returns a list of tag names for the Creator field, and a
// parallel list of values held by those tags.
func (p *Provider) CreatorTags() (tags []string, values [][]string) {
	tags, values = []string{"XMP  dc:creator"}, [][]string{p.dcCreator}
	if len(p.tiffArtist) != 0 {
		tags = append(tags, "XMP  tiff:Artist")
		values = append(values, p.tiffArtist)
	}
	return tags, values
}

// SetCreator sets the value of the Creator field.
func (p *Provider) SetCreator(value string) error {
	p.tiffArtist = nil
	p.rdf.RemoveProperty(artistName)
	if value == "" {
		p.dcCreator = nil
		p.rdf.RemoveProperty(creatorName)
		return nil
	}
	if len(p.dcCreator) == 1 && p.dcCreator[0] == value {
		return nil
	}
	p.dcCreator = []string{value}
	p.rdf.SetProperty(creatorName, makeSeq(p.dcCreator))
	return nil
}
