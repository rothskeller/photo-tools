// adjust-photo-times walks through the numbered images in the current directory
// and adjusts the timestamps of any files whose timestamps are out of bounds.
// Generally these are images that were edited in files that didn't maintain the
// file creation date, or images created later such as title slides.
//
// The arguments to adjust-photo-times determine what is considered "out of
// bounds".  There are two arguments: the starting date of the valid range and
// the ending date of the valid range.  Each one can have the form YYYY,
// YYYY-MM, or YYYY-MM-DD, with the unspecified parts being wildcards.  The
// second argument can be omitted if it is the same as the first.  The most
// common usage is simply listing a single year on the command line, and any
// images dated outside that year are "out of bounds".
//
// When an image's timestamp is out of bounds, it is adjusted to be inline with
// the (validly timestamped) images closest to it in the numbered directory
// order.  By preference, it will be stamped at midnight at the beginning of the
// same day as the next validly stamped image after it.  (This makes it clear
// that the time of day is artificial.)  However, if that would place it out of
// sequence relative to the previous validly stamped image before it, the target
// image will be given the exact same stamp as the next validly stamped image
// after it.  If the target image does not have any validly stamped images after
// it, it will be stamped at 11:59:59 on the same day as the previous validly
// stamped image before it.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	var (
		validStart time.Time
		validEnd   time.Time
		validEndS  string
		images     []string
		stamps     []time.Time
	)
	switch len(os.Args) {
	case 2:
		validEndS = os.Args[1]
	case 3:
		validEndS = os.Args[2]
	default:
		fmt.Fprintf(os.Stderr, "usage: adjust-photo-times valid-start [valid-end]\n")
		os.Exit(2)
	}
	switch len(os.Args[1]) {
	case 4:
		validStart, _ = time.ParseInLocation("2006", os.Args[1], time.Local)
	case 7:
		validStart, _ = time.ParseInLocation("2006-01", os.Args[1], time.Local)
	case 10:
		validStart, _ = time.ParseInLocation("2006-01-02", os.Args[1], time.Local)
	}
	if validStart.IsZero() {
		fmt.Fprintf(os.Stderr, "usage: adjust-photo-times valid-start [valid-end]\n")
		fmt.Fprintf(os.Stderr, "ERROR: valid-start must be YYYY, YYYY-MM, or YYYY-MM-DD\n")
		os.Exit(2)
	}
	switch len(validEndS) {
	case 4:
		validEnd, _ = time.ParseInLocation("2006", validEndS, time.Local)
	case 7:
		validEnd, _ = time.ParseInLocation("2006-01", validEndS, time.Local)
	case 10:
		validEnd, _ = time.ParseInLocation("2006-01-02", validEndS, time.Local)
	}
	if validEnd.IsZero() {
		fmt.Fprintf(os.Stderr, "usage: adjust-photo-times valid-start [valid-end]\n")
		fmt.Fprintf(os.Stderr, "ERROR: valid-end must be YYYY, YYYY-MM, or YYYY-MM-DD\n")
		os.Exit(2)
	}
	switch len(validEndS) {
	case 4:
		validEnd = time.Date(validEnd.Year(), 12, 31, 11, 59, 59, 0, time.Local)
	case 7:
		validEnd = time.Date(validEnd.Year(), validEnd.Month(), 31, 11, 59, 59, 0, time.Local)
	case 10:
		validEnd = time.Date(validEnd.Year(), validEnd.Month(), validEnd.Day(), 11, 59, 59, 0, time.Local)
	}
	if !validStart.Before(validEnd) {
		fmt.Fprintf(os.Stderr, "usage: adjust-photo-times valid-start [valid-end]\n")
		fmt.Fprintf(os.Stderr, "ERROR: valid-start must be before valid-end\n")
		os.Exit(2)
	}
	images, _ = filepath.Glob("[0-9][0-9][0-9]_*")
	j := 0
	for _, im := range images {
		if !strings.HasSuffix(im, ".xmp") {
			images[j] = im
			j++
		}
	}
	images = images[:j]
	sort.Strings(images)
	stamps = make([]time.Time, 0, len(images))
	for _, image := range images {
		var (
			cmd   *exec.Cmd
			dto   []byte
			stamp time.Time
			err   error
		)
		fmt.Print(".")
		cmd = exec.Command("exiftool", "-s3", "-DateTimeOriginal", "-srcfile", "%f.xmp", "-srcfile", "@", image)
		cmd.Stderr = os.Stderr
		if dto, err = cmd.Output(); err != nil {
			fmt.Fprintf(os.Stderr, "\nERROR: date from %s: %s\n", image, err)
			os.Exit(1)
		}
		if len(dto) != 0 {
			if stamp, err = time.ParseInLocation("2006:01:02 15:04:05.000Z\n", string(dto), time.Local); err != nil {
				fmt.Fprintf(os.Stderr, "\nWARNING: date from %s (%s): %s", image, string(dto), err)
			}
		}
		stamps = append(stamps, stamp)
	}
	fmt.Println()
	var lastValid time.Time
	for i, s := range stamps {
		var (
			cmd *exec.Cmd
			err error
		)
		if !s.Before(validStart) && !s.After(validEnd) {
			lastValid = s
			continue
		}
		s = time.Time{}
		for j := i + 1; j < len(stamps); j++ {
			if stamps[j].Before(validStart) || stamps[j].After(validEnd) {
				continue
			}
			s = time.Date(stamps[j].Year(), stamps[j].Month(), stamps[j].Day(), 0, 0, 0, 0, time.Local)
			if s.Before(lastValid) {
				s = stamps[j]
			}
		}
		if s.IsZero() {
			s = time.Date(lastValid.Year(), lastValid.Month(), lastValid.Day(), 23, 59, 59, 0, time.Local)
		}
		cmd = exec.Command("exiftool", "-s3", "-DateTimeOriginal="+s.Format("2006:01:02 15:04:05"), "-srcfile", "%f.xmp", "-srcfile", "@", "-overwrite_original", images[i])
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: set date of %s to %s: %s\n", images[i], s.String(), err)
			os.Exit(1)
		}
		fmt.Printf("%s %s\n", s.Format("2006-01-02 15:04:05"), images[i])
	}
}
