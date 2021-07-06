package tifflike

import "errors"

// rangelist is a sorted list of non-overlapping ranges of the TIFF-like block
// that are filled with data.
type rangelist struct {
	r []uint32
}

// add adds the range [start, end) to the rangelist.  It returns an error if the
// range overlaps with one already consumed on the range list.
func (r *rangelist) add(start, end uint32) error {
	if start > end {
		panic("range starts after range end")
	}
	if start == end {
		return nil
	}
	// Find the lowest existing range that starts after our new one ends.
	idx := 0
	for idx < len(r.r) && r.r[idx] < end {
		idx += 2
	}
	if idx > 0 && start < r.r[idx-1] {
		// The one we're adding starts before the preceding one ends.
		return errors.New("overlapping range")
	}
	if idx > 0 && start == r.r[idx-1] {
		// New range is adjacent to previous range.
		if idx < len(r.r) && end == r.r[idx] {
			// It's also adjacent to the following one, so we can
			// consolidate the preceding, new, and following down to
			// one range.
			r.r[idx-1] = r.r[idx+1]
			r.r = append(r.r[:idx], r.r[idx+2:]...)
		} else {
			// It's not adjacent to the following one, so we extend
			// the preceding one.
			r.r[idx-1] = end
		}
	} else if idx < len(r.r) && end == r.r[idx] {
		// New range is adjacent to following range, so extend it.
		r.r[idx] = start
	} else {
		// New range isn't adjacent to either preceding or following,
		// so we need to insert it.
		r.r = append(r.r, 0, 0)           // ensure correct size
		copy(r.r[idx+2:], r.r[idx:])      // shift contents
		r.r[idx], r.r[idx+1] = start, end // add new range
	}
	return nil
}

// removeTrailer checks to see whether the last range is up against the end of
// the file.  If so, it removes it from the list, and returns the start of that
// range, which becomes our new end of file for rendering purposes.
func (r *rangelist) removeTrailer(eof uint32) uint32 {
	if len(r.r) == 0 || r.r[len(r.r)-1] != eof {
		return eof
	}
	eof = r.r[len(r.r)-2]
	r.r = r.r[:len(r.r)-2]
	return eof
}

// consume finds the largest available range with the requested size, removes
// that size from it, and returns the offset of the consumed area.  It returns
// zero if there is no available range with the requested size.  (Zero is not
// valid because the range from 0 to 8 is never added in any TIFF-like block;
// it's assigned to the header.)
func (r *rangelist) consume(size uint32) (offset uint32) {
	var largestIdx = -1
	var largestSize uint32 = 0
	for idx := 0; idx < len(r.r); idx += 2 {
		sz := r.r[idx+1] - r.r[idx]
		if sz >= size {
			if largestIdx == -1 || largestSize < sz {
				largestIdx = idx
				largestSize = sz
			}
		}
	}
	if largestIdx == -1 {
		return 0
	}
	offset = r.r[largestIdx]
	if largestSize > size {
		r.r[largestIdx] += size
	} else {
		r.r = append(r.r[:largestIdx], r.r[largestIdx+2:]...)
	}
	return offset
}
