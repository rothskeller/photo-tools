package xmp

import (
	"reflect"

	"github.com/rothskeller/photo-tools/metadata/xmp/models/tiff"
)

func (p *XMP) getTIFF() {
	var model *tiff.TiffInfo

	if p != nil && p.doc != nil {
		model = tiff.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.TIFFArtist = model.Artist
	p.xmpDateTimeToMetadata(model.DateTime, &p.TIFFDateTime)
	p.TIFFImageDescription = xmpAltStringToMetadata(model.ImageDescription)
}

func (p *XMP) setTIFF() {
	var (
		model *tiff.TiffInfo
		err   error
	)
	if model, err = tiff.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if p.TIFFArtist != model.Artist {
		model.Artist = p.TIFFArtist
		p.dirty = true
	}
	if dt := p.TIFFDateTime.String(); dt != model.DateTime {
		model.DateTime = dt
		p.dirty = true
	}
	if desc := metadataToXMPAltString(p.TIFFImageDescription); !reflect.DeepEqual(desc, model.ImageDescription) {
		model.ImageDescription = desc
		p.dirty = true
	}
}
