// order-photos copies an ordered set of photos to a destination directory.
// Along the way, it adds sequence numbers to the start of their filenames so
// that the natural sorted order of the directory is the desired order of the
// photos.  It retains the modification times of the photos and any XMP sidecar
// files.
//
// The first argument to the program is the name of a file that defines the
// desired photo order.  Each line in this file names one photo, giving its
// filename relative to the source directory.  Blank lines and # comments are
// ignored.
//
// The second argument to the program is the source directory from which to
// copy the photos, and the third argument is the destination directory into
// which to copy them.  The destination directory will be created if it doesn't
// already exist.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var (
		srcdir  string
		destdir string
		order   *os.File
		scan    *bufio.Scanner
		line    string
		number  int
		err     error
	)
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "usage: order-photos order-file source-directory destination-directory\n")
		os.Exit(2)
	}
	if order, err = os.Open(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	scan = bufio.NewScanner(order)
	srcdir = os.Args[2]
	destdir = os.Args[3]
	if err = os.MkdirAll(destdir, 0777); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if matches, err := filepath.Glob(destdir + "/001_*"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	} else if matches != nil {
		fmt.Fprintf(os.Stderr, "WARNING: numbered images already exist in the destination directory\n")
	}
	for scan.Scan() {
		line = strings.TrimSpace(scan.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		number++
		copyImage(srcdir, line, destdir, number)
	}
	if err = scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", os.Args[1], err)
		os.Exit(1)
	}
}

func copyImage(srcdir, filename, destdir string, number int) {
	var (
		src  *os.File
		dest *os.File
		stat os.FileInfo
		err  error
	)
	if src, err = os.Open(filepath.Join(srcdir, filename)); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if stat, err = src.Stat(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: stat %s: %s\n", filename, err)
		os.Exit(1)
	}
	if dest, err = os.Create(fmt.Sprintf("%s/%03d_%s", destdir, number, filepath.Base(filename))); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if _, err = io.Copy(dest, src); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: copy %s: %s\n", filename, err)
		os.Exit(1)
	}
	src.Close()
	if err = dest.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: close %s: %s\n", filename, err)
		os.Exit(1)
	}
	if err = os.Chtimes(dest.Name(), time.Now(), stat.ModTime()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: chtimes %s: %s\n", filename, err)
		os.Exit(1)
	}
	if strings.HasSuffix(filename, ".xmp") {
		return
	}
	if idx := strings.LastIndexByte(filename, '.'); idx > 0 {
		filename = filename[:idx]
	}
	filename = filename + ".xmp"
	if _, err := os.Stat(filepath.Join(srcdir, filename)); !os.IsNotExist(err) {
		copyImage(srcdir, filename, destdir, number)
	}
}
