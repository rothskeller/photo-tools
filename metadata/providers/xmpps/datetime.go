package xmpps

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var dateCreatedName = rdf.Name{Namespace: nsPS, Name: "DateCreated"}

// getDateTime reads the value of the DateTime field from the RDF.
func (p *Provider) getDateTime() (err error) {
	var s string
	if s, err = getString(p.rdf.Properties, dateCreatedName); err == nil {
		err = p.psDateCreated.Parse(s)
	}
	if err != nil {
		return fmt.Errorf("photoshop:DateCreated: %s", err)
	}
	return nil
}

// DateTime returns the value of the DateTime field.
func (p *Provider) DateTime() (value metadata.DateTime) { return p.psDateCreated }

// DateTimeTags returns a list of tag names for the DateTime field, and
// a parallel list of values held by those tags.
func (p *Provider) DateTimeTags() (tags []string, values []metadata.DateTime) {
	return []string{"XMP  photoshop:DateCreated"}, []metadata.DateTime{p.psDateCreated}
}

// SetDateTime sets the value of the DateTime field.
func (p *Provider) SetDateTime(value metadata.DateTime) error {
	if value.Empty() {
		p.psDateCreated = metadata.DateTime{}
		if _, ok := p.rdf.Properties[dateCreatedName]; ok {
			delete(p.rdf.Properties, dateCreatedName)
			p.dirty = true
		}
		return nil
	}
	if value.Equivalent(p.psDateCreated) {
		return nil
	}
	p.psDateCreated = value
	setString(p.rdf.Properties, dateCreatedName, value.String())
	p.dirty = true
	return nil
}
