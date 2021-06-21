// rename-photos looks for groups of related photos in a directory tree, and
// and renames them to match my naming standards.
package main

import (
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	webdialogs "github.com/rothskeller/photo-tools/web-dialogs"
)

func requestInfo(w http.ResponseWriter, r *http.Request) webdialogs.DialogHandler {
	if len(groupPaths) == 0 {
		return nil
	}
	var paths = groupPaths[0]
	var images = readImages(paths)
	var remaining = map[string]bool{}
	for _, path := range paths {
		remaining[path] = true
	}
	var groups [][]string
	for {
		var found string
		var group []string
		for _, path := range paths {
			if !isXMP(path) && remaining[path] {
				found = path
				break
			}
		}
		if found == "" {
			break
		}
		group = addToGroup(group, found, remaining)
		for _, path := range paths {
			if remaining[path] && sameImage(images[path], images[found]) {
				group = addToGroup(group, path, remaining)
			}
		}
		groups = append(groups, group)
	}
	{
		var strays []string
		for _, path := range paths {
			if remaining[path] {
				strays = append(strays, path)
			}
		}
		if len(strays) > 0 {
			groups = append(groups, strays)
		}
	}
	groupSplit = groups
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := requestTemplate.Execute(w, groups); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	return receiveInfo
}

func isXMP(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".xmp" || ext == ".xmp_original"
}

func readImages(paths []string) (imgs map[string]image.Image) {
	imgs = make(map[string]image.Image)
	for _, path := range paths {
		if fh, err := os.Open(path); err == nil {
			if img, _, err := image.Decode(fh); err == nil {
				imgs[path] = img
			}
			fh.Close()
		}
	}
	return imgs
}

func sameImage(a, b image.Image) bool {
	if a == nil || b == nil {
		return false
	}
	ab := a.Bounds()
	bb := b.Bounds()
	if !ab.Eq(bb) {
		return false
	}
	for x := ab.Min.X; x < ab.Max.X; x++ {
		for y := ab.Min.Y; y < ab.Max.Y; y++ {
			ar, ag, ab, aa := a.At(x, y).RGBA()
			br, bg, bb, ba := b.At(x, y).RGBA()
			if ar != br || ag != bg || ab != bb || aa != ba {
				return false
			}
		}
	}
	return true
}

func addToGroup(group []string, path string, remaining map[string]bool) []string {
	group = append(group, path)
	delete(remaining, path)
	if remaining[path+".xmp"] {
		group = append(group, path+".xmp")
		delete(remaining, path+".xmp")
	}
	if remaining[path+".xmp_original"] {
		group = append(group, path+".xmp_original")
		delete(remaining, path+".xmp_original")
	}
	return group
}

func getCWD() string {
	cwd, _ := os.Getwd()
	return cwd
}

var requestTemplate = template.Must(template.New("request").
	Funcs(map[string]interface{}{"isXMP": isXMP, "cwd": getCWD}).
	Parse(`<!DOCTYPE html><html><body>
  <p>Assign a variant tag to each group, and choose the representative image in each group.
  <br>Leave the variant tag empty for the group with the original.</p>
  <form method="POST">
    {{ range $idx, $group := . }}
      <div style="margin-top:1rem">
        Variant:&nbsp; <input name="v{{ $idx }}">
      </div>
      {{ range $group }}
        <div>
	  <input type="radio" name="g{{ $idx }}" value="{{.}}"{{ if isXMP . }} disabled{{ end }}>
	  <a href="file://{{ cwd }}/{{.}}" target="_blank">{{ . }}</a>
	</div>
      {{ end }}
    {{ end }}
    <div style="margin-top:1rem">
      <input type="submit">
    </div>
  </form>
</body></html>`))
