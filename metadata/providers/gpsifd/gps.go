package gpsifd

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata"
)

// getGPS reads the value of the GPS field from the RDF.
func (p *Provider) getGPS() (err error) {
	latreft := p.ifd.Tag(1)
	latratt := p.ifd.Tag(2)
	longreft := p.ifd.Tag(3)
	longratt := p.ifd.Tag(4)
	if latreft == nil && latratt == nil && longreft == nil && longratt == nil {
		return nil
	}
	if latreft == nil || latratt == nil || longreft == nil || longratt == nil {
		return errors.New("invalid GPS tags")
	}
	var latref, longref string
	var latrat, longrat []uint32
	latref, err = latreft.AsString()
	if err == nil {
		longref, err = longreft.AsString()
	}
	if err == nil {
		latrat, err = latratt.AsRationals()
	}
	if err == nil {
		longrat, err = longratt.AsRationals()
	}
	if err != nil || len(latref) != 1 || len(longref) != 1 || len(latrat) != 6 || len(longrat) != 6 {
		return errors.New("invalid GPS tags")
	}

	var altrefs []byte
	var altref byte
	var altrat []uint32
	altreft := p.ifd.Tag(5)
	altratt := p.ifd.Tag(6)
	if altratt != nil {
		if altreft != nil {
			altrefs, err = altreft.AsBytes()
		}
		if err == nil {
			altrat, err = altratt.AsRationals()
		}
		if err != nil || len(altrefs) > 1 || len(altrat) != 2 {
			return errors.New("invalid GPS tags")
		}
	}
	if len(altrefs) == 1 {
		altref = altrefs[0]
	}
	return p.gpsCoords.ParseEXIF(latref, latrat, longref, longrat, altref, altrat)
}

// GPS returns the value of the GPS field.
func (p *Provider) GPS() (value metadata.GPSCoords) { return p.gpsCoords }

// GPSTags returns a list of tag names for the GPS field, and a parallel list of
// values held by those tags.
func (p *Provider) GPSTags() (tags []string, values []metadata.GPSCoords) {
	return []string{"GPS  GPS*"}, []metadata.GPSCoords{p.gpsCoords}
}

// SetGPS sets the value of the GPS field.
func (p *Provider) SetGPS(value metadata.GPSCoords) (err error) {
	if value.Empty() {
		p.gpsCoords = metadata.GPSCoords{}
		p.ifd.DeleteTag(1)
		p.ifd.DeleteTag(2)
		p.ifd.DeleteTag(3)
		p.ifd.DeleteTag(4)
		p.ifd.DeleteTag(5)
		p.ifd.DeleteTag(6)
		if p.ifd.NextTag(1) == nil {
			p.ifd.DeleteTag(0)
		}
		return nil
	}
	if value.Equivalent(p.gpsCoords) {
		return nil
	}
	p.gpsCoords = value
	latref, lat, longref, long, altref, alt := value.AsEXIF()
	p.ifd.AddTag(1).SetString(latref)
	p.ifd.AddTag(2).SetRationals(lat)
	p.ifd.AddTag(3).SetString(longref)
	p.ifd.AddTag(4).SetRationals(long)
	if alt != nil {
		p.ifd.AddTag(5).SetBytes([]byte{altref})
		p.ifd.AddTag(6).SetRationals(alt)
	} else {
		p.ifd.DeleteTag(5)
		p.ifd.DeleteTag(6)
	}
	return nil
}
