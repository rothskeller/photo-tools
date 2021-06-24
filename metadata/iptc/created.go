package iptc

const (
	idDateCreated uint16 = 0x0237
	idTimeCreated uint16 = 0x023C
)

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
	if err := p.DateTimeCreated.ParseIPTC(date, time); err != nil {
		p.log("DateTimeCreated: %s", err)
	}
}

func (p *IPTC) setDateTimeCreated() {
	if p.DateTimeCreated.Empty() {
		p.deleteDSet(idDateCreated)
		p.deleteDSet(idTimeCreated)
		return
	}
	date, time := p.DateTimeCreated.AsIPTC()
	p.setDSet(idDateCreated, []byte(date))
	p.setDSet(idTimeCreated, []byte(time))
}
