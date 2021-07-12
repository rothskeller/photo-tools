// Package photoshop handles containers of Photoshop Image Resources (PSIRs).
package photoshop

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/rothskeller/photo-tools/metadata"
)

// Photoshop is a handler for containers of Photoshop Image Resources (PSIRs).
type Photoshop struct {
	psirs []PSIR
	dirty bool
}

// A PSIR is a single Photoshop Image Resource.
type PSIR struct {
	id     uint16
	name   string
	reader metadata.Reader
	ps     *Photoshop
}

var psirResourceType = []byte("8BIM")

// Read parses the Photoshop block and returns a handler for it.
func Read(r metadata.Reader) (ps *Photoshop, err error) {
	var (
		offset  int64
		size    int64
		buf     [266]byte
		count   int
		nameend int
		psir    PSIR
		seen    = make(map[uint16]bool)
	)
	ps = new(Photoshop)
	for {
		if count, err = r.ReadAt(buf[:], offset); err == io.EOF && count == 0 {
			return ps, nil
		} else if err != nil && err != io.EOF {
			return nil, fmt.Errorf("Photoshop: %s", err)
		} else if count < 12 || !bytes.Equal(buf[0:4], psirResourceType) {
			return nil, errors.New("Photoshop: invalid PSIR header")
		}
		psir.ps = ps
		psir.id = binary.BigEndian.Uint16(buf[4:6])
		if seen[psir.id] {
			return nil, fmt.Errorf("Photoshop: multiple PSIRs with ID 0x%04x", psir.id)
		}
		seen[psir.id] = true
		if count < 11+int(buf[6]) {
			return nil, fmt.Errorf("Photoshop: incomplete PSIR header")
		}
		nameend = 7 + int(buf[6])
		psir.name = string(buf[7:nameend])
		if nameend%2 == 1 {
			nameend++
		}
		size = int64(binary.BigEndian.Uint32(buf[nameend : nameend+4]))
		offset += int64(nameend) + 4
		psir.reader = io.NewSectionReader(r, offset, size)
		offset += size
		if offset%2 == 1 {
			offset++
		}
		ps.psirs = append(ps.psirs, psir)
	}
}

// Dirty returns whether any PSIRs has been changed.
func (ps *Photoshop) Dirty() bool { return ps.dirty }

// PSIR returns the PSIR with the specified ID, or nil if there is none.
func (ps *Photoshop) PSIR(id uint16) *PSIR {
	for i, psir := range ps.psirs {
		if psir.id == id {
			return &ps.psirs[i]
		}
	}
	return nil
}

// Name returns the name of the PSIR.
func (psir *PSIR) Name() string { return psir.name }

// Reader returns the reader for the PSIR.
func (psir *PSIR) Reader() metadata.Reader { return psir.reader }

// SetReader sets the name and/or reader for a PSIR.
func (psir *PSIR) SetReader(reader metadata.Reader) {
	psir.reader = reader
	psir.ps.dirty = true
}

// AddPSIR adds a PSIR to the Photoshop block.
func (ps *Photoshop) AddPSIR(id uint16, name string, reader metadata.Reader) {
	var firstAfter = -1

	for i, ps := range ps.psirs {
		if ps.id == id {
			panic("adding a redundant PSIR")
		}
		if ps.id > id && firstAfter == -1 {
			firstAfter = i
		}
	}
	if firstAfter == -1 {
		ps.psirs = append(ps.psirs, PSIR{id, name, reader, ps})
	} else {
		ps.psirs = append(ps.psirs, PSIR{})
		copy(ps.psirs[firstAfter+1:], ps.psirs[firstAfter:])
		ps.psirs[firstAfter] = PSIR{id, name, reader, ps}
	}
	ps.dirty = true
}
