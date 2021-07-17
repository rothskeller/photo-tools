// Package photoshop handles containers of Photoshop Image Resources (PSIRs).
package photoshop

import (
	"fmt"
	"io"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers"
)

// Photoshop is a handler for containers of Photoshop Image Resources (PSIRs).
type Photoshop struct {
	psirs []*PSIR
	size  int64
}

var _ containers.Container = (*Photoshop)(nil) // verify interface compliance

// Read reads and parses container structure from the supplied Reader.
func (ps *Photoshop) Read(r metadata.Reader) (err error) {
	var (
		start  int64
		offset int64
		seen   = make(map[uint16]bool)
	)
	if start, err = r.Seek(0, io.SeekCurrent); err != nil {
		panic(err)
	}
	for {
		var psir PSIR

		if err = psir.Read(r); err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			return nil
		}
		if seen[psir.id] {
			return fmt.Errorf("Photoshop: multiple PSIRs with id 0x%x", psir.id)
		}
		seen[psir.id] = true
		ps.psirs = append(ps.psirs, &psir)
		if offset, err = r.Seek(0, io.SeekCurrent); err != nil {
			panic(err)
		}
		if offset%2 != start%2 {
			if _, err = r.Seek(1, io.SeekCurrent); err != nil {
				panic(err)
			}
		}
	}
}

// Empty returns whether the container is empty (and should therefore be omitted
// from the written file, along with whatever tag in the parent container points
// to it).
func (ps *Photoshop) Empty() bool {
	for _, psir := range ps.psirs {
		if !psir.Empty() {
			return false
		}
	}
	return true
}

// Dirty returns whether any PSIRs have been changed.
func (ps *Photoshop) Dirty() bool {
	for _, psir := range ps.psirs {
		if psir.Dirty() {
			return true
		}
	}
	return false
}

// Layout computes the rendered layout of the container, i.e. prepares for a
// call to Write, and returns what the rendered size of the container will be.
func (ps *Photoshop) Layout() int64 {
	ps.size = 0
	for _, psir := range ps.psirs {
		if psir.Empty() {
			continue
		}
		if ps.size%2 == 1 {
			ps.size++
		}
		ps.size += psir.Layout()
	}
	return ps.size
}

// Write writes the rendered container to the specified writer.
func (ps *Photoshop) Write(w io.Writer) (count int, err error) {
	var n int

	for _, psir := range ps.psirs {
		if psir.Empty() {
			continue
		}
		if count%2 == 1 {
			n, err = w.Write([]byte{0})
			count += n
			if err != nil {
				return count, err
			}
		}
		n, err = psir.Write(w)
		count += n
		if err != nil {
			return count, err
		}
	}
	if ps.size != 0 && int(ps.size) != count {
		panic("actual size different from predicted size")
	}
	return count, nil
}

// PSIR returns the PSIR with the specified ID, or nil if there is none.
func (ps *Photoshop) PSIR(id uint16) *PSIR {
	for i, psir := range ps.psirs {
		if psir.id == id {
			return ps.psirs[i]
		}
	}
	return nil
}

// AddPSIR adds a PSIR to the Photoshop block.
func (ps *Photoshop) AddPSIR(id uint16, name string, c containers.Container) {
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
		ps.psirs = append(ps.psirs, &PSIR{id: id, name: name, container: c})
	} else {
		ps.psirs = append(ps.psirs, &PSIR{})
		copy(ps.psirs[firstAfter+1:], ps.psirs[firstAfter:])
		ps.psirs[firstAfter] = &PSIR{id: id, name: name, container: c}
	}
}
