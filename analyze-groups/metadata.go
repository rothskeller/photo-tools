package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var exifcmd *exec.Cmd
var exifin io.WriteCloser
var exifout io.ReadCloser
var exifscan *bufio.Scanner

func startExiftool() {
	exifcmd = exec.Command("exiftool", "-stay_open", "True", "-@", "-")
	exifin, _ = exifcmd.StdinPipe()
	exifout, _ = exifcmd.StdoutPipe()
	exifscan = bufio.NewScanner(exifout)
	exifcmd.Start()
}

func stopExiftool() {
	fmt.Fprint(exifin, "-stay_open\nFalse\n")
}

func getMetadata(path string) (img image.Image, date, author, title, caption, keyword, face, coord, location string) {
	// STEP 1: Compare the images.
	if idx := strings.LastIndexByte(path, '.'); idx > 0 {
		ext := path[idx+1:]
		if strings.HasSuffix(ext, "_original") {
			ext = ext[:len(ext)-9]
		}
		switch ext {
		case "dng", "m4v", "mkv", "mov", "mp4", "tif", "wav":
			img = unreadableImage
		case "gif", "jpg", "png":
			if fh, err := os.Open(path); err == nil {
				if img, _, err = image.Decode(fh); err != nil {
					img = unreadableImage
				}
				fh.Close()
			} else {
				img = unreadableImage
			}
		}
	}
	// Step 2: Get the metadata.
	fmt.Fprintf(exifin, `-j
-G1
-struct
-Artist
-By-line
-Caption-Abstract
-City
-Country-PrimaryLocationCode
-Country-PrimaryLocationName
-CreateDate
-Creator
-Date
-DateCreated
-DateTime
-DateTimeCreated
-DateTimeOriginal
-Description
-GPSLatitude
-GPSLongitude
-GPSPosition
-Headline
-HierarchicalSubject
-ImageDescription
-Keywords
-LocationCreated
-LocationShown
-OffsetTimeOriginal
-PersonInImage
-Province-State
-RegionInfo
-Subject
-Sub-location
-SubSecDateTimeOriginal
-SubSecTimeOriginal
-TimeCreated
-Title
%s
-execute
`, path)
	var buf bytes.Buffer
	for exifscan.Scan() {
		line := exifscan.Text()
		if line == "{ready}" {
			break
		}
		buf.WriteString(line)
		buf.WriteRune('\n')
	}
	dec := json.NewDecoder(&buf)
	dec.DisallowUnknownFields()
	var mlist []Meta
	if err := dec.Decode(&mlist); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}
	if len(mlist) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: wrong number of responses\n")
		return
	}
	var meta = &mlist[0]
	// Render the date(s).
	date = strings.Join(mergeSubstrings(mergeStrings(meta.CompositeDateTimeCreated, meta.CompositeSubSecDateTimeOriginal, meta.ExifIFDDateTimeOriginal, meta.ExifIFDOffsetTimeOriginal, meta.ExifIFDCreateDate, meta.IPTCDateCreated, meta.IPTCTimeCreated, meta.QuickTimeCreateDate, meta.XMPexifDateTimeOriginal, meta.XMPphotoshopDateCreated, meta.XMPtiffDateTime, meta.XMPvideoDateTimeOriginal, meta.XMPxmpCreateDate)), "; ")
	// Render the artist(s).
	author = strings.Join(mergeStrings2([]string{meta.IFD0Artist}, meta.XMPdcCreator), "; ")
	// Render the title(s).
	title = strings.Join(mergeStrings(meta.XMPdcTitle), "; ")
	// Render the caption(s).
	caption = strings.Join(mergeStrings(meta.IFD0ImageDescription, meta.IPTCCaptionAbstract, meta.XMPdcDescription, meta.XMPtiffImageDescription), "; ")
	// Render the keywords.
	keyword = mergeKeywords(meta.IPTCKeywords, meta.XMPdcSubject, strings.Split(meta.XMPpdfKeywords, ", "), meta.XMPlrHierarchicalSubject)
	// Render the faces.
	face = mergeFaces(meta.XMPmwgrsRegionInfo.RegionList)
	// Render the GPS coordinates.
	coord = strings.Join(mergeSubstrings(mergeStrings(meta.CompositeGPSLatitude, meta.CompositeGPSLongitude, meta.CompositeGPSPosition, meta.GPSGPSLatitude, meta.GPSGPSLongitude, meta.XMPexifGPSLatitude, meta.XMPexifGPSLongitude)), "; ")
	// Render the location(s).
	location = strings.Join(mergeStrings(meta.IPTCCity, meta.IPTCCountryPrimaryLocationCode, meta.IPTCCountryPrimaryLocationName, meta.XMPphotoshopCity), "; ")
	return
}

type Meta struct {
	// Date Fields
	CompositeDateTimeCreated        string `json:"Composite:DateTimeCreated"`
	CompositeSubSecDateTimeOriginal string `json:"Composite:SubSecDateTimeOriginal"`
	ExifIFDDateTimeOriginal         string `json:"ExifIFD:DateTimeOriginal"`
	ExifIFDOffsetTimeOriginal       string `json:"ExifIFD:OffsetTimeOriginal"`
	ExifIFDCreateDate               string `json:"ExifIFD:CreateDate"`
	IPTCDateCreated                 string `json:"IPTC:DateCreated"`
	IPTCTimeCreated                 string `json:"IPTC:TimeCreated"`
	QuickTimeCreateDate             string `json:"QuickTime:CreateDate"`
	XMPexifDateTimeOriginal         string `json:"XMP-exif:DateTimeOriginal"`
	XMPphotoshopDateCreated         string `json:"XMP-photoshop:DateCreated"`
	XMPtiffDateTime                 string `json:"XMP-tiff:DateTime"`
	XMPvideoDateTimeOriginal        string `json:"XMP-video:DateTimeOriginal"`
	XMPxmpCreateDate                string `json:"XMP-xmp:CreateDate"`
	// Artist Fields
	IFD0Artist   string   `json:"IFD0:Artist"`
	XMPdcCreator []string `json:"XMP-dc:Creator"`
	// Title Fields
	XMPdcTitle string `json:"XMP-dc:Title"`
	// Caption Fields
	IFD0ImageDescription    string `json:"IFD0:ImageDescription"`
	IPTCCaptionAbstract     string `json:"IPTC:Caption-Abstract"`
	XMPdcDescription        string `json:"XMP-dc:Description"`
	XMPtiffImageDescription string `json:"XMP-tiff:ImageDescription"`
	// Keyword Fields
	IPTCKeywords             []string `json:"IPTC:Keywords"`
	XMPdcSubject             []string `json:"XMP-dc:Subject"`
	XMPlrHierarchicalSubject []string `json:"XMP-lr:HierarchicalSubject"`
	XMPpdfKeywords           string   `json:"XMP-pdf:Keywords"`
	// Face Fields
	XMPmwgrsRegionInfo RegionInfo `json:"XMP-mwg-rs:RegionInfo"`
	// Geolocation Fields
	CompositeGPSLatitude  string `json:"Composite:GPSLatitude"`
	CompositeGPSLongitude string `json:"Composite:GPSLongitude"`
	CompositeGPSPosition  string `json:"Composite:GPSPosition"`
	GPSGPSLatitude        string `json:"GPS:GPSLatitude"`
	GPSGPSLongitude       string `json:"GPS:GPSLongitude"`
	XMPexifGPSLatitude    string `json:"XMP-exif:GPSLatitude"`
	XMPexifGPSLongitude   string `json:"XMP-exif:GPSLongitude"`
	// Location Fields
	IPTCCity                       string `json:"IPTC:City"`
	IPTCCountryPrimaryLocationCode string `json:"IPTC:Country-PrimaryLocationCode"`
	IPTCCountryPrimaryLocationName string `json:"IPTC:Country-PrimaryLocationName"`
	IPTCProvinceState              string `json:"IPTC:Province-State"`
	XMPphotoshopCity               string `json:"XMP-photoshop:City"`
	// Unused Fields
	SourceFile                string
	ExifIFDSubSecTimeOriginal interface{} `json:"ExifIFD:SubSecTimeOriginal"`
}
type RegionInfo struct {
	AppliedToDimensions Dimensions
	RegionList          []RegionStruct
}
type RegionStruct struct {
	Area     AreaStruct
	Name     string
	Rotation float64
	Type     string
}
type Dimensions struct {
	H    float64
	Unit string
	W    float64
}
type AreaStruct struct {
	D, H, W, X, Y float64
	Unit          string
}

func mergeStrings(sl ...string) []string {
	m := make(map[string]struct{})
	for _, s := range sl {
		if t := strings.TrimSpace(s); t != "" {
			m[t] = struct{}{}
		}
	}
	l := make([]string, 0, len(m))
	for s := range m {
		l = append(l, s)
	}
	return l
}

func mergeStrings2(a, b []string) []string {
	m := make(map[string]struct{})
	for _, s := range a {
		if t := strings.TrimSpace(s); t != "" {
			m[t] = struct{}{}
		}
	}
	for _, s := range b {
		if t := strings.TrimSpace(s); t != "" {
			m[t] = struct{}{}
		}
	}
	l := make([]string, 0, len(m))
	for s := range m {
		l = append(l, s)
	}
	return l
}

func mergeKeywords(kwl ...[]string) string {
	flat := make(map[string]struct{})
	hier := make(map[string]struct{})
	for _, kws := range kwl {
		for _, kw := range kws {
			if kw == "" {

			} else if strings.ContainsRune(kw, '|') {
				hier[kw] = struct{}{}
			} else {
				flat[kw] = struct{}{}
			}
		}
	}
	for kw := range hier {
		parts := strings.Split(kw, "|")
		for i, part := range parts {
			delete(hier, strings.Join(parts[:i], "|"))
			delete(flat, part)
		}
	}
	var list []string
	for f := range flat {
		list = append(list, f)
	}
	for h := range hier {
		list = append(list, h)
	}
	sort.Strings(list)
	return strings.Join(list, ", ")
}

func mergeFaces(regions []RegionStruct) string {
	var names []string
	for _, region := range regions {
		if region.Type == "Face" {
			names = append(names, region.Name)
		}
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func mergeSubstrings(list []string) (out []string) {
	for i := 0; i < len(list); i++ {
		found := false
		for j := 0; j < len(list); j++ {
			if i != j && strings.Contains(list[j], list[i]) {
				found = true
				break
			}
		}
		if !found {
			out = append(out, list[i])
		}
	}
	return out
}
