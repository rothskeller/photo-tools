package exif

import (
	"bytes"
	"regexp"
	"strings"
)

const (
	tagDateTime            uint16 = 0x132
	tagSubSecTime          uint16 = 0x9290
	tagOffsetTime          uint16 = 0x9010
	tagDateTimeDigitized   uint16 = 0x9004
	tagSubSecTimeDigitized uint16 = 0x9292
	tagOffsetTimeDigitized uint16 = 0x9012
	tagDateTimeOriginal    uint16 = 0x9003
	tagSubSecTimeOriginal  uint16 = 0x9291
	tagOffsetTimeOriginal  uint16 = 0x9011
)

var (
	dateTimeRE         = regexp.MustCompile(`^\d\d\d\d:\d\d:\d\d \d\d:\d\d:\d\d$`)
	offsetTimeRE       = regexp.MustCompile(`^[-+]\d\d:\d\d$`)
	subSecTimeRE       = regexp.MustCompile(`^\d{1,3}$`)
	combinedDateTimeRE = regexp.MustCompile(`^(\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d)(?:\.(\d{1,3})0*)?([-+]\d\d:\d\d)?$`)
)

// DateTime returns the date and time at which the file was modified, or an
// empty string if they are not recorded or corrupt.  The returned value has the
// form "YYYY-MM-DDTHH:MM:SS", possibly followed by ".sss", possibly followed by
// "±HH:MM".
func (p *EXIF) DateTime() string {
	return p.getDateTimeTagGroup(p.exifIFD, tagDateTime, tagOffsetTime, tagSubSecTime, "")
}

// SetDateTime sets the exif {Date,SubSec,Offset}Time tags based on the provided
// datetime string.
func (p *EXIF) SetDateTime(dto string) {
	p.setDateTimeTagGroup(p.exifIFD, tagDateTime, tagOffsetTime, tagSubSecTime, dto)
}

// DateTimeDigitized returns the date and time at which the picture was scanned,
// or an empty string if they are not recorded or corrupt.  The returned value
// has the form "YYYY-MM-DDTHH:MM:SS", possibly followed by ".sss", possibly
// followed by "±HH:MM".
func (p *EXIF) DateTimeDigitized() string {
	return p.getDateTimeTagGroup(p.exifIFD, tagDateTimeDigitized, tagOffsetTimeDigitized, tagSubSecTimeDigitized, "Digitized")
}

// SetDateTimeDigitized sets the exif {Date,SubSec,Offset}TimeDigitized tags
// based on the provided datetime string.
func (p *EXIF) SetDateTimeDigitized(dto string) {
	p.setDateTimeTagGroup(p.exifIFD, tagDateTimeDigitized, tagOffsetTimeDigitized, tagSubSecTimeDigitized, dto)
}

// DateTimeOriginal returns the date and time at which the picture was taken, or
// an empty string if they are not recorded or corrupt.  The returned value has
// the form "YYYY-MM-DDTHH:MM:SS", possibly followed by ".sss", possibly
// followed by "±HH:MM".
func (p *EXIF) DateTimeOriginal() string {
	return p.getDateTimeTagGroup(p.exifIFD, tagDateTimeOriginal, tagOffsetTimeOriginal, tagSubSecTimeOriginal, "Original")
}

// SetDateTimeOriginal sets the exif {Date,SubSec,Offset}TimeOriginal tags based
// on the provided datetime string.
func (p *EXIF) SetDateTimeOriginal(dto string) {
	p.setDateTimeTagGroup(p.exifIFD, tagDateTimeOriginal, tagOffsetTimeOriginal, tagSubSecTimeOriginal, dto)
}

// getDateTimeTagGroup returns the date and time from an exif tag triplet.  The
// returned value has the form "YYYY-MM-DDTHH:MM:SS", possibly followed by
// ".sss", possibly followed by "±HH:MM".
func (p *EXIF) getDateTimeTagGroup(dtifd *ifdt, dttag, ottag, ssttag uint16, suffix string) string {
	if p == nil || dtifd == nil {
		return ""
	}
	dtot := dtifd.findTag(dttag)
	if dtot == nil {
		return ""
	}
	var dto, ssto, oto string
	dto = p.asciiAt(dtot, "DateTime"+suffix)
	if dto == "" || dto == ":  :     :  :" {
		return ""
	}
	if !dateTimeRE.MatchString(dto) {
		p.log(dtot.offset, "invalid DateTime%s value %q", suffix, dto)
		return ""
	}
	dto = strings.Replace(dto, ":", "-", 2)
	dto = strings.Replace(dto, " ", "T", 1)
	if p.exifIFD == nil {
		return dto
	}
	if sstot := p.exifIFD.findTag(ssttag); sstot != nil {
		ssto = p.asciiAt(sstot, "SubSecTime"+suffix)
		if ssto != "" {
			if subSecTimeRE.MatchString(ssto) {
				dto += "." + ssto
			} else {
				p.log(sstot.offset, "invalid SubSecTime%s value %q", suffix, ssto)
			}
		}
	}
	if otot := p.exifIFD.findTag(ottag); otot != nil {
		oto = p.asciiAt(otot, "OffsetTime"+suffix)
		if oto != "" && oto != ":" {
			if offsetTimeRE.MatchString(oto) {
				if oto == "-00:00" || oto == "+00:00" {
					dto += "Z"
				} else {
					dto += oto
				}
			} else {
				p.log(otot.offset, "invalid OffsetTime%s value %q", suffix, oto)
			}
		}
	}
	return dto
}

// setDateTimeTagGroup sets the exif {Date,SubSec,Offset}Time tags based
// on the provided datetime string.
func (p *EXIF) setDateTimeTagGroup(dtifd *ifdt, dttag, ottag, ssttag uint16, dto string) {
	if p == nil {
		return
	}
	if dto == "" {
		if dtifd != nil {
			p.deleteTag(dtifd, dttag)
		}
		if p.exifIFD != nil {
			p.deleteTag(p.exifIFD, ottag)
			p.deleteTag(p.exifIFD, ssttag)
		}
		return
	}
	parts := combinedDateTimeRE.FindStringSubmatch(dto)
	if parts == nil {
		panic("invalid DateTime value")
	}
	if dtifd != nil {
		dto := strings.Replace(parts[1], "-", ":", 2)
		dto = strings.Replace(dto, "T", " ", 1)
		p.setDateTimeTag(dtifd, dttag, dto)
	}
	if p.exifIFD != nil {
		if parts[2] != "" {
			p.setDateTimeTag(p.exifIFD, ssttag, parts[2])
		} else {
			p.deleteTag(p.exifIFD, ssttag)
		}
		if parts[3] == "Z" {
			p.setDateTimeTag(p.exifIFD, ottag, "+00:00")
		} else if parts[3] != "" {
			p.setDateTimeTag(p.exifIFD, ottag, parts[3])
		} else {
			p.deleteTag(p.exifIFD, ottag)
		}
	}
}

func (p *EXIF) setDateTimeTag(ifd *ifdt, tnum uint16, value string) {
	tag := ifd.findTag(tnum)
	if tag == nil {
		tag = &tagt{tag: tnum, ttype: 2, count: 1, data: []byte{0}}
		p.addTag(ifd, tag)
	}
	encbytes := []byte(value + "\000")
	if !bytes.Equal(encbytes, tag.data) {
		tag.data = encbytes
		tag.count = uint32(len(encbytes))
		ifd.dirty = true
	}
}
