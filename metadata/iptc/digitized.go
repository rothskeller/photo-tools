package iptc

const (
	idDigitalCreationDate uint16 = 0x023E
	idDigitalCreationTime uint16 = 0x023F
)

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
	if err := p.DigitalCreationDateTime.ParseIPTC(date, time); err != nil {
		p.log("DigitalCreationDateTime: %s", err)
	}
	p.saveDigitalCreationDateTime = p.DigitalCreationDateTime
}

func (p *IPTC) setDigitalCreationDateTime() {
	if p.saveDigitalCreationDateTime.Equal(&p.DigitalCreationDateTime) {
		return
	}
	if p.DigitalCreationDateTime.Empty() {
		p.deleteDSet(idDigitalCreationDate)
		p.deleteDSet(idDigitalCreationTime)
		return
	}
	date, time := p.DigitalCreationDateTime.AsIPTC()
	p.setDSet(idDigitalCreationDate, []byte(date))
	p.setDSet(idDigitalCreationTime, []byte(time))
}
