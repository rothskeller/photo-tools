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

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers"
	"github.com/rothskeller/photo-tools/metadata/containers/raw"
)

// An IIM structure represents the entire IIM block.
type IIM struct {
	dsmap         map[uint16][]DataSet
	hash          []byte
	dirty         bool
	size          int64
	hashContainer *raw.Raw
}

var _ containers.Container = (*IIM)(nil) // verify interface compliance

// A DataSet is a single data set within an IIM block.
type DataSet struct {
	ID   uint16
	Data []byte
}

// New creates a new IIM block.
func New() *IIM {
	return &IIM{dsmap: make(map[uint16][]DataSet)}
}

// Read parses the IIM block in the supplied reader.
func (iim *IIM) Read(r metadata.Reader) (err error) {
	var (
		offset int64
		ds     DataSet
		buf    [8]byte
		count  int
		size   uint64
		sum    hash.Hash
	)
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("IIM: %s", err)
	}
	sum = md5.New()
	if _, err = io.Copy(sum, r); err != nil {
		return fmt.Errorf("IIM: %s", err)
	}
	for {
		count, err = r.ReadAt(buf[0:5], offset)
		if err == io.EOF && count == 0 ||
			(count == 1 && buf[0] == 0) ||
			(count == 2 && buf[0] == 0 && buf[1] == 0) {
			// In TIFF files, IPTC data is sometimes stored as LONGs
			// and can have one or two null bytes at the end.
			iim.hash = sum.Sum(nil)
			return nil
		}
		if err != nil {
			return fmt.Errorf("IIM: %s", err)
		}
		offset += 5
		if buf[0] != 0x1C {
			return errors.New("IIM: invalid IIM tag marker")
		}
		ds.ID = binary.BigEndian.Uint16(buf[1:3])
		size = uint64(binary.BigEndian.Uint16(buf[3:5]))
		if size&0x8000 != 0 {
			size -= 0x8000
			if size < 1 || size > 8 {
				return errors.New("IIM: unsupported IIM tag size")
			}
			copy(buf[:], []byte{0, 0, 0, 0, 0, 0, 0, 0})
			if _, err = r.ReadAt(buf[8-size:8], offset); err != nil {
				return fmt.Errorf("IIM: %s", err)
			}
			offset += int64(size)
			size = binary.BigEndian.Uint64(buf[:])
			if size > 0x100000 {
				return errors.New("IIM: unreasonable IIM tag size")
			}
		}
		ds.Data = make([]byte, size)
		if _, err = r.ReadAt(ds.Data, offset); err != nil {
			return fmt.Errorf("IIM: %s", err)
		}
		offset += int64(size)
		iim.dsmap[ds.ID] = append(iim.dsmap[ds.ID], ds)
	}
}

// DataSets returns the list of data sets with the specified ID.  If no such
// data sets exist, it returns an empty list.
func (iim *IIM) DataSets(id uint16) []DataSet {
	return iim.dsmap[id]
}

// SetDataSet puts the provided data set into the IIM block, replacing any other
// data sets with the same ID.  It also marks the IIM block as dirty; callers
// should not call SetDataSet if the data set is unchanged.
func (iim *IIM) SetDataSet(id uint16, data []byte) {
	iim.dsmap[id] = []DataSet{{id, data}}
	iim.dirty = true
}

// SetDataSets puts the provided data sets into the IIM block, replacing any
// other data sets with the same ID.  It also marks the IIM block as dirty;
// callers should not call SetDataSet if the data sets are unchanged.
func (iim *IIM) SetDataSets(id uint16, data [][]byte) {
	var dss = make([]DataSet, len(data))
	for i := range data {
		dss[i].ID = id
		dss[i].Data = data[i]
	}
	iim.dsmap[id] = dss
	iim.dirty = true
}

// RemoveDataSets removes all data sets with the specified ID from the IIM
// block.  If any such data sets existed, the block is marked as dirty.
func (iim *IIM) RemoveDataSets(id uint16) {
	if _, ok := iim.dsmap[id]; ok {
		delete(iim.dsmap, id)
		iim.dirty = true
	}
}

// SetHashContainer sets the container into which the IIM block's hash should be
// written when the block is written.
func (iim *IIM) SetHashContainer(hc *raw.Raw) { iim.hashContainer = hc }

// Empty returns whether the container is empty (and should therefore be omitted
// from the written file, along with whatever tag in the parent container points
// to it).
func (iim *IIM) Empty() bool {
	for _, dss := range iim.dsmap {
		if len(dss) != 0 {
			return false
		}
	}
	return true
}

// Dirty returns whether the IIM block has been changed since it was read.
func (iim *IIM) Dirty() bool { return iim.dirty }

// Layout computes the rendered layout of the container, i.e. prepares for a
// call to Write, and returns what the rendered size of the container will be.
func (iim *IIM) Layout() int64 {
	iim.size = 0
	for _, dss := range iim.dsmap {
		for _, ds := range dss {
			if len(ds.Data) > 0x7FFF {
				iim.size += int64(len(ds.Data)) + 9
			} else {
				iim.size += int64(len(ds.Data)) + 5
			}
		}
	}
	return iim.size
}

// Write writes the IIM block to the specified writer.
func (iim *IIM) Write(w io.Writer) (count int, err error) {
	var (
		buf [9]byte
		mw  io.Writer
		n   int
		ids = make([]uint16, 0, len(iim.dsmap))
	)
	sum := md5.New()
	mw = io.MultiWriter(w, sum)
	for id, ds := range iim.dsmap {
		if len(ds) != 0 {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for _, id := range ids {
		for _, ds := range iim.dsmap[id] {
			buf[0] = 0x1C
			binary.BigEndian.PutUint16(buf[1:3], id)
			if len(ds.Data) > 0x7FFF {
				buf[3] = 0x80
				buf[4] = 0x04
				binary.BigEndian.PutUint32(buf[5:9], uint32(len(ds.Data)))
				n, err = mw.Write(buf[0:9])
				count += n
			} else {
				binary.BigEndian.PutUint16(buf[3:5], uint16(len(ds.Data)))
				n, err = mw.Write(buf[0:5])
				count += n
			}
			if err != nil {
				return count, err
			}
			n, err = mw.Write(ds.Data)
			count += n
			if err != nil {
				return count, err
			}
		}
	}
	iim.hash = sum.Sum(nil)
	if iim.hashContainer != nil {
		iim.hashContainer.SetData(iim.hash)
	}
	if iim.size != 0 && int(iim.size) != count {
		panic("actual size different from predicted size")
	}
	return count, nil
}

// Hash returns the MD5 hash of the IIM block.  It does not reflect any changes
// until after a call to Write.
func (iim *IIM) Hash() []byte { return iim.hash }
