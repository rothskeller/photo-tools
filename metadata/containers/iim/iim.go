// Package iim handles the marshaling and unmarshaling of IPTC IIM metadata.
package iim

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	"sort"
)

// An IIM structure represents the entire IIM block.
type IIM map[uint16][]DataSet

// A DataSet is a single data set within an IIM block.
type DataSet struct {
	ID   uint16
	Data []byte
}

// Reader is the interface that must be satisfied by the parameter to Read.
type Reader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	Size() int64
}

// Read parses the IIM block in the supplied reader.
func Read(r Reader) (iim IIM, md5sum []byte, err error) {
	var (
		offset int64
		ds     DataSet
		buf    [8]byte
		count  int
		size   uint64
		sum    hash.Hash
	)
	iim = make(IIM)
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return nil, nil, fmt.Errorf("IIM: %s", err)
	}
	sum = md5.New()
	if _, err = io.Copy(sum, r); err != nil {
		return nil, nil, fmt.Errorf("IIM: %s", err)
	}
	for {
		count, err = r.ReadAt(buf[0:5], offset)
		if err == io.EOF && count == 0 {
			return iim, sum.Sum(nil), nil
		}
		if err != nil {
			return nil, nil, fmt.Errorf("IIM: %s", err)
		}
		offset += 5
		if buf[0] != 0x1C {
			return nil, nil, errors.New("IIM: invalid IIM tag marker")
		}
		ds.ID = binary.BigEndian.Uint16(buf[1:3])
		size = uint64(binary.BigEndian.Uint16(buf[3:5]))
		if size&0x8000 != 0 {
			size -= 0x8000
			if size < 1 || size > 8 {
				return nil, nil, errors.New("IIM: unsupported IIM tag size")
			}
			copy(buf[:], []byte{0, 0, 0, 0, 0, 0, 0, 0})
			if _, err = r.ReadAt(buf[8-size:8], offset); err != nil {
				return nil, nil, fmt.Errorf("IIM: %s", err)
			}
			offset += int64(size)
			size = binary.BigEndian.Uint64(buf[:])
			if size > 0x100000 {
				return nil, nil, errors.New("IIM: unreasonable IIM tag size")
			}
		}
		ds.Data = make([]byte, size)
		if _, err = r.ReadAt(ds.Data, offset); err != nil {
			return nil, nil, fmt.Errorf("IIM: %s", err)
		}
		offset += int64(size)
		iim[ds.ID] = append(iim[ds.ID], ds)
	}
}

// Render writes the IIM block to the specified writer.
func (iim IIM) Render(w io.Writer) (sum hash.Hash, err error) {
	var (
		buf [9]byte
		mw  io.Writer
		ids = make([]uint16, 0, len(iim))
	)
	sum = md5.New()
	mw = io.MultiWriter(w, sum)
	for id, ds := range iim {
		if len(ds) != 0 {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for _, id := range ids {
		for _, ds := range iim[id] {
			buf[0] = 0x1C
			binary.BigEndian.PutUint16(buf[1:3], id)
			if len(ds.Data) > 0x7FFF {
				buf[3] = 0x80
				buf[4] = 0x04
				binary.BigEndian.PutUint32(buf[5:9], uint32(len(ds.Data)))
				_, err = mw.Write(buf[0:9])
			} else {
				binary.BigEndian.PutUint16(buf[3:5], uint16(len(ds.Data)))
				_, err = mw.Write(buf[0:5])
			}
			if err != nil {
				return nil, err
			}
			if _, err = mw.Write(ds.Data); err != nil {
				return nil, err
			}
		}
	}
	return sum, nil
}