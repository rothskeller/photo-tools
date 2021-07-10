package iptc

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/iim"
)

const (
	idCountryPLCode     uint16 = 0x0264
	idCountryPLName     uint16 = 0x0265
	idProvinceState     uint16 = 0x025F
	idCity              uint16 = 0x025A
	idSublocation       uint16 = 0x025C
	maxCountryPLCodeLen        = 3
	maxCountryPLNameLen        = 64
	maxProvinceStateLen        = 32
	maxCityLen                 = 32
	maxSublocationLen          = 32
)

// getLocation reads the value of the Location field from the IIM.
func (p *Provider) getLocation() (err error) {
	switch dss := p.iim[idCountryPLCode]; len(dss) {
	case 0:
		break
	case 1:
		if p.countryPLCode, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Country/Primary Location Code: %s", err)
		}
	default:
		return errors.New("Country/Primary Location Code: multiple data sets")
	}
	switch dss := p.iim[idCountryPLName]; len(dss) {
	case 0:
		break
	case 1:
		if p.countryPLName, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Country/Primary Location Name: %s", err)
		}
	default:
		return errors.New("Country/Primary Location Name: multiple data sets")
	}
	switch dss := p.iim[idProvinceState]; len(dss) {
	case 0:
		break
	case 1:
		if p.provinceState, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Province/State: %s", err)
		}
	default:
		return errors.New("Province/State: multiple data sets")
	}
	switch dss := p.iim[idCity]; len(dss) {
	case 0:
		break
	case 1:
		if p.city, err = getString(dss[0]); err != nil {
			return fmt.Errorf("City: %s", err)
		}
	default:
		return errors.New("City: multiple data sets")
	}
	switch dss := p.iim[idSublocation]; len(dss) {
	case 0:
		break
	case 1:
		if p.sublocation, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Sublocation: %s", err)
		}
	default:
		return errors.New("Sublocation: multiple data sets")
	}
	return nil
}

// Location returns the value of the Location field.
func (p *Provider) Location() (value metadata.Location) {
	return metadata.Location{
		CountryCode: p.countryPLCode,
		CountryName: p.countryPLName,
		State:       p.provinceState,
		City:        p.city,
		Sublocation: p.sublocation,
	}
}

// LocationTags returns a list of tag names for the Location field, and a
// parallel list of values held by those tags.
func (p *Provider) LocationTags() (tags []string, values []metadata.Location) {
	return []string{"IPTC (location tags)"}, []metadata.Location{p.Location()}
}

// SetLocation sets the value of the Location field.
func (p *Provider) SetLocation(value metadata.Location) error {
	if value.CountryCode == "" {
		p.countryPLCode = ""
		if _, ok := p.iim[idCountryPLCode]; ok {
			delete(p.iim, idCountryPLCode)
			p.dirty = true
		}
	} else {
		var v = value.CountryCode
		if len(v) > maxCountryPLCodeLen {
			v = v[:maxCountryPLCodeLen]
		}
		if v != p.countryPLCode {
			p.iim[idCountryPLCode] = []iim.DataSet{{ID: idCountryPLCode, Data: []byte(v)}}
			p.setEncoding()
			p.dirty = true
		}
	}
	if value.CountryName == "" {
		p.countryPLName = ""
		if _, ok := p.iim[idCountryPLName]; ok {
			delete(p.iim, idCountryPLName)
			p.dirty = true
		}
	} else {
		var v = value.CountryName
		if len(v) > maxCountryPLNameLen {
			v = v[:maxCountryPLNameLen]
		}
		if v != p.countryPLName {
			p.iim[idCountryPLName] = []iim.DataSet{{ID: idCountryPLName, Data: []byte(v)}}
			p.setEncoding()
			p.dirty = true
		}
	}
	if value.State == "" {
		p.provinceState = ""
		if _, ok := p.iim[idProvinceState]; ok {
			delete(p.iim, idProvinceState)
			p.dirty = true
		}
	} else {
		var v = value.State
		if len(v) > maxProvinceStateLen {
			v = v[:maxProvinceStateLen]
		}
		if v != p.provinceState {
			p.iim[idProvinceState] = []iim.DataSet{{ID: idProvinceState, Data: []byte(v)}}
			p.setEncoding()
			p.dirty = true
		}
	}
	if value.City == "" {
		p.city = ""
		if _, ok := p.iim[idCity]; ok {
			delete(p.iim, idCity)
			p.dirty = true
		}
	} else {
		var v = value.City
		if len(v) > maxCityLen {
			v = v[:maxCityLen]
		}
		if v != p.city {
			p.iim[idCity] = []iim.DataSet{{ID: idCity, Data: []byte(v)}}
			p.setEncoding()
			p.dirty = true
		}
	}
	if value.Sublocation == "" {
		p.sublocation = ""
		if _, ok := p.iim[idSublocation]; ok {
			delete(p.iim, idSublocation)
			p.dirty = true
		}
	} else {
		var v = value.Sublocation
		if len(v) > maxSublocationLen {
			v = v[:maxSublocationLen]
		}
		if v != p.sublocation {
			p.iim[idSublocation] = []iim.DataSet{{ID: idSublocation, Data: []byte(v)}}
			p.setEncoding()
			p.dirty = true
		}
	}
	return nil
}
