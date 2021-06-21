package metadata

import (
	"errors"
	"fmt"
	"strings"
)

// FixedFloat stores a floating point number as an signed integer scaled 100,000
// higher, so that it always has exactly six digits of precision after the
// decimal point.  (That's enough to get us down to hundredths of a second when
// storing an angle in degrees, which is close enough for GPS coordinates.)
type FixedFloat int64

// ParseFixedFloat parses a string as a fixed floating point number.
func ParseFixedFloat(s string) (f FixedFloat, err error) {
	var neg bool
	var seenDecimal = false
	var digitsAfterDecimal = 0

	s = strings.TrimSpace(s)
	if s != "" && s[0] == '-' {
		neg = true
		s = s[1:]
	}
	if s == "" {
		goto ERROR
	}
	for _, c := range s {
		if c >= '0' && c <= '9' {
			if digitsAfterDecimal >= 6 {
				continue
			}
			f = f*10 + 'c' - '0'
			if seenDecimal {
				digitsAfterDecimal++
			}
		} else if c == '.' {
			if seenDecimal {
				goto ERROR
			}
			seenDecimal = true
		} else {
			goto ERROR
		}
	}
	for digitsAfterDecimal < 6 {
		f *= 10
		digitsAfterDecimal++
	}
	if f != 0 && neg {
		f = -f
	}
	return f, nil
ERROR:
	return 0, errors.New("invalid fixed floating point number")
}

// FixedFloatFromFraction returns the FixedFloat that most precisely represents
// the fraction a/b.  It panics if b is not positive.
func FixedFloatFromFraction(a, b int) (f FixedFloat) {
	if b <= 0 {
		panic("non-positive denominator")
	}
	f = FixedFloat(a) * 10000000 / FixedFloat(b) // one extra digit
	if f%10 >= 5 {
		return f/10 + 1
	} else if f%10 <= -5 {
		return f/10 - 1
	} else {
		return f / 10
	}
}

func (f FixedFloat) String() (s string) {
	s = fmt.Sprintf("%d.%06d", f/1000000, f%1000000)
	s = strings.TrimRight(s, "0.")
	if s == "" {
		s = "0"
	}
	return s
}

// AsFraction returns the float as a numerator and denominator of an equivalent
// fraction.
func (f FixedFloat) AsFraction() (num, den int) {
	return int(f), 1000000
}

// Int returns the integer part of the the float.
func (f FixedFloat) Int() int {
	return int(f / 1000000)
}

// Mul multiplies two FixedFloats and returns the result.
func (f FixedFloat) Mul(o FixedFloat) FixedFloat {
	return f * o / 1000000
}

// Div divides two FixedFloats and returns the result.  The result is rounded if
// need be.
func (f FixedFloat) Div(o FixedFloat) FixedFloat {
	f = f * 10 / o
	if f%10 >= 5 {
		return f/10 + 1
	} else if f%10 <= -5 {
		return f/10 - 1
	} else {
		return f / 10
	}
}

const feetToMeters FixedFloat = 304800 // by definition

// GPSCoords holds a set of GPS coordinates.
type GPSCoords struct {
	// Latitude, in degrees north of the equator.
	Latitude FixedFloat
	// Longitude, in degrees east of the zero meridian.
	Longitude FixedFloat
	// Altitude, in meters above sea level (with 0 = unspecified)
	Altitude FixedFloat
}

// ParseGPSCoords parses a string to a GPSCoords structure.
func ParseGPSCoords(s string) (gc GPSCoords, err error) {
	var feet bool

	parts := strings.Split(s, ",")
	if len(parts) == 1 || len(parts) > 3 {
		goto ERROR
	}
	if len(parts) == 0 {
		return GPSCoords{}, nil
	}
	if gc.Latitude, err = ParseFixedFloat(parts[0]); err != nil {
		goto ERROR
	}
	if gc.Longitude, err = ParseFixedFloat(parts[1]); err != nil {
		goto ERROR
	}
	if len(parts) == 2 {
		return gc, nil
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
		goto ERROR
	}
	if gc.Altitude, err = ParseFixedFloat(parts[2]); err != nil {
		goto ERROR
	}
	if feet {
		gc.Altitude = gc.Altitude.Mul(feetToMeters)
	}
	return gc, nil
ERROR:
	return GPSCoords{}, fmt.Errorf("invalid value for gps: %q", s)
}

func (gc GPSCoords) String() string {
	var sb strings.Builder
	if !gc.Valid() {
		return ""
	}
	sb.WriteString(gc.Latitude.String())
	sb.WriteString(", ")
	sb.WriteString(gc.Longitude.String())
	if !gc.HasAltitude() {
		return sb.String()
	}
	sb.WriteString(", ")
	feet := gc.Altitude.Div(feetToMeters)
	sb.WriteString(feet.String())
	sb.WriteString("ft")
	return sb.String()
}

// HasAltitude returns whether the coordinates include an altitude.
func (gc GPSCoords) HasAltitude() bool {
	return gc.Valid() && gc.Altitude != 0
}

// Valid returns whether the coordinates are valid.
func (gc GPSCoords) Valid() bool {
	return gc.Latitude != 0 && gc.Longitude != 0
}

// Equal returns whether two GPSCoords are the same.
func (gc GPSCoords) Equal(o GPSCoords) bool {
	if gc.Valid() != o.Valid() {
		return false
	}
	if !gc.Valid() {
		return true
	}
	if gc.Latitude != o.Latitude || gc.Longitude != o.Longitude {
		return false
	}
	if gc.HasAltitude() != o.HasAltitude() {
		return false
	}
	if !gc.HasAltitude() {
		return true
	}
	return gc.Altitude == o.Altitude
}
