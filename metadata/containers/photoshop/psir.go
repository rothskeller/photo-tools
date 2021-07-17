// Package photoshop handles containers of Photoshop Image Resources (PSIRs).
package photoshop

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers"
)

// A PSIR is a single Photoshop Image Resource.
type PSIR struct {
	id        uint16
	name      string
	reader    metadata.Reader
	container containers.Container
	size      int64
}

var _ containers.Container = (*PSIR)(nil) // verify interface compliance

var psirResourceType = []byte("8BIM")

// Read reads and parses the first PSIR in the supplied Reader, leaving the seek
// offset at the end of it.  It returns io.EOF if the reader is at its end and
// there are no more PSIRs left in it.
func (psir *PSIR) Read(r metadata.Reader) (err error) {
	var (
		offset  int64
		size    int64
		buf     [266]byte
		count   int
		nameend int
	)
	if offset, err = r.Seek(0, io.SeekCurrent); err != nil {
		panic(err)
	}
	if count, err = r.Read(buf[:]); err == io.EOF && count == 0 {
		return io.EOF
	} else if err != nil && err != io.EOF {
		return fmt.Errorf("Photoshop: %s", err)
	} else if count < 12 || !bytes.Equal(buf[0:4], psirResourceType) {
		return errors.New("Photoshop: invalid PSIR header")
	}
	psir.id = binary.BigEndian.Uint16(buf[4:6])
	if count < 11+int(buf[6]) {
		return errors.New("Photoshop: incomplete PSIR header")
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
	if _, err = r.Seek(offset, io.SeekStart); err != nil {
		panic(err)
	}
	return nil
}

// Dirty returns whether the contents of the container have been
// changed.
func (psir *PSIR) Dirty() bool {
	if psir == nil || psir.container == nil {
		return false
	}
	return psir.container.Dirty()
}

// Size returns the rendered size of the container, in bytes.
func (psir *PSIR) Size() int64 {
	if psir == nil {
		return 0
	}
	psir.size = 11 + int64(len(psir.name))
	if psir.size%2 == 1 {
		psir.size++
	}
	if psir.container != nil {
		psir.size += psir.container.Size()
	} else {
		psir.size += psir.reader.Size()
	}
	return psir.size
}

// Write writes the rendered container to the specified writer.
func (psir *PSIR) Write(w io.Writer) (count int, err error) {
	var (
		buf [7]byte
		n   int
		n64 int64
	)
	copy(buf[0:4], psirResourceType)
	binary.BigEndian.PutUint16(buf[4:6], psir.id)
	buf[6] = byte(len(psir.name))
	n, err = w.Write(buf[0:7])
	count += n
	if err != nil {
		return count, err
	}
	if psir.name != "" {
		n, err = w.Write([]byte(psir.name))
		count += n
		if err != nil {
			return count, err
		}
	}
	if len(psir.name)%2 == 0 {
		buf[0] = 0
		n, err = w.Write(buf[0:1])
		count += n
		if err != nil {
			return count, err
		}
	}
	if psir.container != nil {
		binary.BigEndian.PutUint32(buf[0:4], uint32(psir.container.Size()))
	} else {
		binary.BigEndian.PutUint32(buf[0:4], uint32(psir.reader.Size()))
	}
	n, err = w.Write(buf[0:4])
	count += n
	if err != nil {
		return count, err
	}
	if psir.container != nil {
		n, err = psir.container.Write(w)
	} else {
		n64, err = io.Copy(w, psir.reader)
		n = int(n64)
	}
	count += n
	if err != nil {
		return count, err
	}
	if psir.size != 0 && int(psir.size) != count {
		panic("actual size different from predicted size")
	}
	return count, err
}

// Name returns the name of the PSIR.
func (psir *PSIR) Name() string { return psir.name }

// Reader returns the reader for the PSIR.
func (psir *PSIR) Reader() metadata.Reader { return psir.reader }

// SetContainer sets the container embedded within a PSIR.
func (psir *PSIR) SetContainer(c containers.Container) {
	psir.container = c
}
