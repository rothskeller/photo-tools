package iptc

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
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
	switch dss := p.iim.DataSets(idCountryPLCode); len(dss) {
	case 0:
		break
	case 1:
		if p.countryPLCode, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Country/Primary Location Code: %s", err)
		}
	default:
		return errors.New("Country/Primary Location Code: multiple data sets")
	}
	switch dss := p.iim.DataSets(idCountryPLName); len(dss) {
	case 0:
		break
	case 1:
		if p.countryPLName, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Country/Primary Location Name: %s", err)
		}
	default:
		return errors.New("Country/Primary Location Name: multiple data sets")
	}
	switch dss := p.iim.DataSets(idProvinceState); len(dss) {
	case 0:
		break
	case 1:
		if p.provinceState, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Province/State: %s", err)
		}
	default:
		return errors.New("Province/State: multiple data sets")
	}
	switch dss := p.iim.DataSets(idCity); len(dss) {
	case 0:
		break
	case 1:
		if p.city, err = getString(dss[0]); err != nil {
			return fmt.Errorf("City: %s", err)
		}
	default:
		return errors.New("City: multiple data sets")
	}
	switch dss := p.iim.DataSets(idSublocation); len(dss) {
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
func (p *Provider) LocationTags() (tags []string, values [][]metadata.Location) {
	return []string{"IPTC (location tags)"}, [][]metadata.Location{{p.Location()}}
}

// SetLocation sets the value of the Location field.
func (p *Provider) SetLocation(value metadata.Location) error {
	if value.CountryCode == "" {
		p.countryPLCode = ""
		p.iim.RemoveDataSets(idCountryPLCode)
	} else {
		var v = value.CountryCode
		if len(v) > maxCountryPLCodeLen {
			v = v[:maxCountryPLCodeLen]
		}
		if v != p.countryPLCode {
			p.iim.SetDataSet(idCountryPLCode, []byte(v))
			p.setEncoding()
		}
	}
	if value.CountryName == "" {
		p.countryPLName = ""
		p.iim.RemoveDataSets(idCountryPLName)
	} else {
		var v = value.CountryName
		if len(v) > maxCountryPLNameLen {
			v = v[:maxCountryPLNameLen]
		}
		if v != p.countryPLName {
			p.iim.SetDataSet(idCountryPLName, []byte(v))
			p.setEncoding()
		}
	}
	if value.State == "" {
		p.provinceState = ""
		p.iim.RemoveDataSets(idProvinceState)
	} else {
		var v = value.State
		if len(v) > maxProvinceStateLen {
			v = v[:maxProvinceStateLen]
		}
		if v != p.provinceState {
			p.iim.SetDataSet(idProvinceState, []byte(v))
			p.setEncoding()
		}
	}
	if value.City == "" {
		p.city = ""
		p.iim.RemoveDataSets(idCity)
	} else {
		var v = value.City
		if len(v) > maxCityLen {
			v = v[:maxCityLen]
		}
		if v != p.city {
			p.iim.SetDataSet(idCity, []byte(v))
			p.setEncoding()
		}
	}
	if value.Sublocation == "" {
		p.sublocation = ""
		p.iim.RemoveDataSets(idSublocation)
	} else {
		var v = value.Sublocation
		if len(v) > maxSublocationLen {
			v = v[:maxSublocationLen]
		}
		if v != p.sublocation {
			p.iim.SetDataSet(idSublocation, []byte(v))
			p.setEncoding()
		}
	}
	return nil
}
