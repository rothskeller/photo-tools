package xmp

import (
	"github.com/rothskeller/photo-tools/metadata/xmp/models/tiff"
	"trimmer.io/go-xmp/xmp"
)

// TIFFArtist returns the list of Artists from the XMP.
func (p *XMP) TIFFArtist() []string {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := tiff.FindModel(p.doc); model != nil {
		if a := model.Artist; a != "" {
			return []string{a}
		}
	}
	return nil
}

// SetTIFFArtist sets the TIFF Artist in the XMP.
func (p *XMP) SetTIFFArtist(artist string) {
	model, err := tiff.MakeModel(p.doc)
	if err != nil {
		p.log("XMP tiff.MakeModel: %s", err)
		return
	}
	model.Artist = artist
}

// TIFFDateTime returns the date and time from the XMP.
func (p *XMP) TIFFDateTime() string {
	if p == nil || p.doc == nil {
		return ""
	}
	if model := tiff.FindModel(p.doc); model != nil {
		if model.DateTime != "" && !dateRE.MatchString(model.DateTime) {
			p.log("TIFFDateTime: invalid value: %q", model.DateTime)
			return ""
		}
		return canonicalDate(model.DateTime)
	}
	return ""
}

// SetTIFFDateTime sets the TIFF DateTime in the XMP.
func (p *XMP) SetTIFFDateTime(dt string) {
	model, err := tiff.MakeModel(p.doc)
	if err != nil {
		p.log("XMP tiff.MakeModel: %s", err)
		return
	}
	model.DateTime = dt
}

// TIFFImageDescription returns the TIFF ImageDescription from the XMP.  It is
// an ordered list of alternatives, each of which has is a [language, value]
// pair.  The first item on the list is the default language.
func (p *XMP) TIFFImageDescription() (descs [][]string) {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := tiff.FindModel(p.doc); model != nil {
		for _, alt := range model.ImageDescription {
			descs = append(descs, []string{alt.Lang, alt.Value})
		}
		return descs
	}
	return nil
}

// SetTIFFImageDescription sets the TIFF ImageDescription in the XMP.  It sets
// only the default language and discards any others.
func (p *XMP) SetTIFFImageDescription(desc string) {
	model, err := tiff.MakeModel(p.doc)
	if err != nil {
		p.log("XMP tiff.MakeModel: %s", err)
		return
	}
	if desc != "" {
		model.ImageDescription = xmp.NewAltString(desc)
	} else {
		model.ImageDescription = nil
	}
}
