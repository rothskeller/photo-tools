package metadata

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const feetToMeters FixedFloat = 304800 // by definition

// ErrParseGPSCoords is the error returned when a string cannot be parsed into
// GPS coordinates.
var ErrParseGPSCoords = errors.New("invalid GPSCoords value")

// GPSCoords holds a set of GPS coordinates.
type GPSCoords struct {
	// Latitude, in degrees north of the equator.
	latitude FixedFloat
	// Longitude, in degrees east of the zero meridian.
	longitude FixedFloat
	// Altitude, in meters above sea level (with 0 = unspecified)
	altitude FixedFloat
}

// Parse sets the value from the input string.  It returns an error if
// the input was invalid.
func (gc *GPSCoords) Parse(s string) (err error) {
	var feet bool

	parts := strings.Split(s, ",")
	if len(parts) == 1 || len(parts) > 3 {
		return ErrParseGPSCoords
	}
	if len(parts) == 0 {
		return nil
	}
	if gc.latitude, err = ParseFixedFloat(parts[0]); err != nil {
		return ErrParseGPSCoords
	}
	if gc.longitude, err = ParseFixedFloat(parts[1]); err != nil {
		return ErrParseGPSCoords
	}
	if len(parts) == 2 {
		return nil
	}
	if strings.HasSuffix(parts[2], "m") {
		parts[2] = parts[2][:len(parts[2])-1]
	} else if strings.HasSuffix(parts[2], "ft") {
		feet = true
		parts[2] = parts[2][:len(parts[2])-2]
	} else if strings.HasSuffix(parts[2], "'") {
		feet = true
		parts[2] = parts[2][:len(parts[2])-1]
	} else {
		return ErrParseGPSCoords
	}
	if gc.altitude, err = ParseFixedFloat(parts[2]); err != nil {
		return ErrParseGPSCoords
	}
	if feet {
		gc.altitude = gc.altitude.Mul(feetToMeters)
	}
	return nil
}

// String returns the value in string form, suitable for input to Parse.
func (gc *GPSCoords) String() string {
	var sb strings.Builder
	if gc.Empty() {
		return ""
	}
	sb.WriteString(gc.latitude.String())
	sb.WriteString(", ")
	sb.WriteString(gc.longitude.String())
	if !gc.HasAltitude() {
		return sb.String()
	}
	sb.WriteString(", ")
	feet := gc.altitude.Div(feetToMeters)
	sb.WriteString(feet.String())
	sb.WriteString("ft")
	return sb.String()
}

// ParseEXIF parses a set of GPS coordinates as represented in EXIF
// It return ErrParseGPSCoords if the input data are invalid.
func (gc *GPSCoords) ParseEXIF(latref string, lat []uint32, longref string, long []uint32, altref byte, alt []uint32) error {
	*gc = GPSCoords{}
	if latref == "" && len(lat) == 0 && longref == "" && len(long) == 0 && altref == 0 && len(alt) == 0 {
		return nil
	}
	if (latref != "N" && latref != "S") || len(lat) != 6 || (longref != "E" && longref != "W") || len(long) != 6 {
		return ErrParseGPSCoords
	}
	if lat[1] <= 0 || lat[3] <= 0 || lat[5] < 0 || long[1] <= 0 || long[3] <= 0 || long[5] < 0 {
		return ErrParseGPSCoords
	}
	if lat[5] == 0 {
		if lat[4] == 0 {
			lat[5] = 1 // illegal, but empirically happens; fix it up
		} else {
			return ErrParseGPSCoords
		}
	}
	if long[5] == 0 {
		if long[4] == 0 {
			long[5] = 1 // illegal, but empirically happens; fix it up
		} else {
			return ErrParseGPSCoords
		}
	}
	gc.latitude = FixedFloatFromFraction(int(lat[0]), int(lat[1]))
	gc.latitude += FixedFloatFromFraction(int(lat[2]), int(lat[3])*60)
	gc.latitude += FixedFloatFromFraction(int(lat[4]), int(lat[5])*3600)
	if latref == "S" {
		gc.latitude = -gc.latitude
	}
	gc.longitude = FixedFloatFromFraction(int(long[0]), int(long[1]))
	gc.longitude += FixedFloatFromFraction(int(long[2]), int(long[3])*60)
	gc.longitude += FixedFloatFromFraction(int(long[4]), int(long[5])*3600)
	if longref == "W" {
		gc.longitude = -gc.longitude
	}
	if altref == 0 && len(alt) == 0 {
		return nil
	}
	if (altref != 0 && altref != 1) || len(alt) != 2 {
		return ErrParseGPSCoords
	}
	gc.altitude = FixedFloatFromFraction(int(alt[0]), int(alt[1]))
	if altref == 1 {
		gc.altitude = -gc.altitude
	}
	return nil
}

// AsEXIF renders a set of GPS coordinates in EXIF metadata form.  Note that
// the ParseEXIF / AsEXIF round trip is not idempotent, because it transforms
// degrees,minutes,seconds into fractional degrees.
func (gc *GPSCoords) AsEXIF() (latref string, lat []uint32, longref string, long []uint32, altref byte, alt []uint32) {
	if gc.Empty() {
		return
	}
	if gc.latitude < 0 {
		latref = "S"
		lat = toEXIFDegrees(-gc.latitude)
	} else {
		latref = "N"
		lat = toEXIFDegrees(gc.latitude)
	}
	if gc.longitude < 0 {
		longref = "W"
		long = toEXIFDegrees(-gc.longitude)
	} else {
		longref = "E"
		long = toEXIFDegrees(gc.longitude)
	}
	if !gc.HasAltitude() {
		return
	}
	if gc.altitude < 0 {
		altref = 1
		alt = toEXIFRational(-gc.altitude)
	} else {
		alt = toEXIFRational(gc.altitude)
	}
	return
}
func toEXIFDegrees(f FixedFloat) (d []uint32) {
	d = make([]uint32, 6)
	copy(d, toEXIFRational(f))
	d[3] = 1
	d[5] = 1
	return d
}
func toEXIFRational(f FixedFloat) (r []uint32) {
	r = make([]uint32, 2)
	r[0] = uint32(f)
	r[1] = uint32(1000000)
	for r[0]%10 == 0 && r[1] > 1 {
		r[0] /= 10
		r[1] /= 10
	}
	return r
}

// ParseXMP parses a set of GPS coordinates as represented in XMP
// It return ErrParseGPSCoords if the input data are invalid.
func (gc *GPSCoords) ParseXMP(lat, long, altref, alt string) (err error) {
	*gc = GPSCoords{}
	if gc.latitude, err = fromXMPAngle(lat, 90); err != nil {
		return err
	}
	if gc.longitude, err = fromXMPAngle(long, 180); err != nil {
		return err
	}
	if alt == "" {
		return nil
	}
	if gc.altitude, err = fromXMPAltitude(altref, alt); err != nil {
		return err
	}
	return nil
}
func fromXMPAngle(xmp string, max int) (f FixedFloat, err error) {
	// According to the spec, angles should have one of the forms DDD,MM.mmX
	// or DDD,MM,SSX where X is the NSEW direction, and mm can have any
	// number of fractional digits.  We will also accept DDD.ddX and ±DDD.dd
	// just in case.
	var (
		neg   bool
		parts []string
	)
	if xmp == "" {
		return 0, nil
	}
	// Detect the sign and remove it from the string.
	if dir := xmp[len(xmp)-1]; strings.ContainsRune("NSEW", rune(dir)) {
		if dir == 'S' || dir == 'W' {
			neg = true
		}
		xmp = xmp[:len(xmp)-1]
	} else if xmp[0] == '-' || xmp[0] == '+' {
		neg = xmp[0] == '-'
		xmp = xmp[1:]
	}
	parts = strings.Split(xmp, ",")
	switch len(parts) {
	case 1: // DDD.dd
		if f, err = ParseFixedFloat(xmp); err != nil {
			return 0, ErrParseGPSCoords
		}
	case 2: // DDD,MM.mm
		var degrees int
		if degrees, err = strconv.Atoi(parts[0]); err != nil {
			return 0, ErrParseGPSCoords
		}
		if f, err = ParseFixedFloat(parts[1]); err != nil {
			return 0, ErrParseGPSCoords
		}
		if f < 0 || f >= FixedFloatFromFraction(60, 1) {
			return 0, ErrParseGPSCoords
		}
		f = f/60 + FixedFloatFromFraction(degrees, 1)
	case 3: // DDD,MM,SS
		var degrees, minutes, seconds int
		if degrees, err = strconv.Atoi(parts[0]); err != nil {
			return 0, ErrParseGPSCoords
		}
		if minutes, err = strconv.Atoi(parts[1]); err != nil || minutes < 0 || minutes >= 60 {
			return 0, ErrParseGPSCoords
		}
		if seconds, err = strconv.Atoi(parts[2]); err != nil || seconds < 0 || seconds >= 60 {
			return 0, ErrParseGPSCoords
		}
		f = FixedFloatFromFraction(degrees, 1) +
			FixedFloatFromFraction(minutes, 60) +
			FixedFloatFromFraction(seconds, 3600)
	default:
		return 0, ErrParseGPSCoords
	}
	if f > FixedFloatFromFraction(max, 1) {
		return 0, ErrParseGPSCoords
	}
	// Apply any sign we detected earlier.
	if neg {
		f = -f
	}
	return f, nil
}
func fromXMPAltitude(ref, alt string) (f FixedFloat, err error) {
	// According to the spec, altitude is expressed as
	// numerator/denominator, with sign indicated by 0 (positive) or 1
	// (negative) in ref.  We will also accept numerator and denominator
	// separated by space, because go-xmp does and they probably had a
	// reason for that.  We will also accept a (possibly signed) floating
	// point in altitude.
	parts := strings.Split(alt, "/")
	if len(parts) == 1 {
		parts = strings.Split(alt, " ")
	}
	switch len(parts) {
	case 1: // float
		if f, err = ParseFixedFloat(alt); err != nil {
			return 0, ErrParseGPSCoords
		}
		switch ref {
		case "", "0":
			break
		case "1":
			if f < 0 {
				return 0, ErrParseGPSCoords
			}
			f = -f
		default:
			return 0, ErrParseGPSCoords
		}
	case 2: // numerator and denominator
		var num, den int
		if num, err = strconv.Atoi(parts[0]); err != nil {
			return 0, ErrParseGPSCoords
		}
		if den, err = strconv.Atoi(parts[1]); err != nil || den < 1 {
			return 0, ErrParseGPSCoords
		}
		f = FixedFloatFromFraction(num, den)
		switch ref {
		case "", "0": // "" isn't legal, but it happens a lot in my library
			break
		case "1":
			f = -f
		default:
			return 0, ErrParseGPSCoords
		}
	default:
		return 0, ErrParseGPSCoords
	}
	return f, nil
}

// AsXMP renders a set of GPS coordinates in XMP metadata form.  Note that
// the ParseXMP / AsXMP round trip is not idempotent, because it transforms
// degrees,minutes,seconds into degrees and fractional minutes.
func (gc *GPSCoords) AsXMP() (lat, long, altref, alt string) {
	if gc.Empty() {
		return
	}
	lat = toXMPAngle(gc.latitude, "N", "S")
	long = toXMPAngle(gc.longitude, "E", "W")
	if gc.HasAltitude() {
		altref, alt = toXMPAltitude(gc.altitude)
	}
	return
}
func toXMPAngle(f FixedFloat, pos, neg string) string {
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
	f -= FixedFloatFromFraction(f.Int(), 1)
	f *= 60
	sb.WriteString(f.String())
	sb.WriteString(suf)
	return sb.String()
}
func toXMPAltitude(f FixedFloat) (ref, alt string) {
	if f < 0 {
		ref = "1"
		f = -f
	} else {
		ref = "0"
	}
	num, den := f.AsFraction()
	for den%10 == 0 && num%10 == 0 {
		num /= 10
		den /= 10
	}
	alt = fmt.Sprintf("%d/%d", num, den)
	return
}

// Empty returns true if the value contains no data.
func (gc *GPSCoords) Empty() bool {
	return gc == nil || gc.latitude == 0 || gc.longitude == 0
}

// HasAltitude returns whether the coordinates include an altitude.
func (gc *GPSCoords) HasAltitude() bool {
	return !gc.Empty() && gc.altitude != 0
}

// Equal returns whether two GPSCoords are the same.
func (gc *GPSCoords) Equal(other *GPSCoords) bool {
	if gc.Empty() != other.Empty() {
		return false
	}
	if gc.Empty() {
		return true
	}
	if gc.latitude != other.latitude || gc.longitude != other.longitude {
		return false
	}
	if gc.HasAltitude() != other.HasAltitude() {
		return false
	}
	if !gc.HasAltitude() {
		return true
	}
	return gc.altitude == other.altitude
}

// Equivalent returns true if the receiver is equal to the argument, to the
// precision of the least precise of the two.  If so, the second return value is
// the more precise of the two.
func (gc *GPSCoords) Equivalent(other *GPSCoords) (bool, *GPSCoords) {
	if gc.Empty() != other.Empty() {
		return false, nil
	}
	if gc.Empty() {
		return true, gc
	}
	// Anything within 0.000003 degrees (in other words, 0.01 second) is
	// considered equivalent.  This handles the inaccuracy of conversions
	// between degrees,minutes,seconds format and fractional degrees.
	if diff := gc.latitude - other.latitude; diff < -3 || diff > 3 {
		return false, nil
	}
	if diff := gc.longitude - other.longitude; diff < -3 || diff > 3 {
		return false, nil
	}
	if other.altitude == 0 {
		return true, gc
	}
	if gc.altitude == 0 {
		return true, other
	}
	// Anything within a tenth of a meter is considered equivalent.
	if diff := gc.altitude - other.altitude; diff < -10000 || diff > 10000 {
		return false, nil
	}
	return true, gc
}
