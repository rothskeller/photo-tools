// Package xmpext contains the "provider" for extended XMP blocks in JPEG (or
// similar) files.  It doesn't actually provide anything; its purpose is to
// ensure that such blocks don't contain any metadata we care about, since we
// can't handle them properly.
package xmpext

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var prohibitedNamespaces = []string{
	"http://purl.org/dc/elements/1.1/",                     // dc
	"http://www.digikam.org/ns/1.0/",                       // digiKam
	"http://ns.adobe.com/exif/1.0/",                        // exif
	"http://iptc.org/std/Iptc4xmpExt/2008-02-29/",          // Iptc4xmpExt
	"http://ns.adobe.com/lightroom/1.0/",                   // lr
	"http://ns.microsoft.com/photo/1.2/",                   // MP
	"http://ns.microsoft.com/photo/1.2/t/RegionInfo#",      // MPRI
	"http://ns.microsoft.com/photo/1.2/t/Region#",          // MPReg
	"http://www.metadataworkinggroup.com/schemas/regions/", // mwg-rs
	"http://ns.adobe.com/photoshop/1.0/",                   // photoshop
	"http://ns.adobe.com/tiff/1.0/",                        // tiff
	"http://ns.adobe.com/xap/1.0/",                         // xmp
}

// A Provider handles data from an XMP/RDF extension metadata block.
type Provider struct {
	metadata.BaseProvider
}

var _ metadata.Provider = (*Provider)(nil) // verify interface compliance

// New creates a new Provider based on the provided RDF block.
func New(rdf *rdf.Packet) (p *Provider, err error) {
	var pnmap = make(map[string]bool)

	for _, pn := range prohibitedNamespaces {
		pnmap[pn] = true
	}
	for _, name := range rdf.Properties() {
		if pnmap[name.Namespace] {
			return nil, errors.New("XMPExt: prohibited namespace")
		}
	}
	return new(Provider), nil
}

// ProviderName is the name for the provider, for debug purposes.
func (p *Provider) ProviderName() string { return "XMP Extension" }
