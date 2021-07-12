package xmp

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
	if lat, err = getString(p.rdf.Property(gpsLatitudeName)); err != nil {
		return fmt.Errorf("exif:GPSLatitude: %s", err)
	}
	if long, err = getString(p.rdf.Property(gpsLongitudeName)); err != nil {
		return fmt.Errorf("exif:GPSLongitude: %s", err)
	}
	if altref, err = getString(p.rdf.Property(gpsAltitudeRefName)); err != nil {
		return fmt.Errorf("exif:GPSAltitudeRef: %s", err)
	}
	if alt, err = getString(p.rdf.Property(gpsAltitudeName)); err != nil {
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
		p.rdf.RemoveProperty(gpsLatitudeName)
		p.rdf.RemoveProperty(gpsLongitudeName)
		p.rdf.RemoveProperty(gpsAltitudeRefName)
		p.rdf.RemoveProperty(gpsAltitudeName)
		return nil
	}
	if value.Equivalent(p.exifGPSCoords) {
		return nil
	}
	p.exifGPSCoords = value
	lat, long, altref, alt := value.AsXMP()
	p.rdf.SetProperty(gpsLatitudeName, makeString(lat))
	p.rdf.SetProperty(gpsLongitudeName, makeString(long))
	p.rdf.SetProperty(gpsAltitudeRefName, makeString(altref))
	p.rdf.SetProperty(gpsAltitudeName, makeString(alt))
	return nil
}
