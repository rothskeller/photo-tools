package strmeta

import (
	"github.com/rothskeller/photo-tools/filefmt"
)

func GetArtist(h filefmt.FileHandler) string {
	if xmp := h.XMP(false); xmp != nil {
		if values := xmp.DCCreator(); len(values) != 0 {
			return values[0]
		}
		if values := xmp.TIFFArtist(); len(values) != 0 {
			return values[0]
		}
	}
	if exif := h.EXIF(); exif != nil {
		if values := exif.Artist(); len(values) != 0 {
			return values[0]
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if values := iptc.Bylines(); len(values) != 0 {
			return values[0]
		}
	}
	return ""
}

func GetArtists(h filefmt.FileHandler) (values []string, change bool) {
	var canonical string

	if xmp := h.XMP(false); xmp != nil {
		values = xmp.DCCreator()
		values = append(values, xmp.TIFFArtist()...)
	}
	if exif := h.EXIF(); exif != nil {
		values = append(values, exif.Artist()...)
	}
	if iptc := h.IPTC(); iptc != nil {
		values = append(values, iptc.Bylines()...)
	}
	if len(values) > 0 {
		canonical = values[0]
	}
	values = removeDuplicates(removeEmpty(values))
	if xmp := h.XMP(false); xmp != nil {
		if !listMatchItem(xmp.DCCreator(), canonical, 0) {
			change = true
		}
		if len(xmp.TIFFArtist()) != 0 {
			change = true
		}
	}
	if exif := h.EXIF(); exif != nil {
		if !listMatchItem(exif.Artist(), canonical, 0) {
			change = true
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if !listMatchItem(iptc.Bylines(), canonical, 32) {
			change = true
		}
	}
	return values, change
}

func SetArtist(h filefmt.FileHandler, artist string) error {
	var artists []string
	if artist != "" {
		artists = append(artists, artist)
	}
	if xmp := h.XMP(len(artists) != 0); xmp != nil {
		xmp.SetDCCreator(artists)
		xmp.SetTIFFArtist("") // deprecated tag, remove
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.SetBylines([]string{artist})
	}
	if exif := h.EXIF(); exif != nil {
		exif.SetArtist([]string{artist})
	}
	return nil
}
