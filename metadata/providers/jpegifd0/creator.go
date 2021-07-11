package jpegifd0

import (
	"fmt"
	"strings"
)

const tagArtist uint16 = 0x13B

// getCreator reads the value of the Creator field from the RDF.
func (p *Provider) getCreator() (err error) {
	tag := p.ifd.Tag(tagArtist)
	if tag == nil {
		return nil
	}
	alist, err := tag.AsString()
	if err != nil {
		return fmt.Errorf("Artist: %s", err)
	}
	// According to the Exif specification, this is a semicolon-separated
	// list of artists; each one may be quoted with quotes if it contains a
	// semicolon or a quote; quotes in the keyword are doubled.  This parser
	// is forgiving: the quotes don't have to surround the entire keyword,
	// and a missing end quote is assumed.
	var abuf string
	var inquotes bool
	var escape bool
	for _, b := range alist {
		switch {
		case escape && b != '"':
			escape = false
			inquotes = false
		case b == '"' && inquotes && escape:
			abuf += `"`
			escape = false
		case b == '"' && inquotes:
			escape = true
		case b == '"':
			inquotes = true
		case b == ';' && !inquotes:
			if t := strings.TrimSpace(abuf); t != "" {
				p.artist = append(p.artist, t)
			}
		default:
			abuf += string(b)
		}
	}
	if t := strings.TrimSpace(abuf); t != "" {
		p.artist = append(p.artist, t)
	}
	return nil
}

// Creator returns the value of the Creator field.
func (p *Provider) Creator() (value string) {
	if len(p.artist) == 0 {
		return ""
	}
	return p.artist[0]
}

// CreatorTags returns a list of tag names for the Creator field, and a
// parallel list of values held by those tags.
func (p *Provider) CreatorTags() (tags []string, values [][]string) {
	return []string{"IFD0 Artist"}, [][]string{p.artist}
}

// SetCreator sets the value of the Creator field.
func (p *Provider) SetCreator(value string) error {
	if value == "" {
		p.artist = nil
		p.ifd.DeleteTag(tagArtist)
		return nil
	}
	if len(p.artist) == 1 && p.artist[0] == value {
		return nil
	}
	p.artist = []string{value}
	if strings.IndexAny(value, `";`) >= 0 {
		value = `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
	}
	p.ifd.AddTag(tagArtist).SetString(value)
	return nil
}
