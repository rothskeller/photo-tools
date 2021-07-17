package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var titleName = rdf.Name{Namespace: nsDC, Name: "title"}

// getTitle reads the value of the Title field from the RDF.
func (p *Provider) getTitle() (err error) {
	if p.dcTitle, err = getAlt(p.rdf.Property(titleName)); err != nil {
		return fmt.Errorf("dc:title: %s", err)
	}
	return nil
}

// Title returns the value of the Title field.
func (p *Provider) Title() (value string) { return p.dcTitle.Default() }

// TitleTags returns a list of tag names for the Title field, and a
// parallel list of values held by those tags.
func (p *Provider) TitleTags() (tags []string, values [][]string) {
	values = [][]string{nil}
	for _, as := range p.dcTitle {
		values[0] = append(values[0], as.Value)
	}
	return []string{"XMP  dc:title"}, values
}

// SetTitle sets the values of the Title field.
func (p *Provider) SetTitle(value string) error {
	if value == "" {
		p.dcTitle = nil
		p.rdf.RemoveProperty(titleName)
		return nil
	}
	if len(p.dcTitle) == 1 && value == p.dcTitle[0].Value {
		return nil
	}
	p.dcTitle = newAltString(value)
	p.rdf.SetProperty(titleName, makeAlt(p.dcTitle))
	return nil
}
