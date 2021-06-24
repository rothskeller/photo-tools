// fix-artist lists all of the artists found in all of the metadata tags in the
// named media files.  Then it takes an artist value on input, and updates all
// of the media files to use that artist tag.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/filefmt"
)

func main() {
	var (
		lines         []string
		table         *tabwriter.Writer
		filesToChange []string
		scan          *bufio.Scanner
		artists       []string
	)
	if len(os.Args) < 2 {
		os.Exit(2)
	}
	table = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	for _, file := range os.Args[1:] {
		handler := filefmt.HandlerFor(file)
		if handler == nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s: unsupported file format\n", file)
			continue
		}
		handler.ReadMetadata()
		if xmp := handler.XMP(false); xmp != nil {
			if vals := xmp.DCCreator(); len(vals) != 0 {
				for i, val := range vals {
					lines = append(lines, val)
					fmt.Fprintf(table, "%d)\t%s\tXMP.dc.creator[%d]\t%s\n", len(lines), file, i, val)
				}
			} else {
				fmt.Fprintf(table, "\t%s\tXMP.dc.creator\t\n", file)
			}
			if a := xmp.TIFFArtist(); a != "" {
				lines = append(lines, a)
				fmt.Fprintf(table, "%d)\t%s\tXMP.tiff.artist\t%s\n", len(lines), file, a)
			}
		}
		if exif := handler.EXIF(); exif != nil {
			if vals := exif.Artist(); len(vals) != 0 {
				for i, val := range vals {
					lines = append(lines, val)
					fmt.Fprintf(table, "%d)\t%s\tEXIF.artist[%d]\t%s\n", len(lines), file, i, val)
				}
			} else {
				fmt.Fprintf(table, "\t%s\tEXIF.artist\t\n", file)
			}
		}
		if iptc := handler.IPTC(); iptc != nil {
			if vals := iptc.Bylines(); len(vals) != 0 {
				for i, val := range vals {
					lines = append(lines, val)
					fmt.Fprintf(table, "%d)\t%s\tIPTC.byline[%d]\t%s\n", len(lines), file, i, val)
				}
			} else {
				fmt.Fprintf(table, "\t%s\tIPTC.byline\t\n", file)
			}
		}
		if problems := handler.Problems(); len(problems) != 0 {
			for _, problem := range handler.Problems() {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file, problem)
			}
		} else {
			filesToChange = append(filesToChange, file)
		}
	}
	table.Flush()
	if len(filesToChange) == 0 {
		os.Exit(0)
	}
RETRY:
	fmt.Println("Enter line # or value to use.  DELETE removes.")
	fmt.Print("? ")
	scan = bufio.NewScanner(os.Stdin)
	if !scan.Scan() {
		os.Exit(0)
	}
	if scan.Text() == "" {
		goto RETRY
	} else if scan.Text() == "DELETE" {
		artists = nil
	} else if line, err := strconv.Atoi(scan.Text()); err == nil {
		if line > 0 && line <= len(lines) {
			artists = []string{lines[line-1]}
		} else {
			fmt.Println("No such line.")
			goto RETRY
		}
	} else {
		artists = []string{scan.Text()}
	}
	for _, file := range filesToChange {
		handler := filefmt.HandlerFor(file)
		dirty := false
		handler.ReadMetadata()
		if xmp := handler.XMP(false); xmp != nil {
			if !equalSS(xmp.DCCreator(), artists) || xmp.TIFFArtist() != "" {
				xmp.SetDCCreator(artists)
				xmp.SetTIFFArtist("") // always delete
				dirty = true
			}
		}
		if exif := handler.EXIF(); exif != nil {
			if !equalSS(exif.Artist(), artists) {
				exif.SetArtist(artists)
				dirty = true
			}
		}
		if iptc := handler.IPTC(); iptc != nil {
			if !equalSS(iptc.Bylines(), artists) {
				iptc.SetBylines(artists)
				dirty = true
			}
		}
		if dirty {
			if err := handler.SaveMetadata(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file, err)
			}
		}
	}
}

func equalSS(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
