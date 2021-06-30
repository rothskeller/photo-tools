package strmeta

// GetCreator returns the highest priority creator value.
func GetCreator(h fileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.DCCreator()) != 0 {
			return xmp.DCCreator()[0]
		}
		if xmp.TIFFArtist() != "" {
			return xmp.TIFFArtist()
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if len(iptc.Bylines()) != 0 {
			return iptc.Bylines()[0]
		}
	}
	if exif := h.EXIF(); exif != nil {
		if len(exif.Artist()) != 0 {
			return exif.Artist()[0]
		}
	}
	return ""
}

// GetCreatorTags returns all of the creator tags and their values.
func GetCreatorTags(h fileHandler) (tags, values []string) {
	if xmp := h.XMP(false); xmp != nil {
		tags, values = tagsForStringList(tags, values, "XMP  dc:Creator", xmp.DCCreator())
		if xmp.TIFFArtist() != "" {
			tags = append(tags, "XMP.tiff.Artist")
			values = append(values, xmp.TIFFArtist())
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		tags, values = tagsForStringList(tags, values, "IPTC Byline", iptc.Bylines())
	}
	if exif := h.EXIF(); exif != nil {
		tags, values = tagsForStringList(tags, values, "EXIF Artist", exif.Artist())
	}
	return tags, values
}

// CheckCreator determines whether the creator is tagged correctly, and is
// consistent with the reference.
func CheckCreator(ref, h fileHandler) (res CheckResult) {
	var value = GetCreator(ref)

	if xmp := h.XMP(false); xmp != nil {
		switch len(xmp.DCCreator()) {
		case 0:
			if value != "" {
				res = ChkIncorrectlyTagged
			}
		case 1:
			if xmp.DCCreator()[0] != value {
				return ChkConflictingValues
			}
		default:
			return ChkConflictingValues
		}
		if xmp.TIFFArtist() != "" && xmp.TIFFArtist() != value {
			return ChkConflictingValues
		} else if xmp.TIFFArtist() != "" {
			res = ChkIncorrectlyTagged
		}
	}
	if i := h.IPTC(); i != nil {
		switch len(i.Bylines()) {
		case 0:
			if value != "" {
				res = ChkIncorrectlyTagged
			}
		case 1:
			if i.Bylines()[0] != value {
				return ChkConflictingValues
			}
		default:
			return ChkConflictingValues
		}
	}
	if exif := h.EXIF(); exif != nil {
		switch len(exif.Artist()) {
		case 0:
			if value != "" {
				res = ChkIncorrectlyTagged
			}
		case 1:
			if exif.Artist()[0] != value {
				return ChkConflictingValues
			}
		default:
			return ChkConflictingValues
		}
	}
	if value != "" && res == 0 {
		return ChkPresent
	}
	if value == "" && res == 0 {
		return ChkExpectedAbsent
	}
	return res
}

// SetCreator sets the creator tags.
func SetCreator(h fileHandler, v string) error {
	var list []string

	if v != "" {
		list = []string{v}
	}
	if xmp := h.XMP(true); xmp != nil {
		if err := xmp.SetDCCreator(list); err != nil {
			return err
		}
		if err := xmp.SetTIFFArtist(""); err != nil { // Always clear deprecated tag
			return err
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if err := iptc.SetBylines(list); err != nil {
			return err
		}
	}
	if exif := h.EXIF(); exif != nil {
		if err := exif.SetArtist(list); err != nil {
			return err
		}
	}
	return nil
}
