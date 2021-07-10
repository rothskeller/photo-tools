package xmptiff

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var artistName = rdf.Name{Namespace: nsTIFF, Name: "Artist"}

// getCreator reads the value of the Creator field from the RDF.
func (p *Provider) getCreator() (err error) {
	var artist string

	if artist, err = getString(p.rdf.Properties, artistName); err != nil {
		// tiff:Artist is supposed to be a single string, but sometimes
		// I see it as a sequence of strings.  We'll accept that.
		if p.tiffArtist, err = getStrings(p.rdf.Properties, artistName); err != nil {
			return fmt.Errorf("tiff:Artist: %s", err)
		}
	} else {
		p.tiffArtist = []string{artist}
	}
	return nil
}

// Creator returns the value of the Creator field.
func (p *Provider) Creator() (value string) {
	if len(p.tiffArtist) == 0 {
		return ""
	}
	return p.tiffArtist[0]
}

// CreatorTags returns a list of tag names for the Creator field, and a
// parallel list of values held by those tags.
func (p *Provider) CreatorTags() (tags []string, values []string) {
	for _, artist := range p.tiffArtist {
		tags = append(tags, "XMP  tiff:Artist")
		values = append(values, artist)
	}
	return tags, values
}

// SetCreator sets the value of the Creator field.
func (p *Provider) SetCreator(value string) error {
	p.tiffArtist = nil
	if _, ok := p.rdf.Properties[artistName]; ok {
		delete(p.rdf.Properties, artistName)
		p.dirty = true
	}
	return nil
}
