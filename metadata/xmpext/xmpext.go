// Package xmpext contains code for handling XMP extension segments in JPEG
// files.
package xmpext

import (
	"bytes"
	"encoding/binary"

	"github.com/rothskeller/photo-tools/metadata/xmp/rdf"
)

// XMPExt is a pseudo-tracker for XMP extension segments in JPEG files.
type XMPExt struct {
	guid     []byte
	buf      []byte
	Problems []string
}

// Parse is called to parse one extension segment.  The first call has a nil
// receiver, and allocates the XMPExt structure; subsequent calls return the
// receiver.  Despite the name (chosen for parallelism with the other segment
// type handlers), this function doesn't actually parse the segment.  That
// waits until all of the segnments have been collected, and Check is called.
func (h *XMPExt) Parse(buf []byte) *XMPExt {
	if h == nil {
		h = new(XMPExt)
	}
	if len(buf) < 40 {
		h.Problems = append(h.Problems, "XMPExt: ill-formed XMP extension block")
		return h
	}
	if h.guid == nil {
		h.guid = buf[:32]
		length := binary.BigEndian.Uint32(buf[32:])
		h.buf = make([]byte, length)
	} else if !bytes.Equal(h.guid, buf[:32]) {
		h.Problems = append(h.Problems, "XMPExt: extension GUID mismatch")
		return h
	} else if len(h.buf) != int(binary.BigEndian.Uint32(buf[32:])) {
		h.Problems = append(h.Problems, "XMPExt: extension size mismatch")
		return h
	}
	offset := binary.BigEndian.Uint32(buf[36:])
	if int(offset)+len(buf)-40 > len(h.buf) {
		h.Problems = append(h.Problems, "XMPExt: extension block exceeds extension size")
		return h
	}
	copy(h.buf[offset:], buf[40:])
	return h
}

// Check looks at the XMP extension block, resulting from merging all of the
// parsed extension segments, to determine whether it can be safely ignored.
// Since this library doesn't know how to write XMP metadata with an extension
// block, the only safe way to handle the file is if the extension block can be
// ignored.  That basically means it doesn't have any of the data we care about
// in it.
func (h *XMPExt) Check() {
	var (
		p   *rdf.Packet
		err error
	)
	if h == nil || len(h.Problems) != 0 {
		return
	}
	if p, err = rdf.ReadPacket(h.buf); err != nil {
		h.Problems = append(h.Problems, "XMPExt: "+err.Error())
		return
	}
	h.buf = nil // free the space, we're done with it
	for name := range p.Properties {
		switch name.Namespace {
		case "http://purl.org/dc/elements/1.1/",
			"http://www.digikam.org/ns/1.0/",
			"http://ns.adobe.com/exif/1.0/",
			"http://iptc.org/std/Iptc4xmpExt/2008-02-29/",
			"http://ns.adobe.com/lightroom/1.0/",
			"http://ns.microsoft.com/photo/1.2/",
			"http://www.metadataworkinggroup.com/schemas/regions/",
			"http://ns.adobe.com/photoshop/1.0/",
			"http://ns.adobe.com/tiff/1.0/",
			"http://ns.adobe.com/xap/1.0/":
			h.Problems = append(h.Problems, "XMPExt: unsupported extension (contains namespaces that should be in main XMP)")
			return
		}
	}
}
