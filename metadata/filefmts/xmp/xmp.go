// Package xmp provides a file format handler for XMP files.
package xmp

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
func Read(fh *os.File) (h *XMP, err error) {
	var (
		line string
		scan = bufio.NewScanner(fh)
	)
	// First, we want to check whether this is an XMP file.  We will look
	// for <?xpacket, <x:xmpmeta, or <rdf:RDF at the beginning of it.
	if _, err = fh.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	for line == "" && scan.Scan() {
		line = strings.TrimSpace(scan.Text())
	}
	if !strings.HasPrefix(line, "<?xpacket") && !strings.HasPrefix(line, "<x:xmpmeta") && !strings.HasPrefix(line, "<rdf:RDF") {
		return nil, nil // not an XMP file
	}
	// It does appear to be an XMP file.  Parse it as RDF.
	h = new(XMP)
	if _, err = fh.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if h.rdf, err = rdf.Read(fh); err != nil {
		return nil, fmt.Errorf("%s: XMP: %s", fh.Name(), err)
	}
	if h.provider, err = xmp.New(h.rdf); err != nil {
		return nil, fmt.Errorf("%s: XMP: %s", fh.Name(), err)
	}
	return h, nil
}

// Provider returns the metadata.Provider for the XMP file.
func (h *XMP) Provider() metadata.Provider { return h.provider }
