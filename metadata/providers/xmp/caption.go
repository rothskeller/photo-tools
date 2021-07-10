package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var descriptionName = rdf.Name{Namespace: nsDC, Name: "description"}

// getCaption reads the value of the Caption field from the RDF.
func (p *Provider) getCaption() (err error) {
	if p.dcDescription, err = getAlt(p.rdf.Properties, descriptionName); err != nil {
		return fmt.Errorf("dc:description: %s", err)
	}
	return nil
}

// Caption returns the value of the Caption field.
func (p *Provider) Caption() (value string) { return p.dcDescription.Default() }

// CaptionTags returns a list of tag names for the Caption field, and a
// parallel list of values held by those tags.
func (p *Provider) CaptionTags() (tags []string, values []string) {
	tags = append(tags, "XML  dc:description")
	if len(p.dcDescription) == 0 {
		return tags, []string{""}
	}
	values = append(values, p.dcDescription[0].Value)
	for _, ai := range p.dcDescription[1:] {
		tags = append(tags, fmt.Sprintf("XMP  dc:description[%s]", ai.Lang))
		values = append(values, ai.Value)
	}
	return tags, values
}

// SetCaption sets the value of the Caption field.
func (p *Provider) SetCaption(value string) error {
	if value == "" {
		p.dcDescription = nil
		if _, ok := p.rdf.Properties[descriptionName]; ok {
			delete(p.rdf.Properties, descriptionName)
			p.dirty = true
		}
		return nil
	}
	if len(p.dcDescription) == 1 && value == p.dcDescription[0].Value {
		return nil
	}
	p.dcDescription = newAltString(value)
	setAlt(p.rdf.Properties, descriptionName, p.dcDescription)
	p.dirty = true
	return nil
}
