package iptc

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/iim"
)

const (
	idByline     uint16 = 0x0250
	maxBylineLen        = 32
)

// getCreator reads the values of the Creator field from the IIM.
func (p *Provider) getCreator() (err error) {
	for _, ds := range p.iim[idByline] {
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
func (p *Provider) CreatorTags() (tags []string, values []string) {
	if len(p.bylines) == 0 {
		return []string{"IPTC  By-line"}, []string{""}
	}
	tags = make([]string, len(p.bylines))
	for i := range p.bylines {
		tags[i] = "IPTC By-Line"
	}
	return tags, p.bylines
}

// SetCreator sets the value of the Creator field.
func (p *Provider) SetCreator(value string) error {
	if value == "" {
		p.bylines = nil
		if _, ok := p.iim[idByline]; ok {
			delete(p.iim, idByline)
			p.dirty = true
		}
		return nil
	}
	if len(value) > maxBylineLen {
		value = value[:maxBylineLen]
	}
	if len(p.bylines) == 1 && p.bylines[0] == value {
		return nil
	}
	p.bylines = []string{value}
	p.iim[idByline] = []iim.DataSet{{ID: idByline, Data: []byte(value)}}
	p.setEncoding()
	p.dirty = true
	return nil
}
