// Package jpeg contains the file format handler for JPEG files.
package jpeg

import (
	"fmt"
	"io"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/iim"
	"github.com/rothskeller/photo-tools/metadata/containers/jpeg"
	"github.com/rothskeller/photo-tools/metadata/containers/photoshop"
	"github.com/rothskeller/photo-tools/metadata/containers/raw"
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
	container *jpeg.JPEG
	exifTIFF  *tiff.TIFF
	psirBlock *photoshop.Photoshop
	iim       *iim.IIM
	xmpRDF    *rdf.Packet
	providers multi.Provider
}

// Read reads the provided file.  It returns nil, nil, if the file is not a JPEG
// file.  It returns an error if the file is a JPEG file but ill-formed, or if a
// read error occurs.  It returns a JPEG file handler for the file if it is read
// successfully.
func Read(r metadata.Reader) (jh *JPEG, err error) {
	var buf [2]byte

	if _, err = r.ReadAt(buf[0:2], 0); err == io.EOF {
		return nil, nil // can't read a signature, assume it's not JPEG
	} else if err != nil {
		return nil, err
	} else if buf[0] != 0xFF || buf[1] != 0xD8 {
		return nil, nil // not a JPEG file
	}
	jh = new(JPEG)
	jh.container = new(jpeg.JPEG)
	if err = jh.container.Read(r); err != nil {
		return nil, err
	}
	if err = jh.readEXIFSegment(); err != nil {
		return nil, err
	}
	if err = jh.readXMPSegments(); err != nil {
		return nil, err
	}
	if err = jh.readPSIRSegment(); err != nil {
		return nil, err
	}
	return jh, nil
}

// Provider returns the metadata.Provider for the JPEG file.
func (jh *JPEG) Provider() metadata.Provider { return jh.providers }

func (jh *JPEG) readEXIFSegment() (err error) {
	var (
		exifSeg          metadata.Reader
		jpegIFD0         *tiff.IFD
		exifIFD          *tiff.IFD
		gpsIFD           *tiff.IFD
		jpegIFD0Provider *jpegifd0.Provider
		exifIFDProvider  *exififd.Provider
		gpsIFDProvider   *gpsifd.Provider
	)
	jh.exifTIFF = new(tiff.TIFF)
	if exifSeg = jh.container.EXIF(); exifSeg != nil {
		if err = jh.exifTIFF.Read(exifSeg); err != nil {
			return fmt.Errorf("EXIF segment: %s", err)
		}
	}
	jh.container.SetEXIFContainer(jh.exifTIFF)
	jpegIFD0 = jh.exifTIFF.IFD0()
	if jpegIFD0Provider, err = jpegifd0.New(jpegIFD0); err != nil {
		return err
	}
	jh.providers = append(jh.providers, jpegIFD0Provider)
	if tag := jpegIFD0.Tag(tagEXIFIFD); tag != nil {
		if exifIFD, err = tag.AsIFD(); err != nil {
			return fmt.Errorf("EXIF IFD: %s", err)
		}
	} else {
		tag = jpegIFD0.AddTag(tagEXIFIFD, 4)
		exifIFD, _ = tag.AddIFD()
	}
	if exifIFDProvider, err = exififd.New(exifIFD, jh.exifTIFF.Encoding()); err != nil {
		return err
	}
	jh.providers = append(jh.providers, exifIFDProvider)
	if tag := jpegIFD0.Tag(tagGPSIFD); tag != nil {
		if gpsIFD, err = tag.AsIFD(); err != nil {
			return fmt.Errorf("GPS IFD: %s", err)
		}
	} else {
		tag = jpegIFD0.AddTag(tagGPSIFD, 4)
		gpsIFD, _ = tag.AddIFD()
	}
	if gpsIFDProvider, err = gpsifd.New(gpsIFD); err != nil {
		return err
	}
	jh.providers = append(jh.providers, gpsIFDProvider)
	return nil
}

func (jh *JPEG) readPSIRSegment() (err error) {
	var (
		psirSeg      metadata.Reader
		iimPSIR      *photoshop.PSIR
		iptcProvider *iptc.Provider
	)
	jh.psirBlock = new(photoshop.Photoshop)
	if psirSeg = jh.container.PSIR(); psirSeg != nil {
		if err = jh.psirBlock.Read(psirSeg); err != nil {
			return fmt.Errorf("PSIR segment: %s", err)
		}
	}
	jh.container.SetPSIRContainer(jh.psirBlock)
	jh.iim = iim.New()
	if iimPSIR = jh.psirBlock.PSIR(psirIDIIM); iimPSIR != nil {
		if err = jh.iim.Read(iimPSIR.Reader()); err != nil {
			return err
		}
		iimPSIR.SetContainer(jh.iim)
	} else {
		jh.psirBlock.AddPSIR(psirIDIIM, "", jh.iim)
	}
	if iptcProvider, err = iptc.New(jh.iim); err != nil {
		return err
	}
	jh.providers = append(jh.providers, iptcProvider)
	return nil
}

func (jh *JPEG) readXMPSegments() (err error) {
	var (
		xmpSeg         metadata.Reader
		xmpExtSeg      metadata.Reader
		xmpExtRDF      *rdf.Packet
		xmpProvider    *xmp.Provider
		xmpExtProvider *xmpext.Provider
	)
	jh.xmpRDF = rdf.New()
	if xmpSeg = jh.container.XMP(); xmpSeg != nil {
		if err = jh.xmpRDF.Read(xmpSeg); err != nil {
			return fmt.Errorf("XMP: %s", err)
		}
	}
	jh.container.SetXMPContainer(jh.xmpRDF)
	if xmpProvider, err = xmp.New(jh.xmpRDF); err != nil {
		return err
	}
	jh.providers = append(jh.providers, xmpProvider)
	xmpExtRDF = rdf.New()
	if xmpExtSeg = jh.container.XMPext(); xmpExtSeg != nil {
		if err = xmpExtRDF.Read(xmpExtSeg); err != nil {
			return fmt.Errorf("XMPExt: %s", err)
		}
	}
	if xmpExtProvider, err = xmpext.New(xmpExtRDF); err != nil {
		return err
	}
	jh.providers = append(jh.providers, xmpExtProvider)
	return nil
}

// Dirty returns whether the metadata from the file have been changed since they
// were read (and therefore need to be saved).
func (jh *JPEG) Dirty() bool { return jh.container.Dirty() }

// Save writes the entire file to the supplied writer, including all revised
// metadata.
func (jh *JPEG) Save(out io.Writer) (err error) {
	if !jh.iim.Empty() {
		var (
			hashRaw  *raw.Raw
			hashPSIR *photoshop.PSIR
		)
		hashRaw = new(raw.Raw)
		if hashPSIR = jh.psirBlock.PSIR(psirIDHash); hashPSIR != nil {
			hashPSIR.SetContainer(hashRaw)
		} else {
			jh.psirBlock.AddPSIR(psirIDHash, "", hashRaw)
		}
		jh.iim.SetHashContainer(hashRaw)
		hashRaw.SetData(make([]byte, 16)) // give it correct size
	}
	jh.container.Layout()
	_, err = jh.container.Write(out)
	return err
}
