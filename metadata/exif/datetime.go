package exif

import (
	"github.com/rothskeller/photo-tools/metadata"
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

// DateTime returns the value of the DateTime tag.
func (p *EXIF) DateTime() metadata.DateTime { return p.dateTime }

func (p *EXIF) getDateTime() {
	p.dateTime = p.getDateTimeTagGroup(p.ifd0, tagDateTime, tagOffsetTime, tagSubSecTime, "")
}

// SetDateTime sets the value of the DateTime tag.
func (p *EXIF) SetDateTime(v metadata.DateTime) error {
	if v.Equal(p.dateTime) {
		return nil
	}
	p.dateTime = v
	p.setDateTimeTagGroup(p.ifd0, tagDateTime, tagOffsetTime, tagSubSecTime, p.dateTime)
	return nil
}

// DateTimeDigitized returns the value of the DateTimeDigitized tag.
func (p *EXIF) DateTimeDigitized() metadata.DateTime { return p.dateTimeDigitized }

func (p *EXIF) getDateTimeDigitized() {
	p.dateTimeDigitized = p.getDateTimeTagGroup(
		p.exifIFD, tagDateTimeDigitized, tagOffsetTimeDigitized, tagSubSecTimeDigitized, "Digitized")
}

// SetDateTimeDigitized sets the value of the DateTimeDigitized tag.
func (p *EXIF) SetDateTimeDigitized(v metadata.DateTime) error {
	if v.Equal(p.dateTimeDigitized) {
		return nil
	}
	p.dateTimeDigitized = v
	if !p.dateTimeDigitized.Empty() && p.exifIFD == nil {
		p.addEXIFIFD()
	}
	p.setDateTimeTagGroup(p.exifIFD, tagDateTimeDigitized, tagOffsetTimeDigitized, tagSubSecTimeDigitized, p.dateTimeDigitized)
	return nil
}

// DateTimeOriginal returns the value of the DateTimeOriginal tag.
func (p *EXIF) DateTimeOriginal() metadata.DateTime { return p.dateTimeOriginal }

func (p *EXIF) getDateTimeOriginal() {
	p.dateTimeOriginal = p.getDateTimeTagGroup(
		p.exifIFD, tagDateTimeOriginal, tagOffsetTimeOriginal, tagSubSecTimeOriginal, "Original")
}

// SetDateTimeOriginal sets the value of the DateTimeOriginal tag.
func (p *EXIF) SetDateTimeOriginal(v metadata.DateTime) error {
	if v.Equal(p.dateTimeOriginal) {
		return nil
	}
	p.dateTimeOriginal = v
	if !p.dateTimeOriginal.Empty() && p.exifIFD == nil {
		p.addEXIFIFD()
	}
	p.setDateTimeTagGroup(p.exifIFD, tagDateTimeOriginal, tagOffsetTimeOriginal, tagSubSecTimeOriginal, p.dateTimeOriginal)
	return nil
}

// getDateTimeTagGroup returns the date and time from an exif tag triplet.
func (p *EXIF) getDateTimeTagGroup(dtifd *ifdt, dttag, ottag, ssttag uint16, suffix string) (dt metadata.DateTime) {
	dtot := dtifd.findTag(dttag)
	if dtot == nil {
		return // empty DateTime
	}
	var dto, ssto, oto string
	dto = p.asciiAt(dtot, "DateTime"+suffix)
	if p.exifIFD != nil {
		if sstot := p.exifIFD.findTag(ssttag); sstot != nil {
			ssto = p.asciiAt(sstot, "SubSecTime"+suffix)
		}
		if otot := p.exifIFD.findTag(ottag); otot != nil {
			oto = p.asciiAt(otot, "OffsetTime"+suffix)
		}
	}
	if err := dt.ParseEXIF(dto, ssto, oto); err != nil {
		p.log(dtot.offset, "DateTime%s: %s", suffix, err)
	}
	return dt
}

// setDateTimeTagGroup sets the exif {Date,SubSec,Offset}Time tags based
// on the provided datetime string.
func (p *EXIF) setDateTimeTagGroup(dtifd *ifdt, dttag, ottag, ssttag uint16, dt metadata.DateTime) {
	if dt.Empty() {
		p.deleteTag(dtifd, dttag)
		p.deleteTag(p.exifIFD, ottag)
		p.deleteTag(p.exifIFD, ssttag)
		return
	}
	dts, ssts, ots := dt.AsEXIF()
	p.setASCIITag(dtifd, dttag, dts)
	if p.exifIFD == nil && (ssts != "" || ots != "") {
		p.addEXIFIFD()
	}
	if ssts != "" {
		p.setASCIITag(p.exifIFD, ssttag, ssts)
	} else {
		p.deleteTag(p.exifIFD, ssttag)
	}
	if ots != "" {
		p.setASCIITag(p.exifIFD, ottag, ots)
	} else {
		p.deleteTag(p.exifIFD, ottag)
	}
}
