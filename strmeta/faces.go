package strmeta

import (
	"errors"
	"sort"
)

// GetFaces returns faces from the highest priority faces tag.
func GetFaces(h fileHandler) (faces []string) {
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.MPFaces) != 0 {
			sort.Strings(xmp.MPFaces)
			return xmp.MPFaces
		} else if len(xmp.MWGRSFaces) != 0 {
			sort.Strings(xmp.MWGRSFaces)
			return xmp.MWGRSFaces
		}
	}
	return nil
}

// GetFacesTags returns all of the faces tags and their values.
func GetFacesTags(h fileHandler) (tags []string, values []string) {
	if xmp := h.XMP(false); xmp != nil {
		for _, face := range xmp.MPFaces {
			tags = append(tags, "XMP.MP:Regions")
			values = append(values, face)
		}
		for _, face := range xmp.MWGRSFaces {
			tags = append(tags, "XMP.mwg-rs:RegionInfo")
			values = append(values, face)
		}
	}
	return tags, values
}

// CheckFaces checks whether the faces are correctly tagged, and are consistent
// with the reference.
func CheckFaces(ref, h fileHandler) (res CheckResult) {
	var value = GetFaces(ref)
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.MPFaces) != 0 {
			sort.Strings(xmp.MPFaces)
			if !stringSliceEqual(value, xmp.MPFaces) {
				return ChkConflictingValues
			}
		} else if len(value) != 0 {
			res = ChkIncorrectlyTagged
		}
		if len(xmp.MWGRSFaces) != 0 {
			sort.Strings(xmp.MWGRSFaces)
			if !stringSliceEqual(value, xmp.MWGRSFaces) {
				return ChkConflictingValues
			}
		} else if len(value) != 0 {
			res = ChkIncorrectlyTagged
		}
	}
	if len(value) != 0 && res == 0 {
		return ChkPresent
	}
	return res
}
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// SetFaces is called to set the face region tags.  It returns "success" if
// there are no actual changes.  Any actual changes cause an error return, since
// the xmp library cannot write those tags.
func SetFaces(h fileHandler, v []string) error {
	sort.Strings(v)
	if xmp := h.XMP(true); xmp != nil {
		sort.Strings(xmp.MPFaces)
		if !stringSliceEqual(v, xmp.MPFaces) {
			return errors.New("this tool cannot change face region tags")
		}
		sort.Strings(xmp.MWGRSFaces)
		if !stringSliceEqual(v, xmp.MWGRSFaces) {
			return errors.New("this tool cannot change face region tags")
		}
	}
	return nil
}
