package xmpexif

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

const (
	pfxEXIF = "exif"
	nsEXIF  = "http://ns.adobe.com/exif/1.0/"
)

// A Provider handles EXIF data mirrored into an XMP/RDF block.
type Provider struct {
	metadata.BaseProvider
	exifDateTimeOriginal  metadata.DateTime
	exifDateTimeDigitized metadata.DateTime
	exifGPSCoords         metadata.GPSCoords
	exifUserComment       altString

	rdf   *rdf.Packet
	dirty bool
}

var _ metadata.Provider = (*Provider)(nil) // verify interface compliance

// New creates a new Provider based on the provided RDF block.
func New(rdf *rdf.Packet) (p *Provider, err error) {
	p = &Provider{rdf: rdf}
	p.rdf.RegisterNamespace(pfxEXIF, nsEXIF)
	if err = p.getCaption(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getDateTime(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getGPS(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	return p, nil
}

// ProviderName is the name for the provider, for debug purposes.
func (p *Provider) ProviderName() string { return "XMP:exif" }
