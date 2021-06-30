package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/exif"
)

// EXIFDateTimeDigitized returns the value of the exif:DateTimeDigitized tag.
func (p *XMP) EXIFDateTimeDigitized() metadata.DateTime { return p.exifDateTimeDigitized }

// EXIFDateTimeOriginal returns the value of the exif:DateTimeOriginal tag.
func (p *XMP) EXIFDateTimeOriginal() metadata.DateTime { return p.exifDateTimeOriginal }

// EXIFGPSCoords returns the values of the exif:GPS* tags.
func (p *XMP) EXIFGPSCoords() metadata.GPSCoords { return p.exifGPSCoords }

// EXIFUserComments returns the values of the exif:UserComment tag.
func (p *XMP) EXIFUserComments() []string { return p.exifUserComments }

func (p *XMP) getEXIF() {
	var model *exif.ExifInfo

	if p != nil && p.doc != nil {
		model = exif.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.xmpDateTimeToMetadata(model.DateTimeDigitized, &p.exifDateTimeDigitized)
	p.xmpDateTimeToMetadata(model.DateTimeOriginal, &p.exifDateTimeOriginal)
	if err := p.exifGPSCoords.ParseXMP(model.GPSLatitude, model.GPSLongitude, model.GPSAltitudeRef, model.GPSAltitude); err != nil {
		p.log("invalid GPS coordinates")
	}
	p.exifUserComments = model.UserComment
}

// SetEXIFDateTimeDigitized sets the value of the exif:DateTimeDigitized tag.
func (p *XMP) SetEXIFDateTimeDigitized(v metadata.DateTime) (err error) {
	var model *exif.ExifInfo

	if model, err = exif.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add exif model to XMP: %s", err)
	}
	if v.Equivalent(p.exifDateTimeDigitized) {
		return nil
	}
	p.exifDateTimeDigitized = v
	model.DateTimeDigitized = v.String()
	p.dirty = true
	return nil
}

// SetEXIFDateTimeOriginal sets the value of the exif:DateTimeOriginal tag.
func (p *XMP) SetEXIFDateTimeOriginal(v metadata.DateTime) (err error) {
	var model *exif.ExifInfo

	if model, err = exif.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add exif model to XMP: %s", err)
	}
	if v.Equivalent(p.exifDateTimeOriginal) {
		return nil
	}
	p.exifDateTimeOriginal = v
	model.DateTimeOriginal = v.String()
	p.dirty = true
	return nil
}

// SetEXIFGPSCoords sets the values of the exif:GPS* tags.
func (p *XMP) SetEXIFGPSCoords(v metadata.GPSCoords) (err error) {
	var model *exif.ExifInfo

	if model, err = exif.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add exif model to XMP: %s", err)
	}
	if v.Equivalent(p.exifGPSCoords) {
		return nil
	}
	p.exifGPSCoords = v
	model.GPSLatitude, model.GPSLongitude, model.GPSAltitudeRef, model.GPSAltitude = v.AsXMP()
	p.dirty = true
	return nil
}

// SetEXIFUserComments sets the values of the exif:UserComment tag.
func (p *XMP) SetEXIFUserComments(v []string) (err error) {
	var model *exif.ExifInfo

	if model, err = exif.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add exif model to XMP: %s", err)
	}
	if stringSliceEqual(v, p.exifUserComments) {
		return nil
	}
	p.exifUserComments = v
	model.UserComment = v
	p.dirty = true
	return nil
}
