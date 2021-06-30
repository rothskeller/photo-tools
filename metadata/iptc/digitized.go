package iptc

import "github.com/rothskeller/photo-tools/metadata"

const (
	idDigitalCreationDate uint16 = 0x023E
	idDigitalCreationTime uint16 = 0x023F
)

// DigitalCreationDateTime returns the values of the Digital Creation Date and
// Digital Creation Time tags.
func (p *IPTC) DigitalCreationDateTime() metadata.DateTime { return p.digitalCreationDateTime }

func (p *IPTC) getDigitalCreationDateTime() {
	var date, time string

	if datedset := p.findDSet(idDigitalCreationDate); datedset != nil {
		date = string(datedset.data)
	} else {
		return
	}
	if timedset := p.findDSet(idDigitalCreationTime); timedset != nil {
		time = string(timedset.data)
	}
	if err := p.digitalCreationDateTime.ParseIPTC(date, time); err != nil {
		p.log("DigitalCreationDateTime: %s", err)
	}
}

// SetDigitalCreationDateTime sets the values of the Digital Creation Date and
// Digital Creation Time tags.
func (p *IPTC) SetDigitalCreationDateTime(v metadata.DateTime) error {
	if v.Equivalent(p.digitalCreationDateTime) {
		return nil
	}
	p.digitalCreationDateTime = v
	if p.digitalCreationDateTime.Empty() {
		p.deleteDSet(idDigitalCreationDate)
		p.deleteDSet(idDigitalCreationTime)
		return nil
	}
	date, time := p.digitalCreationDateTime.AsIPTC()
	p.setDSet(idDigitalCreationDate, []byte(date))
	p.setDSet(idDigitalCreationTime, []byte(time))
	return nil
}
