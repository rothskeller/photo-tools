package xmp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/exif"
)

// EXIFDateTimeOriginal returns the EXIF DateTimeOriginal value from the XMP.
func (p *XMP) EXIFDateTimeOriginal() string {
	if p == nil || p.doc == nil {
		return ""
	}
	if model := exif.FindModel(p.doc); model != nil {
		if model.DateTimeOriginal != "" && !dateRE.MatchString(model.DateTimeOriginal) {
			p.log("EXIFDateTimeOriginal: invalid value: %q", model.DateTimeOriginal)
			return ""
		}
		return canonicalDate(model.DateTimeOriginal)
	}
	return ""
}

// SetEXIFDateTimeOriginal sets the EXIF DateTimeOriginal in the XMP.
func (p *XMP) SetEXIFDateTimeOriginal(dto string) {
	model, err := exif.MakeModel(p.doc)
	if err != nil {
		p.log("XMP exif.MakeModel: %s", err)
		return
	}
	model.DateTimeOriginal = dto
}

// EXIFDateTimeDigitized returns the EXIF DateTimeDigitized value from the XMP.
func (p *XMP) EXIFDateTimeDigitized() string {
	if p == nil || p.doc == nil {
		return ""
	}
	if model := exif.FindModel(p.doc); model != nil {
		if model.DateTimeDigitized != "" && !dateRE.MatchString(model.DateTimeDigitized) {
			p.log("EXIFDateTimeDigitized: invalid value: %q", model.DateTimeDigitized)
			return ""
		}
		return canonicalDate(model.DateTimeDigitized)
	}
	return ""
}

// SetEXIFDateTimeDigitized sets the EXIF DateTimeDigitized in the XMP.
func (p *XMP) SetEXIFDateTimeDigitized(dto string) {
	model, err := exif.MakeModel(p.doc)
	if err != nil {
		p.log("XMP exif.MakeModel: %s", err)
		return
	}
	model.DateTimeDigitized = dto
}

// EXIFGPSCoords returns the EXIF GPS coordinates from the XMP.
func (p *XMP) EXIFGPSCoords() (gc metadata.GPSCoords) {
	model, err := exif.MakeModel(p.doc)
	if err != nil {
		p.log("XMP exif.MakeModel: %s", err)
		return
	}
	if model.GPSLatitude == "" || model.GPSLongitude == "" {
		return metadata.GPSCoords{}
	}
	gc.Latitude = p.exifAngleToFixedFloat(model.GPSLatitude, 90)
	gc.Longitude = p.exifAngleToFixedFloat(model.GPSLongitude, 180)
	if model.GPSAltitude != "" {
		gc.Altitude = p.exifAltitudeToFixedFloat(model.GPSAltitudeRef, model.GPSAltitude)
	}
	// No need to error-check these, because errors return zero, and
	// metadata.GPSCoords treats that as invalid.
	return gc
}

// exifAngleToFixedFloat converts an angle, as specified in EXIF-to-XMP tags,
// into a metadata.FixedFloat.  It logs a problem and returns 0 if the value
// cannot be converted.
func (p *XMP) exifAngleToFixedFloat(exif string, max int) (f metadata.FixedFloat) {
	// According to the spec, GPSLatitude and GPSLongitude should have one
	// of the forms DDD,MM.mmX or DDD,MM,SSX where X is the NSEW direction,
	// and mm can have any number of fractional digits.  We will also accept
	// DDD.ddX and ±DDD.dd just in case.
	var (
		neg   bool
		parts []string
		err   error
	)
	if exif == "" {
		return 0
	}
	// Detect the sign and remove it from the string.
	if dir := exif[len(exif)-1]; strings.ContainsRune("NSEW", rune(dir)) {
		if dir == 'S' || dir == 'W' {
			neg = true
		}
		exif = exif[:len(exif)-1]
	} else if exif[0] == '-' || exif[0] == '+' {
		neg = exif[0] == '-'
		exif = exif[1:]
	}
	parts = strings.Split(exif, ",")
	switch len(parts) {
	case 1: // DDD.dd
		if f, err = metadata.ParseFixedFloat(exif); err != nil {
			goto ERROR
		}
	case 2: // DDD,MM.mm
		var degrees int
		if degrees, err = strconv.Atoi(parts[0]); err != nil {
			goto ERROR
		}
		if f, err = metadata.ParseFixedFloat(parts[1]); err != nil {
			goto ERROR
		}
		if f < 0 || f >= metadata.FixedFloatFromFraction(60, 1) {
			goto ERROR
		}
		f = f/60 + metadata.FixedFloatFromFraction(degrees, 1)
	case 3: // DDD,MM,SS
		var degrees, minutes, seconds int
		if degrees, err = strconv.Atoi(parts[0]); err != nil {
			goto ERROR
		}
		if minutes, err = strconv.Atoi(parts[1]); err != nil || minutes < 0 || minutes >= 60 {
			goto ERROR
		}
		if seconds, err = strconv.Atoi(parts[2]); err != nil || seconds < 0 || seconds >= 60 {
			goto ERROR
		}
		f = metadata.FixedFloatFromFraction(degrees, 1) +
			metadata.FixedFloatFromFraction(minutes, 60) +
			metadata.FixedFloatFromFraction(seconds, 3600)
	default:
		goto ERROR
	}
	if f > metadata.FixedFloatFromFraction(max, 1) {
		goto ERROR
	}
	// Apply any sign we detected earlier.
	if neg {
		f = -f
	}
	return f
ERROR:
	p.log("XMP EXIF: invalid angle value in GPS tag")
	return 0
}

// exifAltitudeToFixedFloat converts an altitude, as specified in EXIF-to-XMP
// tags, into a metadata.FixedFloat.  It logs a problem and returns 0 if the
// value cannot be converted.
func (p *XMP) exifAltitudeToFixedFloat(ref, alt string) (f metadata.FixedFloat) {
	// According to the spec, GPSAltitude is expressed as
	// numerator/denominator, with sign indicated by 0 (positive) or 1
	// (negative) in GPSAltitudeRef.  We will also accept numerator and
	// denominator separated by space, because go-xmp does and they probably
	// had a reason for that.  We will also accept a (possibly signed)
	// floating point in GPSAltitude.
	var (
		parts []string
		err   error
	)
	parts = strings.Split(alt, "/")
	if len(parts) == 1 {
		parts = strings.Split(alt, " ")
	}
	switch len(parts) {
	case 1: // float
		if f, err = metadata.ParseFixedFloat(alt); err != nil {
			goto ERROR
		}
		switch ref {
		case "", "0":
			break
		case "1":
			if f < 0 {
				goto ERROR
			}
			f = -f
		default:
			goto ERROR
		}
	case 2: // numerator and denominator
		var num, den int
		if num, err = strconv.Atoi(parts[0]); err != nil {
			goto ERROR
		}
		if den, err = strconv.Atoi(parts[1]); err != nil || den < 1 {
			goto ERROR
		}
		f = metadata.FixedFloatFromFraction(num, den)
		switch ref {
		case "0":
			break
		case "1":
			f = -f
		default:
			goto ERROR
		}
	default:
		goto ERROR
	}
	return f
ERROR:
	p.log("XMP EXIF: invalid altitude value in GPS tags")
	return 0

}

// SetEXIFGPSCoords sets the EXIF GPS coordinates in the XMP.
func (p *XMP) SetEXIFGPSCoords(gc metadata.GPSCoords) {
	model, err := exif.MakeModel(p.doc)
	if err != nil {
		p.log("XMP exif.MakeModel: %s", err)
		return
	}
	if !gc.Valid() {
		model.GPSAltitude = ""
		model.GPSAltitudeRef = ""
		model.GPSLatitude = ""
		model.GPSLongitude = ""
		return
	}
	model.GPSLatitude = fixedFloatToEXIFAngle(gc.Latitude, "N", "S")
	model.GPSLongitude = fixedFloatToEXIFAngle(gc.Longitude, "E", "W")
	if gc.HasAltitude() {
		if gc.Altitude < 0 {
			model.GPSAltitudeRef = "1"
			model.GPSAltitude = fixedFloatToEXIFRational(-gc.Altitude)
		} else {
			model.GPSAltitudeRef = "0"
			model.GPSAltitude = fixedFloatToEXIFRational(gc.Altitude)
		}
	} else {
		model.GPSAltitude = ""
		model.GPSAltitudeRef = ""
	}
}

func fixedFloatToEXIFAngle(f metadata.FixedFloat, pos, neg string) string {
	var (
		sb  strings.Builder
		suf string
	)
	// We always encode in the DDD,MM.mmX format.
	if f < 0 {
		suf = neg
		f = -f
	} else {
		suf = pos
	}
	fmt.Fprintf(&sb, "%d,", f.Int())
	f -= metadata.FixedFloatFromFraction(f.Int(), 1)
	f *= 60
	sb.WriteString(f.String())
	sb.WriteString(suf)
	return sb.String()
}

func fixedFloatToEXIFRational(f metadata.FixedFloat) string {
	num, den := f.AsFraction()
	for den%10 == 0 && num%10 == 0 {
		num /= 10
		den /= 10
	}
	return fmt.Sprintf("%d/%d", num, den)
}

// EXIFUserComment returns the EXIF UserComment from the XMP.
func (p *XMP) EXIFUserComment() []string {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := exif.FindModel(p.doc); model != nil {
		return model.UserComment
	}
	return nil
}

// SetEXIFUserComment sets the EXIF UserComment in the XMP.
func (p *XMP) SetEXIFUserComment(comment string) {
	model, err := exif.MakeModel(p.doc)
	if err != nil {
		p.log("XMP exif.MakeModel: %s", err)
		return
	}
	if comment != "" {
		model.UserComment = []string{comment}
	} else {
		model.UserComment = nil
	}
}
