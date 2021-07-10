package xmpexif

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	gpsAltitudeName    = rdf.Name{Namespace: nsEXIF, Name: "GPSAltitude"}
	gpsAltitudeRefName = rdf.Name{Namespace: nsEXIF, Name: "GPSAltitudeRef"}
	gpsLatitudeName    = rdf.Name{Namespace: nsEXIF, Name: "GPSLatitude"}
	gpsLongitudeName   = rdf.Name{Namespace: nsEXIF, Name: "GPSLongitude"}
)

// getGPS reads the value of the GPS field from the RDF.
func (p *Provider) getGPS() (err error) {
	var lat, long, altref, alt string
	if lat, err = getString(p.rdf.Properties, gpsLatitudeName); err != nil {
		return fmt.Errorf("exif:GPSLatitude: %s", err)
	}
	if long, err = getString(p.rdf.Properties, gpsLongitudeName); err != nil {
		return fmt.Errorf("exif:GPSLongitude: %s", err)
	}
	if altref, err = getString(p.rdf.Properties, gpsAltitudeRefName); err != nil {
		return fmt.Errorf("exif:GPSAltitudeRef: %s", err)
	}
	if alt, err = getString(p.rdf.Properties, gpsAltitudeName); err != nil {
		return fmt.Errorf("exif:GPSAltitude: %s", err)
	}
	if err = p.exifGPSCoords.ParseXMP(lat, long, altref, alt); err != nil {
		return fmt.Errorf("exif:GPS*: %s", err)
	}
	return nil
}

// GPS returns the value of the GPS field.
func (p *Provider) GPS() (value metadata.GPSCoords) { return p.exifGPSCoords }

// GPSTags returns a list of tag names for the GPS field, and a parallel list of
// values held by those tags.
func (p *Provider) GPSTags() (tags []string, values []metadata.GPSCoords) {
	return []string{"XMP  exif:GPS*"}, []metadata.GPSCoords{p.exifGPSCoords}
}

// SetGPS sets the value of the GPS field.
func (p *Provider) SetGPS(value metadata.GPSCoords) error {
	if value.Empty() {
		p.exifGPSCoords = metadata.GPSCoords{}
		if _, ok := p.rdf.Properties[gpsLatitudeName]; ok {
			delete(p.rdf.Properties, gpsLatitudeName)
			p.dirty = true
		}
		if _, ok := p.rdf.Properties[gpsLongitudeName]; ok {
			delete(p.rdf.Properties, gpsLongitudeName)
			p.dirty = true
		}
		if _, ok := p.rdf.Properties[gpsAltitudeRefName]; ok {
			delete(p.rdf.Properties, gpsAltitudeRefName)
			p.dirty = true
		}
		if _, ok := p.rdf.Properties[gpsAltitudeName]; ok {
			delete(p.rdf.Properties, gpsAltitudeName)
			p.dirty = true
		}
		return nil
	}
	if value.Equivalent(p.exifGPSCoords) {
		return nil
	}
	p.exifGPSCoords = value
	lat, long, altref, alt := value.AsXMP()
	setString(p.rdf.Properties, gpsLatitudeName, lat)
	setString(p.rdf.Properties, gpsLongitudeName, long)
	setString(p.rdf.Properties, gpsAltitudeRefName, altref)
	setString(p.rdf.Properties, gpsAltitudeName, alt)
	p.dirty = true
	return nil
}
