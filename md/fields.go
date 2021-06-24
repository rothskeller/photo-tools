package main

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	strmeta "github.com/rothskeller/photo-tools/strmetao"
)

type field interface {
	name() string
	pluralName() string
	label() string
	multivalued() bool
	langtagged() bool
	newValue() metadata.Metadatum
	get(h filefmt.FileHandler) []metadata.Metadatum
	tags(h filefmt.FileHandler) []metadata.TaggedMetadatum
	set(h filefmt.FileHandler, values []metadata.Metadatum) error
}

type tagvalue struct {
	tag   string
	value metadata.Metadatum
}

type artistField struct{}

func (artistField) name() string                 { return "artist" }
func (artistField) pluralName() string           { return "artist" }
func (artistField) label() string                { return "Artist" }
func (artistField) multivalued() bool            { return false }
func (artistField) langtagged() bool             { return false }
func (artistField) newValue() metadata.Metadatum { return new(metadata.String) }
func (artistField) get(h filefmt.FileHandler) []metadata.Metadatum {
	artist := strmeta.GetArtist(h)
	if !artist.Empty() {
		return []metadata.Metadatum{artist}
	}
	return nil
}
func (artistField) tags(h filefmt.FileHandler) []metadata.TaggedMetadatum {
	return nil
}
func (artistField) set(h filefmt.FileHandler, values []metadata.Metadatum) error {
	panic("not implemented")
}

var _ field = artistField{}
