package exif

func (p *EXIF) getGPSCoords() {
	latreft := p.gpsIFD.findTag(1)
	latratt := p.gpsIFD.findTag(2)
	longreft := p.gpsIFD.findTag(3)
	longratt := p.gpsIFD.findTag(4)
	if latreft == nil && latratt == nil && longreft == nil && longratt == nil {
		return
	}
	if latreft == nil || latreft.ttype != 2 || latreft.count != 2 ||
		latratt == nil || latratt.ttype != 5 || latratt.count != 3 ||
		longreft == nil || longreft.ttype != 2 || longratt.count != 2 ||
		longratt == nil || longratt.ttype != 5 || longratt.count != 3 {
		p.log(p.gpsIFD.offset, "invalid GPS tags")
		return
	}
	latref := p.asciiAt(latreft, "GPSLatitudeRef")
	longref := p.asciiAt(longreft, "GPSLongitudeRef")
	latrat := p.exifRatToUint32(latratt.data)
	longrat := p.exifRatToUint32(longratt.data)

	var altref byte
	var altrat []uint32
	altreft := p.gpsIFD.findTag(5)
	altratt := p.gpsIFD.findTag(6)
	if altreft != nil || altratt != nil {
		if altreft == nil || altreft.ttype != 1 || altreft.count != 1 ||
			altratt == nil || altratt.ttype != 5 || altratt.count != 1 {
			p.log(p.gpsIFD.offset, "invalid GPS tags")
			return
		}
		altref = altreft.data[0]
		altrat = p.exifRatToUint32(altratt.data)
	}
	if err := p.GPSCoords.ParseEXIF(latref, latrat, longref, longrat, altref, altrat); err != nil {
		p.log(p.gpsIFD.offset, err.Error())
	}
}
func (p *EXIF) exifRatToUint32(rat []byte) (u []uint32) {
	u = make([]uint32, len(rat)/4)
	for i := 0; i < len(rat); i += 4 {
		u[i/4] = p.enc.Uint32(rat[i:])
	}
	return u
}

func (p *EXIF) setGPSCoords() {
	if p.GPSCoords.Empty() {
		p.deleteTag(p.gpsIFD, 1)
		p.deleteTag(p.gpsIFD, 2)
		p.deleteTag(p.gpsIFD, 3)
		p.deleteTag(p.gpsIFD, 4)
		p.deleteTag(p.gpsIFD, 5)
		p.deleteTag(p.gpsIFD, 6)
		if p.gpsIFD != nil && len(p.gpsIFD.tags) == 1 && p.gpsIFD.tags[0].tag == 0 {
			p.gpsIFD = nil
			p.deleteTag(p.ifd0, tagGPSIFDOffset)
		}
		return
	}
	if p.gpsIFD == nil {
		p.addGPSIFD()
	}
	latref, latrat, longref, longrat, altref, altrat := p.GPSCoords.AsEXIF()
	p.setASCIITag(p.gpsIFD, 1, latref)
	p.setRationalTag(p.gpsIFD, 2, latrat)
	p.setASCIITag(p.gpsIFD, 3, longref)
	p.setRationalTag(p.gpsIFD, 4, longrat)
	if len(altrat) != 0 {
		p.setByteTag(p.gpsIFD, 5, altref)
		p.setRationalTag(p.gpsIFD, 6, altrat)
	} else {
		p.deleteTag(p.gpsIFD, 5)
		p.deleteTag(p.gpsIFD, 6)
	}
}
