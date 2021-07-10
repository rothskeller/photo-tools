package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var titleName = rdf.Name{Namespace: nsDC, Name: "title"}

// getTitle reads the value of the Title field from the RDF.
func (p *Provider) getTitle() (err error) {
	if p.dcTitle, err = getAlt(p.rdf.Properties, titleName); err != nil {
		return fmt.Errorf("dc:title: %s", err)
	}
	return nil
}

// Title returns the value of the Title field.
func (p *Provider) Title() (value string) { return p.dcTitle.Default() }

// TitleTags returns a list of tag names for the Title field, and a
// parallel list of values held by those tags.
func (p *Provider) TitleTags() (tags []string, values []string) {
	tags = append(tags, "XML  dc:title")
	if len(p.dcTitle) == 0 {
		return tags, []string{""}
	}
	values = append(values, p.dcTitle[0].Value)
	for _, ai := range p.dcTitle[1:] {
		tags = append(tags, fmt.Sprintf("XMP  dc:title[%s]", ai.Lang))
		values = append(values, ai.Value)
	}
	return tags, values
}

// SetTitle sets the values of the Title field.
func (p *Provider) SetTitle(value string) error {
	if value == "" {
		p.dcTitle = nil
		if _, ok := p.rdf.Properties[titleName]; ok {
			delete(p.rdf.Properties, titleName)
			p.dirty = true
		}
		return nil
	}
	if len(p.dcTitle) == 1 && value == p.dcTitle[0].Value {
		return nil
	}
	p.dcTitle = newAltString(value)
	setAlt(p.rdf.Properties, titleName, p.dcTitle)
	p.dirty = true
	return nil
}
