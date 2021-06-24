// Package iptc handles IPTC metadata blocks.
package iptc

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// IPTC is a an IPTC parser and generator.
type IPTC struct {
	Bylines                 []string
	CaptionAbstract         string
	DateTimeCreated         metadata.DateTime
	DigitalCreationDateTime metadata.DateTime
	Keywords                []string
	Location                metadata.Location
	ObjectName              string
	Problems                []string

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
