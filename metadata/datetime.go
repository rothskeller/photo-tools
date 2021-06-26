package metadata

import (
	"errors"
	"strings"
	"time"
)

// DateTime represents a date and time.
type DateTime struct {
	date   string // YYYY-MM-DD
	time   string // HH:MM:SS
	subsec string // zero or more digits
	zone   string // empty, "Z", +HH:MM, or -HH:MM
}

// ErrParseDateTime is the error returned when a string cannot be parsed into a
// DateTime value (or portion thereof).
var ErrParseDateTime = errors.New("invalid DateTime value")

// Parse parses a string into a DateTime.  It returns ErrParseDateTime if the
// string is invalid.
func (dt *DateTime) Parse(s string) error {
	*dt = DateTime{} // clear old data
	if s == "" {
		return nil
	}
	if _, err := time.Parse("2006-01-02", s); err == nil {
		dt.date = s
		dt.time = "00:00:00"
		return nil
	}
	z := s[len(s)-1] == 'Z'
	if z {
		dt.zone = "Z"
		s = s[:len(s)-1]
	}
	if _, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		dt.date = s[:10]
		dt.time = s[11:19]
		if len(s) > 20 {
			dt.subsec = s[20:]
		}
		return nil
	} else if z {
		return ErrParseDateTime
	}
	if _, err := time.Parse("2006-01-02T15:04:05-07:00", s); err == nil {
		dt.date = s[:10]
		dt.time = s[11:19]
		if len(s) > 26 {
			dt.subsec = s[20 : len(s)-6]
		}
		dt.zone = s[len(s)-6:]
		if dt.zone == "-00:00" || dt.zone == "+00:00" {
			dt.zone = "Z"
		}
		return nil
	}
	// The remaining formats aren't strictly legal, but they have been seen
	// in my library so I'm supporting them.
	if _, err := time.Parse("2006-01-02T15:04", s); err == nil {
		dt.date = s[:10]
		dt.time = s[11:17] + ":00"
		return nil
	} else if z {
		return ErrParseDateTime
	}
	if _, err := time.Parse("2006-01-02T15:04-07:00", s); err == nil {
		dt.date = s[:10]
		dt.time = s[11:17] + ":00"
		dt.zone = s[len(s)-6:]
		if dt.zone == "-00:00" || dt.zone == "+00:00" {
			dt.zone = "Z"
		}
		return nil
	}
	return ErrParseDateTime
}

func (dt *DateTime) String() string {
	if dt == nil || dt.date == "" {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(dt.date)
	sb.WriteByte('T')
	sb.WriteString(dt.time)
	if dt.subsec != "" {
		sb.WriteByte('.')
		sb.WriteString(dt.subsec)
	}
	sb.WriteString(dt.zone)
	return sb.String()
}

// ParseEXIF parses a date and time as represented in EXIF metadata.  It returns
// ErrParseDateTime if the input is invalid.
func (dt *DateTime) ParseEXIF(datetime, subsec, offset string) error {
	*dt = DateTime{}
	if datetime == "" || strings.TrimSpace(datetime) == ":  :     :  :" || datetime == "0000:00:00 00:00:00" {
		// That last isn't legal, but it's seen empirically.
		if subsec != "" || offset != "" {
			return ErrParseDateTime
		}
		return nil
	}
	if t, err := time.Parse("2006:01:02 15:04:05", datetime); err == nil {
		dt.date = t.Format("2006-01-02")
		dt.time = datetime[11:19]
	} else {
		return ErrParseDateTime
	}
	if strings.IndexFunc(subsec, func(r rune) bool {
		return r < '0' || r > '9'
	}) >= 0 {
		return ErrParseDateTime
	}
	dt.subsec = subsec
	if offset != "" {
		if _, err := time.Parse("-07:00", offset); err != nil {
			return ErrParseDateTime
		}
		dt.zone = offset
	}
	return nil
}

// AsEXIF returns the data and time as it would be represented in EXIF metadata.
func (dt *DateTime) AsEXIF() (datetime, subsec, offset string) {
	if dt == nil || dt.date == "" {
		return "", "", ""
	}
	datetime = strings.Replace(dt.date, "-", ":", -1) + " " + dt.time
	subsec = dt.subsec
	if dt.zone == "Z" {
		offset = "+00:00"
	} else {
		offset = dt.zone
	}
	return
}

// ParseIPTC parses a date and time as represented in IPTC metadata.  It returns
// ErrParseDateTime if the input is invalid.
func (dt *DateTime) ParseIPTC(date, timev string) error {
	*dt = DateTime{}
	if date == "" || date == "00000000" {
		return nil
	}
	if dval, err := time.Parse("20060102", date); err == nil {
		dt.date = dval.Format("2006-01-02")
	} else {
		return ErrParseDateTime
	}
	if timev == "" {
		dt.time = "00:00:00"
		return nil
	}
	if tval, err := time.Parse("150405", timev); err == nil {
		dt.time = tval.Format("15:04:05")
		return nil
	}
	if tval, err := time.Parse("150405-0700", timev); err == nil {
		dt.time = tval.Format("15:04:05")
		dt.zone = tval.Format("-07:00")
		return nil
	}
	return ErrParseDateTime
}

// AsIPTC returns the data and time as it would be represented in IPTC metadata.
func (dt *DateTime) AsIPTC() (date, timev string) {
	if dt == nil || dt.date == "" {
		return "", ""
	}
	switch dt.zone {
	case "":
		return strings.Replace(dt.date, "-", "", -1), strings.Replace(dt.time, ":", "", -1)
	case "Z":
		return strings.Replace(dt.date, "-", "", -1), strings.Replace(dt.time, ":", "", -1) + "+0000"
	default:
		return strings.Replace(dt.date, "-", "", -1), strings.Replace(dt.time+dt.zone, ":", "", -1)
	}
}

// Empty returns whether the DateTime has a value.
func (dt *DateTime) Empty() bool {
	return dt == nil || dt.date == ""
}

// Equal returns whether two DateTimes are equal.
func (dt *DateTime) Equal(other *DateTime) bool {
	if dt.Empty() != other.Empty() {
		return false
	}
	if dt.Empty() {
		return true
	}
	return dt.date == other.date && dt.time == other.time && dt.subsec == other.subsec && dt.zone == other.zone
}

// Equivalent returns whether two DateTimes are equal, to the precision of the
// least precise of the two.  If so, it returns the more precise one.
func (dt *DateTime) Equivalent(other *DateTime) (bool, *DateTime) {
	if dt.Empty() != other.Empty() {
		return false, nil
	}
	if dt.Empty() {
		return true, dt
	}
	if dt.date != other.date || dt.time != other.time || dt.zone != other.zone {
		return false, nil
	}
	if dt.subsec == other.subsec || other.subsec == "" {
		return true, dt
	}
	if dt.subsec == "" {
		return true, other
	}
	return false, nil
}
