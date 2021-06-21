// md is a program for viewing and editing media file metadata, tailored to the
// metadata conventions in my library.
package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/strmeta"
)

func usage() {
	fmt.Fprint(os.Stderr, `usage: md [fieldname][op][value]... file...
Fields are:
  a(rtist)
  d(atetime)   YYYY-MM-DDTHH:MM:SS.sss±HH:MM
  t(itle)
  c(aption)    (use - to read from stdin)
  k(eywords)   += adds   -= removes   = replaces
  g(ps)        lat,long[,alt]   all signed floats
  l(ocation)   countrycode/state/city/sublocation
  places       equivalent to keyword starting with PLACES
  people       equivalent to keyword starting with PEOPLE
  groups       equivalent to keyword starting with GROUPS
  topics       equivalent to keyword starting with TOPICS
See MANUAL.md for more details.
`)
	os.Exit(2)
}

type operation func(filefmt.FileHandler) error

var datetimeRE = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d(?:T\d\d:\d\d:\d\d(?:\.\d+)?)?(?:[-+]\d\d:\d\d|Z)?$`)
var locationRE = regexp.MustCompile(`^(?:[-+]?\d+(?:\.\d*)?/[-+]?\d+(?:\.\d*)?/(?:[-+]?\d+(?:\.\d*)?/)?)?[A-Z]{3}(?:/-[^/]+)*/(?:[^-/][^/]+)?(?:/-[^/]+)*/((?:[^-/][^/]+)?)(?:/-[^/]+)*(?:/([^-/][^/]+]))?(?:/-[^/]+)*$`)

var changes []operation
var views = map[string]bool{}
var kwviews = map[string]bool{}

func main() {
	var (
		files  []string
		opargs []string
		maxlen int
	)
	// Find the first filename argument in the argument list.  It is the
	// first argument that has a dot, not preceded by an equals sign, and
	// not starting with a plus or minus sign.
	for idx, arg := range os.Args[1:] {
		if arg == "" {
			usage()
		}
		if arg[0] == '-' || arg[0] == '+' {
			continue
		}
		if dot := strings.IndexByte(arg, '.'); dot >= 0 {
			if equal := strings.IndexByte(arg, '='); equal < 0 || equal > dot {
				files = os.Args[idx+1:]
				opargs = os.Args[1 : idx+1]
				break
			}
		}
	}
	if files == nil {
		fmt.Fprintf(os.Stderr, "ERROR: no files found on command line\n")
		usage()
	}
	for _, file := range files {
		if len(file) > maxlen {
			maxlen = len(file)
		}
	}
	// Now parse each of the operations arguments.
	changes = append(changes, strmeta.RecordOldPlaceKeywords)
	for _, arg := range opargs {
		var opidx int
		var field string
		var op byte
		var val string

		switch {
		case arg[0] == '+':
			if len(arg) == 1 {
				fmt.Fprintf(os.Stderr, "ERROR: missing keyword to add\n")
				usage()
			}
			changes = append(changes, func(h filefmt.FileHandler) error {
				return strmeta.AddKeyword(h, arg[1:])
			})
			continue
		case arg[0] == '-':
			if len(arg) == 1 {
				fmt.Fprintf(os.Stderr, "ERROR: missing keyword to remove\n")
				usage()
			}
			changes = append(changes, func(h filefmt.FileHandler) error {
				return strmeta.RemoveKeyword(h, arg[1:])
			})
			continue
		}
		if opidx = strings.IndexByte(arg, '='); opidx >= 0 {
			val = arg[opidx+1:]
		}
		if opidx >= 1 && (arg[opidx-1] == '-' || arg[opidx-1] == '+') {
			opidx--
		}
		if opidx >= 0 {
			op = arg[opidx]
			field = strings.ToUpper(arg[:opidx])
		} else {
			field = strings.ToUpper(arg)
		}
		switch field {
		case "PLACES", "PEOPLE", "GROUPS", "TOPICS":
			switch op {
			case 0:
				kwviews[field+"/"] = true
			case '=':
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.RemoveAllKeywords(h, field+"/"+val)
				})
			case '+':
				if val == "" {
					fmt.Fprintf(os.Stderr, "ERROR: missing keyword to add\n")
					usage()
				}
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.AddKeyword(h, field+"/"+val)
				})
			case '-':
				if val == "" {
					fmt.Fprintf(os.Stderr, "ERROR: missing keyword to remove\n")
					usage()
				}
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.RemoveKeyword(h, field+"/"+val)
				})
			}
		case "KEYWORDS", "KEYWORD", "K":
			switch op {
			case 0:
				kwviews[""] = true
			case '=':
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.RemoveKeyword(h, val)
				})
			case '+':
				if val == "" {
					fmt.Fprintf(os.Stderr, "ERROR: missing keyword to add\n")
					usage()
				}
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.AddKeyword(h, val)
				})
			case '-':
				if val == "" {
					fmt.Fprintf(os.Stderr, "ERROR: missing keyword to remove\n")
					usage()
				}
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.RemoveKeyword(h, val)
				})
			}
		case "ARTIST", "A":
			switch op {
			case 0:
				views["artist"] = true
			case '=':
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.SetArtist(h, val)
				})
			default:
				fmt.Fprintf(os.Stderr, "ERROR: invalid operation on artist: %q\n", arg)
				usage()
			}
		case "DATETIME", "DATE", "D":
			switch op {
			case 0:
				views["datetime"] = true
			case '=':
				if val == "" || datetimeRE.MatchString(val) {
					if len(val) == 10 {
						val += "T00:00:00"
					}
					if strings.HasSuffix(val, "-00:00") || strings.HasSuffix(val, "+00:00") {
						val = val[:len(val)-6] + "Z"
					}
					changes = append(changes, func(h filefmt.FileHandler) error {
						return strmeta.SetDateTime(h, val)
					})
				} else {
					fmt.Fprintf(os.Stderr, "ERROR: invalid value for datetime: %q\n", val)
					usage()
				}
			default:
				fmt.Fprintf(os.Stderr, "ERROR: invalid operation on datetime: %q\n", arg)
				usage()
			}
		case "TITLE", "T":
			switch op {
			case 0:
				views["title"] = true
			case '=':
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.SetTitle(h, val)
				})
			default:
				fmt.Fprintf(os.Stderr, "ERROR: invalid operation on title: %q\n", arg)
				usage()
			}
		case "CAPTION", "C":
			switch op {
			case 0:
				views["caption"] = true
			case '=':
				if val == "-" {
					by, err := io.ReadAll(os.Stdin)
					if err != nil {
						fmt.Fprintf(os.Stderr, "ERROR: stdin: %s\n", err)
						os.Exit(1)
					}
					val = string(by)
				}
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.SetCaption(h, val)
				})
			default:
				fmt.Fprintf(os.Stderr, "ERROR: invalid operation on caption: %q\n", arg)
				usage()
			}
		case "GPS", "G":
			switch op {
			case 0:
				views["gps"] = true
			case '=':
				gc, err := metadata.ParseGPSCoords(val)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: invalid value for gps: %q\n", val)
					usage()
				}
				changes = append(changes, func(h filefmt.FileHandler) error {
					return strmeta.SetGPS(h, gc)
				})
			default:
				fmt.Fprintf(os.Stderr, "ERROR: invalid operation on gps: %q", arg)
				usage()
			}
		case "LOCATION", "L":
			switch op {
			case 0:
				views["location"] = true
			case '=':
				loc := metadata.ParseLocation(val)
				if loc == nil {
					fmt.Fprintf(os.Stderr, "ERROR: invalid value for location: %q\n", val)
					usage()
				}
				changes = append(changes, func(h filefmt.FileHandler) error { return strmeta.SetLocationCaptured(h, loc) })
			default:
				fmt.Fprintf(os.Stderr, "ERROR: invalid operation on location: %q\n", arg)
				usage()
			}
		case "SHOWN", "S":
			switch op {
			case 0:
				views["shown"] = true
			case '=':
				loc := metadata.ParseLocation(val)
				if loc == nil {
					fmt.Fprintf(os.Stderr, "ERROR: invalid value for shown: %q\n", val)
					usage()
				}
				changes = append(changes, func(h filefmt.FileHandler) error { return strmeta.SetLocationShown(h, loc) })
			default:
				fmt.Fprintf(os.Stderr, "ERROR: invalid operation on shown: %q\n", arg)
				usage()
			}
		default:
			if op != 0 {
				fmt.Fprintf(os.Stderr, "ERROR: invalid operation: %q\n", arg)
				usage()
			}
			for _, c := range field {
				switch c {
				case 'A':
					views["artist"] = true
				case 'D':
					views["datetime"] = true
				case 'T':
					views["title"] = true
				case 'C':
					views["caption"] = true
				case 'G':
					views["gps"] = true
				case 'L':
					views["location"] = true
				case 'S':
					views["shown"] = true
				case 'K':
					kwviews[""] = true
				default:
					fmt.Fprintf(os.Stderr, "ERROR: invalid operation: %q\n", arg)
					usage()
				}
			}
		}
	}
	if len(changes) == 1 {
		changes = nil // remove the record old place keywords
	}
	if len(changes) == 0 && len(views) == 0 && len(kwviews) == 0 {
		views["artist"] = true
		views["caption"] = true
		views["datetime"] = true
		views["gps"] = true
		views["location"] = true
		views["shown"] = true
		views["title"] = true
		kwviews[""] = true
	} else if len(changes) != 0 {
		changes = append(changes, strmeta.RemoveOldPlaceKeywords)
	}
	// Walk through the list of files.
FILES:
	for _, file := range files {
		handler := filefmt.HandlerFor(file)
		if handler == nil {
			fmt.Fprintf(os.Stderr, "ERROR: no file format handler for %s\n", file)
			continue
		}
		problems := handler.ReadMetadata()
		if len(problems) != 0 {
			for _, p := range problems {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file, p)
			}
			continue
		}
		for _, change := range changes {
			if err := change(handler); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file, err)
				continue FILES
			}
		}
		if len(changes) != 0 {
			if xmp := handler.XMP(false); xmp != nil {
				// Get rid of namespaces for random software
				// that are added by digiKam.
				xmp.RemoveNamespace("acdsee", "http://ns.acdsee.com/iptc/1.0/")
				xmp.RemoveNamespace("video", "http://www.video/")
			}
			if err := handler.SaveMetadata(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file, err)
				continue
			}
		}
		if len(changes) == 0 && len(kwviews) == 0 && len(views) == 1 && views["caption"] && len(files) == 1 {
			caption := strmeta.GetCaption(handler)
			fmt.Print(caption)
			if !strings.HasSuffix(caption, "\n") {
				fmt.Println()
			}
			continue
		}
		showFields(file, handler, len(files), maxlen)
		if problems := handler.Problems(); len(problems) != 0 {
			for _, p := range problems {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file, p)
			}
		}
	}
}
