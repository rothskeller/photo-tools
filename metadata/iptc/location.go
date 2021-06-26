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

func (p *IPTC) getLocation() {
	if dset := p.findDSet(idCountryPLCode); dset != nil {
		p.CountryPLCode = strings.TrimSpace(p.decodeString(dset.data, "CountryPLCode"))
		p.saveCountryPLCode = p.CountryPLCode
	}
	if dset := p.findDSet(idCountryPLName); dset != nil {
		p.CountryPLName = strings.TrimSpace(p.decodeString(dset.data, "CountryPLName"))
		p.saveCountryPLName = p.CountryPLName
	}
	if dset := p.findDSet(idProvinceState); dset != nil {
		p.ProvinceState = strings.TrimSpace(p.decodeString(dset.data, "ProvinceState"))
		p.saveProvinceState = p.ProvinceState
	}
	if dset := p.findDSet(idCity); dset != nil {
		p.City = strings.TrimSpace(p.decodeString(dset.data, "City"))
		p.saveCity = p.City
	}
	if dset := p.findDSet(idSublocation); dset != nil {
		p.Sublocation = strings.TrimSpace(p.decodeString(dset.data, "Sublocation"))
		p.saveSublocation = p.Sublocation
	}
}

func (p *IPTC) setLocation() {
	if stringEqualMax(p.CountryPLCode, p.saveCountryPLCode, MaxCountryPLCodeLen) &&
		stringEqualMax(p.CountryPLName, p.saveCountryPLName, MaxCountryPLNameLen) &&
		stringEqualMax(p.ProvinceState, p.saveProvinceState, MaxProvinceStateLen) &&
		stringEqualMax(p.City, p.saveCity, MaxCityLen) &&
		stringEqualMax(p.Sublocation, p.saveSublocation, MaxSublocationLen) {
		return
	}
	if p.CountryPLCode != "" {
		p.setDSet(idCountryPLCode, []byte(applyMax(p.CountryPLCode, MaxCountryPLCodeLen)))
	} else {
		p.deleteDSet(idCountryPLCode)
	}
	if p.CountryPLName != "" {
		p.setDSet(idCountryPLName, []byte(applyMax(p.CountryPLName, MaxCountryPLNameLen)))
	} else {
		p.deleteDSet(idCountryPLName)
	}
	if p.ProvinceState != "" {
		p.setDSet(idProvinceState, []byte(applyMax(p.ProvinceState, MaxProvinceStateLen)))
	} else {
		p.deleteDSet(idProvinceState)
	}
	if p.City != "" {
		p.setDSet(idCity, []byte(applyMax(p.City, MaxCityLen)))
	} else {
		p.deleteDSet(idCity)
	}
	if p.Sublocation != "" {
		p.setDSet(idSublocation, []byte(applyMax(p.Sublocation, MaxSublocationLen)))
	} else {
		p.deleteDSet(idSublocation)
	}
}
