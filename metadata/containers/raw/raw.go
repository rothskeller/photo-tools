// Package raw provides a "container" for a raw byte array.
package raw

import (
	"bytes"
	"io"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers"
)

// Raw is a container for a raw byte array.
type Raw struct {
	data  []byte
	dirty bool
}

var _ containers.Container = (*Raw)(nil) // verify interface compliance

// Read reads and parses the container structure from the supplied Reader.
func (raw *Raw) Read(r metadata.Reader) (err error) {
	raw.data = make([]byte, r.Size())
	_, err = r.ReadAt(raw.data, 0)
	return err
}

// Dirty returns whether the contents of the container have been
// changed.
func (raw *Raw) Dirty() bool { return raw.dirty }

// Size returns the rendered size of the container, in bytes.
func (raw *Raw) Size() int64 { return int64(len(raw.data)) }

// Write writes the rendered container to the specified writer.
func (raw *Raw) Write(w io.Writer) (n int, err error) { return w.Write(raw.data) }

// Data returns the data in the container.
func (raw *Raw) Data() (data []byte) {
	data = append([]byte{}, raw.data...)
	return data
}

// SetData sets the data in the container.
func (raw *Raw) SetData(data []byte) {
	if bytes.Equal(raw.data, data) {
		return
	}
	raw.data = append([]byte{}, data...)
	raw.dirty = true
}
