package exif

import "github.com/rothskeller/photo-tools/metadata"

// GPSCoords returns the GPS coordinates from the EXIF tags.
func (p *EXIF) GPSCoords() metadata.GPSCoords { return p.gpsCoords }

func (p *EXIF) getGPSCoords() {
	var err error

	latreft := p.gpsIFD.Tag(1)
	latratt := p.gpsIFD.Tag(2)
	longreft := p.gpsIFD.Tag(3)
	longratt := p.gpsIFD.Tag(4)
	if latreft == nil && latratt == nil && longreft == nil && longratt == nil {
		return
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
		p.log("invalid GPS tags")
		return
	}

	var altrefs []byte
	var altref byte
	var altrat []uint32
	altreft := p.gpsIFD.Tag(5)
	altratt := p.gpsIFD.Tag(6)
	if altratt != nil {
		if altreft != nil {
			altrefs, err = altreft.AsBytes()
		}
		if err == nil {
			altrat, err = altratt.AsRationals()
		}
		if err != nil || len(altrefs) > 1 || len(altrat) != 2 {
			p.log("invalid GPS tags")
			return
		}
	}
	if len(altrefs) == 1 {
		altref = altrefs[0]
	}
	if err := p.gpsCoords.ParseEXIF(latref, latrat, longref, longrat, altref, altrat); err != nil {
		p.log(err.Error())
	}
}

// SetGPSCoords sets the values of the GPS coordinate tags.
func (p *EXIF) SetGPSCoords(v metadata.GPSCoords) error {
	if v.Equivalent(p.gpsCoords) {
		return nil
	}
	p.gpsCoords = v
	if p.gpsCoords.Empty() {
		p.gpsIFD.DeleteTag(1)
		p.gpsIFD.DeleteTag(2)
		p.gpsIFD.DeleteTag(3)
		p.gpsIFD.DeleteTag(4)
		p.gpsIFD.DeleteTag(5)
		p.gpsIFD.DeleteTag(6)
		if p.gpsIFD != nil && p.gpsIFD.NextTag(1) == nil {
			p.ifd0.DeleteTag(tagGPSIFDOffset)
		}
		return nil
	}
	if p.gpsIFD == nil {
		p.addGPSIFD()
	}
	latref, latrat, longref, longrat, altref, altrat := p.gpsCoords.AsEXIF()
	p.gpsIFD.AddTag(1).SetString(latref)
	p.gpsIFD.AddTag(2).SetRationals(latrat)
	p.gpsIFD.AddTag(3).SetString(longref)
	p.gpsIFD.AddTag(4).SetRationals(longrat)
	if len(altrat) != 0 {
		p.gpsIFD.AddTag(5).SetBytes([]byte{altref})
		p.gpsIFD.AddTag(6).SetRationals(altrat)
	} else {
		p.gpsIFD.DeleteTag(5)
		p.gpsIFD.DeleteTag(6)
	}
	return nil
}
