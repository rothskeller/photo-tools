package exif

import (
	"bytes"
	"strings"
)

const tagArtist uint16 = 0x13B

// Artist returns the list of people in the Artist tag, if any.
func (p *EXIF) Artist() (artists []string) {
	if p == nil || p.ifd0 == nil {
		return nil
	}
	tag := p.ifd0.findTag(tagArtist)
	if tag == nil {
		return nil
	}
	alist := p.asciiAt(tag, "Artist")
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
				artists = append(artists, t)
			}
		default:
			abuf += string(b)
		}
	}
	if t := strings.TrimSpace(abuf); t != "" {
		artists = append(artists, t)
	}
	return artists
}

// SetArtist sets the EXIF Artist tag.
func (p *EXIF) SetArtist(artists []string) {
	if p == nil || p.ifd0 == nil {
		// We're not going to add an EXIF block just for this.
		return
	}
	if len(artists) == 0 {
		p.deleteTag(p.ifd0, tagArtist)
		return
	}
	tag := p.ifd0.findTag(tagArtist)
	if tag == nil {
		tag = &tagt{tag: tagArtist, ttype: 2, count: 1, data: []byte{0}}
		p.addTag(p.ifd0, tag)
	}
	var encoded = make([]string, 0, len(artists))
	for _, a := range artists {
		if strings.IndexAny(a, `";`) >= 0 {
			a = `"` + strings.ReplaceAll(a, `"`, `""`) + `"`
		}
		encoded = append(encoded, a)
	}
	encbytes := []byte(strings.Join(encoded, "; ") + "\000")
	if !bytes.Equal(encbytes, tag.data) {
		tag.data = encbytes
		tag.count = uint32(len(encbytes))
		p.ifd0.dirty = true
	}
}
