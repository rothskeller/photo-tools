package exif

import (
	"bytes"
	"encoding/binary"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

const tagUserComment uint16 = 0x9286

var charsetASCII = []byte("ASCII\000\000\000")
var charsetUnicode = []byte("UNICODE\000")
var charsetUnknown = []byte("\000\000\000\000\000\000\000\000")

// UserComment returns the UserComment tag, if provided.
func (p *EXIF) UserComment() string {
	if p == nil || p.exifIFD == nil {
		return ""
	}
	idt := p.exifIFD.findTag(tagUserComment)
	if idt == nil {
		return ""
	}
	if idt.ttype != 7 || idt.count < 8 {
		p.log(idt.doff, "UserComment is ill-formed")
		return ""
	}
	switch {
	case bytes.Equal(idt.data[:8], charsetASCII):
		return string(idt.data[8:])
	case bytes.Equal(idt.data[:8], charsetUnicode):
		// By the spec, this should be UCS-2.  In practice it may
		// actually be UTF-16.  Reading it as UTF-16 handles both cases.
		var enc encoding.Encoding
		if p.enc == binary.BigEndian {
			enc = unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
		} else {
			enc = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		}
		u8, err := enc.NewDecoder().String(string(idt.data[8:]))
		if err != nil || strings.ContainsRune(u8, utf8.RuneError) {
			p.log(idt.doff, "UserComment is invalid UTF-16, ignoring")
			return ""
		}
		return u8
	case bytes.Equal(idt.data[:8], charsetUnknown):
		// There's a decent chance this is actually UTF-8.  Let's try it.
		var s = string(idt.data[8:])
		if !utf8.ValidString(s) {
			p.log(idt.doff, "UserComment is in unknown character set, ignoring")
			return ""
		}
		return s
	default:
		p.log(idt.doff, "UserComment is in unknown character set, ignoring")
		return ""
	}
}

// SetUserComment sets the UserComment tag.
func (p *EXIF) SetUserComment(comment string) {
	if p == nil || p.exifIFD == nil {
		return
	}
	if comment == "" {
		p.deleteTag(p.exifIFD, tagUserComment)
		return
	}
	tag := p.exifIFD.findTag(tagUserComment)
	if tag == nil {
		tag = &tagt{tag: tagUserComment, ttype: 7, count: 1, data: []byte{0}}
		p.addTag(p.exifIFD, tag)
	}
	var encbytes []byte
	if strings.IndexFunc(comment, func(r rune) bool {
		return r >= utf8.RuneSelf
	}) < 0 {
		encbytes = make([]byte, len(comment)+8)
		copy(encbytes, charsetASCII)
		copy(encbytes[8:], comment)
	} else {
		var enc encoding.Encoding
		if p.enc == binary.BigEndian {
			enc = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		} else {
			enc = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		}
		u16, err := enc.NewEncoder().String(comment)
		if err != nil {
			panic("can't encode comment into UTF-16?")
		}
		encbytes = make([]byte, len(u16)+8)
		copy(encbytes, charsetUnicode)
		copy(encbytes[8:], u16)
	}
	if !bytes.Equal(encbytes, tag.data) {
		tag.data = encbytes
		tag.count = uint32(len(encbytes))
		p.exifIFD.dirty = true
	}
}
