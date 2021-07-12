package tiffifd0

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/tiff"
)

// A Provider handles metadata in the IFD0 of a TIFF file.
type Provider struct {
	metadata.BaseProvider
	artist           string
	dateTime         metadata.DateTime
	imageDescription string

	ifd *tiff.IFD
}

var _ metadata.Provider = (*Provider)(nil) // verify interface compliance

// New creates a new Provider based on the provided IFD.
func New(ifd *tiff.IFD) (p *Provider, err error) {
	p = &Provider{ifd: ifd}
	if err = p.getCaption(); err != nil {
		return nil, fmt.Errorf("JPEG IFD0: %s", err)
	}
	if err = p.getCreator(); err != nil {
		return nil, fmt.Errorf("JPEG IFD0: %s", err)
	}
	if err = p.getDateTime(); err != nil {
		return nil, fmt.Errorf("JPEG IFD0: %s", err)
	}
	return p, nil
}

// ProviderName is the name for the provider, for debug purposes.
func (p *Provider) ProviderName() string { return "TIFF IFD0" }
