package xmp

import (
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
	p.TIFFArtist = model.Artist
	p.xmpDateTimeToMetadata(model.DateTime, &p.TIFFDateTime)
	p.TIFFImageDescription = model.ImageDescription
}

func (p *XMP) setTIFF() {
	var (
		model *tiff.TiffInfo
		dt    metadata.DateTime
		err   error
	)
	if model, err = tiff.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if p.TIFFArtist != model.Artist {
		model.Artist = p.TIFFArtist
		p.dirty = true
	}
	p.xmpDateTimeToMetadata(model.DateTime, &dt)
	if eq, _ := dt.Equivalent(&p.TIFFDateTime); !eq {
		model.DateTime = p.TIFFDateTime.String()
		p.dirty = true
	}
	if !metadata.EqualAltStrings(p.TIFFImageDescription, model.ImageDescription) {
		model.ImageDescription = p.TIFFImageDescription
		p.dirty = true
	}
}
