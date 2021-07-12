package iptc

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
)

const (
	idDateCreated         uint16 = 0x0237
	idTimeCreated         uint16 = 0x023C
	idDigitalCreationDate uint16 = 0x023E
	idDigitalCreationTime uint16 = 0x023F
)

// getDateTime reads the values of the DateTime field from the IIM.
func (p *Provider) getDateTime() (err error) {
	var date, time string

	switch dss := p.iim.DataSets(idDateCreated); len(dss) {
	case 0:
		break
	case 1:
		if date, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Date Created: %s", err)
		}
	default:
		return errors.New("Date Created: multiple data sets")
	}
	if date != "" {
		switch dss := p.iim.DataSets(idTimeCreated); len(dss) {
		case 0:
			break
		case 1:
			if time, err = getString(dss[0]); err != nil {
				return fmt.Errorf("Time Created: %s", err)
			}
		default:
			return errors.New("Time Created: multiple data sets")
		}
	}
	if err = p.dateTimeCreated.ParseIPTC(date, time); err != nil {
		return fmt.Errorf("Date/Time Created: %s", err)
	}
	switch dss := p.iim.DataSets(idDigitalCreationDate); len(dss) {
	case 0:
		date = ""
	case 1:
		if date, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Digital Creation Date: %s", err)
		}
	default:
		return errors.New("Digital Creation Date: multiple data sets")
	}
	if date != "" {
		switch dss := p.iim.DataSets(idDigitalCreationTime); len(dss) {
		case 0:
			time = ""
		case 1:
			if time, err = getString(dss[0]); err != nil {
				return fmt.Errorf("Digital Creation Time: %s", err)
			}
		default:
			return errors.New("Digital Creation Time: multiple data sets")
		}
	}
	if err = p.digitalCreationDateTime.ParseIPTC(date, time); err != nil {
		return fmt.Errorf("Digital Creation Date/Time: %s", err)
	}
	return nil
}

// DateTime returns the value of the DateTime field.
func (p *Provider) DateTime() (value metadata.DateTime) {
	if !p.dateTimeCreated.Empty() {
		return p.dateTimeCreated
	}
	return p.digitalCreationDateTime // which may be empty
}

// DateTimeTags returns a list of tag names for the DateTime field, and
// a parallel list of values held by those tags.
func (p *Provider) DateTimeTags() (tags []string, values []metadata.DateTime) {
	tags = append(tags, "IPTC Date/Time Created")
	values = append(values, p.dateTimeCreated)
	if !p.digitalCreationDateTime.Empty() {
		tags = append(tags, "IPTC Digital Creation Date/Time")
		values = append(values, p.digitalCreationDateTime)
	}
	return tags, values
}

// SetDateTime sets the value of the DateTime field.
func (p *Provider) SetDateTime(value metadata.DateTime) error {
	p.digitalCreationDateTime = metadata.DateTime{}
	p.iim.RemoveDataSets(idDigitalCreationDate)
	p.iim.RemoveDataSets(idDigitalCreationDate)
	if value.Empty() {
		p.dateTimeCreated = metadata.DateTime{}
		p.iim.RemoveDataSets(idDateCreated)
		p.iim.RemoveDataSets(idDateCreated)
		return nil
	}
	if value.Equivalent(p.dateTimeCreated) {
		return nil
	}
	p.dateTimeCreated = value
	date, time := value.AsIPTC()
	p.iim.SetDataSet(idDateCreated, []byte(date))
	p.iim.SetDataSet(idTimeCreated, []byte(time))
	return nil
}
