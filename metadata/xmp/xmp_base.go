package xmp

import (
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
		err   error
	)
	if model, err = xmpbase.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if d := p.XMPCreateDate.String(); d != model.CreateDate {
		model.CreateDate = d
		p.dirty = true
	}
	if d := p.XMPMetadataDate.String(); d != model.MetadataDate {
		model.MetadataDate = d
		p.dirty = true
	}
	if d := p.XMPModifyDate.String(); d != model.ModifyDate {
		model.ModifyDate = d
		p.dirty = true
	}
}
