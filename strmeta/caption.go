package strmeta

import (
	"github.com/rothskeller/photo-tools/filefmt"
)

func GetCaption(h filefmt.FileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if alternatives := xmp.DCDescription(); len(alternatives) != 0 {
			return alternatives[0][1]
		}
		if values := xmp.EXIFUserComment(); len(values) != 0 {
			return values[0]
		}
		if alternatives := xmp.TIFFImageDescription(); len(alternatives) != 0 {
			return alternatives[0][1]
		}
	}
	if exif := h.EXIF(); exif != nil {
		if value := exif.UserComment(); value != "" {
			return value
		}
		if value := exif.ImageDescription(); value != "" {
			return value
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if value := iptc.CaptionAbstract(); value != "" {
			return value
		}
	}
	return ""
}

func GetCaptions(h filefmt.FileHandler) (values []string, change bool) {
	var canonical string

	if xmp := h.XMP(false); xmp != nil {
		for _, alt := range xmp.DCDescription() {
			values = append(values, alt[1])
		}
		values = append(values, xmp.EXIFUserComment()...)
		for _, alt := range xmp.TIFFImageDescription() {
			values = append(values, alt[1])
		}
	}
	if exif := h.EXIF(); exif != nil {
		values = append(values, exif.ImageDescription())
		values = append(values, exif.UserComment())
	}
	if iptc := h.IPTC(); iptc != nil {
		values = append(values, iptc.CaptionAbstract())
	}
	if len(values) > 0 {
		canonical = values[0]
	}
	values = removeDuplicates(removeEmpty(values))
	if xmp := h.XMP(false); xmp != nil {
		if alts := xmp.DCDescription(); len(alts) == 0 && canonical != "" {
			change = true
		} else if len(alts) > 1 || (len(alts) == 1 && (alts[0][1] != canonical || alts[0][0] != "")) {
			change = true
		}
		if len(xmp.EXIFUserComment()) != 0 {
			change = true
		}
		if alts := xmp.TIFFImageDescription(); len(alts) == 0 && canonical != "" {
			change = true
		} else if len(alts) > 1 || (len(alts) == 1 && (alts[0][1] != canonical || alts[0][0] != "")) {
			change = true
		}
	}
	if exif := h.EXIF(); exif != nil {
		if exif.ImageDescription() != canonical {
			change = true
		}
		if exif.UserComment() != "" {
			change = true
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if !equalMaxLen(canonical, iptc.CaptionAbstract(), 2000) {
			change = true
		}
	}
	return values, change
}

func SetCaption(h filefmt.FileHandler, caption string) error {
	if xmp := h.XMP(caption != ""); xmp != nil {
		xmp.SetDCDescription(caption)
		xmp.SetEXIFUserComment("") // always remove
		xmp.SetTIFFImageDescription(caption)
	}
	if exif := h.EXIF(); exif != nil {
		exif.SetImageDescription(caption)
		exif.SetUserComment("") // always remove
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.SetCaptionAbstract(caption)
	}
	return nil
}
