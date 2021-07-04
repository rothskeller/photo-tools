package strmeta

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/iptc"
)

// GetCaption returns the highest priority caption value.
func GetCaption(h fileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if def := xmp.DCDescription().Default(); def != "" {
			return def
		}
		if def := xmp.EXIFUserComment().Default(); def != "" {
			return def
		}
		if def := xmp.TIFFImageDescription().Default(); def != "" {
			return def
		}
	}
	if exif := h.EXIF(); exif != nil {
		if exif.UserComment() != "" {
			return exif.UserComment()
		}
		if exif.ImageDescription() != "" {
			return exif.ImageDescription()
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if iptc.CaptionAbstract() != "" {
			return iptc.CaptionAbstract()
		}
	}
	return ""
}

// GetCaptionTags returns all of the caption tags and their values.
func GetCaptionTags(h fileHandler) (tags, values []string) {
	if xmp := h.XMP(false); xmp != nil {
		tags, values = tagsForAltString(tags, values, "XMP  dc:Description", xmp.DCDescription())
		tags, values = tagsForAltString(tags, values, "XMP  exif:UserComment", xmp.EXIFUserComment())
		tags, values = tagsForAltString(tags, values, "XMP  tiff:ImageDescription", xmp.TIFFImageDescription())
	}
	if exif := h.EXIF(); exif != nil {
		if exif.UserComment() != "" {
			tags = append(tags, "EXIF UserComment")
			values = append(values, exif.UserComment())
		}
		tags = append(tags, "EXIF ImageDescription")
		values = append(values, exif.ImageDescription())
	}
	if iptc := h.IPTC(); iptc != nil {
		tags = append(tags, "IPTC CaptionAbstract")
		values = append(values, iptc.CaptionAbstract())
	}
	return tags, values
}

// CheckCaption determines whether the caption is tagged correctly.
func CheckCaption(h fileHandler) (res CheckResult) {
	var value = GetCaption(h)
	if xmp := h.XMP(false); xmp != nil {
		switch len(xmp.DCDescription()) {
		case 0:
			if value != "" {
				res = ChkIncorrectlyTagged
			}
		case 1:
			if xmp.DCDescription()[0].Value != value {
				return ChkConflictingValues
			}
		default:
			return ChkConflictingValues
		}
		switch len(xmp.EXIFUserComment()) {
		case 0:
			break
		case 1:
			if xmp.EXIFUserComment()[0].Value != value {
				return ChkConflictingValues
			}
			res = ChkIncorrectlyTagged
		default:
			return ChkConflictingValues
		}
		switch len(xmp.TIFFImageDescription()) {
		case 0:
			if value != "" {
				res = ChkIncorrectlyTagged
			}
		case 1:
			if xmp.TIFFImageDescription()[0].Value != value {
				return ChkConflictingValues
			}
		default:
			return ChkConflictingValues
		}
	}
	if exif := h.EXIF(); exif != nil {
		if exif.UserComment() != "" && exif.UserComment() != value {
			return ChkConflictingValues
		} else if exif.UserComment() != "" {
			res = ChkIncorrectlyTagged
		}
		if exif.ImageDescription() != "" && exif.ImageDescription() != value {
			return ChkConflictingValues
		} else if exif.ImageDescription() == "" && value != "" {
			res = ChkIncorrectlyTagged
		}
	}
	if i := h.IPTC(); i != nil {
		if i.CaptionAbstract() != "" && !stringEqualMax(value, i.CaptionAbstract(), iptc.MaxCaptionAbstractLen) {
			return ChkConflictingValues
		} else if i.CaptionAbstract() == "" && value != "" {
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
		if err := xmp.SetDCDescription(as); err != nil {
			return err
		}
		if err := xmp.SetEXIFUserComments(nil); err != nil { // Always clear unwanted tag
			return err
		}
		if err := xmp.SetTIFFImageDescription(as); err != nil {
			return err
		}
	}
	if exif := h.EXIF(); exif != nil {
		if err := exif.SetUserComment(""); err != nil { // Always clear unwanted tag
			return err
		}
		if err := exif.SetImageDescription(v); err != nil {
			return err
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if err := iptc.SetCaptionAbstract(v); err != nil {
			return err
		}
	}
	return nil
}
