package iptc

import (
	"strings"
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

// CountryPLCode returns the value of the Country/Primary Location Code tag.
func (p *IPTC) CountryPLCode() string { return p.countryPLCode }

// CountryPLName returns the value of the Country/Primary Location Name tag.
func (p *IPTC) CountryPLName() string { return p.countryPLName }

// ProvinceState returns the value of the Province/State tag.
func (p *IPTC) ProvinceState() string { return p.provinceState }

// City returns the value of the City tag.
func (p *IPTC) City() string { return p.city }

// Sublocation returns the value of the Sub-location tag.
func (p *IPTC) Sublocation() string { return p.sublocation }

func (p *IPTC) getLocation() {
	if dset := p.findDSet(idCountryPLCode); dset != nil {
		p.countryPLCode = strings.TrimSpace(p.decodeString(dset.data, "CountryPLCode"))
	}
	if dset := p.findDSet(idCountryPLName); dset != nil {
		p.countryPLName = strings.TrimSpace(p.decodeString(dset.data, "CountryPLName"))
	}
	if dset := p.findDSet(idProvinceState); dset != nil {
		p.provinceState = strings.TrimSpace(p.decodeString(dset.data, "ProvinceState"))
	}
	if dset := p.findDSet(idCity); dset != nil {
		p.city = strings.TrimSpace(p.decodeString(dset.data, "City"))
	}
	if dset := p.findDSet(idSublocation); dset != nil {
		p.sublocation = strings.TrimSpace(p.decodeString(dset.data, "Sublocation"))
	}
}

// SetCountryPLCode sets the value of the Country/Primary Location Code tag.
func (p *IPTC) SetCountryPLCode(v string) error {
	if stringEqualMax(v, p.countryPLCode, MaxCountryPLCodeLen) {
		return nil
	}
	p.countryPLCode = applyMax(v, MaxCountryPLCodeLen)
	if p.countryPLCode == "" {
		p.deleteDSet(idCountryPLCode)
	} else {
		p.setDSet(idCountryPLCode, []byte(p.countryPLCode))
	}
	return nil
}

// SetCountryPLName sets the value of the Country/Primary Location Name tag.
func (p *IPTC) SetCountryPLName(v string) error {
	if stringEqualMax(v, p.countryPLName, MaxCountryPLNameLen) {
		return nil
	}
	p.countryPLName = applyMax(v, MaxCountryPLNameLen)
	if p.countryPLName == "" {
		p.deleteDSet(idCountryPLName)
	} else {
		p.setDSet(idCountryPLName, []byte(p.countryPLName))
	}
	return nil
}

// SetProvinceState sets the value of the Province/State tag.
func (p *IPTC) SetProvinceState(v string) error {
	if stringEqualMax(v, p.provinceState, MaxProvinceStateLen) {
		return nil
	}
	p.provinceState = applyMax(v, MaxProvinceStateLen)
	if p.provinceState == "" {
		p.deleteDSet(idProvinceState)
	} else {
		p.setDSet(idProvinceState, []byte(p.provinceState))
	}
	return nil
}

// SetCity sets the value of the City tag.
func (p *IPTC) SetCity(v string) error {
	if stringEqualMax(v, p.city, MaxCityLen) {
		return nil
	}
	p.city = applyMax(v, MaxCityLen)
	if p.city == "" {
		p.deleteDSet(idCity)
	} else {
		p.setDSet(idCity, []byte(p.city))
	}
	return nil
}

// SetSublocation sets the value of the Sub-location tag.
func (p *IPTC) SetSublocation(v string) error {
	if stringEqualMax(v, p.sublocation, MaxSublocationLen) {
		return nil
	}
	p.sublocation = applyMax(v, MaxSublocationLen)
	if p.sublocation == "" {
		p.deleteDSet(idSublocation)
	} else {
		p.setDSet(idSublocation, []byte(p.sublocation))
	}
	return nil
}
