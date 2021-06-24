package strmeta

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
)

// GetShown returns the highest priority location shown values.  If it returns
// two values, they are the same location in two languages, the first one
// English and the second one the native language of the location.
func GetShown(h fileHandler) []metadata.Location {
	xmp := h.XMP(false)
	if xmp == nil || len(xmp.IPTCLocationsShown) == 0 {
		return nil
	}
	var loc = xmp.IPTCLocationsShown[0]
	var en *metadata.Location
	for i := range loc {
		if loc[i].Lang == "en" {
			en = &loc[i]
			break
		}
	}
	if en != nil {
		var other *metadata.Location
		for i := range loc {
			if &loc[i] != en {
				other = &loc[i]
			}
		}
		if other != nil {
			return []metadata.Location{*en, *other}
		}
		return []metadata.Location{*en}
	}
	if len(loc) >= 2 && loc[0].Lang == "" {
		return loc[:2]
	}
	if len(loc) >= 1 {
		return loc[:1]
	}
	return nil
}

// GetShownTags returns all of the location shown tags and their values.
func GetShownTags(h fileHandler) (tags []string, values []metadata.Location) {
	if xmp := h.XMP(false); xmp != nil {
		for i, locs := range xmp.IPTCLocationsShown {
			label := fmt.Sprintf("XMP.iptc:LocationShown[%d]", i)
			switch len(locs) {
			case 0:
				tags = append(tags, label)
				values = append(values, metadata.Location{})
			case 1:
				if locs[0].Lang == "" {
					tags = append(tags, label)
				} else {
					tags = append(tags, fmt.Sprintf("%s[%s]", label, locs[0].Lang))
				}
				values = append(values, locs[0])
			default:
				for _, loc := range locs {
					tags = append(tags, fmt.Sprintf("XMP.iptc:LocationCreated[%s]", loc.Lang))
					values = append(values, loc)
				}
			}
		}
		if len(xmp.IPTCLocationsShown) == 0 {
			tags = append(tags, "XMP.iptc:LocationShown")
			values = append(values, metadata.Location{})
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		tags = append(tags, "IPTC.Location")
		values = append(values, iptc.Location)
	}
	return tags, values
}

// SetShown sets the location shown tags.  There can be at most two values.  If
// there are two values, they are the same location in two languages, the first
// one English and the second one the native language of the location.
func SetShown(h fileHandler, v []metadata.Location) error {
	if xmp := h.XMP(true); xmp != nil {
		xmp.IPTCLocationsShown = [][]metadata.Location{v}
	}
	return nil
}
