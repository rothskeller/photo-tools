// Package iptc handles IPTC metadata blocks.
package iptc

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// IPTC is a an IPTC parser and generator.
type IPTC struct {
	Bylines                 []string
	CaptionAbstract         string
	City                    string
	CountryPLCode           string
	CountryPLName           string
	DateTimeCreated         metadata.DateTime
	DigitalCreationDateTime metadata.DateTime
	Keywords                []string
	ObjectName              string
	ProvinceState           string
	Sublocation             string
	Problems                []string

	saveBylines                 []string
	saveCaptionAbstract         string
	saveCity                    string
	saveCountryPLCode           string
	saveCountryPLName           string
	saveDateTimeCreated         metadata.DateTime
	saveDigitalCreationDateTime metadata.DateTime
	saveKeywords                []string
	saveObjectName              string
	saveProvinceState           string
	saveSublocation             string

	offset uint32
	buf    []byte
	psir   []*psirt
	dsets  []*dsett
	dirty  bool
}

type psirt struct {
	offset uint32
	id     uint16
	name   string
	buf    []byte
}

type dsett struct {
	offset uint32
	id     uint16
	data   []byte
}

func stringEqualMax(a, b string, max int) bool {
	if a == b {
		return true
	}
	if len(a) == max && len(b) > max && a == b[:max] {
		return true
	}
	if len(b) == max && len(a) > max && b == a[:max] {
		return true
	}
	return false
}

func stringSliceEqualMax(a, b []string, max int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !stringEqualMax(a[i], b[i], max) {
			return false
		}
	}
	return true
}
