package metadata

import (
	"errors"
	"fmt"
	"strings"
)

// FixedFloat stores a floating point number as an signed integer scaled
// 1,000,000 higher, so that it always has exactly six digits of precision after
// the decimal point.  (That's enough to get us down to hundredths of a second
// when storing an angle in degrees, which is close enough for GPS coordinates.)
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
			f = f*10 + FixedFloat(c) - '0'
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

// FixedFloatFromFloat returns the FixedFloat that most precisely represents the
// supplied floating point number.
func FixedFloatFromFloat(v float64) (f FixedFloat) {
	f = FixedFloat(v*10000000.0 + 5.0)
	f /= 10
	return f
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
	frac := f % 1000000
	if frac < 0 {
		frac = -frac
	}
	s = fmt.Sprintf("%d.%06d", f/1000000, frac)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
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

// AsFloat64 returns the float as a (possibly imprecise) Go float value.
func (f FixedFloat) AsFloat64() float64 {
	return float64(f) / 1000000.0
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
	f = f * 10000000 / o
	if f%10 >= 5 {
		return f/10 + 1
	} else if f%10 <= -5 {
		return f/10 - 1
	} else {
		return f / 10
	}
}
