package exif

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

const tagUserComment uint16 = 0x9286

var charsetASCII = []byte("ASCII\000\000\000")
var charsetUnicode = []byte("UNICODE\000")
var charsetUnknown = []byte("\000\000\000\000\000\000\000\000")

// UserComment returns the value of the UserComment tag.
func (p *EXIF) UserComment() string { return p.userComment }

func (p *EXIF) getUserComment() {
	idt := p.exifIFD.findTag(tagUserComment)
	if idt == nil {
		return
	}
	if idt.ttype != 7 || idt.count < 8 {
		p.log(idt.doff, "UserComment is ill-formed")
		return
	}
	switch {
	case bytes.Equal(idt.data[:8], charsetASCII):
		// I've found that comments are often padded with nulls (or
		// even composed entirely of them).  We don't want those.
		p.userComment = strings.TrimRight(string(idt.data[8:]), "\000")
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
			return
		}
		p.userComment = u8
	case bytes.Equal(idt.data[:8], charsetUnknown):
		// There's a decent chance this is actually UTF-8.  Let's try it.
		var s = string(idt.data[8:])
		if !utf8.ValidString(s) {
			p.log(idt.doff, "UserComment is in unknown character set, ignoring")
			return
		}
		p.userComment = s
	default:
		p.log(idt.doff, "UserComment is in unknown character set, ignoring")
	}
}

// SetUserComment sets the value of the UserComment tag.
func (p *EXIF) SetUserComment(v string) error {
	if v == p.userComment {
		return nil
	}
	p.userComment = v
	if p.exifIFD == nil && p.userComment != "" {
		p.addEXIFIFD()
	}
	if p.userComment == "" {
		p.deleteTag(p.exifIFD, tagUserComment)
		return nil
	}
	tag := p.exifIFD.findTag(tagUserComment)
	if tag == nil {
		tag = &tagt{tag: tagUserComment, ttype: 7, count: 1, data: []byte{0}}
		p.addTag(p.exifIFD, tag)
	}
	var encbytes []byte
	if s := p.userComment; strings.IndexFunc(s, func(r rune) bool {
		return r >= utf8.RuneSelf
	}) < 0 {
		encbytes = make([]byte, len(s)+8)
		copy(encbytes, charsetASCII)
		copy(encbytes[8:], []byte(s))
	} else {
		var enc encoding.Encoding
		if p.enc == binary.BigEndian {
			enc = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		} else {
			enc = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		}
		u16, err := enc.NewEncoder().String(s)
		if err != nil {
			return errors.New("can't encode comment into UTF-16 for EXIF")
		}
		encbytes = make([]byte, len(u16)+8)
		copy(encbytes, charsetUnicode)
		copy(encbytes[8:], u16)
	}
	tag.data = encbytes
	tag.count = uint32(len(encbytes))
	p.exifIFD.dirty = true
	return nil
}
