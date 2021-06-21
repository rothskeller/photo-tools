package main

import (
	"fmt"
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

func showFields(file string, h filefmt.FileHandler, filecount, filelen int) {
	var outfn func(field, value string, change bool)

	if filecount == 1 {
		outfn = printFieldConsistValue
	} else if len(kwviews) != 0 && len(views) != 0 {
		fmt.Printf("=== %s ===\n", file)
		outfn = printFieldConsistValue
	} else if len(views) > 1 {
		fmt.Printf("=== %s ===\n", file)
		outfn = printFieldConsistValue
	} else {
		outfn = func(field, value string, change bool) {
			printFileConsistValue(file, filelen, value, change)
		}
	}
	if views["title"] {
		values, change := strmeta.GetTitles(h)
		for _, val := range values {
			outfn("Title", val, change)
		}
		if len(values) == 0 {
			outfn("Title", "", change)
		}
	}
	if views["datetime"] {
		values, change := strmeta.GetDatesTimes(h)
		for _, val := range values {
			outfn("DateTime", val, change)
		}
		if len(values) == 0 {
			outfn("DateTime", "", change)
		}
	}
	if views["gps"] {
		values, change := strmeta.GetGPSs(h)
		for _, val := range values {
			outfn("Lat/Long", val.String(), change)
		}
		if len(values) == 0 {
			outfn("Lat/Long", "", change)
		}
	}
	if views["location"] {
		values, change := strmeta.GetLocationsCaptured(h)
		for _, val := range values {
			outfn("Location", val.String(), change)
		}
		if len(values) == 0 {
			outfn("Location", "", change)
		}
	}
	if views["shown"] {
		values, change := strmeta.GetLocationsShown(h)
		for _, val := range values {
			outfn("Shown", val.String(), change)
		}
		if len(values) == 0 {
			outfn("Shown", "", change)
		}
	}
	if views["artist"] {
		values, change := strmeta.GetArtists(h)
		for _, val := range values {
			outfn("Artist", val, change)
		}
		if len(values) == 0 {
			outfn("Artist", "", change)
		}
	}
	if len(kwviews) == 1 {
		for prefix := range kwviews {
			values, change := strmeta.GetKeywords(h, prefix, true)
			for _, val := range values {
				outfn("Keyword", val, change)
			}
			if len(values) == 0 {
				outfn("Keyword", "", change)
			}
		}
	}
	if len(kwviews) > 1 {
		values, change := strmeta.GetKeywords(h, "", true)
		shown := false
	VALUES:
		for _, val := range values {
			for prefix := range kwviews {
				if strings.HasPrefix(val, prefix) {
					outfn("Keyword", val, change)
					shown = true
					continue VALUES
				}
			}
		}
		if !shown {
			outfn("Keyword", "", change)
		}
	}
	if views["caption"] {
		values, change := strmeta.GetCaptions(h)
		for _, val := range values {
			outfn("Caption", val, change)
		}
		if len(values) == 0 {
			outfn("Caption", "", change)
		}
	}
}

var changeFlag = map[bool]string{
	false: "   ",
	true:  "CHG",
}

func printFieldConsistValue(field, value string, change bool) {
	fmt.Printf("%-8s %s %s\n", field, changeFlag[change], escapeString(value))
}

func printFileConsistValue(file string, filelen int, value string, change bool) {
	fmt.Printf("%-*s %s %s\n", filelen, file, changeFlag[change], escapeString(value))
}

func escapeString(s string) string {
	return strings.Replace(strings.Replace(s, "\\", "\\\\", -1), "\n", "\\n", -1)
}
