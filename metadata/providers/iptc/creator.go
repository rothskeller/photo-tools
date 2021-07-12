package iptc

import (
	"fmt"
)

const (
	idByline     uint16 = 0x0250
	maxBylineLen        = 32
)

// getCreator reads the values of the Creator field from the IIM.
func (p *Provider) getCreator() (err error) {
	for _, ds := range p.iim.DataSets(idByline) {
		if byline, err := getString(ds); err == nil {
			p.bylines = append(p.bylines, byline)
		} else {
			return fmt.Errorf("By-line: %s", err)
		}
	}
	return nil
}

// Creator returns the value of the Creator field.
func (p *Provider) Creator() (value string) {
	if len(p.bylines) == 0 {
		return ""
	}
	return p.bylines[0]
}

// CreatorTags returns a list of tag names for the Creator field, and a
// parallel list of values held by those tags.
func (p *Provider) CreatorTags() (tags []string, values [][]string) {
	return []string{"IPTC  By-line"}, [][]string{p.bylines}
}

// SetCreator sets the value of the Creator field.
func (p *Provider) SetCreator(value string) error {
	if value == "" {
		p.bylines = nil
		p.iim.RemoveDataSets(idByline)
		return nil
	}
	if len(value) > maxBylineLen {
		value = value[:maxBylineLen]
	}
	if len(p.bylines) == 1 && p.bylines[0] == value {
		return nil
	}
	p.bylines = []string{value}
	p.iim.SetDataSet(idByline, []byte(value))
	p.setEncoding()
	return nil
}
