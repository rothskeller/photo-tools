package iptc

import "github.com/rothskeller/photo-tools/metadata"

const (
	idDateCreated uint16 = 0x0237
	idTimeCreated uint16 = 0x023C
)

// DateTimeCreated returns the values of the Date Created and Time Created tags.
func (p *IPTC) DateTimeCreated() metadata.DateTime { return p.dateTimeCreated }

func (p *IPTC) getDateTimeCreated() {
	var date, time string

	if datedset := p.findDSet(idDateCreated); datedset != nil {
		date = string(datedset.data)
	} else {
		return
	}
	if timedset := p.findDSet(idTimeCreated); timedset != nil {
		time = string(timedset.data)
	}
	if err := p.dateTimeCreated.ParseIPTC(date, time); err != nil {
		p.log("DateTimeCreated: %s", err)
	}
}

// SetDateTimeCreated sets the values of the Date Created and Time Created tags.
func (p *IPTC) SetDateTimeCreated(v metadata.DateTime) error {
	if v.Equivalent(p.dateTimeCreated) {
		return nil
	}
	p.dateTimeCreated = v
	if p.dateTimeCreated.Empty() {
		p.deleteDSet(idDateCreated)
		p.deleteDSet(idTimeCreated)
		return nil
	}
	date, time := p.dateTimeCreated.AsIPTC()
	p.setDSet(idDateCreated, []byte(date))
	p.setDSet(idTimeCreated, []byte(time))
	return nil
}
