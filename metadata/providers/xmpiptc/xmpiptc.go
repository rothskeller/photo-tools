package xmpiptc

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

const (
	pfxIPTC = "Iptc4xmpExt"
	nsIPTC  = "http://iptc.org/std/Iptc4xmpExt/2008-02-29/"
)

// A Provider handles IPTC data mirrored in an XMP/RDF block.
type Provider struct {
	metadata.BaseProvider
	iptcLocationCreated location
	iptcLocationsShown  []location

	rdf   *rdf.Packet
	dirty bool
}

var _ metadata.Provider = (*Provider)(nil) // verify interface compliance

// New creates a new Provider based on the provided RDF block.
func New(rdf *rdf.Packet) (p *Provider, err error) {
	p = &Provider{rdf: rdf}
	p.rdf.RegisterNamespace(pfxIPTC, nsIPTC)
	if err = p.getLocation(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	return p, nil
}

// ProviderName is the name for the provider, for debug purposes.
func (p *Provider) ProviderName() string { return "XMP:iptc" }
