package strmeta

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// GetCaption returns the highest priority caption value.
func GetCaption(h fileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.DCDescription) != 0 {
			return xmp.DCDescription[0].Value
		}
		if len(xmp.EXIFUserComments) != 0 {
			return xmp.EXIFUserComments[0]
		}
		if len(xmp.TIFFImageDescription) != 0 {
			return xmp.TIFFImageDescription[0].Value
		}
	}
	if exif := h.EXIF(); exif != nil {
		if exif.UserComment != "" {
			return exif.UserComment
		}
		if exif.ImageDescription != "" {
			return exif.ImageDescription
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if iptc.CaptionAbstract != "" {
			return iptc.CaptionAbstract
		}
	}
	return ""
}

// GetCaptionTags returns all of the caption tags and their values.
func GetCaptionTags(h fileHandler) (tags, values []string) {
	if xmp := h.XMP(false); xmp != nil {
		tags, values = tagsForLangStrings(tags, values, "XMP.dc:Description", xmp.DCDescription)
		for _, v := range xmp.EXIFUserComments {
			tags = append(tags, "XMP.exif:UserComment")
			values = append(values, v)
		}
		tags, values = tagsForLangStrings(tags, values, "XMP.tiff:ImageDescription", xmp.TIFFImageDescription)
	}
	if exif := h.EXIF(); exif != nil {
		if exif.UserComment != "" {
			tags = append(tags, "EXIF.UserComment")
			values = append(values, exif.UserComment)
		}
		tags = append(tags, "EXIF.ImageDescription")
		values = append(values, exif.ImageDescription)
	}
	if iptc := h.IPTC(); iptc != nil {
		tags = append(tags, "IPTC.CaptionAbstract")
		values = append(values, iptc.CaptionAbstract)
	}
	return tags, values
}

// SetCaption sets the caption tags.
func SetCaption(h fileHandler, v string) error {
	var list []metadata.LangString

	if v != "" {
		list = []metadata.LangString{{Lang: "", Value: v}}
	}
	if xmp := h.XMP(true); xmp != nil {
		xmp.DCDescription = list
		xmp.EXIFUserComments = nil // Always clear unwanted tag
		xmp.TIFFImageDescription = list
	}
	if exif := h.EXIF(); exif != nil {
		exif.UserComment = "" // Always clear unwanted tag
		exif.ImageDescription = v
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.CaptionAbstract = v
	}
	return nil
}
