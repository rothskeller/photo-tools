package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	xmpbase "github.com/rothskeller/photo-tools/metadata/xmp/models/xmp_base"
)

// XMPCreateDate returns the value of the xmp:CreateDate tag.
func (p *XMP) XMPCreateDate() metadata.DateTime { return p.xmpCreateDate }

// XMPMetadataDate returns the value of the xmp:MetadataDate tag.
func (p *XMP) XMPMetadataDate() metadata.DateTime { return p.xmpMetadataDate }

// XMPModifyDate returns the value of the xmp:ModifyDate tag.
func (p *XMP) XMPModifyDate() metadata.DateTime { return p.xmpModifyDate }

func (p *XMP) getXMP() {
	var model *xmpbase.XmpBase

	if p != nil && p.doc != nil {
		model = xmpbase.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.xmpDateTimeToMetadata(model.CreateDate, &p.xmpCreateDate)
	p.xmpDateTimeToMetadata(model.MetadataDate, &p.xmpMetadataDate)
	p.xmpDateTimeToMetadata(model.ModifyDate, &p.xmpModifyDate)
}

// SetXMPCreateDate sets the value of the xmp:CreateDate tag.
func (p *XMP) SetXMPCreateDate(v metadata.DateTime) (err error) {
	var model *xmpbase.XmpBase

	if model, err = xmpbase.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add xmp model to XMP: %s", err)
	}
	if v.Equivalent(p.xmpCreateDate) {
		return nil
	}
	p.xmpCreateDate = v
	model.CreateDate = v.String()
	p.dirty = true
	return nil
}

// SetXMPMetadataDate sets the value of the xmp:MetadataDate tag.
func (p *XMP) SetXMPMetadataDate(v metadata.DateTime) (err error) {
	var model *xmpbase.XmpBase

	if model, err = xmpbase.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add xmp model to XMP: %s", err)
	}
	if v.Equivalent(p.xmpMetadataDate) {
		return nil
	}
	p.xmpMetadataDate = v
	model.MetadataDate = v.String()
	p.dirty = true
	return nil
}

// SetXMPModifyDate sets the value of the xmp:ModifyDate tag.
func (p *XMP) SetXMPModifyDate(v metadata.DateTime) (err error) {
	var model *xmpbase.XmpBase

	if model, err = xmpbase.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add xmp model to XMP: %s", err)
	}
	if v.Equivalent(p.xmpModifyDate) {
		return nil
	}
	p.xmpModifyDate = v
	model.ModifyDate = v.String()
	p.dirty = true
	return nil
}
