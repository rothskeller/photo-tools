// Package tiff contains the file format handler for TIFF files.
package tiff

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/iim"
	"github.com/rothskeller/photo-tools/metadata/containers/photoshop"
	"github.com/rothskeller/photo-tools/metadata/containers/raw"
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
	container *tiff.TIFF
	tiffIFD0  *tiff.IFD
	psirBlock *photoshop.Photoshop
	psirIIM   *iim.IIM
	providers multi.Provider
}

var tiffHeaderLE = []byte{0x49, 0x49, 0x2A, 0x00}
var tiffHeaderBE = []byte{0x4D, 0x4D, 0x00, 0x2A}

// Read reads the provided file.  It returns nil, nil, if the file is not a TIFF
// file.  It returns an error if the file is a TIFF file but ill-formed, or if a
// read error occurs.  It returns a TIFF file handler for the file if it is read
// successfully.
func Read(r metadata.Reader) (h *TIFF, err error) {
	var buf [4]byte

	if _, err = r.ReadAt(buf[0:4], 0); err == io.EOF {
		return nil, nil // can't read a signature, assume it's not TIFF
	} else if err != nil {
		return nil, err
	} else if !bytes.Equal(buf[0:4], tiffHeaderBE) && !bytes.Equal(buf[0:4], tiffHeaderLE) {
		return nil, nil // not a JPEG file
	}
	h = new(TIFF)
	h.container = new(tiff.TIFF)
	if err = h.container.Read(r); err != nil {
		return nil, err
	}
	if err = h.readIFD0(); err != nil {
		return nil, err
	}
	if err = h.readEXIFIFD(); err != nil {
		return nil, err
	}
	if err = h.readGPSIFD(); err != nil {
		return nil, err
	}
	if err = h.readXMPTag(); err != nil {
		return nil, err
	}
	if err = h.readIPTCTag(); err != nil {
		return nil, err
	}
	if err = h.readPSIRIFD(); err != nil {
		return nil, err
	}
	return h, nil
}

// Provider returns the metadata.Provider for the JPEG file.
func (h *TIFF) Provider() metadata.Provider { return h.providers }

func (h *TIFF) readIFD0() (err error) {
	var tiffIFD0Provider *tiffifd0.Provider

	if h.tiffIFD0 = h.container.IFD0(); h.tiffIFD0 == nil {
		return errors.New("TIFF: no IFD0")
	}
	if tiffIFD0Provider, err = tiffifd0.New(h.tiffIFD0); err != nil {
		return err
	}
	h.providers = append(h.providers, tiffIFD0Provider)
	return nil
}

func (h *TIFF) readXMPTag() (err error) {
	var (
		xmpRDF      *rdf.Packet
		xmpProvider *xmp.Provider
	)
	if tag := h.tiffIFD0.Tag(tagXMP); tag != nil {
		var r metadata.Reader
		if r, err = tag.AsUnknownReader(); err != nil {
			if r, err = tag.AsBytesReader(); err != nil {
				return fmt.Errorf("XMP: %s", err)
			}
		}
		xmpRDF = rdf.New()
		if err = xmpRDF.Read(r); err != nil {
			return fmt.Errorf("XMP: %s", err)
		}
		tag.SetContainer(xmpRDF)
		if xmpProvider, err = xmp.New(xmpRDF); err != nil {
			return err
		}
		h.providers = append(h.providers, xmpProvider)
	}
	return nil
}

func (h *TIFF) readIPTCTag() (err error) {
	var (
		iimc         *iim.IIM
		iptcProvider *iptc.Provider
	)
	if tag := h.tiffIFD0.Tag(tagIPTC); tag != nil {
		// Empirically, tag type can be UNKNOWN or LONG.
		var r metadata.Reader
		if r, err = tag.AsUnknownReader(); err != nil {
			if r, err = tag.AsLongReader(); err != nil {
				return fmt.Errorf("IPTC: %s", err)
			}
		}
		iimc = new(iim.IIM)
		if err = iimc.Read(r); err != nil {
			return fmt.Errorf("IPTC: %s", err)
		}
		tag.SetContainer(iimc)
		if iptcProvider, err = iptc.New(iimc); err != nil {
			return err
		}
		h.providers = append(h.providers, iptcProvider)
	}
	return nil
}

func (h *TIFF) readPSIRIFD() (err error) {
	var (
		iimPSIR          *photoshop.PSIR
		psirIPTCProvider *iptc.Provider
	)
	if tag := h.tiffIFD0.Tag(tagPSIR); tag != nil {
		var r metadata.Reader
		// Empirically, tag type could be either UNKNOWN or BYTE.
		if r, err = tag.AsUnknownReader(); err != nil {
			if r, err = tag.AsBytesReader(); err != nil {
				return fmt.Errorf("PSIR: %s", err)
			}
		}
		h.psirBlock = new(photoshop.Photoshop)
		if err = h.psirBlock.Read(r); err != nil {
			return fmt.Errorf("PSIR block: %s", err)
		}
		tag.SetContainer(h.psirBlock)
		if iimPSIR = h.psirBlock.PSIR(psirIDIIM); iimPSIR == nil {
			return nil
		}
		h.psirIIM = new(iim.IIM)
		if err = h.psirIIM.Read(iimPSIR.Reader()); err != nil {
			return fmt.Errorf("PSIR: %s", err)
		}
		iimPSIR.SetContainer(h.psirIIM)
		if psirIPTCProvider, err = iptc.New(h.psirIIM); err != nil {
			return fmt.Errorf("PSIR: %s", err)
		}
		h.providers = append(h.providers, psirIPTCProvider)
	}
	return nil
}

func (h *TIFF) readEXIFIFD() (err error) {
	var (
		exifIFD         *tiff.IFD
		exifIFDProvider *exififd.Provider
	)
	if tag := h.tiffIFD0.Tag(tagEXIFIFD); tag != nil {
		if exifIFD, err = tag.AsIFD(); err != nil {
			return fmt.Errorf("EXIF IFD: %s", err)
		}
		if exifIFDProvider, err = exififd.New(exifIFD, h.container.Encoding()); err != nil {
			return err
		}
		h.providers = append(h.providers, exifIFDProvider)
	}
	return nil
}

func (h *TIFF) readGPSIFD() (err error) {
	var (
		gpsIFD         *tiff.IFD
		gpsIFDProvider *gpsifd.Provider
	)
	if tag := h.tiffIFD0.Tag(tagGPSIFD); tag != nil {
		if gpsIFD, err = tag.AsIFD(); err != nil {
			return fmt.Errorf("GPS IFD: %s", err)
		}
		if gpsIFDProvider, err = gpsifd.New(gpsIFD); err != nil {
			return err
		}
		h.providers = append(h.providers, gpsIFDProvider)
	}
	return nil
}

// Dirty returns whether the metadata from the file have been changed since they
// were read (and therefore need to be saved).
func (h *TIFF) Dirty() bool {
	return h.container.Dirty() ||
		(h.psirIIM != nil && h.psirIIM.Dirty())
}

// Save writes the entire file to the supplied writer, including all revised
// metadata.
func (h *TIFF) Save(out io.Writer) (err error) {
	if h.psirIIM != nil {
		var (
			hashRaw  *raw.Raw
			hashPSIR *photoshop.PSIR
		)
		hashRaw = new(raw.Raw)
		if hashPSIR = h.psirBlock.PSIR(psirIDHash); hashPSIR != nil {
			hashPSIR.SetContainer(hashRaw)
		} else {
			h.psirBlock.AddPSIR(psirIDHash, "", hashRaw)
		}
		h.psirIIM.SetHashContainer(hashRaw)
		hashRaw.SetData(make([]byte, 16)) // give it correct size
	}
	_, err = h.container.Write(out)
	return err
}
