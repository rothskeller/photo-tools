// Package photoshop handles containers of Photoshop Image Resources (PSIRs).
package photoshop

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Photoshop is a handler for containers of Photoshop Image Resources (PSIRs).
type Photoshop struct {
	psirs []PSIR
}

// A PSIR is a single Photoshop Image Resource.
type PSIR struct {
	ID     uint16
	Name   string
	Reader Reader
}

// Reader is the interface that must be honored by the parameter to Read.
type Reader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	Size() int64
}

var psirResourceType = []byte("8BIM")

// Read parses the Photoshop block and returns a handler for it.
func Read(r Reader) (ps *Photoshop, err error) {
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
		psir.ID = binary.BigEndian.Uint16(buf[4:6])
		if seen[psir.ID] {
			return nil, fmt.Errorf("Photoshop: multiple PSIRs with ID 0x%04x", psir.ID)
		}
		seen[psir.ID] = true
		if count < 11+int(buf[6]) {
			return nil, fmt.Errorf("Photoshop: incomplete PSIR header")
		}
		nameend = 7 + int(buf[6])
		psir.Name = string(buf[7:nameend])
		if nameend%2 == 1 {
			nameend++
		}
		size = int64(binary.BigEndian.Uint32(buf[nameend : nameend+4]))
		offset += int64(nameend) + 4
		psir.Reader = io.NewSectionReader(r, offset, size)
		offset += size
		if offset%2 == 1 {
			offset++
		}
	}
}

// PSIR returns the PSIR with the specified ID, or nil if there is none.
func (ps *Photoshop) PSIR(id uint16) *PSIR {
	for i, psir := range ps.psirs {
		if psir.ID == id {
			return &ps.psirs[i]
		}
	}
	return nil
}

// AddPSIR adds a PSIR to the Photoshop block.
func (ps *Photoshop) AddPSIR(psir *PSIR) {
	var firstAfter = -1

	for i, ps := range ps.psirs {
		if ps.ID == psir.ID {
			panic("adding a redundant PSIR")
		}
		if ps.ID > psir.ID && firstAfter == -1 {
			firstAfter = i
		}
	}
	if firstAfter == -1 {
		ps.psirs = append(ps.psirs, *psir)
	} else {
		ps.psirs = append(ps.psirs, PSIR{})
		copy(ps.psirs[firstAfter+1:], ps.psirs[firstAfter:])
		ps.psirs[firstAfter] = *psir
	}
}
