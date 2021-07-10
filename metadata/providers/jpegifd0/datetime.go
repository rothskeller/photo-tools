package jpegifd0

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
)

const tagDateTime uint16 = 0x132

// getDateTime reads the value of the DateTime field from the RDF.
func (p *Provider) getDateTime() (err error) {
	dtot := p.ifd.Tag(tagDateTime)
	if dtot == nil {
		return nil
	}
	var dto string
	if dto, err = dtot.AsString(); err != nil {
		return fmt.Errorf("DateTime: %s", err)
	}
	if err = p.dateTime.ParseEXIF(dto, "", ""); err != nil {
		return fmt.Errorf("DateTime: %s", err)
	}
	return nil
}

// DateTime returns the value of the DateTime field.
func (p *Provider) DateTime() (value metadata.DateTime) { return p.dateTime }

// DateTimeTags returns a list of tag names for the DateTime field, and
// a parallel list of values held by those tags.
func (p *Provider) DateTimeTags() (tags []string, values []metadata.DateTime) {
	if p.dateTime.Empty() {
		return nil, nil
	}
	return []string{"IFD0 DateTime"}, []metadata.DateTime{p.dateTime}
}

// SetDateTime sets the value of the DateTime field.
func (p *Provider) SetDateTime(value metadata.DateTime) error {
	p.dateTime = metadata.DateTime{}
	p.ifd.DeleteTag(tagDateTime)
	return nil
}
