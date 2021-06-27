// Package xmp handles XMP metadata blocks.
package xmp

import (
	"bytes"
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"trimmer.io/go-xmp/xmp"
)

// XMP is a an XMP parser and generator.
type XMP struct {
	DCCreator             []string
	DCDescription         metadata.AltString
	DCSubject             []string
	DCTitle               metadata.AltString
	DigiKamTagsList       []metadata.Keyword
	EXIFDateTimeOriginal  metadata.DateTime
	EXIFDateTimeDigitized metadata.DateTime
	EXIFGPSCoords         metadata.GPSCoords
	EXIFUserComments      []string
	IPTCLocationCreated   metadata.Location
	IPTCLocationsShown    []metadata.Location
	LRHierarchicalSubject []metadata.Keyword
	PSDateCreated         metadata.DateTime
	TIFFArtist            string
	TIFFDateTime          metadata.DateTime
	TIFFImageDescription  metadata.AltString
	XMPCreateDate         metadata.DateTime
	XMPMetadataDate       metadata.DateTime
	XMPModifyDate         metadata.DateTime
	Problems              []string

	doc   *xmp.Document
	dirty bool
}

// New creates a new XMP block, to be added to a media file (or sidecar) that
// doesn't already have one.
func New() (p *XMP) {
	p = new(XMP)
	p.doc = xmp.NewDocument()
	p.dirty = true
	return p
}

// Parse parses an XMP block and returns the parse results.
func Parse(buf []byte) (p *XMP) {
	var err error

	p = new(XMP)
	if p.doc, err = xmp.Read(bytes.NewReader(buf)); err != nil {
		p.log("XMP parse error: %s", err)
		p.doc = nil // just to be sure
	}
	p.getDC()
	p.getDigiKam()
	p.getEXIF()
	p.getIPTC()
	p.getLR()
	p.getPS()
	p.getTIFF()
	p.getXMP()
	return p
}

// RemoveNamespace removes an entire XML namespace from the XMP block.
func (p *XMP) RemoveNamespace(label, uri string) {
	p.doc.RemoveNamespace(xmp.NewNamespace(label, uri, nil))
}

// Dirty returns whether the XMP data have changed and need to be saved.
func (p *XMP) Dirty() bool {
	if p == nil || len(p.Problems) != 0 {
		return false
	}
	p.setDC()
	p.setDigiKam()
	p.setEXIF()
	p.setIPTC()
	p.setLR()
	p.setPS()
	p.setTIFF()
	p.setXMP()
	return p.dirty
}

// Render renders and returns the encoded XMP block, reflecting the data that
// was read, as subsequently modified by any SetXXX calls.
func (p *XMP) Render() ([]byte, error) {
	var buf bytes.Buffer
	if len(p.Problems) != 0 {
		panic("XMP Render with parse problems")
	}
	p.setDC()
	p.setDigiKam()
	p.setEXIF()
	p.setIPTC()
	p.setLR()
	p.setPS()
	p.setTIFF()
	p.setXMP()
	if err := xmp.NewEncoder(&buf).Encode(p.doc); err != nil {
		return nil, fmt.Errorf("XMP.Encode: %s", err)
	}
	return buf.Bytes(), nil
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
