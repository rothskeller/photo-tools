package xmp

import (
	"github.com/rothskeller/photo-tools/metadata/xmp/models/ps"
)

func (p *XMP) getPS() {
	var model *ps.PhotoshopInfo

	if p != nil && p.doc != nil {
		model = ps.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.xmpDateTimeToMetadata(model.DateCreated, &p.PSDateCreated)
}

func (p *XMP) setPS() {
	var (
		model *ps.PhotoshopInfo
		err   error
	)
	if model, err = ps.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if dc := p.PSDateCreated.String(); dc != model.DateCreated {
		model.DateCreated = dc
		p.dirty = true
	}
}
