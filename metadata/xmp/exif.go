package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/exif"
	"trimmer.io/go-xmp/xmp"
)

func (p *XMP) getEXIF() {
	var model *exif.ExifInfo

	if p != nil && p.doc != nil {
		model = exif.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.xmpDateTimeToMetadata(model.DateTimeDigitized, &p.EXIFDateTimeDigitized)
	p.xmpDateTimeToMetadata(model.DateTimeOriginal, &p.EXIFDateTimeOriginal)
	p.xmpEXIFGPSCoordsToMetadata(model.GPSLatitude, model.GPSLongitude, model.GPSAltitudeRef, model.GPSAltitude, &p.EXIFGPSCoords)
	p.EXIFUserComments = model.UserComment
}

func (p *XMP) setEXIF() {
	var (
		model *exif.ExifInfo
		dt    metadata.DateTime
		gps   metadata.GPSCoords
		err   error
	)
	if model, err = exif.MakeModel(p.doc); err != nil {
		panic(err)
	}
	p.xmpDateTimeToMetadata(model.DateTimeDigitized, &dt)
	if eq, _ := dt.Equivalent(&p.EXIFDateTimeDigitized); !eq {
		model.DateTimeDigitized = p.EXIFDateTimeDigitized.String()
		p.dirty = true
	}
	p.xmpDateTimeToMetadata(model.DateTimeOriginal, &dt)
	if eq, _ := dt.Equivalent(&p.EXIFDateTimeOriginal); !eq {
		model.DateTimeOriginal = p.EXIFDateTimeOriginal.String()
		p.dirty = true
	}
	// GPS coordinate transformations involve rounding error, hence the
	// roundabout equivalency check.
	p.xmpEXIFGPSCoordsToMetadata(model.GPSLatitude, model.GPSLongitude, model.GPSAltitudeRef, model.GPSAltitude, &gps)
	if eq, _ := gps.Equivalent(&p.EXIFGPSCoords); !eq {
		model.GPSLatitude, model.GPSLongitude, model.GPSAltitudeRef, model.GPSAltitude = p.EXIFGPSCoords.AsXMP()
		p.dirty = true
	}
	if !stringSliceEqual(xmp.StringArray(p.EXIFUserComments), model.UserComment) {
		model.UserComment = p.EXIFUserComments
		p.dirty = true
	}
}

func (p *XMP) xmpEXIFGPSCoordsToMetadata(lat, long, altref, alt string, m *metadata.GPSCoords) {
	if err := m.ParseXMP(lat, long, altref, alt); err != nil {
		p.log("invalid GPS coordinates")
	}
}
