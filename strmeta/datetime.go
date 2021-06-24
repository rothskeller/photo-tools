package strmeta

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// GetDateTime returns the highest priority date/time value.
func GetDateTime(h fileHandler) metadata.DateTime {
	if exif := h.EXIF(); exif != nil {
		if !exif.DateTimeOriginal.Empty() {
			return exif.DateTimeOriginal
		}
		if !exif.DateTimeDigitized.Empty() {
			return exif.DateTimeDigitized
		}
		if !exif.DateTime.Empty() {
			return exif.DateTime
		}
	}
	if xmp := h.XMP(false); xmp != nil {
		if !xmp.EXIFDateTimeOriginal.Empty() {
			return xmp.EXIFDateTimeOriginal
		}
		if !xmp.EXIFDateTimeDigitized.Empty() {
			return xmp.EXIFDateTimeDigitized
		}
		if !xmp.PSDateCreated.Empty() {
			return xmp.PSDateCreated
		}
		if !xmp.XMPCreateDate.Empty() {
			return xmp.XMPCreateDate
		}
		if !xmp.TIFFDateTime.Empty() {
			return xmp.TIFFDateTime
		}
		if !xmp.XMPModifyDate.Empty() {
			return xmp.XMPModifyDate
		}
		if !xmp.XMPMetadataDate.Empty() {
			return xmp.XMPMetadataDate
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if !iptc.DateTimeCreated.Empty() {
			return iptc.DateTimeCreated
		}
		if !iptc.DigitalCreationDateTime.Empty() {
			return iptc.DigitalCreationDateTime
		}
	}
	return metadata.DateTime{}
}

// GetDateTimeTags returns all of the date/time tags and their values.
func GetDateTimeTags(h fileHandler) (tags []string, values []metadata.DateTime) {
	if exif := h.EXIF(); exif != nil {
		tags = append(tags, "EXIF.DateTimeOriginal")
		values = append(values, exif.DateTimeOriginal)
		if !exif.DateTimeDigitized.Empty() {
			tags = append(tags, "EXIF.DateTimeDigitized")
			values = append(values, exif.DateTimeDigitized)
		}
		if !exif.DateTime.Empty() {
			tags = append(tags, "EXIF.DateTime")
			values = append(values, exif.DateTime)
		}
	}
	if xmp := h.XMP(false); xmp != nil {
		tags = append(tags, "XMP.exif:DateTimeOriginal")
		values = append(values, xmp.EXIFDateTimeOriginal)
		if !xmp.EXIFDateTimeDigitized.Empty() {
			tags = append(tags, "XMP.exif:DateTimeDigitized")
			values = append(values, xmp.EXIFDateTimeDigitized)
		}
		tags = append(tags, "XMP.ps:DateCreated")
		values = append(values, xmp.PSDateCreated)
		tags = append(tags, "XMP.xmp:CreateDate")
		values = append(values, xmp.XMPCreateDate)
		if !xmp.TIFFDateTime.Empty() {
			tags = append(tags, "XMP.tiff:DateTime")
			values = append(values, xmp.TIFFDateTime)
		}
		if !xmp.XMPModifyDate.Empty() {
			tags = append(tags, "XMP.xmp:ModifyDate")
			values = append(values, xmp.XMPModifyDate)
		}
		if !xmp.XMPMetadataDate.Empty() {
			tags = append(tags, "XMP.xmp:MetadataDate")
			values = append(values, xmp.XMPMetadataDate)
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		tags = append(tags, "IPTC.DateTimeCreated")
		values = append(values, iptc.DateTimeCreated)
		if !iptc.DigitalCreationDateTime.Empty() {
			tags = append(tags, "IPTC.DigitalCreationDateTime")
			values = append(values, iptc.DigitalCreationDateTime)
		}
	}
	return tags, values
}

// SetDateTime sets the date/time tags.
func SetDateTime(h fileHandler, v metadata.DateTime) error {
	if exif := h.EXIF(); exif != nil {
		exif.DateTimeOriginal = v
		exif.DateTimeDigitized = metadata.DateTime{} // Always clear unwanted tag
		exif.DateTime = metadata.DateTime{}          // Always clear unwanted tag
	}
	if xmp := h.XMP(true); xmp != nil {
		xmp.EXIFDateTimeOriginal = v
		xmp.PSDateCreated = v
		xmp.XMPCreateDate = v
		xmp.EXIFDateTimeDigitized = metadata.DateTime{} // Always clear unwanted tag
		xmp.TIFFDateTime = metadata.DateTime{}          // Always clear unwanted tag
		xmp.XMPModifyDate = metadata.DateTime{}         // Always clear unwanted tag
		xmp.XMPMetadataDate = metadata.DateTime{}       // Always clear unwanted tag
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.DateTimeCreated = v
		iptc.DigitalCreationDateTime = metadata.DateTime{} // Always clear unwanted tag
	}
	return nil
}
