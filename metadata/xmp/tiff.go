package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/tiff"
	"trimmer.io/go-xmp/xmp"
)

// TIFFArtist returns the value of the tiff:Artist tag.
func (p *XMP) TIFFArtist() string { return p.tiffArtist }

// TIFFDateTime returns the value of the tiff:DateTime tag.
func (p *XMP) TIFFDateTime() metadata.DateTime { return p.tiffDateTime }

// TIFFImageDescription returns the value of the tiff:ImageDescription tag.
func (p *XMP) TIFFImageDescription() metadata.AltString { return p.tiffImageDescription }

func (p *XMP) getTIFF() {
	var model *tiff.TiffInfo

	if p != nil && p.doc != nil {
		model = tiff.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.tiffArtist = model.Artist
	p.xmpDateTimeToMetadata(model.DateTime, &p.tiffDateTime)
	p.tiffImageDescription = model.ImageDescription
}

// SetTIFFArtist sets the value of the tiff:Artist tag.
func (p *XMP) SetTIFFArtist(v string) (err error) {
	var model *tiff.TiffInfo

	if model, err = tiff.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add tiff model to XMP: %s", err)
	}
	if v == p.tiffArtist {
		return nil
	}
	p.tiffArtist = v
	model.Artist = v
	p.dirty = true
	return nil
}

// SetTIFFDateTime sets the value of the tiff:DateTime tag.
func (p *XMP) SetTIFFDateTime(v metadata.DateTime) (err error) {
	var model *tiff.TiffInfo

	if model, err = tiff.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add tiff model to XMP: %s", err)
	}
	if v.Equivalent(p.tiffDateTime) {
		return nil
	}
	p.tiffDateTime = v
	model.DateTime = v.String()
	p.dirty = true
	return nil
}

// SetTIFFImageDescription sets the value of the tiff:ImageDescription tag.
func (p *XMP) SetTIFFImageDescription(v xmp.AltString) (err error) {
	var model *tiff.TiffInfo

	if model, err = tiff.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add tiff model to XMP: %s", err)
	}
	if metadata.EqualAltStrings(v, p.tiffImageDescription) {
		return nil
	}
	p.tiffImageDescription = v
	model.ImageDescription = v
	p.dirty = true
	return nil
}
