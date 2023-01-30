package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

const (
	pfxDC      = "dc"
	nsDC       = "http://purl.org/dc/elements/1.1/"
	pfxDigiKam = "digiKam"
	nsDigiKam  = "http://www.digikam.org/ns/1.0/"
	pfxEXIF    = "exif"
	nsEXIF     = "http://ns.adobe.com/exif/1.0/"
	pfxIPTC    = "Iptc4xmpExt"
	nsIPTC     = "http://iptc.org/std/Iptc4xmpExt/2008-02-29/"
	pfxLR      = "lr"
	nsLR       = "http://ns.adobe.com/lightroom/1.0/"
	pfxMP      = "MP"
	nsMP       = "http://ns.microsoft.com/photo/1.2/"
	pfxMPRI    = "MPRI"
	nsMPRI     = "http://ns.microsoft.com/photo/1.2/t/RegionInfo#"
	pfxMPReg   = "MPReg"
	nsMPReg    = "http://ns.microsoft.com/photo/1.2/t/Region#"
	pfxMWGRS   = "mwg-rs"
	nsMWGRS    = "http://www.metadataworkinggroup.com/schemas/regions/"
	pfxPS      = "photoshop"
	nsPS       = "http://ns.adobe.com/photoshop/1.0/"
	pfxTIFF    = "tiff"
	nsTIFF     = "http://ns.adobe.com/tiff/1.0/"
	pfxXMP     = "xmp"
	nsXMP      = "http://ns.adobe.com/xap/1.0/"
)

// A Provider handles data from an XMP/RDF block — but only the native XMP data,
// not the namespaces that mirror EXIF, IPTC, and/or Photoshop data.
type Provider struct {
	metadata.BaseProvider
	dcCreator               []string
	dcDescription           altString
	dcSubject               []string
	dcTitle                 altString
	digiKamTagsList         []metadata.HierValue
	exifDateTimeOriginal    metadata.DateTime
	exifDateTimeDigitized   metadata.DateTime
	exifGPSCoords           metadata.GPSCoords
	exifUserComment         altString
	iptcLocationCreated     location
	iptcLocationsShown      []location
	lrHierarchicalSubject   []metadata.HierValue
	mpRegPersonDisplayNames []string
	mwgrsNames              []string
	psDateCreated           metadata.DateTime
	tiffArtist              []string
	tiffDateTime            metadata.DateTime
	tiffImageDescription    altString
	tiffOrientation         metadata.Orientation
	xmpCreateDate           metadata.DateTime
	xmpMetadataDate         metadata.DateTime
	xmpModifyDate           metadata.DateTime

	rdf *rdf.Packet
}

var _ metadata.Provider = (*Provider)(nil) // verify interface compliance

// New creates a new Provider based on the provided RDF block.
func New(rdf *rdf.Packet) (p *Provider, err error) {
	p = &Provider{rdf: rdf}
	p.rdf.RegisterNamespace(pfxDC, nsDC)
	p.rdf.RegisterNamespace(pfxDigiKam, nsDigiKam)
	p.rdf.RegisterNamespace(pfxEXIF, nsEXIF)
	p.rdf.RegisterNamespace(pfxIPTC, nsIPTC)
	p.rdf.RegisterNamespace(pfxLR, nsLR)
	p.rdf.RegisterNamespace(pfxMP, nsMP)
	p.rdf.RegisterNamespace(pfxMPRI, nsMPRI)
	p.rdf.RegisterNamespace(pfxMPReg, nsMPReg)
	p.rdf.RegisterNamespace(pfxMWGRS, nsMWGRS)
	p.rdf.RegisterNamespace(pfxPS, nsPS)
	p.rdf.RegisterNamespace(pfxTIFF, nsTIFF)
	p.rdf.RegisterNamespace(pfxXMP, nsXMP)
	if err = p.getCaption(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getCreator(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getDateTime(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getFaces(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getGPS(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getGroups(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getKeywords(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getLocation(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getOrientation(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getPeople(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getPlaces(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getTitle(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if err = p.getTopics(); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	return p, nil
}

// ProviderName is the name for the provider, for debug purposes.
func (p *Provider) ProviderName() string { return "XMP" }
