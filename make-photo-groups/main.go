// OBSOLETE: this program was used for a one-off task on 2021-05-01.  I'm not
// deleting it yet, in case the results are wrong and I need to revise and redo
// it.  But it can safely be deleted after that.
//
// make-photo-groups searches an album directory and generates a list of photo
// groups.  Each group consists of a single original image, and all of the
// variants and other derived files related to it.  It outputs a file named
// "groups", listing the groups.  The groups are separated by blank lines.  Each
// group has one or more lines containing a relative path from the album root.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var prefixRE = regexp.MustCompile(`^[0-9]{2,3}[-_ ]`)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: make-photo-groups directory\n")
		os.Exit(2)
	}
	if err := os.Chdir(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	paths := make(map[string][]string)
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := filepath.Base(path)
		name = strings.ReplaceAll(name, "_original", "")
	RETRY:
		if idx := strings.LastIndexByte(name, '.'); idx > 0 {
			ext := name[idx+1:]
			switch ext {
			case "xmp":
				name = name[:idx]
				goto RETRY
			case "dng", "gif", "jpg", "m4v", "mkv", "mov", "mp4", "png", "tif", "wav":
				name = name[:idx]
			default:
				return nil
			}
		} else {
			return nil
		}
		name = prefixRE.ReplaceAllLiteralString(name, "")
		if strings.HasPrefix(name, "LRE_") {
			name = name[4:]
		}
		if idx := strings.IndexByte(name, '.'); idx > 0 {
			name = name[:idx]
		}
		paths[name] = append(paths[name], path)
		return nil
	})
	names := make([]string, 0, len(paths))
	for name := range paths {
		names = append(names, name)
	}
	sort.Strings(names)
	j := 0
	for _, n := range names {
		if j > 0 && len(names[j-1]) > 4 && strings.HasPrefix(n, names[j-1]) {
			paths[names[j-1]] = append(paths[names[j-1]], paths[n]...)
			fmt.Printf("%s: merging %s into %s\n", os.Args[1], n, names[j-1])
		} else {
			names[j] = n
			j++
		}
	}
	names = names[:j]
	fh, _ := os.Create("groups")
	for _, n := range names {
		p := paths[n]
		sort.Slice(p, func(i, j int) bool {
			ic := strings.Count(p[i], "/")
			jc := strings.Count(p[j], "/")
			if ic != jc {
				return ic < jc
			}
			return p[i] < p[j]
		})
		for _, f := range p {
			fmt.Fprintln(fh, f)
		}
		fmt.Fprintln(fh)
	}
	fh.Close()
}
