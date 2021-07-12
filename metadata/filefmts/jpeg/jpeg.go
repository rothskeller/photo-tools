// Package jpeg contains the file format handler for JPEG files.
package jpeg

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/iim"
	"github.com/rothskeller/photo-tools/metadata/containers/jpeg"
	"github.com/rothskeller/photo-tools/metadata/containers/photoshop"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
	"github.com/rothskeller/photo-tools/metadata/containers/tiff"
	"github.com/rothskeller/photo-tools/metadata/providers/exififd"
	"github.com/rothskeller/photo-tools/metadata/providers/gpsifd"
	"github.com/rothskeller/photo-tools/metadata/providers/iptc"
	"github.com/rothskeller/photo-tools/metadata/providers/jpegifd0"
	"github.com/rothskeller/photo-tools/metadata/providers/multi"
	"github.com/rothskeller/photo-tools/metadata/providers/xmp"
	"github.com/rothskeller/photo-tools/metadata/providers/xmpext"
)

const (
	tagEXIFIFD uint16 = 0x8769
	tagGPSIFD  uint16 = 0x8825
	psirIDIIM  uint16 = 0x404
	psirIDHash uint16 = 0x425
)

// JPEG is a JPEG file handler.
type JPEG struct {
	container        *jpeg.JPEG
	exifSeg          jpeg.SegmentReader
	psirSeg          jpeg.SegmentReader
	xmpSeg           jpeg.SegmentReader
	xmpExtSeg        jpeg.SegmentReader
	exifTIFF         *tiff.TIFF
	jpegIFD0         *tiff.IFD
	exifIFD          *tiff.IFD
	gpsIFD           *tiff.IFD
	psirBlock        *photoshop.Photoshop
	iimPSIR          *photoshop.PSIR
	hashPSIR         *photoshop.PSIR
	iim              *iim.IIM
	xmpRDF           *rdf.Packet
	xmpExtRDF        *rdf.Packet
	jpegIFD0Provider *jpegifd0.Provider
	exifIFDProvider  *exififd.Provider
	gpsIFDProvider   *gpsifd.Provider
	iptcProvider     *iptc.Provider
	xmpProvider      *xmp.Provider
	xmpExtProvider   *xmpext.Provider
	providers        multi.Provider
}

// Read reads the provided file.  It returns nil, nil, if the file is not a JPEG
// file.  It returns an error if the file is a JPEG file but ill-formed, or if a
// read error occurs.  It returns a JPEG file handler for the file if it is read
// successfully.
func Read(fh *os.File) (jh *JPEG, err error) {
	var (
		buf  [2]byte
		size int64
	)
	if _, err = fh.ReadAt(buf[0:2], 0); err == io.EOF {
		return nil, nil // can't read a signature, assume it's not JPEG
	} else if err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	} else if buf[0] != 0xFF || buf[1] != 0xD8 {
		return nil, nil // not a JPEG file
	}
	if size, err = fh.Seek(0, io.SeekEnd); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	jh = new(JPEG)
	if jh.container, err = jpeg.Read(io.NewSectionReader(fh, 0, size)); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if err = jh.readEXIFSegment(); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if err = jh.readXMPSegments(); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if err = jh.readPSIRSegment(); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	return jh, nil
}

// Provider returns the metadata.Provider for the JPEG file.
func (jh *JPEG) Provider() metadata.Provider { return jh.providers }

func (jh *JPEG) readEXIFSegment() (err error) {
	if jh.exifSeg = jh.container.EXIF(); jh.exifSeg == nil {
		return nil
	}
	if jh.exifTIFF, err = tiff.Read(jh.exifSeg); err != nil {
		return fmt.Errorf("EXIF segment: %s", err)
	}
	if jh.jpegIFD0 = jh.exifTIFF.IFD0(); jh.jpegIFD0 == nil {
		return errors.New("EXIF segment: no IFD0")
	}
	if jh.jpegIFD0Provider, err = jpegifd0.New(jh.jpegIFD0); err != nil {
		return err
	}
	jh.providers = append(jh.providers, jh.jpegIFD0Provider)
	if tag := jh.jpegIFD0.Tag(tagEXIFIFD); tag != nil {
		if jh.exifIFD, err = tag.AsIFD(); err != nil {
			return fmt.Errorf("EXIF IFD: %s", err)
		}
		if jh.exifIFDProvider, err = exififd.New(jh.exifIFD, jh.exifTIFF.Encoding()); err != nil {
			return err
		}
		jh.providers = append(jh.providers, jh.exifIFDProvider)
	}
	if tag := jh.jpegIFD0.Tag(tagGPSIFD); tag != nil {
		if jh.gpsIFD, err = tag.AsIFD(); err != nil {
			return fmt.Errorf("GPS IFD: %s", err)
		}
		if jh.gpsIFDProvider, err = gpsifd.New(jh.gpsIFD); err != nil {
			return err
		}
		jh.providers = append(jh.providers, jh.gpsIFDProvider)
	}
	return nil
}

func (jh *JPEG) readPSIRSegment() (err error) {
	if jh.psirSeg = jh.container.PSIR(); jh.psirSeg == nil {
		return nil
	}
	if jh.psirBlock, err = photoshop.Read(jh.psirSeg); err != nil {
		return fmt.Errorf("PSIR segment: %s", err)
	}
	if jh.iimPSIR = jh.psirBlock.PSIR(psirIDIIM); jh.iimPSIR == nil {
		return nil
	}
	if jh.iim, _, err = iim.Read(jh.iimPSIR.Reader()); err != nil {
		return err
	}
	if jh.iptcProvider, err = iptc.New(jh.iim); err != nil {
		return err
	}
	jh.providers = append(jh.providers, jh.iptcProvider)
	jh.hashPSIR = jh.psirBlock.PSIR(psirIDHash)
	return nil
}

func (jh *JPEG) readXMPSegments() (err error) {
	if jh.xmpSeg = jh.container.XMP(); jh.xmpSeg != nil {
		if jh.xmpRDF, err = rdf.Read(jh.xmpSeg); err != nil {
			return fmt.Errorf("XMP: %s", err)
		}
		if jh.xmpProvider, err = xmp.New(jh.xmpRDF); err != nil {
			return err
		}
		jh.providers = append(jh.providers, jh.xmpProvider)
	}
	if jh.xmpExtSeg = jh.container.XMPext(); jh.xmpExtSeg != nil {
		if jh.xmpExtRDF, err = rdf.Read(jh.xmpExtSeg); err != nil {
			return fmt.Errorf("XMPExt: %s", err)
		}
		if jh.xmpExtProvider, err = xmpext.New(jh.xmpExtRDF); err != nil {
			return err
		}
		jh.providers = append(jh.providers, jh.xmpExtProvider)
	}
	return nil
}
