package gpsifd

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/tiff"
)

// A Provider handles metadata in a GPS IFD.
type Provider struct {
	metadata.BaseProvider
	gpsCoords metadata.GPSCoords

	ifd *tiff.IFD
}

var _ metadata.Provider = (*Provider)(nil) // verify interface compliance

// New creates a new Provider based on the provided IFD.
func New(ifd *tiff.IFD) (p *Provider, err error) {
	p = &Provider{ifd: ifd}
	if err = p.getGPS(); err != nil {
		return nil, fmt.Errorf("GPS IFD: %s", err)
	}
	return p, nil
}

// ProviderName is the name for the provider, for debug purposes.
func (p *Provider) ProviderName() string { return "GPS IFD" }
