// Package xmp provides a file format handler for XMP files.
package xmp

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
	"github.com/rothskeller/photo-tools/metadata/providers/xmp"
)

// XMP is a file handler for XMP files.
type XMP struct {
	rdf      *rdf.Packet
	provider *xmp.Provider
}

// Read reads the provided file.  It returns nil, nil, if the file is not an XMP
// file.  It returns an error if the file is a XMP file but ill-formed, or if a
// read error occurs.  It returns a XMP file handler for the file if it is read
// successfully.
func Read(r metadata.Reader) (h *XMP, err error) {
	var (
		line string
		scan = bufio.NewScanner(r)
	)
	// First, we want to check whether this is an XMP file.  We will look
	// for <?xpacket, <x:xmpmeta, or <rdf:RDF at the beginning of it.
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	for line == "" && scan.Scan() {
		line = strings.TrimSpace(scan.Text())
	}
	if !strings.HasPrefix(line, "<?xpacket") && !strings.HasPrefix(line, "<x:xmpmeta") && !strings.HasPrefix(line, "<rdf:RDF") {
		if _, err = r.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		return nil, nil // not an XMP file
	}
	// It does appear to be an XMP file.  Parse it as RDF.
	h = new(XMP)
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	h.rdf = rdf.New()
	if err = h.rdf.Read(r); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	if h.provider, err = xmp.New(h.rdf); err != nil {
		return nil, fmt.Errorf("XMP: %s", err)
	}
	return h, nil
}

// Provider returns the metadata.Provider for the XMP file.
func (h *XMP) Provider() metadata.Provider { return h.provider }

// Dirty returns whether the metadata from the file have been changed
// since they were read (and therefore need to be saved).
func (h *XMP) Dirty() bool { return h.rdf.Dirty() }

// Save writes the entire file to the supplied writer, including all
// revised metadata.
func (h *XMP) Save(out io.Writer) (err error) {
	_, err = h.rdf.Write(out)
	return err
}
