package xmptiff

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var dateTimeName = rdf.Name{Namespace: nsTIFF, Name: "DateTime"}

// getDateTime reads the value of the DateTime field from the RDF.
func (p *Provider) getDateTime() (err error) {
	var s string
	if s, err = getString(p.rdf.Properties, dateTimeName); err == nil {
		err = p.tiffDateTime.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("tiff:DateTime: %s", err)
	}
	return nil
}

// DateTime returns the value of the DateTime field.
func (p *Provider) DateTime() (value metadata.DateTime) { return p.tiffDateTime }

// DateTimeTags returns a list of tag names for the DateTime field, and
// a parallel list of values held by those tags.
func (p *Provider) DateTimeTags() (tags []string, values []metadata.DateTime) {
	if p.tiffDateTime.Empty() {
		return nil, nil
	}
	return []string{"XMP  tiff:DateTime"}, []metadata.DateTime{p.tiffDateTime}
}

// SetDateTime sets the value of the DateTime field.
func (p *Provider) SetDateTime(value metadata.DateTime) error {
	p.tiffDateTime = metadata.DateTime{}
	if _, ok := p.rdf.Properties[dateTimeName]; ok {
		delete(p.rdf.Properties, dateTimeName)
		p.dirty = true
	}
	return nil
}
