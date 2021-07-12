package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	exifDateTimeDigitizedName = rdf.Name{Namespace: nsEXIF, Name: "DateTimeDigitized"}
	exifDateTimeOriginalName  = rdf.Name{Namespace: nsEXIF, Name: "DateTimeOriginal"}
	psDateCreatedName         = rdf.Name{Namespace: nsPS, Name: "DateCreated"}
	tiffDateTimeName          = rdf.Name{Namespace: nsTIFF, Name: "DateTime"}
	xmpCreateDateName         = rdf.Name{Namespace: nsXMP, Name: "CreateDate"}
	xmpMetadataDateName       = rdf.Name{Namespace: nsXMP, Name: "MetadataDate"}
	xmpModifyDateName         = rdf.Name{Namespace: nsXMP, Name: "ModifyDate"}
)

// getDateTime reads the value of the DateTime field from the RDF.
func (p *Provider) getDateTime() (err error) {
	var s string
	if s, err = getString(p.rdf.Property(exifDateTimeDigitizedName)); err == nil {
		err = p.exifDateTimeDigitized.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("exif:DateTimeDigitized: %s", err)
	}
	if s, err = getString(p.rdf.Property(exifDateTimeOriginalName)); err == nil {
		err = p.exifDateTimeOriginal.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("exif:DateTimeOriginal: %s", err)
	}
	if s, err = getString(p.rdf.Property(psDateCreatedName)); err == nil {
		err = p.psDateCreated.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("photoshop:DateCreated: %s", err)
	}
	if s, err = getString(p.rdf.Property(tiffDateTimeName)); err == nil {
		err = p.tiffDateTime.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("tiff:DateTime: %s", err)
	}
	if s, err = getString(p.rdf.Property(xmpCreateDateName)); err == nil {
		err = p.xmpCreateDate.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("xmp:CreateDate: %s", err)
	}
	if s, err = getString(p.rdf.Property(xmpMetadataDateName)); err == nil {
		err = p.xmpMetadataDate.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("xmp:MetadataDate: %s", err)
	}
	if s, err = getString(p.rdf.Property(xmpModifyDateName)); err == nil {
		err = p.xmpModifyDate.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("xmp:ModifyDate: %s", err)
	}
	return nil
}

// DateTime returns the value of the DateTime field.
func (p *Provider) DateTime() (value metadata.DateTime) {
	if !p.xmpCreateDate.Empty() {
		return p.xmpCreateDate
	}
	if !p.psDateCreated.Empty() {
		return p.psDateCreated
	}
	if !p.exifDateTimeOriginal.Empty() {
		return p.exifDateTimeOriginal
	}
	if !p.xmpModifyDate.Empty() {
		return p.xmpModifyDate
	}
	if !p.xmpMetadataDate.Empty() {
		return p.xmpMetadataDate
	}
	if !p.tiffDateTime.Empty() {
		return p.tiffDateTime
	}
	if !p.exifDateTimeDigitized.Empty() {
		return p.exifDateTimeDigitized
	}
	return metadata.DateTime{}
}

// DateTimeTags returns a list of tag names for the DateTime field, and
// a parallel list of values held by those tags.
func (p *Provider) DateTimeTags() (tags []string, values []metadata.DateTime) {
	tags = append(tags, "XMP  xmp:CreateDate", "XMP  photoshop:DateCreated", "XMP  exif:DateTimeOriginal")
	values = append(values, p.xmpCreateDate, p.psDateCreated, p.exifDateTimeOriginal)
	if !p.xmpModifyDate.Empty() {
		tags = append(tags, "XMP  xmp:ModifyDate")
		values = append(values, p.xmpModifyDate)
	}
	if !p.xmpMetadataDate.Empty() {
		tags = append(tags, "XMP  xmp:MetadataDate")
		values = append(values, p.xmpMetadataDate)
	}
	if !p.tiffDateTime.Empty() {
		tags = append(tags, "XMP  tiff:DateTime")
		values = append(values, p.tiffDateTime)
	}
	if !p.exifDateTimeDigitized.Empty() {
		tags = append(tags, "XMP  exif:DateTimeDigitized")
		values = append(values, p.exifDateTimeDigitized)
	}
	return tags, values
}

// SetDateTime sets the value of the DateTime field.
func (p *Provider) SetDateTime(value metadata.DateTime) error {
	p.xmpModifyDate = metadata.DateTime{}
	p.rdf.RemoveProperty(xmpModifyDateName)
	p.xmpMetadataDate = metadata.DateTime{}
	p.rdf.RemoveProperty(xmpMetadataDateName)
	p.tiffDateTime = metadata.DateTime{}
	p.rdf.RemoveProperty(tiffDateTimeName)
	p.exifDateTimeDigitized = metadata.DateTime{}
	p.rdf.RemoveProperty(exifDateTimeDigitizedName)
	if value.Empty() {
		p.xmpCreateDate = metadata.DateTime{}
		p.rdf.RemoveProperty(xmpCreateDateName)
		p.psDateCreated = metadata.DateTime{}
		p.rdf.RemoveProperty(psDateCreatedName)
		p.exifDateTimeOriginal = metadata.DateTime{}
		p.rdf.RemoveProperty(exifDateTimeOriginalName)
		return nil
	}
	if !value.Equivalent(p.xmpCreateDate) {
		p.xmpCreateDate = value
		p.rdf.SetProperty(xmpCreateDateName, makeString(value.String()))
	}
	if !value.Equivalent(p.psDateCreated) {
		p.psDateCreated = value
		p.rdf.SetProperty(psDateCreatedName, makeString(value.String()))
	}
	if !value.Equivalent(p.exifDateTimeOriginal) {
		p.exifDateTimeOriginal = value
		p.rdf.SetProperty(exifDateTimeOriginalName, makeString(value.String()))
	}
	return nil
}
