package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	readXMP     = flag.Bool("read-xmp", false, "read XMP sidecar file")
	creator     = flag.String("creator", "", "person who took the photo")
	mine        = flag.Bool("mine", false, "-creator 'Steve Roth'")
	title       = flag.String("title", "", "title")
	description = flag.String("desc", "", "description (or @filename)")
	adjustTime  timeOffset
	timeZone    timeOffset
	date        dateRange
	gps         gpsCoords
	loc         location
	people      opList
	keywords    opList
	exiftool    EXIFTool
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: photo-meta [options] image-file...`)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Var(&adjustTime, "adjust-time", "adjust time by offset (±HH:MM[:SS.sss])")
	flag.Var(&timeZone, "tz", "time zone offset (±HH:MM[:SS.sss])")
	flag.Var(&date, "date", "date and time (range)")
	flag.Var(&gps, "gps", "latitude, longitude[, altitude]")
	flag.Var(&loc, "loc", "countrycode:[state]:[city][:sublocation]")
	flag.Var(&people, "person", "[+-=]person,...")
	flag.Var(&keywords, "kw", "[+-=]keyword,...")
	flag.Parse()
	if strings.HasPrefix(*description, "@") {
		if by, err := ioutil.ReadFile((*description)[1:]); err == nil {
			*description = string(by)
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "ERROR: no images specified")
		usage()
	}
	if err := exiftool.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := exiftool.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}()
	for _, arg := range flag.Args() {
		processImage(arg)
	}
}

type timeOffset time.Duration

func (to *timeOffset) Set(s string) error {
	var (
		neg bool
		t   time.Time
		err error
	)
	if s == "" {
		goto ERROR
	}
	if s[0] == '-' {
		neg = true
	} else if s[0] != '-' {
		goto ERROR
	}
	if t, err = time.Parse("15:04:05", s[1:]); err != nil {
		if t, err = time.Parse("15:04", s[1:]); err != nil {
			goto ERROR
		}
	}
	*to = timeOffset(t.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)))
	if neg {
		*to = -*to
	}
	return nil
ERROR:
	return errors.New("invalid time offset format")
}

func (to timeOffset) String() string {
	var sb strings.Builder
	if to < 0 {
		sb.WriteByte('-')
		to = -to
	} else {
		sb.WriteByte('+')
	}
	fmt.Fprintf(&sb, "%02d", to/timeOffset(time.Hour))
	to %= timeOffset(time.Hour)
	fmt.Fprintf(&sb, ":%02d", to/timeOffset(time.Minute))
	to %= timeOffset(time.Minute)
	if to != 0 {
		fmt.Fprintf(&sb, ":%02d", to/timeOffset(time.Second))
		to %= timeOffset(time.Second)
	}
	if to != 0 {
		fmt.Fprintf(&sb, ".%03d", to/timeOffset(time.Millisecond))
	}
	return sb.String()
}

type dateRange struct {
	start time.Time
	end   time.Time
}

func (dr *dateRange) Set(s string) error {
	var (
		slash  int
		shastz bool
		ehastz bool
		se     time.Time
		es     time.Time
		err    error
	)
	if slash = strings.IndexByte(s, '/'); slash < 0 {
		if dr.start, dr.end, shastz, err = parseTime(s); err != nil {
			return err
		}
		if !shastz {
			return errors.New("-date with time requires time zone")
		}
		return nil
	}
	if dr.start, se, shastz, err = parseTime(s[:slash]); err != nil {
		return err
	}
	if es, dr.end, ehastz, err = parseTime(s[slash+1:]); err != nil {
		return err
	}
	if es.Before(se) {
		return errors.New("-date range out of order")
	}
	if shastz && ehastz && se.Location() != es.Location() {
		return errors.New("-date range must have one time zone")
	} else if shastz && !ehastz {
		dr.end = dr.end.In(se.Location())
	} else if ehastz && !shastz {
		dr.start = dr.start.In(es.Location())
	} else if !shastz && !ehastz {
		return errors.New("-date with time requires time zone")
	}
	return nil
}

func parseTime(s string) (st, et time.Time, tz bool, err error) {
	if st, err = time.Parse("2006-01-02T15:04:05.000Z07:00", s); err == nil {
		return st, st, true, nil
	}
	if st, err = time.ParseInLocation("2006-01-02T15:04:05.000", s, time.Local); err == nil {
		return st, st, false, nil
	}
	if st, err = time.Parse("2006-01-02T15:04:05Z07:00", s); err == nil {
		et = st.Add(999 * time.Millisecond)
		return st, et, true, nil
	}
	if st, err = time.ParseInLocation("2006-01-02T15:04:05", s, time.Local); err == nil {
		et = st.Add(999 * time.Millisecond)
		return st, et, false, nil
	}
	if st, err = time.Parse("2006-01-02T15:04Z07:00", s); err == nil {
		et = st.Add(59*time.Second + 999*time.Millisecond)
		return st, et, true, nil
	}
	if st, err = time.ParseInLocation("2006-01-02T15:04", s, time.Local); err == nil {
		et = st.Add(59*time.Second + 999*time.Millisecond)
		return st, et, false, nil
	}
	if st, err = time.Parse("2006-01-02", s); err == nil {
		et = time.Date(st.Year(), st.Month(), st.Day(), 23, 59, 59, 999000000, st.Location())
		return st, et, true, nil
	}
	if st, err = time.Parse("2006-01", s); err == nil {
		et = time.Date(st.Year(), st.Month()+1, 1, 0, 0, 0, 0, st.Location()).Add(-time.Millisecond)
		return st, et, true, nil
	}
	if st, err = time.Parse("2006", s); err == nil {
		et = time.Date(st.Year()+1, 1, 1, 0, 0, 0, 0, st.Location()).Add(-time.Millisecond)
		return st, et, true, nil
	}
	return st, et, false, err
}

func (dr *dateRange) String() string {
	if dr.start.IsZero() {
		return ""
	}
	ss := dr.start.Format("2006-01-02T15:04:05.000-07:00")
	if dr.start.Equal(dr.end) {
		return ss
	}
	es := dr.end.Format("2006-01-02T15:04:05.000-07:00")
	const min = "0000-01-01T00:00:00.000-07:00"
	var equalUntil int
	for _, idx := range []int{4, 7, 10, 16, 19} {
		if ss[:idx] == es[:idx] {
			equalUntil = idx
		}
	}
	var max = time.Date(dr.end.Year(), dr.end.Month()+1, 1, 0, 0, 0, 0, dr.end.Location()).Add(-time.Millisecond).Format("2006-01-02T15:04:05.000-07:00")
	var minmaxAfter int
	for _, idx := range []int{23, 19, 16, 10, 7, 4} {
		if ss[idx:] == min[idx:] && es[idx:] == max[idx:] {
			minmaxAfter = idx
		}
	}
	switch {
	case equalUntil == minmaxAfter && equalUntil <= 10:
		return ss[:equalUntil]
	case equalUntil == minmaxAfter:
		return ss[:equalUntil] + ss[23:]
	case minmaxAfter <= 10:
		return ss[:minmaxAfter] + "/" + es[:minmaxAfter]
	default:
		return ss[:minmaxAfter] + ss[23:] + "/" + es[:minmaxAfter] + es[23:]
	}
}

type gpsCoords struct {
	lat      float64
	long     float64
	alt      float64
	altValid bool
}

var sepRE = regexp.MustCompile(`\s*,\s*`)

func (g *gpsCoords) Set(s string) (err error) {
	var parts = sepRE.Split(s, -1)
	if len(parts) < 2 || len(parts) > 3 {
		goto ERROR
	}
	if g.lat, err = strconv.ParseFloat(parts[0], 64); err != nil || g.lat < -90 || g.lat > 90 {
		goto ERROR
	}
	if g.long, err = strconv.ParseFloat(parts[1], 64); err != nil || g.long < -180 || g.long > 180 {
		goto ERROR
	}
	if len(parts) == 3 {
		if g.alt, err = strconv.ParseFloat(parts[2], 64); err != nil {
			goto ERROR
		}
		g.altValid = true
	}
	return nil
ERROR:
	return errors.New("invalid GPS coordinate format")
}

func (g *gpsCoords) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%f,%f", g.lat, g.long)
	if g.altValid {
		fmt.Fprintf(&sb, ",%f", g.alt)
	}
	return sb.String()
}

type location struct {
	countrycode string
	state       string
	city        string
	sublocation string
}

func (l *location) Set(s string) error {
	var parts = strings.SplitN(s, ":", 4)
	if len(parts) < 3 {
		return errors.New("invalid location format")
	}
	l.countrycode = strings.TrimSpace(parts[0])
	l.state = strings.TrimSpace(parts[1])
	l.city = strings.TrimSpace(parts[2])
	if len(parts) == 4 {
		l.sublocation = strings.TrimSpace(parts[3])
	}
	if l.countrycode == "" || (l.city == "" && l.sublocation == "") {
		return errors.New("invalid location format")
	}
	return nil
}

func (l *location) String() string {
	if l.sublocation == "" {
		return strings.Join([]string{l.countrycode, l.state, l.city}, ":")
	}
	return strings.Join([]string{l.countrycode, l.state, l.city, l.sublocation}, ":")
}

type opList struct {
	op   byte
	list []string
}

func (ol *opList) Set(s string) error {
	if s == "" {
		return errors.New("invalid list")
	}
	if s[0] == '+' || s[0] == '-' || s[0] == '=' {
		ol.op = s[0]
		s = s[1:]
	} else {
		ol.op = '='
	}
	if s == "" {
		return errors.New("invalid list")
	}
	ol.list = sepRE.Split(s, -1)
	for i := range ol.list {
		ol.list[i] = strings.TrimSpace(ol.list[i])
	}
	return nil
}

func (ol *opList) String() string {
	return string(ol.op) + strings.Join(ol.list, ",")
}
