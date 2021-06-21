// rename-photos looks for groups of related photos in a directory tree, and
// and renames them to match my naming standards.
package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	webdialogs "github.com/rothskeller/photo-tools/web-dialogs"
)

var prefixRE = regexp.MustCompile(`^[0-9]{2,3}[-_ ]`)
var allDigitsRE = regexp.MustCompile(`^[0-9]*$`)

var groupNames []string
var groupPaths [][]string
var groupSplit [][]string
var renames map[string]string

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: rename-photos [directory]\n")
		os.Exit(2)
	}
	if len(os.Args) == 2 {
		if err := os.Chdir(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	getGroups()
	webdialogs.Main(requestInfo)
}

func getGroups() {
	var paths = make(map[string][]string)
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
	var names = make([]string, 0, len(paths))
	for name := range paths {
		names = append(names, name)
	}
	sort.Strings(names)
	j := 0
	for _, n := range names {
		if j > 0 && len(names[j-1]) > 4 && strings.HasPrefix(n, names[j-1]) && !allDigitsRE.MatchString(n[len(names[j-1]):]) {
			paths[names[j-1]] = append(paths[names[j-1]], paths[n]...)
		} else {
			names[j] = n
			j++
		}
	}
	groupNames = names[:j]
	groupPaths = make([][]string, len(groupNames))
	for i, name := range groupNames {
		groupPaths[i] = paths[name]
		sort.Strings(groupPaths[i])
	}
}
