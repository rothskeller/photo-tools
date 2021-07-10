package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var creatorName = rdf.Name{Namespace: nsDC, Name: "creator"}

// getCreator reads the value of the Creator field from the RDF.
func (p *Provider) getCreator() (err error) {
	if p.dcCreator, err = getStrings(p.rdf.Properties, creatorName); err != nil {
		return fmt.Errorf("dc:creator: %s", err)
	}
	return nil
}

// Creator returns the value of the Creator field.
func (p *Provider) Creator() (value string) {
	if len(p.dcCreator) == 0 {
		return ""
	}
	return p.dcCreator[0]
}

// CreatorTags returns a list of tag names for the Creator field, and a
// parallel list of values held by those tags.
func (p *Provider) CreatorTags() (tags []string, values []string) {
	if len(p.dcCreator) == 0 {
		return []string{"XMP  dc:creator"}, []string{""}
	}
	tags = make([]string, len(p.dcCreator))
	for i := range p.dcCreator {
		tags[i] = "XMP  dc:creator"
	}
	return tags, p.dcCreator
}

// SetCreator sets the value of the Creator field.
func (p *Provider) SetCreator(value string) error {
	if value == "" {
		p.dcCreator = nil
		if _, ok := p.rdf.Properties[creatorName]; ok {
			delete(p.rdf.Properties, creatorName)
			p.dirty = true
		}
		return nil
	}
	if len(p.dcCreator) == 1 && p.dcCreator[0] == value {
		return nil
	}
	p.dcCreator = []string{value}
	setSeq(p.rdf.Properties, creatorName, p.dcCreator)
	p.dirty = true
	return nil
}
