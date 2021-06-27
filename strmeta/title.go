package strmeta

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/iptc"
)

// GetTitle returns the highest priority title value.
func GetTitle(h fileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if def := xmp.DCTitle.Default(); def != "" {
			return def
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if iptc.ObjectName != "" {
			return iptc.ObjectName
		}
	}
	return ""
}

// GetTitleTags returns all of the title tags and their values.
func GetTitleTags(h fileHandler) (tags, values []string) {
	if xmp := h.XMP(false); xmp != nil {
		tags, values = tagsForAltString(tags, values, "XMP.dc:Title", xmp.DCTitle)
	}
	if iptc := h.IPTC(); iptc != nil {
		tags = append(tags, "IPTC.ObjectName")
		values = append(values, iptc.ObjectName)
	}
	return tags, values
}

// CheckTitle determines whether the title is tagged correctly, and is
// consistent with the reference.
func CheckTitle(ref, h fileHandler) (res CheckResult) {
	var value = GetTitle(ref)

	if xmp := h.XMP(false); xmp != nil {
		switch len(xmp.DCTitle) {
		case 0:
			if value != "" {
				res = ChkIncorrectlyTagged
			}
		case 1:
			if value != xmp.DCTitle[0].Value {
				return ChkConflictingValues
			}
		default:
			return ChkConflictingValues
		}
	}
	if i := h.IPTC(); i != nil {
		if !stringEqualMax(value, i.ObjectName, iptc.MaxObjectNameLen) {
			if i.ObjectName != "" {
				return ChkConflictingValues
			}
			res = ChkIncorrectlyTagged
		}
	}
	if value != "" && res == 0 {
		return ChkPresent
	}
	return res
}

// SetTitle sets the title tags.
func SetTitle(h fileHandler, v string) error {
	var as metadata.AltString

	if v != "" {
		as = metadata.NewAltString(v)
	}
	if xmp := h.XMP(true); xmp != nil {
		xmp.DCTitle = as
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.ObjectName = v
	}
	return nil
}
