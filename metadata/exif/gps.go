package exif

import (
	"bytes"

	"github.com/rothskeller/photo-tools/metadata"
)

// GPSCoords returns the GPS coordinates where the picture was taken.
func (p *EXIF) GPSCoords() (gc metadata.GPSCoords) {
	if p == nil || p.gpsIFD == nil {
		return metadata.GPSCoords{}
	}
	gc.Latitude = p.exifToFixedFloat(1, 2)
	gc.Longitude = p.exifToFixedFloat(3, 4)
	gc.Altitude = p.exifToFixedFloat(5, 6)
	return gc
}

// SetGPSCoords sets or removes the GPS coordinates where the picture was taken.
func (p *EXIF) SetGPSCoords(gc metadata.GPSCoords) {
	if p == nil || p.ifd0 == nil {
		// We're not going to add an EXIF block just for this.
		return
	}
	if !gc.Valid() {
		// Remove any existing GPS tags.
		if p.gpsIFD != nil {
			p.deleteTag(p.gpsIFD, 1)
			p.deleteTag(p.gpsIFD, 2)
			p.deleteTag(p.gpsIFD, 3)
			p.deleteTag(p.gpsIFD, 4)
			p.deleteTag(p.gpsIFD, 5)
			p.deleteTag(p.gpsIFD, 6)
			// If this is the last tag in the block, other than the
			// GPS version, remove the whole block.
			if len(p.gpsIFD.tags) == 1 && p.gpsIFD.tags[0].tag == 0 {
				p.gpsIFD = nil
				p.deleteTag(p.ifd0, tagGPSIFDOffset)
			}
		}
		return
	}
	if p.gpsIFD == nil {
		// We need to add GPS tags, and there's no GPS block, so add the
		// block.
		p.gpsIFD = new(ifdt)
		p.addTag(p.gpsIFD, &tagt{
			tag:   0, // GPS Version
			ttype: 1, // BYTE
			count: 4,
			data:  []byte{2, 3, 0, 0},
		})
		// We also need to make sure there's a place to store the offset
		// of the new GPS block.
		p.addTag(p.ifd0, &tagt{
			tag:   tagGPSIFDOffset,
			ttype: 4, // LONG
			count: 1,
			data:  []byte{0, 0, 0, 0}, // filled in later
		})
	}
	p.fixedFloatToExif(gc.Latitude, 1, 2, 'N', 'S')
	p.fixedFloatToExif(gc.Longitude, 3, 4, 'E', 'W')
	if gc.HasAltitude() {
		p.fixedFloatToExif(gc.Altitude, 5, 6, 0, 1)
	} else {
		p.deleteTag(p.gpsIFD, 5)
		p.deleteTag(p.gpsIFD, 6)
	}
}

// exifToFixedFloat takes tag numbers for a reference tag and a tag containing
// rational values, and returns the result converted to a single FixedFloat with
// appropriate sign.  It returns 0 if the tags are missing// or invalid.
func (p *EXIF) exifToFixedFloat(ref, rat uint16) (val metadata.FixedFloat) {
	var neg bool

	reft := p.gpsIFD.findTag(ref)
	ratt := p.gpsIFD.findTag(rat)
	if reft == nil || ratt == nil {
		return 0
	}
	if reft.ttype == 2 && reft.count == 2 && ratt.ttype == 5 && ratt.count == 3 {
		if reft.data[0] == 'W' || reft.data[0] == 'S' {
			neg = true
		}
	} else if reft.ttype == 1 && reft.count == 1 && ratt.ttype == 5 && ratt.count == 1 {
		if reft.data[0] == 1 {
			neg = true
		}
	} else {
		p.log(reft.offset, "invalid GPS tags")
		return 0
	}
	val = p.exifRatToFixedFloat(ratt.data)
	if neg {
		val = -val
	}
	return val
}

// exifRatToFixedFloat converts the data from an EXIF tag containing rational
// values, and returns the result converted to a single FixedFloat.
func (p *EXIF) exifRatToFixedFloat(data []byte) (val metadata.FixedFloat) {
	denmult := 1
	for len(data) != 0 {
		num := p.enc.Uint32(data)
		den := p.enc.Uint32(data[4:])
		if den == 0 {
			return val
		}
		val += metadata.FixedFloatFromFraction(int(num), int(den)*denmult)
		denmult *= 60
	}
	return val
}

func (p *EXIF) fixedFloatToExif(val metadata.FixedFloat, reft, ratt uint16, pos, neg byte) {
	// First, set the reference tag (i.e, the sign of the value).
	reftag := p.gpsIFD.findTag(reft)
	if reftag == nil {
		if pos > 0 { // Latitude or Longitude
			reftag = &tagt{tag: reft, ttype: 2, count: 2, data: []byte{0, 0}}
		} else {
			reftag = &tagt{tag: reft, ttype: 1, count: 1, data: []byte{0}}
		}
		p.addTag(p.gpsIFD, reftag)
	}
	var refbyte byte
	if val >= 0 {
		refbyte = pos
	} else {
		refbyte = neg
		val = -val
	}
	if refbyte != reftag.data[0] {
		reftag.data[0] = refbyte
		p.gpsIFD.dirty = true
	}
	// Next, set the value tag.
	rattag := p.gpsIFD.findTag(ratt)
	if rattag == nil {
		if pos > 0 { // Latitude or Longitude
			rattag = &tagt{tag: ratt, ttype: 5, count: 3, data: make([]byte, 24)}
		} else {
			rattag = &tagt{tag: ratt, ttype: 5, count: 1, data: make([]byte, 8)}
		}
		p.addTag(p.gpsIFD, rattag)
	} else {
		// Before converting the FixedFloat to EXIF, let's first try
		// converting the EXIF to a FixedFloat and see if they match.
		// That allows us to avoid converting an EXIF degree, minutes,
		// seconds/100 triplet into a less accurate degrees/1000000, 0,
		// 0 triplet if we don't need to.
		trial := p.exifRatToFixedFloat(rattag.data)
		if trial == val {
			return
		}
	}
	data := make([]byte, len(rattag.data))
	if pos > 0 {
		// For latitude and longitude, set the second and third
		// rationals to 0/1, since we do not use them.
		p.enc.PutUint32(data[8:], 0)
		p.enc.PutUint32(data[12:], 1)
		p.enc.PutUint32(data[16:], 0)
		p.enc.PutUint32(data[20:], 1)
	}
	p.enc.PutUint32(data, uint32(val))
	p.enc.PutUint32(data[4:], 1000000)
	// If the data changed, mark the IFD dirty.
	if !bytes.Equal(data, rattag.data) {
		rattag.data = data
		p.gpsIFD.dirty = true
	}
}
