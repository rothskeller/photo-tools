package jpegifd0

import (
	"fmt"
)

const tagImageDescription uint16 = 0x10E

// getCaption reads the value of the Caption field from the RDF.
func (p *Provider) getCaption() (err error) {
	if tag := p.ifd.Tag(tagImageDescription); tag != nil {
		if p.imageDescription, err = tag.AsString(); err != nil {
			return fmt.Errorf("ImageDescription: %s", err)
		}
	}
	return nil
}

// Caption returns the value of the Caption field.
func (p *Provider) Caption() (value string) { return p.imageDescription }

// CaptionTags returns a list of tag names for the Caption field, and a
// parallel list of values held by those tags.
func (p *Provider) CaptionTags() (tags []string, values []string) {
	return []string{"IFD0 ImageDescription"}, []string{p.imageDescription}
}

// SetCaption sets the value of the Caption field.
func (p *Provider) SetCaption(value string) error {
	if value == "" {
		p.imageDescription = ""
		p.ifd.DeleteTag(tagImageDescription)
		return nil
	}
	if value == p.imageDescription {
		return nil
	}
	p.imageDescription = value
	p.ifd.AddTag(tagImageDescription).SetString(value)
	return nil
}
