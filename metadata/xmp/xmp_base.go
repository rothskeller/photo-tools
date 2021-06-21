package xmp

import (
	xmpbase "github.com/rothskeller/photo-tools/metadata/xmp/models/xmp_base"
)

// CreateDate returns the creation date from the XMP.
func (p *XMP) CreateDate() string {
	if p == nil || p.doc == nil {
		return ""
	}
	if model := xmpbase.FindModel(p.doc); model != nil {
		if model.CreateDate != "" && !dateRE.MatchString(model.CreateDate) {
			p.log("CreateDate: invalid value: %q", model.CreateDate)
			return ""
		}
		return canonicalDate(model.CreateDate)
	}
	return ""
}

// SetCreateDate sets the creation date in the XMP.
func (p *XMP) SetCreateDate(v string) {
	model, err := xmpbase.MakeModel(p.doc)
	if err != nil {
		p.log("XMP xmpbase.MakeModel: %s", err)
		return
	}
	model.CreateDate = v
}

// ModifyDate returns the modification date from the XMP.
func (p *XMP) ModifyDate() string {
	if p == nil || p.doc == nil {
		return ""
	}
	if model := xmpbase.FindModel(p.doc); model != nil {
		if model.ModifyDate != "" && !dateRE.MatchString(model.ModifyDate) {
			p.log("ModifyDate: invalid value: %q", model.ModifyDate)
			return ""
		}
		return canonicalDate(model.ModifyDate)
	}
	return ""
}

// SetModifyDate sets the modification date in the XMP.
func (p *XMP) SetModifyDate(v string) {
	model, err := xmpbase.MakeModel(p.doc)
	if err != nil {
		p.log("XMP xmpbase.MakeModel: %s", err)
		return
	}
	model.ModifyDate = v
}

// MetadataDate returns the metadata date from the XMP.
func (p *XMP) MetadataDate() string {
	if p == nil || p.doc == nil {
		return ""
	}
	if model := xmpbase.FindModel(p.doc); model != nil {
		if model.MetadataDate != "" && !dateRE.MatchString(model.MetadataDate) {
			p.log("MetadataDate: invalid value: %q", model.MetadataDate)
			return ""
		}
		return canonicalDate(model.MetadataDate)
	}
	return ""
}

// SetMetadataDate sets the metadata date in the XMP.
func (p *XMP) SetMetadataDate(v string) {
	model, err := xmpbase.MakeModel(p.doc)
	if err != nil {
		p.log("XMP xmpbase.MakeModel: %s", err)
		return
	}
	model.MetadataDate = v
}
