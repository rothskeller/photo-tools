package xmp

import (
	"reflect"

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
		err   error
	)
	if model, err = exif.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if dtd := p.EXIFDateTimeDigitized.String(); dtd != model.DateTimeDigitized {
		model.DateTimeDigitized = dtd
		p.dirty = true
	}
	if dto := p.EXIFDateTimeOriginal.String(); dto != model.DateTimeOriginal {
		model.DateTimeOriginal = dto
		p.dirty = true
	}
	if lat, long, altref, alt := p.EXIFGPSCoords.AsXMP(); lat != model.GPSLatitude || long != model.GPSLongitude ||
		altref != model.GPSAltitudeRef || alt != model.GPSAltitude {
		model.GPSLatitude, model.GPSLongitude, model.GPSAltitudeRef, model.GPSAltitude = lat, long, altref, alt
		p.dirty = true
	}
	if !reflect.DeepEqual(xmp.StringArray(p.EXIFUserComments), model.UserComment) {
		model.UserComment = p.EXIFUserComments
		p.dirty = true
	}
}

func (p *XMP) xmpEXIFGPSCoordsToMetadata(lat, long, altref, alt string, m *metadata.GPSCoords) {
	if err := m.ParseXMP(lat, long, altref, alt); err != nil {
		p.log("invalid GPS coordinates")
	}
}
