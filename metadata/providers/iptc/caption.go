package iptc

import (
	"errors"
	"fmt"
)

const (
	idCaptionAbstract     uint16 = 0x0278
	maxCaptionAbstractLen        = 2000
)

// getCaption reads the value of the Caption field from the IIM.
func (p *Provider) getCaption() (err error) {
	switch dss := p.iim.DataSets(idCaptionAbstract); len(dss) {
	case 0:
		break
	case 1:
		if p.captionAbstract, err = getString(dss[0]); err != nil {
			return fmt.Errorf("Caption/Abstract: %s", err)
		}
	default:
		return errors.New("Caption/Abstract: multiple data sets")
	}
	return nil
}

// Caption returns the value of the Caption field.
func (p *Provider) Caption() (value string) { return p.captionAbstract }

// CaptionTags returns a list of tag names for the Caption field, and a
// parallel list of values held by those tags.
func (p *Provider) CaptionTags() (tags []string, values [][]string) {
	if p.captionAbstract == "" {
		return []string{"IPTC Caption/Abstract"}, [][]string{nil}
	}
	return []string{"IPTC Caption/Abstract"}, [][]string{{p.captionAbstract}}
}

// SetCaption sets the value of the Caption field.
func (p *Provider) SetCaption(value string) error {
	if value == "" {
		p.captionAbstract = ""
		p.iim.RemoveDataSets(idCaptionAbstract)
		return nil
	}
	if len(value) > maxCaptionAbstractLen {
		value = value[:maxCaptionAbstractLen]
	}
	if value == p.captionAbstract {
		return nil
	}
	p.captionAbstract = value
	p.iim.SetDataSet(idCaptionAbstract, []byte(value))
	p.setEncoding()
	return nil
}
