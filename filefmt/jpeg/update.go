package jpeg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"

	"github.com/rothskeller/photo-tools/metadata/exif"
	"github.com/rothskeller/photo-tools/metadata/iptc"
	"github.com/rothskeller/photo-tools/metadata/xmp"
)

// SaveMetadata rewrites the JPEG file with the supplied metadata.
func (h *JPEG) SaveMetadata() (err error) {
	var (
		tempfn string
		ifh    *os.File
		ofh    *os.File
		in     *bufio.Reader
		out    *bufio.Writer
		offset uint32
		s      *segment
	)
	if len(h.problems) != 0 {
		panic("JPEG UpdateMetadata after parse failures")
	}
	if !h.exif.Dirty() && !h.iptc.Dirty() && !h.xmp.Dirty() {
		return nil
	}
	if h.xmp == nil {
		h.xmp = xmp.New()
	}
	if ifh, err = os.Open(h.path); err != nil {
		return err
	}
	defer ifh.Close()
	in = bufio.NewReader(ifh)
	tempfn = filepath.Dir(h.path) + "/." + filepath.Base(h.path) + ".TEMP"
	if ofh, err = os.Create(tempfn); err != nil {
		return err
	}
	defer ofh.Close()
	defer os.Remove(tempfn)
	out = bufio.NewWriter(ofh)

	for {
		if s, offset, err = h.readSegment(offset, in); err != nil {
			return err
		}
		if s == nil || s.marker == 0xDA {
			break
		}
		switch s.marker {
		case 0xD8: // SOI
			if err = writeSegment(out, s); err == nil {
				if err = writeEXIFSegment(out, h.exif); err == nil {
					if err = writeXMPSegment(out, h.xmp); err == nil {
						err = writeIPTCSegment(out, h.iptc)
					}
				}
			}
		case 0xE1:
			err = maybeWriteAPP1Segment(out, s)
		case 0xED:
			err = maybeWriteAPP13Segment(out, s)
		default:
			err = writeSegment(out, s)
		}
		if err != nil {
			return err
		}
	}
	if err = writeSegment(out, s); err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	if err = out.Flush(); err != nil {
		return err
	}
	if err = ofh.Close(); err != nil {
		return err
	}
	if err = os.Rename(tempfn, h.path); err != nil {
		return err
	}
	return nil
}

type segment struct {
	marker byte
	size   uint16
	buf    []byte
}

func (h *JPEG) readSegment(offset uint32, in io.Reader) (s *segment, outoff uint32, err error) {
	var (
		sbuf [2]byte
	)
	s = new(segment)

RESTART:
	if _, err = in.Read(sbuf[:1]); err != nil && err != io.EOF {
		return nil, offset, err
	}
	if err == io.EOF {
		return nil, offset, nil
	}
	if sbuf[0] != 0xFF {
		offset++
		goto RESTART
	}
	offset++
PADDING:
	if _, err = in.Read(sbuf[:1]); err != nil {
		return nil, offset, err
	}
	offset++
	s.marker = sbuf[0]
	switch s.marker {
	case 0xFF:
		goto PADDING
	case 0x00:
		goto RESTART
	case 0x01, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9, 0xDA:
		// Most of these are non-segment markers.  0xDA is a real
		// segment marker but it marks the point after which there can
		// be no more metadata, so we return it as a non-segment marker
		// and the calling code handles the rest of the file
		// differently.
		return s, offset, nil
	}
	if _, err = io.ReadFull(in, sbuf[:]); err != nil {
		return nil, offset, err
	}
	offset += 2
	s.size = binary.BigEndian.Uint16(sbuf[:]) - 2
	if s.size < 1 {
		goto RESTART
	}
	s.buf = make([]byte, s.size)
	if _, err = io.ReadFull(in, s.buf); err != nil {
		return nil, offset, err
	}
	offset += uint32(s.size)
	return s, offset, nil
}

func writeSegment(out io.Writer, s *segment) (err error) {
	var sbuf [2]byte

	sbuf[0] = 0xFF
	sbuf[1] = s.marker
	if _, err = out.Write(sbuf[:]); err != nil {
		return err
	}
	if s.size == 0 {
		return nil
	}
	binary.BigEndian.PutUint16(sbuf[:], s.size+2)
	if _, err = out.Write(sbuf[:]); err != nil {
		return err
	}
	_, err = out.Write(s.buf)
	return err
}

func writeEXIFSegment(out io.Writer, e *exif.EXIF) (err error) {
	var sbuf [2]byte
	var buf = e.Render(0xFFF7)

	sbuf[0] = 0xFF
	sbuf[1] = 0xE1
	if _, err = out.Write(sbuf[:]); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(sbuf[:], uint16(len(buf)+8))
	if _, err = out.Write(sbuf[:]); err != nil {
		return err
	}
	if _, err = out.Write([]byte("Exif\000\000")); err != nil {
		return err
	}
	_, err = out.Write(buf)
	return err
}

func writeXMPSegment(out io.Writer, x *xmp.XMP) (err error) {
	var sbuf [2]byte
	var buf []byte

	if buf, err = x.Render(); err != nil {
		return err
	}
	sbuf[0] = 0xFF
	sbuf[1] = 0xE1
	if _, err = out.Write(sbuf[:]); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(sbuf[:], uint16(len(buf)+31))
	if _, err = out.Write(sbuf[:]); err != nil {
		return err
	}
	if _, err = out.Write([]byte("http://ns.adobe.com/xap/1.0/\000")); err != nil {
		return err
	}
	_, err = out.Write(buf)
	return err
}

func writeIPTCSegment(out io.Writer, i *iptc.IPTC) (err error) {
	var sbuf [2]byte
	var buf = i.Render()

	sbuf[0] = 0xFF
	sbuf[1] = 0xED
	if _, err = out.Write(sbuf[:]); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(sbuf[:], uint16(len(buf)+16))
	if _, err = out.Write(sbuf[:]); err != nil {
		return err
	}
	if _, err = out.Write([]byte("Photoshop 3.0\000")); err != nil {
		return err
	}
	_, err = out.Write(buf)
	return err
}

func maybeWriteAPP1Segment(out io.Writer, s *segment) (err error) {
	if s.size >= 6 && bytes.HasPrefix(s.buf, []byte("Exif\000\000")) {
		return nil
	} else if s.size >= 29 && bytes.HasPrefix(s.buf, []byte("http://ns.adobe.com/xap/1.0/\000")) {
		return nil
	}
	return writeSegment(out, s)
}

func maybeWriteAPP13Segment(out io.Writer, s *segment) (err error) {
	if s.size >= 14 && bytes.HasPrefix(s.buf, []byte("Photoshop 3.0\000")) {
		return nil
	}
	return writeSegment(out, s)
}
