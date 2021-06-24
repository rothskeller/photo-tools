package iptc

import (
	"strings"
	"unicode/utf8"
)

const (
	idCountryPLCode uint16 = 0x0264
	idCountryPLName uint16 = 0x0265
	idProvinceState uint16 = 0x025F
	idCity          uint16 = 0x025A
	idSublocation   uint16 = 0x025C
)

// Maximum lengths of various fields.
const (
	MaxCountryPLCodeLen = 3
	MaxCountryPLNameLen = 64
	MaxProvinceStateLen = 32
	MaxCityLen          = 32
	MaxSublocationLen   = 32
)

func (p *IPTC) getLocation() {
	var ccode, cname, state, city, subloc string

	if dset := p.findDSet(idCountryPLCode); dset != nil {
		if utf8.Valid(dset.data) {
			ccode = strings.TrimSpace(string(dset.data))
		} else {
			p.log("ignoring non-UTF8 Country/Primary Location Code")
			return
		}
	}
	if dset := p.findDSet(idCountryPLName); dset != nil {
		if utf8.Valid(dset.data) {
			cname = strings.TrimSpace(string(dset.data))
		} else {
			p.log("ignoring non-UTF8 IPTC Country/Primary Location Name")
		}
	}
	if dset := p.findDSet(idProvinceState); dset != nil {
		if utf8.Valid(dset.data) {
			state = strings.TrimSpace(string(dset.data))
		} else {
			p.log("ignoring non-UTF8 IPTC Province/State")
		}
	}
	if dset := p.findDSet(idCity); dset != nil {
		if utf8.Valid(dset.data) {
			city = strings.TrimSpace(string(dset.data))
		} else {
			p.log("ignoring non-UTF8 IPTC City")
		}
	}
	if dset := p.findDSet(idSublocation); dset != nil {
		if utf8.Valid(dset.data) {
			subloc = strings.TrimSpace(string(dset.data))
		} else {
			p.log("ignoring non-UTF8 IPTC Sub-location")
		}
	}
	if err := p.Location.ParseComponents(ccode, cname, state, city, subloc); err != nil {
		p.log("Location: %s", err)
	}
}

func (p *IPTC) setLocation() {
	if p.Location.Empty() {
		p.deleteDSet(idCountryPLCode)
		p.deleteDSet(idCountryPLName)
		p.deleteDSet(idProvinceState)
		p.deleteDSet(idCity)
		p.deleteDSet(idSublocation)
		return
	}
	p.setDSet(idCountryPLCode, []byte(applyMax(p.Location.CountryCode, MaxCountryPLCodeLen)))
	p.setDSet(idCountryPLName, []byte(applyMax(p.Location.CountryName, MaxCountryPLNameLen)))
	p.setDSet(idProvinceState, []byte(applyMax(p.Location.State, MaxProvinceStateLen)))
	p.setDSet(idCity, []byte(applyMax(p.Location.City, MaxCityLen)))
	p.setDSet(idSublocation, []byte(applyMax(p.Location.Sublocation, MaxSublocationLen)))
}
