package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mailru/easyjson/jlexer"
)

type imageMetadata struct {
	time        dateRange
	gps         gpsCoords
	loc         location
	creator     string
	people      []string
	title       string
	description string
	keywords    []string

	hasEXIFDateTimeOriginal   bool
	hasEXIFOffsetTimeOriginal bool
	hasIPTCDateCreated        bool
	hasEXIFGPS                bool
	hasEXIFArtist             bool
	hasIPTCByline             bool
	hasIPTCHeadline           bool
	hasEXIFImageDescription   bool
	hasIPTCCaptionAbstract    bool
	hasIPTCKeywords           bool
}

func readMetadata(file string) (im *imageMetadata, err error) {
	switch {
	case strings.HasSuffix(file, ".jpg"):
		return readJPEGMetadata(file)
	default:
		return nil, fmt.Errorf("%s: unknown file format", file)
	}
}

func readJPEGMetadata(file string) (im *imageMetadata, err error) {
	var json []string
	json, err = exiftool.Run(
		"-coordFormat", "%+.6f",
		"-groupNames",
		"-json", "-struct",
		"-EXIF:DateTimeOriginal",
		"-EXIF:OffsetTimeOriginal",
		"-EXIF:SubSecTimeOriginal",
		"-IPTC:DateCreated",
		"-IPTC:TimeCreated",
		"-XMP:Date",
		"-EXIF:GPSLatitude",
		"-EXIF:GPSLatitudeRef",
		"-EXIF:GPSLongitude",
		"-EXIF:GPSLongitudeRef",
		"-EXIF:GPSAltitude",
		"-EXIF:GPSAltitudeRef",
		"-IPTC:Country-PrimaryLocationCode",
		"-IPTC:Country-PrimaryLocationName",
		"-IPTC:Province-State",
		"-IPTC:City",
		"-IPTC:Sub-location",
		"-XMP:LocationCreated",
		"-EXIF:Artist",
		"-IPTC:By-line",
		"-XMP:Creator",
		"-XMP:PersonInImage",
		"-IPTC:Headline",
		"-XMP:Title",
		"-EXIF:ImageDescription",
		"-IPTC:Caption-Abstract",
		"-XMP:Description",
		"-IPTC:Keywords",
		"-XMP:Subject",
		file,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR (%s): %s\n", file, err)
		return
	}
	if len(json) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR (%s): %d lines from exiftool, expected 1\n", file, len(json))
		return
	}
	if im, err = decodeMetadata(json[0]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR (%s): exiftool output: %s\n", file, err)
		return
	}
	if *readXMP {
		var xmpfn = file[:len(file)-3] + "xmp"
		if _, err = os.Stat(xmpfn); err == nil {
			var xmpim *imageMetadata
			if xmpim, err = readXMPMetadata(xmpfn); err != nil {
				return nil, err
			}
			mergeMetadata(im, xmpim)
		}
	}
	return im, nil
}

func readXMPMetadata(file string) (im *imageMetadata, err error) {
	return nil, nil
}

func mergeMetadata(im1, im2 *imageMetadata) {}

func decodeMetadata(json string) (im *imageMetadata, err error) {
	var (
		in       jlexer.Lexer
		nanosecs time.Duration
		tz       = time.Local
	)
	in.Data = []byte(json)
	im = new(imageMetadata)
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		switch key {
		case "EXIF:DateTimeOriginal":
			var t time.Time
			t, err = time.ParseInLocation("2006:01:02 15:04:05", in.String(), tz)
			t = t.Add(nanosecs)
			im.time.start, im.time.end = t, t
			im.hasEXIFDateTimeOriginal = true
			in.AddError(err)
		case "EXIF:OffsetTimeOriginal":
			var t time.Time
			t, err = time.Parse("-07:00", in.String())
			tz = t.Location()
			if im.hasEXIFDateTimeOriginal {
				im.time.start = im.time.start.In(tz)
				im.time.end = im.time.end.In(tz)
			}
			im.hasEXIFOffsetTimeOriginal = true
			in.AddError(err)
		case "EXIF:SubSecTimeOriginal":
			var ns int
			var s = in.String() + "000000000"
			ns, err = strconv.Atoi(s[:9])
			nanosecs = time.Duration(ns) * time.Nanosecond
			if im.hasEXIFDateTimeOriginal {
				im.time.start = im.time.start.Add(nanosecs)
				im.time.end = im.time.end.Add(nanosecs)
			}
			in.AddError(err)
		case "IPTC:DateCreated":
		case "IPTC:TimeCreated":
		case "XMP:Date":
		case "EXIF:GPSLatitude":
		case "EXIF:GPSLatitudeRef":
		case "EXIF:GPSLongitude":
		case "EXIF:GPSLongitudeRef":
		case "EXIF:GPSAltitude":
		case "EXIF:GPSAltitudeRef":
		case "IPTC:Country-PrimaryLocationCode":
		case "IPTC:Country-PrimaryLocationName":
		case "IPTC:Province-State":
		case "IPTC:City":
		case "IPTC:Sub-location":
		case "XMP:LocationCreated":
		case "EXIF:Artist":
		case "IPTC:By-line":
		case "XMP:Creator":
		case "XMP:PersonInImage":
		case "IPTC:Headline":
		case "XMP:Title":
		case "EXIF:ImageDescription":
		case "IPTC:Caption-Abstract":
		case "XMP:Description":
		case "IPTC:Keywords":
		case "XMP:Subject":
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	in.Consumed()
	return im, in.Error()
}
