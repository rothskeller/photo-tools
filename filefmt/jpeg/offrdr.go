package jpeg

import "io"

// OffsetReader tracks the offset into a reader as reads progress.
type OffsetReader struct {
	r   io.Reader
	off int
}

// NewOffsetReader creates a new OffsetReader for the underlying reader.
func NewOffsetReader(r io.Reader) *OffsetReader {
	return &OffsetReader{r, 0}
}

// Read reads from the OffsetReader.
func (o *OffsetReader) Read(buf []byte) (n int, err error) {
	n, err = o.r.Read(buf)
	o.off += n
	return n, err
}

// ReadByte reads a byte from the OffsetReader.
func (o *OffsetReader) ReadByte() (b byte, err error) {
	var n int
	if r, ok := o.r.(interface{ ReadByte() (byte, error) }); ok {
		b, err = r.ReadByte()
		if err == nil {
			n = 1
		}
	} else {
		var buf [1]byte
		n, err = o.r.Read(buf[:])
		b = buf[0]
	}
	o.off += n
	return b, err
}

// Offset returns the current offset into the underlying reader.
func (o *OffsetReader) Offset() int {
	return o.off
}
