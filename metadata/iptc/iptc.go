// Package iptc handles IPTC metadata blocks.
package iptc

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// IPTC is a an IPTC parser and generator.
type IPTC struct {
	Bylines                 []*metadata.String
	CaptionAbstract         *metadata.String
	DateTimeCreated         *metadata.DateTime
	DigitalCreationDateTime *metadata.DateTime
	Keywords                []*metadata.String
	Location                *metadata.Location
	ObjectName              *metadata.String
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
