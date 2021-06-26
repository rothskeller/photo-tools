// Package xmpext contains code for handling XMP extension segments in JPEG
// files.
package xmpext

import (
	"bytes"
	"encoding/binary"

	"trimmer.io/go-xmp/xmp"
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
		doc *xmp.Document
		ns  xmp.NamespaceList
		err error
	)
	if h == nil || len(h.Problems) != 0 {
		return
	}
	if doc, err = xmp.Read(bytes.NewReader(h.buf)); err != nil {
		h.Problems = append(h.Problems, "XMPExt: "+err.Error())
		return
	}
	ns = doc.Namespaces()
	if ns.ContainsName("dc") ||
		ns.ContainsName("digiKam") ||
		ns.ContainsName("exif") ||
		ns.ContainsName("IptcxmpExt") ||
		ns.ContainsName("lr") ||
		ns.ContainsName("photoshop") ||
		ns.ContainsName("tiff") ||
		ns.ContainsName("xmp") {
		h.Problems = append(h.Problems, "XMPExt: unsupported extension (contains namespaces that should be in main XMP)")
	}
	h.buf = nil // free the space, we're done with it
}
