package exif

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/tifflike"
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
func (p *EXIF) getDateTimeTagGroup(dtifd *tifflike.IFD, dttag, ottag, ssttag uint16, suffix string) (dt metadata.DateTime) {
	var err error

	dtot := dtifd.Tag(dttag)
	if dtot == nil {
		return // empty DateTime
	}
	var dto, ssto, oto string
	if dto, err = dtot.AsString(); err != nil {
		p.log("DateTime%s: %s", suffix, err)
	}
	if p.exifIFD != nil {
		if sstot := p.exifIFD.Tag(ssttag); sstot != nil {
			if ssto, err = sstot.AsString(); err != nil {
				p.log("SubSecTime%s: %s", suffix, err)
			}
		}
		if otot := p.exifIFD.Tag(ottag); otot != nil {
			if oto, err = otot.AsString(); err != nil {
				p.log("OffsetTime%s: %s", suffix, err)
			}
		}
	}
	if err := dt.ParseEXIF(dto, ssto, oto); err != nil {
		p.log("DateTime%s: %s", suffix, err)
	}
	return dt
}

// setDateTimeTagGroup sets the exif {Date,SubSec,Offset}Time tags based
// on the provided datetime string.
func (p *EXIF) setDateTimeTagGroup(dtifd *tifflike.IFD, dttag, ottag, ssttag uint16, dt metadata.DateTime) {
	if dt.Empty() {
		dtifd.DeleteTag(dttag)
		p.exifIFD.DeleteTag(ottag)
		p.exifIFD.DeleteTag(ssttag)
		return
	}
	dts, ssts, ots := dt.AsEXIF()
	dtifd.AddTag(dttag).SetString(dts)
	if p.exifIFD == nil && (ssts != "" || ots != "") {
		p.addEXIFIFD()
	}
	if ssts != "" {
		p.exifIFD.AddTag(ssttag).SetString(ssts)
	} else {
		p.exifIFD.DeleteTag(ssttag)
	}
	if ots != "" {
		p.exifIFD.AddTag(ottag).SetString(ots)
	} else {
		p.exifIFD.DeleteTag(ottag)
	}
}
