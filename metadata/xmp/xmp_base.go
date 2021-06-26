package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
	xmpbase "github.com/rothskeller/photo-tools/metadata/xmp/models/xmp_base"
)

func (p *XMP) getXMP() {
	var model *xmpbase.XmpBase

	if p != nil && p.doc != nil {
		model = xmpbase.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.xmpDateTimeToMetadata(model.CreateDate, &p.XMPCreateDate)
	p.xmpDateTimeToMetadata(model.MetadataDate, &p.XMPMetadataDate)
	p.xmpDateTimeToMetadata(model.ModifyDate, &p.XMPModifyDate)
}

func (p *XMP) setXMP() {
	var (
		model *xmpbase.XmpBase
		dt    metadata.DateTime
		err   error
	)
	if model, err = xmpbase.MakeModel(p.doc); err != nil {
		panic(err)
	}
	p.xmpDateTimeToMetadata(model.CreateDate, &dt)
	if eq, _ := dt.Equivalent(&p.XMPCreateDate); !eq {
		model.CreateDate = p.XMPCreateDate.String()
		p.dirty = true
	}
	p.xmpDateTimeToMetadata(model.MetadataDate, &dt)
	if eq, _ := dt.Equivalent(&p.XMPMetadataDate); !eq {
		model.MetadataDate = p.XMPMetadataDate.String()
		p.dirty = true
	}
	p.xmpDateTimeToMetadata(model.ModifyDate, &dt)
	if eq, _ := dt.Equivalent(&p.XMPModifyDate); !eq {
		model.ModifyDate = p.XMPModifyDate.String()
		p.dirty = true
	}
}
