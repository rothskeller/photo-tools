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

func (p *EXIF) getDateTime() {
	p.DateTime = p.getDateTimeTagGroup(p.ifd0, tagDateTime, tagOffsetTime, tagSubSecTime, "")
}

func (p *EXIF) setDateTime() {
	p.setDateTimeTagGroup(p.ifd0, tagDateTime, tagOffsetTime, tagSubSecTime, p.DateTime)
}

func (p *EXIF) getDateTimeDigitized() {
	p.DateTimeDigitized = p.getDateTimeTagGroup(
		p.exifIFD, tagDateTimeDigitized, tagOffsetTimeDigitized, tagSubSecTimeDigitized, "Digitized")
}

func (p *EXIF) setDateTimeDigitized() {
	if !p.DateTimeDigitized.Empty() && p.exifIFD == nil {
		p.addEXIFIFD()
	}
	p.setDateTimeTagGroup(p.exifIFD, tagDateTimeDigitized, tagOffsetTimeDigitized, tagSubSecTimeDigitized, p.DateTimeDigitized)
}

func (p *EXIF) getDateTimeOriginal() {
	p.DateTimeOriginal = p.getDateTimeTagGroup(
		p.exifIFD, tagDateTimeOriginal, tagOffsetTimeOriginal, tagSubSecTimeOriginal, "Original")
}

func (p *EXIF) setDateTimeOriginal() {
	if !p.DateTimeOriginal.Empty() && p.exifIFD == nil {
		p.addEXIFIFD()
	}
	p.setDateTimeTagGroup(p.exifIFD, tagDateTimeOriginal, tagOffsetTimeOriginal, tagSubSecTimeOriginal, p.DateTimeOriginal)
}

// getDateTimeTagGroup returns the date and time from an exif tag triplet.  The
// returned value has the form "YYYY-MM-DDTHH:MM:SS", possibly followed by
// ".sss", possibly followed by "±HH:MM".
func (p *EXIF) getDateTimeTagGroup(dtifd *ifdt, dttag, ottag, ssttag uint16, suffix string) (dt *metadata.DateTime) {
	dtot := dtifd.findTag(dttag)
	if dtot == nil {
		return nil
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
	dt = new(metadata.DateTime)
	if err := dt.ParseEXIF(dto, ssto, oto); err != nil {
		p.log(dtot.offset, "DateTime%s: %s", suffix, err)
	}
	return dt
}

// setDateTimeTagGroup sets the exif {Date,SubSec,Offset}Time tags based
// on the provided datetime string.
func (p *EXIF) setDateTimeTagGroup(dtifd *ifdt, dttag, ottag, ssttag uint16, dt *metadata.DateTime) {
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
