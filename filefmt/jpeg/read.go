package jpeg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"os"

	"github.com/rothskeller/photo-tools/metadata/exif"
	"github.com/rothskeller/photo-tools/metadata/xmp"
)

// ReadMetadata reads the metadata from a JPEG file.  It returns a list of
// problems found while reading.
func (h *JPEG) ReadMetadata() {
	var (
		fh     *os.File
		in     *OffsetReader
		marker byte
		size   uint16
		ok     bool
		err    error
		header = make([]byte, 2)
	)
	// Open the file and get name and modification time.
	if fh, err = os.Open(h.path); err != nil {
		h.problems = []string{err.Error()}
		return
	}
	defer fh.Close()
	// Check that it's really a JPEG file.
	if _, err = fh.Read(header); err != nil {
		h.problems = []string{err.Error()}
		return
	}
	if header[0] != 0xFF || header[1] != 0xD8 {
		h.problems = []string{"missing JPEG file header (not really a JPEG?)"}
		return
	}
	if _, err := fh.Seek(0, 0); err != nil {
		h.problems = []string{err.Error()}
		return
	}
	// Walk through each of the segments of the file, until we see image
	// data (which means there will be no more metadata).
	in = NewOffsetReader(bufio.NewReader(fh))
	for {
		mstart := in.Offset()
		marker, size, ok = h.readMarker(in)
		if marker == 0x00 || marker == 0xDA || !ok { // EOF or image data segment
			break
		}
		if size == 0 { // non-segment marker
			continue
		}
		buf := make([]byte, size)
		if _, err = io.ReadFull(in, buf); err != nil {
			h.problems = []string{err.Error()}
			return
		}
		switch marker {
		case 0xE0: // APP0 segment, with JFIF or JFXX metadata
			if len(buf) >= 5 && bytes.HasPrefix(buf, []byte("JFIF\000")) {
				h.jfif = append(h.jfif, segment{marker: marker, size: size, buf: buf})
			} else if len(buf) >= 5 && bytes.HasPrefix(buf, []byte("JFXX\000")) {
				h.jfif = append(h.jfif, segment{marker: marker, size: size, buf: buf})
			}
		case 0xE1: // APP1 segment, with EXIF or XMP metadata
			if len(buf) >= 6 && bytes.HasPrefix(buf, []byte("Exif\000\000")) {
				h.exif = exif.Parse(buf[6:], uint32(mstart+6))
			} else if len(buf) >= 29 && bytes.HasPrefix(buf, []byte("http://ns.adobe.com/xap/1.0/\000")) {
				h.xmp = xmp.Parse(buf[29:])
			} else if len(buf) >= 35 && bytes.HasPrefix(buf, []byte("http://ns.adobe.com/xmp/extension/\000")) {
				h.xmpext = h.xmpext.Parse(buf[35:])
			}
		case 0xED: // APP13 segment, with IPTC metadata
			if len(buf) >= 14 && bytes.HasPrefix(buf, []byte("Photoshop 3.0\000")) {
				h.iptc = h.iptc.Parse(buf[14:], uint32(mstart+14))
			}
		}
	}
	h.iptc.Check()
	h.xmpext.Check()
}

// readMarker reads a marker from the JPEG file, including the size of the
// segment started by the marker.  It returns the marker byte and the segment
// size (not self-inclusive).  If it cannot parse the file, it skips over the
// offensive part with a log message.  It returns an error only if the file
// cannot be read, or ends unexpectedly.
func (h *JPEG) readMarker(in *OffsetReader) (marker byte, size uint16, ok bool) {
	var (
		sbuf [2]byte
		err  error
	)

RESTART:
	if marker, err = in.ReadByte(); err != nil && err != io.EOF {
		h.problems = append(h.problems, err.Error())
		return 0, 0, false
	}
	if err == io.EOF {
		return 0, 0, true
	}
	if marker != 0xFF {
		h.log("[%x] 0x%x found where marker (0xFF) expected", in.Offset(), marker)
		goto RESTART
	}
PADDING:
	if marker, err = in.ReadByte(); err != nil {
		h.problems = append(h.problems, err.Error())
		return 0, 0, false
	}
	switch marker {
	case 0xFF:
		goto PADDING
	case 0x00:
		h.log("[%x] 0xFF00 found where marker (0xFFnn) expected", in.Offset())
		goto RESTART
	case 0x01, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9:
		return marker, 0, true
	}
	if _, err = io.ReadFull(in, sbuf[:]); err != nil {
		h.problems = append(h.problems, err.Error())
		return 0, 0, false
	}
	size = binary.BigEndian.Uint16(sbuf[:])
	if size < 2 {
		h.log("[%x] invalid segment size", in.Offset())
		goto RESTART
	}
	return marker, size - 2, true
}
