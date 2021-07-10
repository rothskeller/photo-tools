package exififd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

const tagUserComment uint16 = 0x9286

var charsetASCII = []byte("ASCII\000\000\000")
var charsetUnicode = []byte("UNICODE\000")
var charsetUnknown = []byte("\000\000\000\000\000\000\000\000")

// getCaption reads the value of the Caption field from the RDF.
func (p *Provider) getCaption() (err error) {
	tag := p.ifd.Tag(tagUserComment)
	if tag == nil {
		return nil
	}
	data, err := tag.AsUnknown()
	if err != nil {
		return fmt.Errorf("UserComment: %s", err)
	}
	if len(data) < 8 {
		return errors.New("UserComment: wrong length")
	}
	switch {
	case bytes.Equal(data[:8], charsetASCII):
		// I've found that comments are often padded with nulls (or
		// even composed entirely of them).  We don't want those.
		p.userComment = strings.TrimRight(string(data[8:]), "\000")
	case bytes.Equal(data[:8], charsetUnicode):
		// By the spec, this should be UCS-2.  In practice it may
		// actually be UTF-16.  Reading it as UTF-16 handles both cases.
		var enc encoding.Encoding
		if p.enc == binary.BigEndian {
			enc = unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
		} else {
			enc = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		}
		u8, err := enc.NewDecoder().String(string(data[8:]))
		if err != nil || strings.ContainsRune(u8, utf8.RuneError) {
			return errors.New("UserComment: invalid UTF-16")
		}
		p.userComment = u8
	case bytes.Equal(data[:8], charsetUnknown):
		// There's a decent chance this is actually UTF-8.  Let's try it.
		var s = string(data[8:])
		if !utf8.ValidString(s) {
			return errors.New("UserComment: unknown character set")
		}
		p.userComment = s
	default:
		return errors.New("UserComment: unknown character set")
	}
	return nil
}

// Caption returns the value of the Caption field.
func (p *Provider) Caption() (value string) { return p.userComment }

// CaptionTags returns a list of tag names for the Caption field, and a
// parallel list of values held by those tags.
func (p *Provider) CaptionTags() (tags []string, values []string) {
	if p.userComment == "" {
		return nil, nil
	}
	return []string{"EXIF UserComment"}, []string{p.userComment}
}

// SetCaption sets the value of the Caption field.
func (p *Provider) SetCaption(value string) error {
	p.userComment = ""
	p.ifd.DeleteTag(tagUserComment)
	return nil
}
