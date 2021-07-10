package exififd

import (
	"encoding/binary"
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/tiff"
)

// A Provider handles metadata in an EXIF IFD.
type Provider struct {
	metadata.BaseProvider
	dateTimeDigitized metadata.DateTime
	dateTimeOriginal  metadata.DateTime
	userComment       string

	ifd *tiff.IFD
	enc binary.ByteOrder
}

var _ metadata.Provider = (*Provider)(nil) // verify interface compliance

// New creates a new Provider based on the provided IFD.
func New(ifd *tiff.IFD, enc binary.ByteOrder) (p *Provider, err error) {
	p = &Provider{ifd: ifd, enc: enc}
	if err = p.getCaption(); err != nil {
		return nil, fmt.Errorf("EXIF IFD: %s", err)
	}
	if err = p.getDateTime(); err != nil {
		return nil, fmt.Errorf("EXIF IFD: %s", err)
	}
	return p, nil
}

// ProviderName is the name for the provider, for debug purposes.
func (p *Provider) ProviderName() string { return "EXIF IFD" }
