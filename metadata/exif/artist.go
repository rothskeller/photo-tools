package exif

import (
	"strings"
)

const tagArtist uint16 = 0x13B

func (p *EXIF) getArtist() {
	tag := p.ifd0.findTag(tagArtist)
	if tag == nil {
		return
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
				p.Artist = append(p.Artist, t)
			}
		default:
			abuf += string(b)
		}
	}
	if t := strings.TrimSpace(abuf); t != "" {
		p.Artist = append(p.Artist, t)
	}
}

func (p *EXIF) setArtist() {
	if len(p.Artist) == 0 {
		p.deleteTag(p.ifd0, tagArtist)
		return
	}
	var encoded = make([]string, 0, len(p.Artist))
	for _, a := range p.Artist {
		if strings.IndexAny(a, `";`) >= 0 {
			a = `"` + strings.ReplaceAll(a, `"`, `""`) + `"`
		}
		encoded = append(encoded, a)
	}
	p.setASCIITag(p.ifd0, tagArtist, strings.Join(encoded, "; "))
}
