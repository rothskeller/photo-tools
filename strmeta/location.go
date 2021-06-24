package strmeta

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
)

// GetLocation returns the highest priority location captured values.  If it
// returns two values, they are the same location in two languages, the first
// one English and the second one the native language of the location.
func GetLocation(h fileHandler) []metadata.Location {
	if xmp := h.XMP(false); xmp != nil {
		var en *metadata.Location
		for i := range xmp.IPTCLocationCreated {
			if xmp.IPTCLocationCreated[i].Lang == "en" {
				en = &xmp.IPTCLocationCreated[i]
				break
			}
		}
		if en != nil {
			var other *metadata.Location
			for i := range xmp.IPTCLocationCreated {
				if &xmp.IPTCLocationCreated[i] != en {
					other = &xmp.IPTCLocationCreated[i]
				}
			}
			if other != nil {
				return []metadata.Location{*en, *other}
			}
			return []metadata.Location{*en}
		}
		if len(xmp.IPTCLocationCreated) >= 2 && xmp.IPTCLocationCreated[0].Lang == "" {
			return xmp.IPTCLocationCreated[:2]
		}
		if len(xmp.IPTCLocationCreated) >= 1 {
			return xmp.IPTCLocationCreated[:1]
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if !iptc.Location.Empty() {
			return []metadata.Location{iptc.Location}
		}
	}
	return nil
}

// GetLocationTags returns all of the location captured tags and their values.
func GetLocationTags(h fileHandler) (tags []string, values []metadata.Location) {
	if xmp := h.XMP(false); xmp != nil {
		switch len(xmp.IPTCLocationCreated) {
		case 0:
			tags = append(tags, "XMP.iptc:LocationCreated")
			values = append(values, metadata.Location{})
		case 1:
			if xmp.IPTCLocationCreated[0].Lang == "" {
				tags = append(tags, "XMP.iptc:LocationCreated")
			} else {
				tags = append(tags, fmt.Sprintf("XMP.iptc:LocationCreated[%s]", xmp.IPTCLocationCreated[0].Lang))
			}
			values = append(values, xmp.IPTCLocationCreated[0])
		default:
			for _, loc := range xmp.IPTCLocationCreated {
				tags = append(tags, fmt.Sprintf("XMP.iptc:LocationCreated[%s]", loc.Lang))
				values = append(values, loc)
			}
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		tags = append(tags, "IPTC.Location")
		values = append(values, iptc.Location)
	}
	return tags, values
}

// SetLocation sets the location captured tags.  There can be at most two
// values.  If there are two values, they are the same location in two
// languages, the first one English and the second one the native language of
// the location.
func SetLocation(h fileHandler, v []metadata.Location) error {
	if xmp := h.XMP(true); xmp != nil {
		xmp.IPTCLocationCreated = v
	}
	if iptc := h.IPTC(); iptc != nil {
		if len(v) != 0 {
			iptc.Location = v[0]
		} else {
			iptc.Location = metadata.Location{}
		}
	}
	return nil
}
