package strmeta

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
)

// GetGPS returns the canonical GPS coordinates, in string form.
func GetGPS(h filefmt.FileHandler) metadata.GPSCoords {
	if exif := h.EXIF(); exif != nil {
		if gc := exif.GPSCoords(); gc.Valid() {
			return gc
		}
	}
	if xmp := h.XMP(false); xmp != nil {
		if gc := xmp.EXIFGPSCoords(); gc.Valid() {
			return gc
		}
	}
	return metadata.GPSCoords{}
}

// GetGPSs returns all tagged GPS coordinate sets, in string form; the first one
// is the canonical one.  If the metadata would be changed by calling SetGPS
// with the canonical value, the returned change flag is true.
func GetGPSs(h filefmt.FileHandler) (values []metadata.GPSCoords, change bool) {
	canonical := GetGPS(h)
	if !canonical.Valid() {
		return nil, false
	}
	values = append(values, canonical)
	if exif := h.EXIF(); exif != nil {
		if gc := exif.GPSCoords(); !gc.Valid() {
			return values, true
		}
	}
	if xmp := h.XMP(false); xmp != nil {
		if gc := xmp.EXIFGPSCoords(); !canonical.Equal(gc) {
			if gc.Valid() {
				values = append(values, gc)
			}
			return values, true
		}
	}
	return values, false
}

// SetGPS sets the GPS coordinates.  vals must be a slice of 0, 2, or 3 floats.
// 0 removes the GPS coordinates.  2 vals are latitude and longitude.  A third
// val is altitude in meters.
func SetGPS(h filefmt.FileHandler, gc metadata.GPSCoords) error {
	if exif := h.EXIF(); exif != nil {
		exif.SetGPSCoords(gc)
	}
	if xmp := h.XMP(gc.Valid()); xmp != nil {
		xmp.SetEXIFGPSCoords(gc)
	}
	return nil
}
