package strmeta

import "github.com/rothskeller/photo-tools/metadata"

// GetGPSCoords returns the highest priority GPS coordinates value.
func GetGPSCoords(h fileHandler) metadata.GPSCoords {
	if xmp := h.XMP(false); xmp != nil {
		if !xmp.EXIFGPSCoords.Empty() {
			return xmp.EXIFGPSCoords
		}
	}
	if exif := h.EXIF(); exif != nil {
		if !exif.GPSCoords.Empty() {
			return exif.GPSCoords
		}
	}
	return metadata.GPSCoords{}
}

// GetGPSCoordsTags returns all of the GPS coordinates tags and their values.
func GetGPSCoordsTags(h fileHandler) (tags []string, values []metadata.GPSCoords) {
	if xmp := h.XMP(false); xmp != nil {
		tags = append(tags, "XMP  exif:GPSCoords")
		values = append(values, xmp.EXIFGPSCoords)
	}
	if exif := h.EXIF(); exif != nil {
		tags = append(tags, "EXIF GPSCoords")
		values = append(values, exif.GPSCoords)
	}
	return tags, values
}

// CheckGPSCoords determines whether the GPS coordinates are tagged correctly,
// and are consistent with the reference.
func CheckGPSCoords(ref, h fileHandler) (res CheckResult) {
	var value = GetGPSCoords(ref)

	if xmp := h.XMP(false); xmp != nil {
		if !xmp.EXIFGPSCoords.Empty() {
			if !value.Equivalent(&xmp.EXIFGPSCoords) {
				return ChkConflictingValues
			}
		} else if !value.Empty() {
			res = ChkIncorrectlyTagged
		}
	}
	if exif := h.EXIF(); exif != nil {
		if !exif.GPSCoords.Empty() {
			if !value.Equivalent(&exif.GPSCoords) {
				return ChkConflictingValues
			}
		} else if !value.Empty() {
			res = ChkIncorrectlyTagged
		}
	}
	if !value.Empty() && res == 0 {
		return ChkPresent
	}
	if value.Empty() && res == 0 {
		return ChkExpectedAbsent
	}
	return res
}

// SetGPSCoords sets the GPS coordinates tags.
func SetGPSCoords(h fileHandler, v metadata.GPSCoords) error {
	if xmp := h.XMP(true); xmp != nil {
		xmp.EXIFGPSCoords = v
	}
	if exif := h.EXIF(); exif != nil {
		exif.GPSCoords = v
	}
	return nil
}
