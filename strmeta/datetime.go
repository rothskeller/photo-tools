package strmeta

import (
	"regexp"

	"github.com/rothskeller/photo-tools/filefmt"
)

func GetDateTime(h filefmt.FileHandler) string {
	// digiKam's algorithm for reading timesstamps is complex, involving a
	// weighted priority system looking for how many tags contain a
	// particular timestamp.  I'm not going to try to duplicate it here.
	// The key point is that it gives priority to EXIF first, then XMP, and
	// lastly IPTC.  And within the XMP, it's the EXIF-mirroring tags first,
	// then the XMP-native tags, and finally the IPTC-mirroring tags.
	if exif := h.EXIF(); exif != nil {
		if dto := exif.DateTimeOriginal(); dto != "" {
			return dto
		}
		if dtd := exif.DateTimeDigitized(); dtd != "" {
			return dtd
		}
		if dt := exif.DateTime(); dt != "" {
			return dt
		}
	}
	if xmp := h.XMP(false); xmp != nil {
		if dto := xmp.EXIFDateTimeOriginal(); dto != "" {
			return dto
		}
		if dtd := xmp.EXIFDateTimeDigitized(); dtd != "" {
			return dtd
		}
		if dc := xmp.PSDateCreated(); dc != "" {
			return dc
		}
		if cd := xmp.CreateDate(); cd != "" {
			return cd
		}
		if dt := xmp.TIFFDateTime(); dt != "" {
			return dt
		}
		if md := xmp.ModifyDate(); md != "" {
			return md
		}
		if md := xmp.MetadataDate(); md != "" {
			return md
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if dtc := iptc.DateTimeCreated(); dtc != "" {
			return dtc
		}
		if ddt := iptc.DigitalCreationDateTime(); ddt != "" {
			return ddt
		}
	}
	return ""
}

func GetDatesTimes(h filefmt.FileHandler) (values []string, change bool) {
	var canonical = GetDateTime(h)
	if canonical == "" {
		return nil, false
	}
	values = []string{canonical}
	if exif := h.EXIF(); exif != nil {
		values = addDTValue(values, &change, exif.DateTimeOriginal(), true, true)
		values = addDTValue(values, &change, exif.DateTimeDigitized(), true, false)
		values = addDTValue(values, &change, exif.DateTime(), true, false)
	}
	if xmp := h.XMP(false); xmp != nil {
		values = addDTValue(values, &change, xmp.EXIFDateTimeOriginal(), true, true)
		values = addDTValue(values, &change, xmp.CreateDate(), true, true)
		values = addDTValue(values, &change, xmp.EXIFDateTimeDigitized(), true, false)
		values = addDTValue(values, &change, xmp.PSDateCreated(), true, false)
		values = addDTValue(values, &change, xmp.TIFFDateTime(), true, false)
		values = addDTValue(values, &change, xmp.ModifyDate(), true, false)
		values = addDTValue(values, &change, xmp.MetadataDate(), true, false)
	}
	if iptc := h.IPTC(); iptc != nil {
		values = addDTValue(values, &change, iptc.DateTimeCreated(), false, true)
		values = addDTValue(values, &change, iptc.DigitalCreationDateTime(), false, false)
	}
	return values, change
}

var dateTimePartsRE = regexp.MustCompile(`^(\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d)(\.\d+)?([+-]\d\d:\d\d|Z)?$`)

func addDTValue(values []string, change *bool, add string, subsec, wanted bool) []string {
	println("addDTValue", add, *change)
	defer func() {
		println("=>", len(values), *change)
	}()
	if add == "" {
		if wanted {
			*change = true
		}
		return values
	}
	if !wanted {
		*change = true
	}
	ap := dateTimePartsRE.FindStringSubmatch(add)
	for idx, v := range values {
		var mod bool

		vp := dateTimePartsRE.FindStringSubmatch(v)
		if vp == nil {
			println(v)
		}
		if vp[1] != ap[1] {
			continue
		}
		// If one has subseconds and the other doesn't, apply them.
		if vp[2] != "" && ap[2] == "" {
			ap[2] = vp[2]
			mod = true
		} else if vp[2] == "" && ap[2] != "" {
			vp[2] = ap[2]
			mod = true
		}
		// If one has a time zone and the other doesn't, apply it.
		if vp[3] != "" && ap[3] == "" {
			ap[3] = vp[3]
			mod = true
		} else if vp[3] == "" && ap[3] != "" {
			vp[3] = ap[3]
			mod = true
		}
		if vp[2] == ap[2] && vp[3] == ap[3] {
			v = vp[1] + ap[2] + vp[3]
			values[idx] = v
			if mod {
				*change = true
			}
			return values
		}
	}
	values = append(values, add)
	*change = true
	return values
}

func SetDateTime(h filefmt.FileHandler, dt string) error {
	if exif := h.EXIF(); exif != nil {
		exif.SetDateTimeOriginal(dt)
		exif.SetDateTime("")
		exif.SetDateTimeDigitized("")
	}
	if xmp := h.XMP(false); xmp != nil {
		xmp.SetEXIFDateTimeOriginal(dt)
		xmp.SetCreateDate(dt)
		xmp.SetEXIFDateTimeDigitized("")
		xmp.SetPSDateCreated("")
		xmp.SetTIFFDateTime("")
		xmp.SetModifyDate("")
		xmp.SetMetadataDate("")
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.SetDateTimeCreated(dt)
		iptc.SetDigitalCreationDateTime("")
	}
	return nil
}
