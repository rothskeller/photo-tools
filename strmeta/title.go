package strmeta

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// GetTitle returns the highest priority title value.
func GetTitle(h fileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.DCTitle) != 0 {
			return xmp.DCTitle[0].Value
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
		tags, values = tagsForLangStrings(tags, values, "XMP.dc:Title", xmp.DCTitle)
	}
	if iptc := h.IPTC(); iptc != nil {
		tags = append(tags, "IPTC.ObjectName")
		values = append(values, iptc.ObjectName)
	}
	return tags, values
}

// SetTitle sets the title tags.
func SetTitle(h fileHandler, v string) error {
	var list []metadata.LangString

	if v != "" {
		list = []metadata.LangString{{Lang: "", Value: v}}
	}
	if xmp := h.XMP(true); xmp != nil {
		xmp.DCTitle = list
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.ObjectName = v
	}
	return nil
}
