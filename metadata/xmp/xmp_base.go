package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
)

const nsXMP = "http://ns.adobe.com/xap/1.0/"
const pfxXMP = "xmp"

// XMPCreateDate returns the value of the xmp:CreateDate tag.
func (p *XMP) XMPCreateDate() metadata.DateTime { return p.xmpCreateDate }

// XMPMetadataDate returns the value of the xmp:MetadataDate tag.
func (p *XMP) XMPMetadataDate() metadata.DateTime { return p.xmpMetadataDate }

// XMPModifyDate returns the value of the xmp:ModifyDate tag.
func (p *XMP) XMPModifyDate() metadata.DateTime { return p.xmpModifyDate }

func (p *XMP) getXMP() {
	p.xmpDateTimeToMetadata(p.getString(p.rdf.Properties, pfxXMP, nsXMP, "CreateDate"), &p.xmpCreateDate)
	p.xmpDateTimeToMetadata(p.getString(p.rdf.Properties, pfxXMP, nsXMP, "MetadataDate"), &p.xmpMetadataDate)
	p.xmpDateTimeToMetadata(p.getString(p.rdf.Properties, pfxXMP, nsXMP, "ModifyDate"), &p.xmpModifyDate)
}

// SetXMPCreateDate sets the value of the xmp:CreateDate tag.
func (p *XMP) SetXMPCreateDate(v metadata.DateTime) (err error) {
	if v.Equivalent(p.xmpCreateDate) {
		return nil
	}
	p.xmpCreateDate = v
	p.setString(p.rdf.Properties, nsXMP, "CreateDate", v.String())
	return nil
}

// SetXMPMetadataDate sets the value of the xmp:MetadataDate tag.
func (p *XMP) SetXMPMetadataDate(v metadata.DateTime) (err error) {
	if v.Equivalent(p.xmpMetadataDate) {
		return nil
	}
	p.xmpMetadataDate = v
	p.setString(p.rdf.Properties, nsXMP, "MetadataDate", v.String())
	return nil
}

// SetXMPModifyDate sets the value of the xmp:ModifyDate tag.
func (p *XMP) SetXMPModifyDate(v metadata.DateTime) (err error) {
	if v.Equivalent(p.xmpModifyDate) {
		return nil
	}
	p.xmpModifyDate = v
	p.setString(p.rdf.Properties, nsXMP, "ModifyDate", v.String())
	return nil
}
