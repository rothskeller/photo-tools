// Package xmp handles XMP metadata blocks.
package xmp

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"trimmer.io/go-xmp/xmp"
)

var dateRE = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d(?:\.\d+)?(?:[-+]\d\d:\d\d|Z)?$`)

// XMP is a an XMP parser and generator.
type XMP struct {
	doc      *xmp.Document
	problems []string
}

// New creates a new XMP block, to be added to a media file (or sidecar) that
// doesn't already have one.
func New() (p *XMP) {
	p = new(XMP)
	p.doc = xmp.NewDocument()
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
	return p
}

// RemoveNamespace removes an entire XML namespace from the XMP block.
func (p *XMP) RemoveNamespace(label, uri string) {
	p.doc.RemoveNamespace(xmp.NewNamespace(label, uri, nil))
}

// Render renders and returns the encoded XMP block, reflecting the data that
// was read, as subsequently modified by any SetXXX calls.
func (p *XMP) Render() ([]byte, error) {
	var buf bytes.Buffer
	if len(p.problems) != 0 {
		panic("XMP Render with parse problems")
	}
	if err := xmp.NewEncoder(&buf).Encode(p.doc); err != nil {
		return nil, fmt.Errorf("XMP.Encode: %s", err)
	}
	return buf.Bytes(), nil
}

// Problems returns the accumulated list of problems.
func (p *XMP) Problems() []string {
	if p == nil || p.doc == nil {
		return nil
	}
	return p.problems
}

// CanUpdate returns whether the XMP block can be rewritten safely.
func (p *XMP) CanUpdate() bool {
	return len(p.problems) == 0
}

func (p *XMP) log(f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	p.problems = append(p.problems, s)
}

func canonicalDate(date string) string {
	if strings.HasSuffix(date, "-00:00") || strings.HasSuffix(date, "+00:00") {
		return date[:len(date)-6] + "Z"
	}
	return date
}
