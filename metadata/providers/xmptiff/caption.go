package xmptiff

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var imageDescriptionName = rdf.Name{Namespace: nsTIFF, Name: "ImageDescription"}

// getCaption reads the value of the Caption field from the RDF.
func (p *Provider) getCaption() (err error) {
	if p.tiffImageDescription, err = getAlt(p.rdf.Properties, imageDescriptionName); err != nil {
		return fmt.Errorf("tiff:ImageDescription: %s", err)
	}
	return nil
}

// Caption returns the value of the Caption field.
func (p *Provider) Caption() (value string) { return p.tiffImageDescription.Default() }

// CaptionTags returns a list of tag names for the Caption field, and a
// parallel list of values held by those tags.
func (p *Provider) CaptionTags() (tags []string, values [][]string) {
	vlist := make([]string, len(p.tiffImageDescription))
	for i := range p.tiffImageDescription {
		vlist[i] = p.tiffImageDescription[i].Value
	}
	return []string{"XML  tiff:ImageDescription"}, [][]string{vlist}
}

// SetCaption sets the value of the Caption field.
func (p *Provider) SetCaption(value string) error {
	if value == "" {
		p.tiffImageDescription = nil
		if _, ok := p.rdf.Properties[imageDescriptionName]; ok {
			delete(p.rdf.Properties, imageDescriptionName)
			p.dirty = true
		}
		return nil
	}
	if len(p.tiffImageDescription) == 1 && value == p.tiffImageDescription[0].Value {
		return nil
	}
	p.tiffImageDescription = newAltString(value)
	setAlt(p.rdf.Properties, imageDescriptionName, p.tiffImageDescription)
	p.dirty = true
	return nil
}
