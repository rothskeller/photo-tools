package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
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
		dc    metadata.DateTime
		err   error
	)
	if model, err = ps.MakeModel(p.doc); err != nil {
		panic(err)
	}
	p.xmpDateTimeToMetadata(model.DateCreated, &dc)
	if eq, _ := dc.Equivalent(&p.PSDateCreated); !eq {
		model.DateCreated = p.PSDateCreated.String()
		p.dirty = true
	}
}
