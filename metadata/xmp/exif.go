package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
)

const nsEXIF = "http://ns.adobe.com/exif/1.0/"
const pfxEXIF = "exif"

// EXIFDateTimeDigitized returns the value of the exif:DateTimeDigitized tag.
func (p *XMP) EXIFDateTimeDigitized() metadata.DateTime { return p.exifDateTimeDigitized }

// EXIFDateTimeOriginal returns the value of the exif:DateTimeOriginal tag.
func (p *XMP) EXIFDateTimeOriginal() metadata.DateTime { return p.exifDateTimeOriginal }

// EXIFGPSCoords returns the values of the exif:GPS* tags.
func (p *XMP) EXIFGPSCoords() metadata.GPSCoords { return p.exifGPSCoords }

// EXIFUserComments returns the values of the exif:UserComment tag.
func (p *XMP) EXIFUserComments() []string { return p.exifUserComments }

func (p *XMP) getEXIF() {
	p.xmpDateTimeToMetadata(p.getString(p.rdf.Properties, pfxEXIF, nsEXIF, "DateTimeDigitized"), &p.exifDateTimeDigitized)
	p.xmpDateTimeToMetadata(p.getString(p.rdf.Properties, pfxEXIF, nsEXIF, "DateTimeOriginal"), &p.exifDateTimeOriginal)
	if err := p.exifGPSCoords.ParseXMP(
		p.getString(p.rdf.Properties, pfxEXIF, nsEXIF, "GPSLatitude"),
		p.getString(p.rdf.Properties, pfxEXIF, nsEXIF, "GPSLongitude"),
		p.getString(p.rdf.Properties, pfxEXIF, nsEXIF, "GPSAltitudeRef"),
		p.getString(p.rdf.Properties, pfxEXIF, nsEXIF, "GPSAltitude"),
	); err != nil {
		p.log("invalid GPS coordinates")
	}
	p.exifUserComments = p.getStrings(p.rdf.Properties, pfxEXIF, nsEXIF, "UserComment")
	p.rdf.RegisterNamespace(pfxEXIF, nsEXIF)
}

// SetEXIFDateTimeDigitized sets the value of the exif:DateTimeDigitized tag.
func (p *XMP) SetEXIFDateTimeDigitized(v metadata.DateTime) (err error) {
	if v.Equivalent(p.exifDateTimeDigitized) {
		return nil
	}
	p.exifDateTimeDigitized = v
	p.setString(p.rdf.Properties, nsEXIF, "DateTimeDigitized", v.String())
	return nil
}

// SetEXIFDateTimeOriginal sets the value of the exif:DateTimeOriginal tag.
func (p *XMP) SetEXIFDateTimeOriginal(v metadata.DateTime) (err error) {
	if v.Equivalent(p.exifDateTimeOriginal) {
		return nil
	}
	p.exifDateTimeOriginal = v
	p.setString(p.rdf.Properties, nsEXIF, "DateTimeOriginal", v.String())
	return nil
}

// SetEXIFGPSCoords sets the values of the exif:GPS* tags.
func (p *XMP) SetEXIFGPSCoords(v metadata.GPSCoords) (err error) {
	if v.Equivalent(p.exifGPSCoords) {
		return nil
	}
	p.exifGPSCoords = v
	lat, long, altref, alt := v.AsXMP()
	p.setString(p.rdf.Properties, nsEXIF, "GPSLatitude", lat)
	p.setString(p.rdf.Properties, nsEXIF, "GPSLongitude", long)
	p.setString(p.rdf.Properties, nsEXIF, "GPSAltitudeRef", altref)
	p.setString(p.rdf.Properties, nsEXIF, "GPSAltitude", alt)
	return nil
}

// SetEXIFUserComments sets the values of the exif:UserComment tag.
func (p *XMP) SetEXIFUserComments(v []string) (err error) {
	if stringSliceEqual(v, p.exifUserComments) {
		return nil
	}
	p.exifUserComments = v
	p.setBag(p.rdf.Properties, nsEXIF, "UserComment", v)
	return nil
}
