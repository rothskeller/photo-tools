package iptc

import (
	"bytes"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/iim"
	"golang.org/x/text/encoding/charmap"
)

const idCodedCharacterSet uint16 = 0x015A

// A Provider handles data from an IPTC IIM block.
type Provider struct {
	metadata.BaseProvider
	bylines                 []string
	captionAbstract         string
	city                    string
	countryPLCode           string
	countryPLName           string
	dateTimeCreated         metadata.DateTime
	digitalCreationDateTime metadata.DateTime
	keywords                []string
	objectName              string
	provinceState           string
	sublocation             string

	iim   iim.IIM
	dirty bool
}

var utf8Escape1 = []byte{0x1B, 0x25, 0x47}
var utf8Escape2 = []byte{0x1B, 0x25, 0x2F, 0x49}

var _ metadata.Provider = (*Provider)(nil) // verify interface compliance

// New creates a new Provider based on the provided IIM block.
func New(iim iim.IIM) (p *Provider, err error) {
	p = &Provider{iim: iim}
	if err = p.verifyEncoding(); err != nil {
		return nil, fmt.Errorf("IPTC IIM: %s", err)
	}
	if err = p.getCaption(); err != nil {
		return nil, fmt.Errorf("IPTC IIM: %s", err)
	}
	if err = p.getCreator(); err != nil {
		return nil, fmt.Errorf("IPTC IIM: %s", err)
	}
	if err = p.getDateTime(); err != nil {
		return nil, fmt.Errorf("IPTC IIM: %s", err)
	}
	if err = p.getKeywords(); err != nil {
		return nil, fmt.Errorf("IPTC IIM: %s", err)
	}
	if err = p.getLocation(); err != nil {
		return nil, fmt.Errorf("IPTC IIM: %s", err)
	}
	if err = p.getTitle(); err != nil {
		return nil, fmt.Errorf("IPTC IIM: %s", err)
	}
	return p, nil
}

// ProviderName is the name for the provider, for debug purposes.
func (p *Provider) ProviderName() string { return "IPTC" }

func (p *Provider) verifyEncoding() (err error) {
	switch dss := p.iim[idCodedCharacterSet]; len(dss) {
	case 0:
		break
	case 1:
		if !bytes.Equal(dss[0].Data, utf8Escape1) && !bytes.Equal(dss[0].Data, utf8Escape2) {
			return errors.New("Coded Character Set: not UTF-8")
		}
	default:
		return errors.New("Coded Character Set: multiple data sets")
	}
	return nil
}

// setEncoding adds the record that defines the encoding as UTF-8.  It is called
// whenever a data set with a string value is changed.
func (p *Provider) setEncoding() {
	p.iim[idCodedCharacterSet] = []iim.DataSet{{ID: idCodedCharacterSet, Data: utf8Escape1}}
}

func getString(ds iim.DataSet) (string, error) {
	if utf8.Valid(ds.Data) {
		return string(ds.Data), nil
	}
	if by, err := charmap.ISO8859_1.NewDecoder().Bytes(ds.Data); err == nil {
		return string(by), nil
	}
	return "", errors.New("can't determine character encoding")
}
