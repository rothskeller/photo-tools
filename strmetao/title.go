package strmeta

import (
	"github.com/rothskeller/photo-tools/filefmt"
)

func GetTitle(h filefmt.FileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if alternatives := xmp.DCTitle(); len(alternatives) != 0 {
			return alternatives[0][1]
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if value := iptc.ObjectName(); value != "" {
			return value
		}
	}
	return ""
}

func GetTitles(h filefmt.FileHandler) (values []string, change bool) {
	var canonical string

	if xmp := h.XMP(false); xmp != nil {
		for _, alt := range xmp.DCTitle() {
			values = append(values, alt[1])
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		values = append(values, iptc.ObjectName())
	}
	if len(values) > 0 {
		canonical = values[0]
	}
	values = removeDuplicates(removeEmpty(values))
	if xmp := h.XMP(false); xmp != nil {
		if alts := xmp.DCTitle(); len(alts) == 0 && canonical != "" {
			change = true
		} else if len(alts) > 1 || (len(alts) == 1 && (alts[0][1] != canonical || alts[0][0] != "")) {
			change = true
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if !equalMaxLen(canonical, iptc.ObjectName(), 64) {
			change = true
		}
	}
	return values, change
}

func SetTitle(h filefmt.FileHandler, title string) error {
	if xmp := h.XMP(title != ""); xmp != nil {
		xmp.SetDCTitle(title)
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.SetObjectName(title)
	}
	return nil
}
