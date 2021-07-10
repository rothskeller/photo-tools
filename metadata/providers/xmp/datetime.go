package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	createDateName   = rdf.Name{Namespace: nsXMP, Name: "CreateDate"}
	metadataDateName = rdf.Name{Namespace: nsXMP, Name: "MetadataDate"}
	modifyDateName   = rdf.Name{Namespace: nsXMP, Name: "ModifyDate"}
)

// getDateTime reads the value of the DateTime field from the RDF.
func (p *Provider) getDateTime() (err error) {
	var s string
	if s, err = getString(p.rdf.Properties, createDateName); err == nil {
		err = p.xmpCreateDate.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("xmp:CreateDate: %s", err)
	}
	if s, err = getString(p.rdf.Properties, metadataDateName); err == nil {
		err = p.xmpMetadataDate.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("xmp:MetadataDate: %s", err)
	}
	if s, err = getString(p.rdf.Properties, modifyDateName); err == nil {
		err = p.xmpModifyDate.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("xmp:ModifyDate: %s", err)
	}
	return nil
}

// DateTime returns the value of the DateTime field.
func (p *Provider) DateTime() (value metadata.DateTime) {
	if !p.xmpCreateDate.Empty() {
		return p.xmpCreateDate
	}
	if !p.xmpModifyDate.Empty() {
		return p.xmpModifyDate
	}
	return p.xmpMetadataDate // which may be empty
}

// DateTimeTags returns a list of tag names for the DateTime field, and
// a parallel list of values held by those tags.
func (p *Provider) DateTimeTags() (tags []string, values []metadata.DateTime) {
	tags = append(tags, "XMP  xmp:CreateDate")
	values = append(values, p.xmpCreateDate)
	if !p.xmpModifyDate.Empty() {
		tags = append(tags, "XMP  xmp:ModifyDate")
		values = append(values, p.xmpModifyDate)
	}
	if !p.xmpMetadataDate.Empty() {
		tags = append(tags, "XMP  xmp:MetadataDate")
		values = append(values, p.xmpMetadataDate)
	}
	return tags, values
}

// SetDateTime sets the value of the DateTime field.
func (p *Provider) SetDateTime(value metadata.DateTime) error {
	p.xmpModifyDate = metadata.DateTime{}
	if _, ok := p.rdf.Properties[modifyDateName]; ok {
		delete(p.rdf.Properties, modifyDateName)
		p.dirty = true
	}
	p.xmpMetadataDate = metadata.DateTime{}
	if _, ok := p.rdf.Properties[metadataDateName]; ok {
		delete(p.rdf.Properties, metadataDateName)
		p.dirty = true
	}
	if value.Empty() {
		p.xmpCreateDate = metadata.DateTime{}
		if _, ok := p.rdf.Properties[createDateName]; ok {
			delete(p.rdf.Properties, createDateName)
			p.dirty = true
		}
		return nil
	}
	if value.Equivalent(p.xmpCreateDate) {
		return nil
	}
	p.xmpCreateDate = value
	setString(p.rdf.Properties, createDateName, value.String())
	p.dirty = true
	return nil
}
