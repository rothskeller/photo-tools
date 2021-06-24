package xmp

import (
	"reflect"

	"github.com/rothskeller/photo-tools/metadata"
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
	if model.Artist != "" {
		p.TIFFArtist = metadata.NewString(model.Artist)
	}
	p.TIFFDateTime = p.xmpDateTimeToMetadata(model.DateTime)
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
	if artist := p.TIFFArtist.String(); artist != model.Artist {
		model.Artist = artist
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
