package strmeta

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/iptc"
)

// GetCaption returns the highest priority caption value.
func GetCaption(h fileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if def := xmp.DCDescription.Default(); def != "" {
			return def
		}
		if len(xmp.EXIFUserComments) != 0 {
			return xmp.EXIFUserComments[0]
		}
		if def := xmp.TIFFImageDescription.Default(); def != "" {
			return def
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
		tags, values = tagsForAltString(tags, values, "XMP.dc:Description", xmp.DCDescription)
		for _, v := range xmp.EXIFUserComments {
			tags = append(tags, "XMP.exif:UserComment")
			values = append(values, v)
		}
		tags, values = tagsForAltString(tags, values, "XMP.tiff:ImageDescription", xmp.TIFFImageDescription)
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

// CheckCaption determines whether the caption is tagged correctly, and is
// consistent with the reference.
func CheckCaption(ref, h fileHandler) (res CheckResult) {
	var value = GetCaption(ref)
	if xmp := h.XMP(false); xmp != nil {
		switch len(xmp.DCDescription) {
		case 0:
			if value != "" {
				res = ChkIncorrectlyTagged
			}
		case 1:
			if xmp.DCDescription[0].Value != value {
				return ChkConflictingValues
			}
		default:
			return ChkConflictingValues
		}
		switch len(xmp.EXIFUserComments) {
		case 0:
			break
		case 1:
			if xmp.EXIFUserComments[0] != value {
				return ChkConflictingValues
			}
			res = ChkIncorrectlyTagged
		default:
			return ChkConflictingValues
		}
		switch len(xmp.TIFFImageDescription) {
		case 0:
			if value != "" {
				res = ChkIncorrectlyTagged
			}
		case 1:
			if xmp.TIFFImageDescription[0].Value != value {
				return ChkConflictingValues
			}
		default:
			return ChkConflictingValues
		}
	}
	if exif := h.EXIF(); exif != nil {
		if exif.UserComment != "" && exif.UserComment != value {
			return ChkConflictingValues
		} else if exif.UserComment != "" {
			res = ChkIncorrectlyTagged
		}
		if exif.ImageDescription != "" && exif.ImageDescription != value {
			return ChkConflictingValues
		} else if exif.ImageDescription == "" && value != "" {
			res = ChkIncorrectlyTagged
		}
	}
	if i := h.IPTC(); i != nil {
		if i.CaptionAbstract != "" && !stringEqualMax(value, i.CaptionAbstract, iptc.MaxCaptionAbstractLen) {
			return ChkConflictingValues
		} else if i.CaptionAbstract == "" && value != "" {
			res = ChkIncorrectlyTagged
		}
	}
	if value != "" && res == 0 {
		return ChkPresent
	}
	return res
}

// SetCaption sets the caption tags.
func SetCaption(h fileHandler, v string) error {
	var as metadata.AltString

	if v != "" {
		as = metadata.NewAltString(v)
	}
	if xmp := h.XMP(true); xmp != nil {
		xmp.DCDescription = as
		xmp.EXIFUserComments = nil // Always clear unwanted tag
		xmp.TIFFImageDescription = as
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
