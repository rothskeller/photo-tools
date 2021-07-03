package strmeta

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/filefmt"
)

// GetFaces returns the names of face regions from the highest priority face
// region tag.
func GetFaces(h fileHandler) []string {
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.MPRegPersonDisplayNames()) != 0 {
			return xmp.MPRegPersonDisplayNames()
		}
		return xmp.MWGRSNames()
	}
	return nil
}

// GetFaceTags returns all of the face region tags and their values.
func GetFaceTags(h filefmt.FileHandler) (tags []string, values []string) {
	if xmp := h.XMP(false); xmp != nil {
		for _, face := range xmp.MPRegPersonDisplayNames() {
			tags = append(tags, "XMP  MP:Regions")
			values = append(values, face)
		}
		for _, face := range xmp.MWGRSNames() {
			tags = append(tags, "XMP  mwg-rs:RegionInfo")
			values = append(values, face)
		}
	}
	return tags, values
}

// CheckFaces determines whether the face regions are tagged correctly, and are
// consistent with the reference.
func CheckFaces(ref, h filefmt.FileHandler) (res CheckResult) {
	var value = GetFaces(ref)

	if xmp := h.XMP(false); xmp != nil {
		if !stringSliceEqualUnordered(value, xmp.MPRegPersonDisplayNames()) {
			if len(xmp.MPRegPersonDisplayNames()) != 0 {
				return ChkConflictingValues
			}
			res = ChkIncorrectlyTagged
		}
		if !stringSliceEqualUnordered(value, xmp.MWGRSNames()) {
			if len(xmp.MWGRSNames()) != 0 {
				return ChkConflictingValues
			}
			res = ChkIncorrectlyTagged
		}
	}
	// We also want to call them inconsistently tagged if there are missing
	// person tags for them.
	var people = getFilteredKeywords(ref, personPredicate, false)
	for _, face := range value {
		var found = false
		for _, kw := range people {
			if face == kw[1] {
				found = true
				break
			}
		}
		if !found {
			res = ChkIncorrectlyTagged
			break
		}
	}
	if res == 0 && len(value) != 0 {
		res = ChkPresent
	}
	return res
}

// SetFaces sets the face region tags.  It can only remove existing tags; it
// can't add new ones.
func SetFaces(h filefmt.FileHandler, v []string) error {
	var (
		list  []string
		faces = make(map[string]bool)
	)
	if xmp := h.XMP(false); xmp != nil {
		// First, the MPReg tag.
		for _, face := range v {
			faces[face] = true
		}
		for _, face := range xmp.MPRegPersonDisplayNames() {
			if _, ok := faces[face]; ok {
				list = append(list, face)
				faces[face] = false
			}
		}
		if err := xmp.SetMPRegPersonDisplayNames(list); err != nil {
			return err
		}
		// Now repeat the exact same thing for the mwg-rs tag.
		list = list[:0]
		for _, face := range v {
			faces[face] = true
		}
		for _, face := range xmp.MWGRSNames() {
			if _, ok := faces[face]; ok {
				list = append(list, face)
				faces[face] = false
			}
		}
		if err := xmp.SetMWGRSNames(list); err != nil {
			return err
		}
		// We're happy if the faces we were asked to set were seen in
		// either of the two tags; it doesn't have to be both.  But
		// if we didn't see a face on either one, we should flag that.
		for face, unseen := range faces {
			if unseen {
				return fmt.Errorf("cannot add face region for %q", face)
			}
		}
	} else if len(v) != 0 {
		return errors.New("cannot add face regions")
	}
	return nil
}
