package strmeta

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/iptc"
)

// GetLocationCaptured returns the canonical value for location captured.
func GetLocationCaptured(h filefmt.FileHandler) *metadata.Location {
	if xmp := h.XMP(false); xmp != nil {
		if locs := xmp.IPTCLocationsCreated(); len(locs) != 0 {
			return locs[0]
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if loc := iptc.Location(); loc.Valid() {
			return loc
		}
	}
	return nil
}

// GetLocationsCaptured returns all values for location captured, and a flag
// indicating whether setting the canonical value would result in a tag change.
func GetLocationsCaptured(h filefmt.FileHandler) (values []*metadata.Location, change bool) {
	canonical := GetLocationCaptured(h)
	if canonical == nil {
		return nil, false
	}
	if xmp := h.XMP(false); xmp != nil {
		values = append(values, xmp.IPTCLocationsCreated()...)
		if len(values) != 1 {
			change = true
		}
	}
	if iptch := h.IPTC(); iptch != nil {
		loc := iptch.Location()
		if !loc.Valid() {
			return values, true
		}
		if canonical.CountryCode != loc.CountryCode ||
			!equalMaxLen(canonical.CountryName, loc.CountryName, iptc.MaxCountryPLNameLen) ||
			!equalMaxLen(canonical.State, loc.State, iptc.MaxProvinceStateLen) ||
			!equalMaxLen(canonical.City, loc.City, iptc.MaxCityLen) ||
			!equalMaxLen(canonical.Sublocation, loc.Sublocation, iptc.MaxSublocationLen) {
			values = append(values, loc)
			return values, true
		}
	}
	return values, change
}

// SetLocationCaptured sets the location captured.
func SetLocationCaptured(h filefmt.FileHandler, loc *metadata.Location) error {
	if xmp := h.XMP(loc.Valid()); xmp != nil {
		xmp.SetIPTCLocationCreated(loc)
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.SetLocation(loc)
	}
	return nil
}

// GetLocationShown returns the canonical value for location shown.
func GetLocationShown(h filefmt.FileHandler) *metadata.Location {
	if xmp := h.XMP(false); xmp != nil {
		if locs := xmp.IPTCLocationsShown(); len(locs) != 0 {
			return locs[0]
		}
	}
	return nil
}

// GetLocationsShown returns all values for location shown, and a flag
// indicating whether setting the canonical value would result in a tag change.
func GetLocationsShown(h filefmt.FileHandler) (values []*metadata.Location, change bool) {
	if xmp := h.XMP(false); xmp != nil {
		values = xmp.IPTCLocationsShown()
		return values, len(values) > 1
	}
	return nil, false
}

// SetLocationShown sets the location shown.
func SetLocationShown(h filefmt.FileHandler, loc *metadata.Location) error {
	if xmp := h.XMP(loc.Valid()); xmp != nil {
		xmp.SetIPTCLocationShown(loc)
	}
	return nil
}
