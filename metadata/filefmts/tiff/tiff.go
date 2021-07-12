// Package tiff contains the file format handler for TIFF files.
package tiff

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/iim"
	"github.com/rothskeller/photo-tools/metadata/containers/photoshop"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
	"github.com/rothskeller/photo-tools/metadata/containers/tiff"
	"github.com/rothskeller/photo-tools/metadata/providers/exififd"
	"github.com/rothskeller/photo-tools/metadata/providers/gpsifd"
	"github.com/rothskeller/photo-tools/metadata/providers/iptc"
	"github.com/rothskeller/photo-tools/metadata/providers/multi"
	"github.com/rothskeller/photo-tools/metadata/providers/tiffifd0"
	"github.com/rothskeller/photo-tools/metadata/providers/xmp"
)

const (
	tagEXIFIFD uint16 = 0x8769
	tagGPSIFD  uint16 = 0x8825
	tagIPTC    uint16 = 0x83BB
	tagPSIR    uint16 = 0x8649
	tagXMP     uint16 = 0x02BC
	psirIDIIM  uint16 = 0x404
	psirIDHash uint16 = 0x425
)

// TIFF is a TIFF file handler.
type TIFF struct {
	container        *tiff.TIFF
	tiffIFD0         *tiff.IFD
	xmpIFD           *tiff.IFD
	tiffIFD0Provider *tiffifd0.Provider
	xmpRDF           *rdf.Packet
	xmpProvider      *xmp.Provider
	iim              *iim.IIM
	iptcProvider     *iptc.Provider
	psirBlock        *photoshop.Photoshop
	iimPSIR          *photoshop.PSIR
	psirIIM          *iim.IIM
	psirIPTCProvider *iptc.Provider
	hashPSIR         *photoshop.PSIR
	exifIFD          *tiff.IFD
	exifIFDProvider  *exififd.Provider
	gpsIFD           *tiff.IFD
	gpsIFDProvider   *gpsifd.Provider
	providers        multi.Provider
}

var tiffHeaderLE = []byte{0x49, 0x49, 0x2A, 0x00}
var tiffHeaderBE = []byte{0x4D, 0x4D, 0x00, 0x2A}

// Read reads the provided file.  It returns nil, nil, if the file is not a TIFF
// file.  It returns an error if the file is a TIFF file but ill-formed, or if a
// read error occurs.  It returns a TIFF file handler for the file if it is read
// successfully.
func Read(fh *os.File) (h *TIFF, err error) {
	var (
		buf  [4]byte
		size int64
	)
	if _, err = fh.ReadAt(buf[0:4], 0); err == io.EOF {
		return nil, nil // can't read a signature, assume it's not TIFF
	} else if err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	} else if !bytes.Equal(buf[0:4], tiffHeaderBE) && !bytes.Equal(buf[0:4], tiffHeaderLE) {
		return nil, nil // not a JPEG file
	}
	if size, err = fh.Seek(0, io.SeekEnd); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	h = new(TIFF)
	if h.container, err = tiff.Read(io.NewSectionReader(fh, 0, size)); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if err = h.readIFD0(); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if err = h.readEXIFIFD(); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if err = h.readGPSIFD(); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if err = h.readXMPTag(); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if err = h.readIPTCTag(); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	if err = h.readPSIRIFD(); err != nil {
		return nil, fmt.Errorf("%s: %s", fh.Name(), err)
	}
	return h, nil
}

// Provider returns the metadata.Provider for the JPEG file.
func (h *TIFF) Provider() metadata.Provider { return h.providers }

func (h *TIFF) readIFD0() (err error) {
	if h.tiffIFD0 = h.container.IFD0(); h.tiffIFD0 == nil {
		return errors.New("TIFF: no IFD0")
	}
	if h.tiffIFD0Provider, err = tiffifd0.New(h.tiffIFD0); err != nil {
		return err
	}
	h.providers = append(h.providers, h.tiffIFD0Provider)
	return nil
}

func (h *TIFF) readXMPTag() (err error) {
	if tag := h.tiffIFD0.Tag(tagXMP); tag != nil {
		var r metadata.Reader
		if r, err = tag.AsUnknownReader(); err != nil {
			if r, err = tag.AsBytesReader(); err != nil {
				return fmt.Errorf("XMP: %s", err)
			}
		}
		if h.xmpRDF, err = rdf.Read(r); err != nil {
			return fmt.Errorf("XMP: %s", err)
		}
		if h.xmpProvider, err = xmp.New(h.xmpRDF); err != nil {
			return err
		}
		h.providers = append(h.providers, h.xmpProvider)
	}
	return nil
}

func (h *TIFF) readIPTCTag() (err error) {
	if tag := h.tiffIFD0.Tag(tagIPTC); tag != nil {
		// Empirically, tag type can be UNKNOWN or LONG.
		var r metadata.Reader
		if r, err = tag.AsUnknownReader(); err != nil {
			if r, err = tag.AsLongReader(); err != nil {
				return fmt.Errorf("IPTC: %s", err)
			}
		}
		if h.iim, _, err = iim.Read(r); err != nil {
			return fmt.Errorf("IPTC: %s", err)
		}
		if h.iptcProvider, err = iptc.New(h.iim); err != nil {
			return err
		}
		h.providers = append(h.providers, h.iptcProvider)
	}
	return nil
}

func (h *TIFF) readPSIRIFD() (err error) {
	if tag := h.tiffIFD0.Tag(tagPSIR); tag != nil {
		var r metadata.Reader
		// Empirically, tag type could be either UNKNOWN or BYTE.
		if r, err = tag.AsUnknownReader(); err != nil {
			if r, err = tag.AsBytesReader(); err != nil {
				return fmt.Errorf("PSIR: %s", err)
			}
		}
		if h.psirBlock, err = photoshop.Read(r); err != nil {
			return fmt.Errorf("PSIR block: %s", err)
		}
		if h.iimPSIR = h.psirBlock.PSIR(psirIDIIM); h.iimPSIR == nil {
			return nil
		}
		if h.psirIIM, _, err = iim.Read(h.iimPSIR.Reader()); err != nil {
			return fmt.Errorf("PSIR: %s", err)
		}
		if h.psirIPTCProvider, err = iptc.New(h.psirIIM); err != nil {
			return fmt.Errorf("PSIR: %s", err)
		}
		h.providers = append(h.providers, h.psirIPTCProvider)
		h.hashPSIR = h.psirBlock.PSIR(psirIDHash)
	}
	return nil
}

func (h *TIFF) readEXIFIFD() (err error) {
	if tag := h.tiffIFD0.Tag(tagEXIFIFD); tag != nil {
		if h.exifIFD, err = tag.AsIFD(); err != nil {
			return fmt.Errorf("EXIF IFD: %s", err)
		}
		if h.exifIFDProvider, err = exififd.New(h.exifIFD, h.container.Encoding()); err != nil {
			return err
		}
		h.providers = append(h.providers, h.exifIFDProvider)
	}
	return nil
}

func (h *TIFF) readGPSIFD() (err error) {
	if tag := h.tiffIFD0.Tag(tagGPSIFD); tag != nil {
		if h.gpsIFD, err = tag.AsIFD(); err != nil {
			return fmt.Errorf("GPS IFD: %s", err)
		}
		if h.gpsIFDProvider, err = gpsifd.New(h.gpsIFD); err != nil {
			return err
		}
		h.providers = append(h.providers, h.gpsIFDProvider)
	}
	return nil
}
