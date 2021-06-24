package strmeta

// GetCreator returns the highest priority creator value.
func GetCreator(h fileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.DCCreator) != 0 {
			return xmp.DCCreator[0]
		}
		if xmp.TIFFArtist != "" {
			return xmp.TIFFArtist
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if len(iptc.Bylines) != 0 {
			return iptc.Bylines[0]
		}
	}
	if exif := h.EXIF(); exif != nil {
		if len(exif.Artist) != 0 {
			return exif.Artist[0]
		}
	}
	return ""
}

// GetCreatorTags returns all of the creator tags and their values.
func GetCreatorTags(h fileHandler) (tags, values []string) {
	if xmp := h.XMP(false); xmp != nil {
		tags, values = tagsForStringList(tags, values, "XMP.dc:Creator", xmp.DCCreator)
		if xmp.TIFFArtist != "" {
			tags = append(tags, "XMP.tiff.Artist")
			values = append(values, xmp.TIFFArtist)
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		tags, values = tagsForStringList(tags, values, "IPTC.Byline", iptc.Bylines)
	}
	if exif := h.EXIF(); exif != nil {
		tags, values = tagsForStringList(tags, values, "EXIF.Artist", exif.Artist)
	}
	return tags, values
}

// SetCreator sets the creator tags.
func SetCreator(h fileHandler, v string) error {
	var list []string

	if v != "" {
		list = []string{v}
	}
	if xmp := h.XMP(true); xmp != nil {
		xmp.DCCreator = list
		xmp.TIFFArtist = "" // Always clear deprecated tag
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.Bylines = list
	}
	if exif := h.EXIF(); exif != nil {
		exif.Artist = list
	}
	return nil
}
