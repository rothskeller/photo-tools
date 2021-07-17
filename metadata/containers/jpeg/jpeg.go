// Package jpeg handles marshaling and unmarshaling of JPEG file segments.
package jpeg

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers"
)

const (
	markerJFIF   byte = 0xE0
	markerJFXX   byte = 0xE0
	markerEXIF   byte = 0xE1
	markerXMP    byte = 0xE1
	markerXMPext byte = 0xE1
	markerPSIR   byte = 0xED
)

var (
	nsJFIF   = []byte("JFIF\000")
	nsJFXX   = []byte("JFXX\000")
	nsEXIF   = []byte("Exif\000\000")
	nsXMP    = []byte("http://ns.adobe.com/xap/1.0/\000")
	nsXMPext = []byte("http://ns.adobe.com/xmp/extension/\000")
	nsPSIR   = []byte("Photoshop 3.0\000")
)

// A JPEG is a container of Segments.
type JPEG struct {
	start  *segmentGroup
	jfif   []*segmentGroup // multiple namespaces, so can't be one group
	exif   *segmentGroup
	xmp    *segmentGroup
	xmpext *segmentGroup
	psir   *segmentGroup
	others []*segmentGroup
	end    *segmentGroup
	size   int64
}

var _ containers.Container = (*JPEG)(nil) // verify interface compliance

// Read creates a new JPEG container handler, reading the specified reader.  It
// returns an error if the container is ill-formed or unreadable.
func (jpeg *JPEG) Read(r metadata.Reader) (err error) {
	var seg *segmentGroup

	seg = new(segmentGroup)
	if err = seg.Read(r); err == io.EOF || seg.marker != 0xD8 {
		return errors.New("JPEG: not a jpeg file")
	} else if err != nil {
		return fmt.Errorf("JPEG: %s", err)
	}
	jpeg.start = seg
	for seg.marker != 0xDA {
		seg = new(segmentGroup)
		if err = seg.Read(r); err != nil {
			return fmt.Errorf("JPEG: %s", err)
		}
		switch {
		case seg.marker == 0xDA:
			jpeg.end = seg
		case seg.namespace == nil:
			jpeg.others = append(jpeg.others, seg)
		case bytes.Equal(seg.namespace, nsJFIF):
			jpeg.jfif = append(jpeg.jfif, seg)
		case bytes.Equal(seg.namespace, nsJFXX):
			jpeg.jfif = append(jpeg.jfif, seg)
		case bytes.Equal(seg.namespace, nsEXIF):
			jpeg.exif = jpeg.exif.merge(seg, 0)
		case bytes.Equal(seg.namespace, nsXMP):
			jpeg.xmp = jpeg.xmp.merge(seg, 0)
		case bytes.Equal(seg.namespace, nsXMPext):
			jpeg.xmpext = jpeg.xmpext.merge(seg, 40)
		case bytes.Equal(seg.namespace, nsPSIR):
			jpeg.psir = jpeg.psir.merge(seg, 0)
		}
	}
	return nil
}

// Empty returns whether the container is empty (and should therefore be omitted
// from the written file, along with whatever tag in the parent container points
// to it).
func (jpeg *JPEG) Empty() bool { return false } // JPEGs are never empty

// Dirty returns whether any of the JPEG segments have been changed.
func (jpeg *JPEG) Dirty() bool {
	return jpeg.exif.Dirty() || jpeg.xmp.Dirty() || jpeg.psir.Dirty()
}

// Layout computes the rendered layout of the container, i.e. prepares for a
// call to Write, and returns what the rendered size of the container will be.
func (jpeg *JPEG) Layout() int64 {
	jpeg.size = jpeg.start.Layout()
	for _, seg := range jpeg.jfif {
		jpeg.size += seg.Layout()
	}
	jpeg.size += jpeg.exif.Layout()
	jpeg.size += jpeg.xmp.Layout()
	jpeg.size += jpeg.xmpext.Layout()
	jpeg.size += jpeg.psir.Layout()
	for _, seg := range jpeg.others {
		jpeg.size += seg.Layout()
	}
	jpeg.size += jpeg.end.Layout()
	return jpeg.size
}

// Write renders a JPEG file to the specified writer.
func (jpeg *JPEG) Write(w io.Writer) (count int, err error) {
	var n int

	n, err = jpeg.start.Write(w)
	count += n
	if err != nil {
		return count, err
	}
	for _, seg := range jpeg.jfif {
		n, err = seg.Write(w)
		count += n
		if err != nil {
			return count, err
		}
	}
	n, err = jpeg.exif.Write(w)
	count += n
	if err != nil {
		return count, err
	}
	n, err = jpeg.xmp.Write(w)
	count += n
	if err != nil {
		return count, err
	}
	n, err = jpeg.xmpext.Write(w)
	count += n
	if err != nil {
		return count, err
	}
	n, err = jpeg.psir.Write(w)
	count += n
	if err != nil {
		return count, err
	}
	for _, seg := range jpeg.others {
		n, err = seg.Write(w)
		count += n
		if err != nil {
			return count, err
		}
	}
	n, err = jpeg.end.Write(w)
	count += n
	if err != nil {
		return count, err
	}
	if jpeg.size != 0 && int(jpeg.size) != count {
		panic("actual size different from predicted size")
	}
	return count, err
}

// EXIF returns the contents of the EXIF segment, if any.
func (jpeg *JPEG) EXIF() metadata.Reader {
	if jpeg.exif != nil {
		return jpeg.exif.reader
	}
	return nil
}

// XMP returns the contents of the XMP segment, if any.
func (jpeg *JPEG) XMP() metadata.Reader {
	var (
		buf  [1]byte
		size int64
	)
	if jpeg.xmp == nil {
		return nil
	}
	// Many XMP segments in my library have extraneous null bytes at the
	// end, which the RDF parser can't handle.  Detect and remove them.
	size = jpeg.xmp.reader.Size()
	for size > 0 {
		jpeg.xmp.reader.ReadAt(buf[:], size-1)
		if buf[0] == 0 {
			size--
		} else {
			break
		}
	}
	if size < jpeg.xmp.reader.Size() {
		return io.NewSectionReader(jpeg.xmp.reader, 0, size)
	}
	return jpeg.xmp.reader
}

// XMPext returns the contents of the XMP extension segment, if any.
func (jpeg *JPEG) XMPext() metadata.Reader {
	if jpeg.xmpext != nil {
		return jpeg.xmpext.reader
	}
	return nil
}

// PSIR returns the contents of the PSIR segment, if any.
func (jpeg *JPEG) PSIR() metadata.Reader {
	if jpeg.psir != nil {
		return jpeg.psir.reader
	}
	return nil
}

// SetEXIFContainer sets the contents of the EXIF segment to those provided by
// the supplied container.
func (jpeg *JPEG) SetEXIFContainer(c containers.Container) {
	if jpeg.exif == nil {
		jpeg.exif = &segmentGroup{
			marker:    markerEXIF,
			namespace: nsEXIF,
		}
	}
	jpeg.exif.container = c
}

// SetXMPContainer sets the contents of the XMP segment to those provided by the
// supplied container.
func (jpeg *JPEG) SetXMPContainer(c containers.Container) {
	if jpeg.xmp == nil {
		jpeg.xmp = &segmentGroup{
			marker:    markerXMP,
			namespace: nsXMP,
		}
	}
	jpeg.xmp.container = c
}

// SetPSIRContainer sets the contents of the PSIR segment to those provided by the
// supplied container.
func (jpeg *JPEG) SetPSIRContainer(c containers.Container) {
	if jpeg.psir == nil {
		jpeg.psir = &segmentGroup{
			marker:    markerPSIR,
			namespace: nsPSIR,
		}
	}
	jpeg.psir.container = c
}
