package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
)

const nsTIFF = "http://ns.adobe.com/tiff/1.0/"
const pfxTIFF = "tiff"

// TIFFArtist returns the value of the tiff:Artist tag.
func (p *XMP) TIFFArtist() string { return p.tiffArtist }

// TIFFDateTime returns the value of the tiff:DateTime tag.
func (p *XMP) TIFFDateTime() metadata.DateTime { return p.tiffDateTime }

// TIFFImageDescription returns the value of the tiff:ImageDescription tag.
func (p *XMP) TIFFImageDescription() metadata.AltString { return p.tiffImageDescription }

func (p *XMP) getTIFF() {
	p.tiffArtist = p.getString(p.rdf.Properties, pfxTIFF, nsTIFF, "Artist")
	p.xmpDateTimeToMetadata(p.getString(p.rdf.Properties, pfxTIFF, nsTIFF, "DateTime"), &p.tiffDateTime)
	p.tiffImageDescription = p.getAlt(p.rdf.Properties, pfxTIFF, nsTIFF, "ImageDescription")
	p.rdf.RegisterNamespace(pfxTIFF, nsTIFF)
}

// SetTIFFArtist sets the value of the tiff:Artist tag.
func (p *XMP) SetTIFFArtist(v string) (err error) {
	if v == p.tiffArtist {
		return nil
	}
	p.tiffArtist = v
	p.setString(p.rdf.Properties, nsTIFF, "Artist", v)
	return nil
}

// SetTIFFDateTime sets the value of the tiff:DateTime tag.
func (p *XMP) SetTIFFDateTime(v metadata.DateTime) (err error) {
	if v.Equivalent(p.tiffDateTime) {
		return nil
	}
	p.tiffDateTime = v
	p.setString(p.rdf.Properties, nsTIFF, "DateTime", v.String())
	return nil
}

// SetTIFFImageDescription sets the value of the tiff:ImageDescription tag.
func (p *XMP) SetTIFFImageDescription(v metadata.AltString) (err error) {
	if v.Equal(p.tiffImageDescription) {
		return nil
	}
	p.tiffImageDescription = v
	p.setAlt(p.rdf.Properties, nsTIFF, "ImageDescription", v)
	return nil
}
