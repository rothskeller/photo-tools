package exififd

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
)

const (
	tagSubSecTime          uint16 = 0x9290
	tagOffsetTime          uint16 = 0x9010
	tagDateTimeDigitized   uint16 = 0x9004
	tagSubSecTimeDigitized uint16 = 0x9292
	tagOffsetTimeDigitized uint16 = 0x9012
	tagDateTimeOriginal    uint16 = 0x9003
	tagSubSecTimeOriginal  uint16 = 0x9291
	tagOffsetTimeOriginal  uint16 = 0x9011
)

// getDateTime reads the value of the DateTime field from the RDF.
func (p *Provider) getDateTime() (err error) {
	p.dateTimeDigitized, err = p.getDateTimeTagGroup(
		tagDateTimeDigitized, tagOffsetTimeDigitized, tagSubSecTimeDigitized, "Digitized")
	if err == nil {
		p.dateTimeOriginal, err = p.getDateTimeTagGroup(
			tagDateTimeOriginal, tagOffsetTimeOriginal, tagSubSecTimeOriginal, "Original")
	}
	return err
}
func (p *Provider) getDateTimeTagGroup(dttag, ottag, ssttag uint16, suffix string) (dt metadata.DateTime, err error) {
	dtot := p.ifd.Tag(dttag)
	if dtot == nil {
		return metadata.DateTime{}, nil
	}
	var dto, ssto, oto string
	if dto, err = dtot.AsString(); err != nil {
		return metadata.DateTime{}, fmt.Errorf("DateTime%s: %s", suffix, err)
	}
	if p.ifd != nil {
		if sstot := p.ifd.Tag(ssttag); sstot != nil {
			if ssto, err = sstot.AsString(); err != nil {
				return metadata.DateTime{}, fmt.Errorf("SubSecTime%s: %s", suffix, err)
			}
		}
		if otot := p.ifd.Tag(ottag); otot != nil {
			if oto, err = otot.AsString(); err != nil {
				return metadata.DateTime{}, fmt.Errorf("OffsetTime%s: %s", suffix, err)
			}
		}
	}
	if err := dt.ParseEXIF(dto, ssto, oto); err != nil {
		return metadata.DateTime{}, fmt.Errorf("DateTime%s: %s", suffix, err)
	}
	return dt, nil
}

// DateTime returns the value of the DateTime field.
func (p *Provider) DateTime() (value metadata.DateTime) {
	if !p.dateTimeOriginal.Empty() {
		return p.dateTimeOriginal
	}
	return p.dateTimeDigitized // which may be empty
}

// DateTimeTags returns a list of tag names for the DateTime field, and
// a parallel list of values held by those tags.
func (p *Provider) DateTimeTags() (tags []string, values []metadata.DateTime) {
	tags = append(tags, "EXIF DateTimeOriginal*")
	values = append(values, p.dateTimeOriginal)
	if !p.dateTimeDigitized.Empty() {
		tags = append(tags, "EXIF DateTimeDigitized*")
		values = append(values, p.dateTimeDigitized)
	}
	return tags, values
}

// SetDateTime sets the value of the DateTime field.
func (p *Provider) SetDateTime(value metadata.DateTime) error {
	p.ifd.DeleteTag(tagSubSecTime)
	p.ifd.DeleteTag(tagOffsetTime)
	p.dateTimeDigitized = metadata.DateTime{}
	p.ifd.DeleteTag(tagDateTimeDigitized)
	p.ifd.DeleteTag(tagSubSecTimeDigitized)
	p.ifd.DeleteTag(tagOffsetTimeDigitized)
	if value.Empty() {
		p.dateTimeOriginal = metadata.DateTime{}
		p.ifd.DeleteTag(tagDateTimeOriginal)
		p.ifd.DeleteTag(tagSubSecTimeOriginal)
		p.ifd.DeleteTag(tagOffsetTimeOriginal)
		return nil
	}
	if value.Equivalent(p.dateTimeOriginal) {
		return nil
	}
	p.dateTimeOriginal = value
	dto, ssto, oto := value.AsEXIF()
	p.ifd.AddTag(tagDateTimeOriginal).SetString(dto)
	if ssto != "" {
		p.ifd.AddTag(tagSubSecTimeOriginal).SetString(ssto)
	} else {
		p.ifd.DeleteTag(tagSubSecTimeOriginal)
	}
	if oto != "" {
		p.ifd.AddTag(tagOffsetTimeOriginal).SetString(oto)
	} else {
		p.ifd.DeleteTag(tagOffsetTimeOriginal)
	}
	return nil
}
