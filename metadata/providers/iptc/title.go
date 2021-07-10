package iptc

import (
	"errors"
	"fmt"
)

// MaxObjectNameLen is the maximum length of the Object Name entry.
const (
	idObjectName     uint16 = 0x0205
	maxObjectNameLen        = 64
)

// getTitle reads the value of the Title field from the RDF.
func (p *Provider) getTitle() (err error) {
	switch dss := p.iim[idObjectName]; len(dss) {
	case 0:
		break
	case 1:
		if p.objectName, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Object Name: %s", err)
		}
	default:
		return errors.New("Object Name: multiple data sets")
	}
	return nil
}

// Title returns the value of the Title field.
func (p *Provider) Title() (value string) { return p.objectName }

// TitleTags returns a list of tag names for the Title field, and a
// parallel list of values held by those tags.
func (p *Provider) TitleTags() (tags []string, values []string) {
	if p.objectName == "" {
		return nil, nil
	}
	return []string{"IPTC Object Name"}, []string{p.objectName}
}

// SetTitle sets the values of the Title field.
func (p *Provider) SetTitle(value string) error {
	p.objectName = ""
	if _, ok := p.iim[idObjectName]; ok {
		delete(p.iim, idObjectName)
		p.setEncoding()
		p.dirty = true
	}
	return nil
}
