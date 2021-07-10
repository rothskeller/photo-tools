package xmpexif

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	dateTimeDigitizedName = rdf.Name{Namespace: nsEXIF, Name: "DateTimeDigitized"}
	dateTimeOriginalName  = rdf.Name{Namespace: nsEXIF, Name: "DateTimeOriginal"}
)

// getDateTime reads the value of the DateTime field from the RDF.
func (p *Provider) getDateTime() (err error) {
	var s string
	if s, err = getString(p.rdf.Properties, dateTimeDigitizedName); err == nil {
		err = p.exifDateTimeDigitized.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("exif:DateTimeDigitized: %s", err)
	}
	if s, err = getString(p.rdf.Properties, dateTimeOriginalName); err == nil {
		err = p.exifDateTimeOriginal.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("exif:DateTimeOriginal: %s", err)
	}
	return nil
}

// DateTime returns the value of the DateTime field.
func (p *Provider) DateTime() (value metadata.DateTime) {
	if !p.exifDateTimeOriginal.Empty() {
		return p.exifDateTimeOriginal
	}
	return p.exifDateTimeDigitized // which may be empty
}

// DateTimeTags returns a list of tag names for the DateTime field, and
// a parallel list of values held by those tags.
func (p *Provider) DateTimeTags() (tags []string, values []metadata.DateTime) {
	tags = append(tags, "XMP  exif:DateTimeOriginal")
	values = append(values, p.exifDateTimeOriginal)
	if !p.exifDateTimeDigitized.Empty() {
		tags = append(tags, "XMP  exif:DateTimeDigitized")
		values = append(values, p.exifDateTimeDigitized)
	}
	return tags, values
}

// SetDateTime sets the value of the DateTime field.
func (p *Provider) SetDateTime(value metadata.DateTime) error {
	p.exifDateTimeDigitized = metadata.DateTime{}
	if _, ok := p.rdf.Properties[dateTimeDigitizedName]; ok {
		delete(p.rdf.Properties, dateTimeDigitizedName)
		p.dirty = true
	}
	if value.Empty() {
		p.exifDateTimeOriginal = metadata.DateTime{}
		if _, ok := p.rdf.Properties[dateTimeOriginalName]; ok {
			delete(p.rdf.Properties, dateTimeOriginalName)
			p.dirty = true
		}
		return nil
	}
	if value.Equivalent(p.exifDateTimeOriginal) {
		return nil
	}
	p.exifDateTimeOriginal = value
	setString(p.rdf.Properties, dateTimeOriginalName, value.String())
	p.dirty = true
	return nil
}
