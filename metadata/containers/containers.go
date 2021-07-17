// Package containers defines the interface that all containers must satisfy.
package containers

import (
	"io"

	"github.com/rothskeller/photo-tools/metadata"
)

// Container is the interface that all containers must satisfy.  Note that, in
// addition to these methods, each container will provide container-specific
// methods for reading and changing the contents of the container.
//
// Containers must take care not to consume arbitrarily large amounts of memory.
// If the container's data can be arbitrarily large, it should not be cached in
// memory; instead the container should remember the structure only, and read
// the actual data piece by piece when requested.  Similarly, the Size and Write
// operations should avoid building an arbitrarily large copy of the data in
// memory while writing, instead copying data piece by piece from the Reader
// directly to the Writer.
type Container interface {
	// Read reads and parses the container structure from the supplied
	// Reader.  The reader will continue to be used after Read returns, and
	// must remain open and usable as long as the Container is in scope.
	Read(r metadata.Reader) error
	// Empty returns whether the container is empty (and should therefore
	// be omitted from the written file, along with whatever tag in the
	// parent container points to it).
	Empty() bool
	// Dirty returns whether the contents of the container have been
	// changed.
	Dirty() bool
	// Layout computes the rendered layout of the container, i.e. prepares
	// for a call to Write, and returns what the rendered size of the
	// container will be.
	Layout() int64
	// Write writes the rendered container to the specified writer.
	Write(w io.Writer) (n int, err error)
}
