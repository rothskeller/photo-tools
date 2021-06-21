package iptc

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/rothskeller/photo-tools/metadata"
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

// Location returns the IPTC location.
func (p *IPTC) Location() (loc *metadata.Location) {
	loc = new(metadata.Location)
	if dset := p.findDSet(idCountryPLCode); dset != nil {
		if utf8.Valid(dset.data) {
			loc.CountryCode = strings.TrimSpace(string(dset.data))
		} else {
			p.problems = append(p.problems, "ignoring non-UTF8 IPTC Country/Primary Location Code")
		}
	}
	if dset := p.findDSet(idCountryPLName); dset != nil {
		if utf8.Valid(dset.data) {
			loc.CountryName = strings.TrimSpace(string(dset.data))
		} else {
			p.problems = append(p.problems, "ignoring non-UTF8 IPTC Country/Primary Location Name")
		}
	}
	if dset := p.findDSet(idProvinceState); dset != nil {
		if utf8.Valid(dset.data) {
			loc.State = strings.TrimSpace(string(dset.data))
		} else {
			p.problems = append(p.problems, "ignoring non-UTF8 IPTC Province/State")
		}
	}
	if dset := p.findDSet(idCity); dset != nil {
		if utf8.Valid(dset.data) {
			loc.City = strings.TrimSpace(string(dset.data))
		} else {
			p.problems = append(p.problems, "ignoring non-UTF8 IPTC City")
		}
	}
	if dset := p.findDSet(idSublocation); dset != nil {
		if utf8.Valid(dset.data) {
			loc.Sublocation = strings.TrimSpace(string(dset.data))
		} else {
			p.problems = append(p.problems, "ignoring non-UTF8 IPTC Sub-location")
		}
	}
	return loc
}

// SetLocation sets the IPTC location.
func (p *IPTC) SetLocation(loc *metadata.Location) {
	if p == nil {
		return
	}
	if !loc.Valid() {
		p.deleteDSet(idCountryPLCode)
		p.deleteDSet(idCountryPLName)
		p.deleteDSet(idProvinceState)
		p.deleteDSet(idCity)
		p.deleteDSet(idSublocation)
		return
	}
	p.setLocationPart(idCountryPLCode, MaxCountryPLCodeLen, loc.CountryCode)
	p.setLocationPart(idCountryPLName, MaxCountryPLNameLen, loc.CountryName)
	p.setLocationPart(idProvinceState, MaxProvinceStateLen, loc.State)
	p.setLocationPart(idCity, MaxCityLen, loc.City)
	p.setLocationPart(idSublocation, MaxSublocationLen, loc.Sublocation)
}
func (p *IPTC) setLocationPart(id uint16, max int, val string) {
	dset := p.findDSet(id)
	encoded := []byte(applyMax(val, max))
	if dset != nil {
		if !bytes.Equal(encoded, dset.data) {
			dset.data = encoded
			p.dirty = true
		}
	} else {
		p.dsets = append(p.dsets, &dsett{0, id, encoded})
		p.dirty = true
	}
}
