package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rothskeller/photo-tools/md/operations"
)

/* FILE FORMAT

The remembered file set and targeted subset are stored in $HOME/.md.  The first
line of the file contains the result of os.Getwd when the file was stored; if
this does not match the current value of os.Getwd, the file is disregarded.

Subsequent lines contain names of files in the remembered file set, in order as
given on the command line.  (Order is important for the "copy" operation.)
Those which are in the targeted subset appear bare; those which are not are
prefixed with a pound sign (#).
*/

var dir string
var savefilename string

func init() {
	var err error

	if dir, err = os.Getwd(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: getwd: %s\n", err)
		os.Exit(1)
	}
	savefilename = filepath.Join(os.Getenv("HOME"), ".md")
}

func readMDFile() []string {
	var (
		by    []byte
		lines []string
		err   error
	)
	if by, err = os.ReadFile(savefilename); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return nil
	}
	lines = strings.Split(string(by), "\n")
	if lines[0] != dir {
		os.Remove(savefilename)
		return nil
	}
	lines = lines[1:]
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func writeMDFile(lines []string) {
	var (
		fh  *os.File
		err error
	)
	if fh, err = os.Create(savefilename); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %\n", err)
		return
	}
	defer fh.Close()
	fmt.Fprintln(fh, dir)
	for _, line := range lines {
		fmt.Fprintln(fh, line)
	}
}

func getTargetedFiles() []string {
	lines := readMDFile()
	j := 0
	for _, line := range lines {
		if line[0] != '#' {
			lines[j] = line
			j++
		}
	}
	return lines[:j]
}

func restoreFullSet() (files []string, err error) {
	files = readMDFile()
	if len(files) == 0 {
		return nil, errors.New("no remembered file set")
	}
	for i := range files {
		files[i] = strings.TrimLeft(files[i], "#")
	}
	writeMDFile(files)
	return files, nil
}

func getFirstBatch() (files []string, err error) {
	files = readMDFile()
	if len(files) == 0 {
		return nil, errors.New("no remembered file set")
	}
	batches, _ := splitIntoBatches(files)
	return selectBatch(files, batches[0])
}

func getNextBatch() (files []string, err error) {
	files = readMDFile()
	if len(files) == 0 {
		return nil, errors.New("no remembered file set")
	}
	batches, current := splitIntoBatches(files)
	if current < 0 {
		return nil, errors.New("not in batch mode")
	}
	if current == len(batches)-1 {
		return nil, errors.New("already on last batch")
	}
	return selectBatch(files, batches[current+1])
}

func getPrevBatch() (files []string, err error) {
	files = readMDFile()
	if len(files) == 0 {
		return nil, errors.New("no remembered file set")
	}
	batches, current := splitIntoBatches(files)
	if current < 0 {
		return nil, errors.New("not in batch mode")
	}
	if current == 0 {
		return nil, errors.New("already on first batch")
	}
	return selectBatch(files, batches[current-1])
}

func splitIntoBatches(files []string) (batches [][]string, current int) {
	var (
		fnames   []string
		selected = make(map[string]bool)
	)
	// Act on a copy of the files so we don't change their order.
	fnames = make([]string, len(files))
	for i := range files {
		if files[i][0] == '#' {
			fnames[i] = files[i][1:]
		} else {
			fnames[i] = files[i]
			selected[files[i]] = true
		}
		fnames[i] = strings.TrimLeft(files[i], "#")
	}
	sort.Slice(fnames, func(i, j int) bool { return batchSortFn(fnames[i], fnames[j]) })
	batches = append(batches, []string{})
	for _, fname := range fnames {
		if len(batches[len(batches)-1]) == 0 || basename(batches[len(batches)-1][0]) == basename(fname) {
			batches[len(batches)-1] = append(batches[len(batches)-1], fname)
		} else {
			batches = append(batches, []string{fname})
		}
	}
	current = -1
	for i, batch := range batches {
		for j, fname := range batch {
			if selected[fname] {
				if j == 0 && current == -1 {
					current = i
				} else if j != 0 && current == i {
					// nothing
				} else {
					return batches, -1
				}
			} else if current == i {
				return batches, -1
			}
		}
	}
	return batches, current
}

func selectBatch(files, batch []string) ([]string, error) {
	var bmap = make(map[string]bool)
	for _, fname := range batch {
		bmap[fname] = true
	}
	for i := range files {
		if files[i][0] == '#' && bmap[files[i][1:]] {
			files[i] = files[i][1:]
		} else if files[i][0] != '#' && !bmap[files[i]] {
			files[i] = "#" + files[i]
		}
	}
	writeMDFile(files)
	return batch, nil
}

func selectSubset() (selected []string, err error) {
	var (
		files  []string
		line   string
		nums   []int
		scan   *bufio.Scanner
		selmap map[string]bool
		seen   = make(map[int]bool)
	)
	files = readMDFile()
	if len(files) == 0 {
		return nil, errors.New("no remembered file set")
	}
	for i, file := range files {
		fmt.Printf("%3d %s\n", i+1, strings.TrimLeft(file, "#"))
	}
	// Repeat reading lines from stdin until we get a valid answer.
	scan = bufio.NewScanner(os.Stdin)
RETRY:
	selected = nil
	selmap = make(map[string]bool)
	fmt.Printf("Select? ")
	if !scan.Scan() {
		return nil, scan.Err()
	}
	if line = scan.Text(); line == "" {
		return getTargetedFiles(), nil
	}
	if nums = operations.ParseNumberList(line); len(nums) == 0 {
		goto RETRY
	}
	for _, num := range nums {
		if num < 1 || num > len(files) {
			fmt.Printf("ERROR: no such file number %d\n", num)
			goto RETRY
		}
		if seen[num] {
			fmt.Printf("ERROR: file number %d listed twice\n", num)
			goto RETRY
		}
		seen[num] = true
		selected = append(selected, strings.TrimLeft(files[num-1], "#"))
		selmap[strings.TrimLeft(files[num-1], "#")] = true
	}
	for i := range files {
		if files[i][0] == '#' && selmap[files[i][1:]] {
			files[i] = files[i][1:]
		} else if files[i][0] != '#' && !selmap[files[i]] {
			files[i] = "#" + files[i]
		}
	}
	writeMDFile(files)
	return selected, nil
}

// Files are sorted by basename first, then by variant (the part between the
// first dot and the extension, if any), then by directory, and finally by
// extension.  This ensures, for example, that "foo.jpg" comes before
// "foo.aaa.jpg".
func batchSortFn(a, b string) bool {
	var adir, abase, avariant, aext string
	var bdir, bbase, bvariant, bext string
	adir = filepath.Dir(a)
	bdir = filepath.Dir(b)
	abase = filepath.Base(a)
	bbase = filepath.Base(b)
	if strings.HasPrefix(abase, "#") {
	}
	if strings.HasSuffix(abase, ".xmp") {
		abase, aext = abase[:len(abase)-4], ".xmp"
	}
	if strings.HasSuffix(bbase, ".xmp") {
		bbase, bext = bbase[:len(bbase)-4], ".xmp"
	}
	if ext := filepath.Ext(abase); ext != "" {
		aext = ext + aext
		abase = abase[:len(abase)-len(ext)]
	}
	if ext := filepath.Ext(bbase); ext != "" {
		bext = ext + bext
		bbase = bbase[:len(bbase)-len(ext)]
	}
	if idx := strings.IndexByte(abase, '.'); idx >= 0 {
		abase, avariant = abase[:idx], abase[idx:]
	}
	if idx := strings.IndexByte(bbase, '.'); idx >= 0 {
		bbase, bvariant = bbase[:idx], bbase[idx:]
	}
	switch {
	case abase < bbase:
		return true
	case bbase < abase:
		return false
	case avariant < bvariant:
		return true
	case bvariant < avariant:
		return false
	case adir < bdir:
		return true
	case bdir < adir:
		return false
	default:
		return aext < bext
	}
}

// basename returns the path name, shorn of any directory information and
// anything after the first period (including the period itself).
func basename(path string) string {
	path = filepath.Base(path)
	if idx := strings.IndexByte(path, '.'); idx >= 0 {
		path = path[:idx]
	}
	return path
}
