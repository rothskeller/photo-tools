// OBSOLETE: this program was used for a one-off task on 2021-05-01.  I'm not
// deleting it yet, in case the results are wrong and I need to revise and redo
// it.  But it can safely be deleted after that.
//
// gen-order examines an album directory and generates an "order" file within
// it.  The "order" file lists the basenames of the images in the album that
// have been selected for display, in display order.  For albums that have
// subdirectories, the "order" file includes the path to the image, relative to
// the album root.
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: gen-order directory\n")
		os.Exit(2)
	}
	if err := os.Chdir(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	var order []string
	dir, _ := os.Getwd()
	if _, err := os.Stat("Mylio"); err == nil && dir < "2005-09" {
		order = genOrder("Mylio")
	} else if _, err := os.Stat("exported"); err == nil {
		order = genOrder("exported")
	} else {
		// No source for an order, so make sure there isn't one.
		os.Remove("order")
		os.Exit(0)
	}
	order = resolveOrder(order)
	fh, _ := os.Create("order")
	for _, o := range order {
		fmt.Fprintln(fh, o)
	}
	fh.Close()
}

var prefixRE = regexp.MustCompile(`^[0-9]{2,3}[-_ ]`)

// genOrder builds the ordered list of basenames of selected images.
func genOrder(root string) (order []string) {
	dir, _ := os.Getwd()
	fh, err := os.Open(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s/%s: %s\n", dir, root, err)
		return nil
	}
	defer fh.Close()
	ents, err := fh.ReadDir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s/%s: %s\n", dir, root, err)
		return nil
	}
	sort.Slice(ents, func(i, j int) bool {
		return ents[i].Name() < ents[j].Name()
	})
	var seen = make(map[string]bool)
	for _, ent := range ents {
		if ent.IsDir() {
			sub := genOrder(filepath.Join(root, ent.Name()))
			for _, s := range sub {
				order = append(order, filepath.Join(ent.Name(), s))
			}
			continue
		}
		name := ent.Name()
		name = strings.ReplaceAll(name, "_original", "")
		if strings.HasSuffix(name, ".xmp") {
			continue
		}
		if idx := strings.LastIndexByte(name, '.'); idx > 0 {
			ext := name[idx+1:]
			switch ext {
			case "dng", "gif", "jpg", "m4v", "mkv", "mov", "mp4", "png", "tif", "wav":
			default:
				fmt.Fprintf(os.Stderr, "WARNING: skipping extension %s\n", ext)
				continue
			}
			name = name[:idx]
		} else {
			continue
		}
		name = prefixRE.ReplaceAllLiteralString(name, "")
		if strings.HasPrefix(name, "LRE_") {
			name = name[4:]
		}
		if seen[name] {
			continue
		}
		seen[name] = true
		order = append(order, name)
	}
	return order
}

// resolveOrder resolves each of the basenames in the ordered list to include
// the subdirectory containing the canonical version of that image.
func resolveOrder(order []string) (out []string) {
	// First, build a map from basename to resolved name.
	var nmap = make(map[string]map[string]struct{})
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if b := filepath.Base(path); b == "Mylio" || b == "exported" {
			return fs.SkipDir
		}
		if err != nil {
			return nil
		}
		name := d.Name()
		if strings.Contains(name, "_original") {
			return nil
		}
		if strings.HasSuffix(name, ".xmp") {
			return nil
		}
		if idx := strings.LastIndexByte(name, '.'); idx > 0 {
			name = name[:idx]
		} else {
			return nil
		}
		if strings.HasPrefix(name, "LRE_") {
			name = name[4:]
		}
		if nmap[name] == nil {
			nmap[name] = make(map[string]struct{})
		}
		nmap[name][filepath.Join(filepath.Dir(path), name)] = struct{}{}
		return nil
	})
	// Then, replace each item in order with its resolved name.
	for _, o := range order {
		nm := nmap[filepath.Base(o)]
		var nlist []string
		for n := range nm {
			if strings.ContainsRune(o, '/') {
				if filepath.Dir(o) != filepath.Dir(n) {
					continue
				}
			}
			nlist = append(nlist, n)
		}
		if len(nlist) == 0 {
			fmt.Fprintf(os.Stderr, "WARNING: can't resolve %s\n", o)
			out = append(out, o)
			continue
		}
		if len(nlist) > 1 {
			fmt.Fprintf(os.Stderr, "WARNING: multiple resolutions for %s\n", o)
		}
		out = append(out, nlist...)
	}
	return out
}
