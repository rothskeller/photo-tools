// Package xmp handles XMP metadata blocks.
package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/rdf"
)

// XMP is a an XMP parser and generator.
type XMP struct {
	dcCreator               []string
	dcDescription           metadata.AltString
	dcSubject               []string
	dcTitle                 metadata.AltString
	digiKamTagsList         []metadata.Keyword
	exifDateTimeOriginal    metadata.DateTime
	exifDateTimeDigitized   metadata.DateTime
	exifGPSCoords           metadata.GPSCoords
	exifUserComments        []string
	iptcLocationCreated     metadata.Location
	iptcLocationsShown      []metadata.Location
	lrHierarchicalSubject   []metadata.Keyword
	mpRegPersonDisplayNames []string
	mwgrsNames              []string
	psDateCreated           metadata.DateTime
	tiffArtist              string
	tiffDateTime            metadata.DateTime
	tiffImageDescription    metadata.AltString
	xmpCreateDate           metadata.DateTime
	xmpMetadataDate         metadata.DateTime
	xmpModifyDate           metadata.DateTime
	Problems                []string

	rdf   *rdf.Packet
	dirty bool
}

// New creates a new XMP block, to be added to a media file (or sidecar) that
// doesn't already have one.
func New() (p *XMP) {
	p = new(XMP)
	p.rdf = rdf.NewPacket()
	p.dirty = true
	return p
}

// Parse parses an XMP block and returns the parse results.
func Parse(buf []byte) (p *XMP) {
	var err error

	p = new(XMP)
	// I have found that some XMP blocks are null-terminated, for no
	// obvious reason.  The XML parser doesn't like that.
	for len(buf) != 0 && buf[len(buf)-1] == 0 {
		buf = buf[:len(buf)-1]
	}
	if p.rdf, err = rdf.ReadPacket(buf); err != nil {
		p.log("XMP parse error: %s", err)
		return p
	}
	p.getDC()
	p.getDigiKam()
	p.getEXIF()
	p.getIPTC()
	p.getLR()
	p.getMWGRS()
	p.getMP()
	p.getPS()
	p.getTIFF()
	p.getXMP()
	return p
}

// RemoveNamespace removes an entire XML namespace from the XMP block.
func (p *XMP) RemoveNamespace(_, uri string) {
	for name := range p.rdf.Properties {
		if name.Namespace == uri {
			delete(p.rdf.Properties, name)
		}
	}
}

// Dirty returns whether the XMP data have changed and need to be saved.
func (p *XMP) Dirty() bool {
	if p == nil || len(p.Problems) != 0 {
		return false
	}
	return p.dirty
}

// Render renders and returns the encoded XMP block, reflecting the data that
// was read, as subsequently modified by any SetXXX calls.
func (p *XMP) Render() (buf []byte, err error) {
	if len(p.Problems) != 0 {
		panic("XMP Render with parse problems")
	}
	if buf, err = p.rdf.Render(); err != nil {
		return nil, fmt.Errorf("XMP.Render: %s", err)
	}
	return buf, nil
}

func (p *XMP) log(f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	p.Problems = append(p.Problems, "XMP: "+s)
}

func (p *XMP) xmpDateTimeToMetadata(x string, m *metadata.DateTime) {
	if err := m.Parse(x); err != nil {
		p.log("invalid DateTime value")
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
