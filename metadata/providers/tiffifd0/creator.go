package tiffifd0

import (
	"fmt"
)

const tagArtist uint16 = 0x13B

// getCreator reads the value of the Creator field from the IFD.
func (p *Provider) getCreator() (err error) {
	tag := p.ifd.Tag(tagArtist)
	if tag == nil {
		return nil
	}
	p.artist, err = tag.AsString()
	if err != nil {
		return fmt.Errorf("Artist: %s", err)
	}
	return nil
}

// Creator returns the value of the Creator field.
func (p *Provider) Creator() (value string) { return p.artist }

// CreatorTags returns a list of tag names for the Creator field, and a
// parallel list of values held by those tags.
func (p *Provider) CreatorTags() (tags []string, values [][]string) {
	return []string{"IFD0 Artist"}, [][]string{{p.artist}}
}

// SetCreator sets the value of the Creator field.
func (p *Provider) SetCreator(value string) error {
	if value == "" {
		p.artist = ""
		p.ifd.DeleteTag(tagArtist)
		return nil
	}
	if p.artist == value {
		return nil
	}
	p.artist = value
	p.ifd.AddTag(tagArtist).SetString(value)
	return nil
}
