package strmeta

import (
	"errors"
	"fmt"
	"sort"

	"github.com/rothskeller/photo-tools/filefmt"
)

// GetFaces returns the names of face regions.  Since we can't add regions, we
// don't use the highest priority tag; we use the union of both tags.
func GetFaces(h fileHandler) (faces []string) {
	if xmp := h.XMP(false); xmp != nil {
		var fmap = make(map[string]bool)
		for _, face := range xmp.MPRegPersonDisplayNames() {
			fmap[face] = true
		}
		for _, face := range xmp.MWGRSNames() {
			fmap[face] = true
		}
		faces = make([]string, 0, len(fmap))
		for face := range fmap {
			faces = append(faces, face)
		}
		sort.Strings(faces)
		return faces
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
	var local = GetFaces(h)
	if !stringSliceEqual(value, local) {
		if len(local) != 0 {
			return ChkConflictingValues
		}
		return ChkIncorrectlyTagged
	}
	if len(value) != 0 {
		return ChkPresent
	}
	return ChkOptionalAbsent
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
